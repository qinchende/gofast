// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package httpx

import (
	"bytes"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/core/lang"
	"net/http"
	"sync"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对标准 http.ResponseWriter 的包裹，加入对响应的状态管理
const (
	defaultStatus      = 0
	errAlreadyRendered = "ResWarp: already committed. "
)

// 自定义 ResponseWriter, 对标准库的一层包裹处理，需要对返回的数据做缓存，做到更灵活的控制。
// 实现接口 ResponseWriter
// 思想：通过加锁的方式，控制不同Goroutine对内存的竞争
type ResponseWrap struct {
	http.ResponseWriter               // Raw http.ResponseWriter
	header              http.Header   // Response Header
	respLock            sync.Mutex    // render locker
	dataBuf             *bytes.Buffer // render data（指针，每次申请新的data buffer。因为fst.Context用sync.Pool）
	status              int16         // HttpStatus
	committed           bool          // 防止重复render的标记
	isTimeout           bool          // 是否是因为超时触发的render
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (w *ResponseWrap) Header() http.Header {
	if w.header == nil {
		w.header = w.ResponseWriter.Header()
	}
	return w.header
}

func (w *ResponseWrap) Status() int {
	return int(w.status)
}

func (w *ResponseWrap) IsTimeout() bool {
	return w.isTimeout
}

// 数据长度
func (w *ResponseWrap) DataSize() int {
	return w.dataBuf.Len()
}

// 当前已写的数据内容
func (w *ResponseWrap) WrittenData() []byte {
	return w.dataBuf.Bytes()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (w *ResponseWrap) Reset(res http.ResponseWriter) {
	w.header = nil
	w.committed = false
	w.isTimeout = false
	w.ResponseWriter = res
	w.status = defaultStatus
	w.dataBuf = new(bytes.Buffer) // 申请新的内存
}

// 重置返回结果（没有最终response的情况下，可以重置返回内容）
func (w *ResponseWrap) Flush() bool {
	w.respLock.Lock()
	defer w.respLock.Unlock()

	if w.committed {
		if !w.isTimeout {
			logx.Warn().Msg(errAlreadyRendered + "Can't Flush.")
		}
		return false
	}
	w.status = defaultStatus
	w.dataBuf.Reset()
	return true
}

// Gin: 只会改变这里的w.status值，而不会改变response给客户端的状态了。（这没有多大意义，GoFast做出改变）
// GoFast: 没有提交之前可以无限次的改变，最终返回最后一次设置的值
func (w *ResponseWrap) WriteHeader(newStatus int) {
	w.respLock.Lock()
	defer w.respLock.Unlock()

	if w.committed {
		if !w.isTimeout {
			logx.Warn().MsgF("%sCan't WriteHeader from %d to %d.", errAlreadyRendered, w.status, newStatus)
		}
		return
	}

	if w.status != int16(newStatus) && w.status != defaultStatus {
		logx.Warn().MsgF("Response status already %d, but now change to %d.", w.status, newStatus)
	}
	w.status = int16(newStatus)
}

// 最后都要通过这个函数Render所有数据
// 问题1: 是否要避免 double render?
// 答：目前不需要管这个事，调用多少次Write，就往返回流写入多少数据。double render是前段业务逻辑的问题，开发应该主动避免。
func (w *ResponseWrap) Write(data []byte) (n int, err error) {
	w.respLock.Lock()
	defer w.respLock.Unlock()

	if w.committed {
		if !w.isTimeout {
			logx.Warn().Msg(errAlreadyRendered + "Can't Write.")
		}
		return 0, nil
	}
	n, err = w.dataBuf.Write(data)
	return
}

func (w *ResponseWrap) WriteString(s string) (n int, err error) {
	w.respLock.Lock()
	defer w.respLock.Unlock()

	if w.committed {
		if !w.isTimeout {
			logx.Warn().MsgF(errAlreadyRendered + "Can't WriteString.")
		}
		return 0, nil
	}
	n, err = w.dataBuf.WriteString(s)
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Render才会真的往返回通道写数据，Render只执行一次
func (w *ResponseWrap) Send() (n int, err error) {
	if w.tryToCommit("Can't Send.") == false {
		return
	}
	if w.status == defaultStatus {
		w.status = http.StatusOK
	}
	n, err = w.realResp()
	w.respLock.Unlock() // 因为tryToCommit没有解锁

	if err != nil {
		logx.Trace().MsgF("realSend error: %s", err)
	}
	return
}

// 这个主要用于严重错误的时候，特殊状态的返回
// 如果还没有render，强制返回服务器错误，中断其它返回。否则啥也不做。
func (w *ResponseWrap) SendHijack(resStatus int, data []byte) (n int) {
	if w.tryToCommit("Can't Hijack.") == false {
		return
	}
	w.resetResponse(resStatus, data)
	n, err := w.realResp()
	w.respLock.Unlock() // 因为tryToCommit没有解锁

	if err != nil {
		logx.Trace().MsgF("realSend error: %s", err)
	}
	return
}

// 强制跳转
func (w *ResponseWrap) SendHijackRedirect(req *http.Request, resStatus int, redirectUrl string) {
	if w.tryToCommit("Can't Hijack Redirect.") == false {
		return
	}
	w.resetResponse(resStatus, lang.ToBytes(redirectUrl))
	http.Redirect(w, req, redirectUrl, resStatus)
	w.respLock.Unlock() // 因为tryToCommit没有解锁
}

// 超时协程调用
func (w *ResponseWrap) SendByTimeoutGoroutine(resStatus int, data []byte) bool {
	w.isTimeout = true
	if w.tryToCommit("Can't Send by timeout goroutine.") == false {
		return false
	}
	w.resetResponse(resStatus, data)
	_, err := w.realResp()
	w.respLock.Unlock() // 因为tryToCommit没有解锁

	if err != nil {
		logx.Trace().MsgF("realSend error: %s", err)
	}
	return true
}

// 打劫成功，强制改写返回结果
func (w *ResponseWrap) resetResponse(resStatus int, data []byte) {
	w.status = int16(resStatus)
	w.dataBuf.Reset()
	_, _ = w.dataBuf.Write(data)
}

// NOTE: 要避免 double render。只执行第一次Render的结果，后面的Render直接丢弃
func (w *ResponseWrap) tryToCommit(tip string) bool {
	w.respLock.Lock()
	if w.committed {
		w.respLock.Unlock()
		if !w.isTimeout {
			logx.Warn().Msg(errAlreadyRendered + tip)
		}
		return false
	}
	w.committed = true
	return true // Note: Important! 此时没有解锁，需要在调用外部解锁
}

// NOTE：调用此方法才是真正意义上的 对请求Response，之后再无法更改Response的结果
func (w *ResponseWrap) realResp() (n int, err error) {
	w.ResponseWriter.WriteHeader(int(w.status))
	n, err = w.ResponseWriter.Write(w.dataBuf.Bytes())
	return
}
