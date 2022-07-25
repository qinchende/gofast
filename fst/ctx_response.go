// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"bytes"
	"github.com/qinchende/gofast/logx"
	"net/http"
	"sync"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对标准 http.ResponseWriter 的包裹，加入对响应的状态管理
// var errAlreadyRendered = errors.New("ResponseWrap: already send")
const (
	defaultStatus      = 0
	errAlreadyRendered = "ResponseWrap: already committed. "
)

// 自定义 ResponseWriter, 对标准库的一层包裹处理，需要对返回的数据做缓存，做到更灵活的控制。
// 实现接口 ResponseWriter
type ResponseWrap struct {
	http.ResponseWriter

	mu        sync.Mutex
	status    int
	dataBuf   *bytes.Buffer // 记录响应的数据，用于框架统一封装之后的打印信息等场景
	committed bool
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (w *ResponseWrap) Reset(res http.ResponseWriter) {
	w.committed = false
	w.ResponseWriter = res
	w.status = defaultStatus
	w.dataBuf = new(bytes.Buffer)
}

// TODO：这是个问题，如何重置已经被写入的 Header 值
func (w *ResponseWrap) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Gin: 只会改变这里的w.status值，而不会改变response给客户端的状态了。（这没有多大意义，GoFast做出改变）
// GoFast: 打印警告日志，不改变变量。也就是只能被调用一次。
func (w *ResponseWrap) WriteHeader(newStatus int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.committed {
		logx.WarnF("Response status %d committed, the status %d is useless.", w.status, newStatus)
		return
	}

	if w.status <= defaultStatus {
		w.status = newStatus
	} else {
		logx.WarnF("Response status already %d, can't change to %d.", w.status, newStatus)
	}
}

// 最后都要通过这个函数Render所有数据
// 问题1: 是否要避免 double render?
// 答：目前不需要管这个事，调用多少次Write，就往返回流写入多少数据。double render是前段业务逻辑的问题，开发应该主动避免。
func (w *ResponseWrap) Write(data []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.committed {
		logx.Warn(errAlreadyRendered, "Can't Write.")
		return 0, nil
	}
	n, err = w.dataBuf.Write(data)
	return
}

func (w *ResponseWrap) WriteString(s string) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.committed {
		logx.Warn(errAlreadyRendered, "Can't WriteString.")
		return 0, nil
	}
	n, err = w.dataBuf.WriteString(s)
	return
}

func (w *ResponseWrap) Status() int {
	return w.status
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
// Render才会真的往返回通道写数据，Render只执行一次
func (w *ResponseWrap) Send() (n int, err error) {
	w.mu.Lock()
	// 不允许 double render
	if w.committed {
		logx.Warn(errAlreadyRendered, "Can't Send.")
		w.mu.Unlock()
		return 0, nil
	}
	w.committed = true
	w.mu.Unlock()

	if w.status == defaultStatus {
		w.status = http.StatusOK
	}
	w.ResponseWriter.WriteHeader(w.status)
	n, err = w.ResponseWriter.Write(w.dataBuf.Bytes())
	return
}

// 这个主要用于严重错误的时候，特殊状态的返回
// 如果还没有render，强制返回服务器错误，中断其它返回。否则啥也不做。
func (w *ResponseWrap) SendHijack(resStatus int, data []byte) (n int) {
	w.mu.Lock()
	// 已经render，无法打劫，啥也不做
	if w.committed {
		w.mu.Unlock()
		logx.Warn(errAlreadyRendered, "Can't Hijack.")
		return 0
	}
	w.committed = true
	w.mu.Unlock()

	// 打劫成功，强制改写返回结果
	w.status = resStatus
	w.dataBuf.Reset()
	_, _ = w.dataBuf.Write(data)

	w.ResponseWriter.WriteHeader(w.status)
	n, err := w.ResponseWriter.Write(w.dataBuf.Bytes())
	if err != nil {
		logx.ErrorStackF("SendHijack ResponseWriter error: %s", err)
	}
	return
}
