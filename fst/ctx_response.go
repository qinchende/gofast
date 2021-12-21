// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"bytes"
	"errors"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/bytesconv"
	"net/http"
	"sync"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对标准 http.ResponseWriter 的包裹，加入对响应的状态管理
const (
	notWritAnyData = 0
	defaultStatus  = http.StatusOK
)

var errAlreadyRendered = errors.New("ResponseWriter: already rendered")

// 自定义 ResponseWriter, 对标准库的一层包裹处理，需要对返回的数据做缓存，做到更灵活的控制。
// 实现接口 ResponseWriter
type ResWriterWrap struct {
	http.ResponseWriter
	mu       sync.Mutex
	rendered bool
	status   int
	dataSize int
	dataBuf  *bytes.Buffer // 记录响应的数据，用于框架统一封装之后的打印信息等场景
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (w *ResWriterWrap) Reset(res http.ResponseWriter) {
	w.ResponseWriter = res
	w.status = defaultStatus    // 默认返回200 OK
	w.dataSize = notWritAnyData // 一定要初始化为-1，因为0代表已设置好返回状态
	w.dataBuf = new(bytes.Buffer)
	w.rendered = false
}

func (w *ResWriterWrap) Header() http.Header {
	return w.ResponseWriter.Header()
}

// 在没有调用 WriteHeaderNow() 之前，设置status code都是可以的，会对最终response起作用
// add by sdx 2021.08.25
// 否则：
// Gin: 只会改变这里的w.status值，而不会改变response给客户端的状态了。（这没有多大意义，GoFast做出改变）
// GoFast: 打印警告日志，不改变变量。
func (w *ResWriterWrap) WriteHeader(newStatus int) {
	// 设置不一样的状态时要做一定处理，否则啥也不做。
	if w.status != newStatus {
		if w.status == -1 {
			logx.DebugPrint("[WARNING] HTTP status %d rendered, so status %d is useless.", w.status, newStatus)
		} else {
			logx.DebugPrint("[WARNING] HTTP status %d, now change to %d.", w.status, newStatus)
			w.status = newStatus
		}
	} else {
		logx.DebugPrint("[WARNING] HTTP status %d rendered, so status %d is useless.", w.status, newStatus)
	}
}

// 最后都要通过这个函数Render所有数据
// 问题1: 是否要避免 double render?
// 答：目前不需要管这个事，调用多少次Write，就往返回流写入多少数据。double render是前段业务逻辑的问题，开发应该主动避免。
func (w *ResWriterWrap) Write(data []byte) (n int, err error) {
	n, err = w.dataBuf.Write(data)
	w.dataSize += n
	return
}

func (w *ResWriterWrap) WriteString(s string) (n int, err error) {
	n, err = w.dataBuf.WriteString(s)
	w.dataSize += n
	return
}

func (w *ResWriterWrap) Status() int {
	return w.status
}

// 数据长度
func (w *ResWriterWrap) Size() int {
	return w.dataSize
}

// 当前已写的数据内容
func (w *ResWriterWrap) WrittenBytes() []byte {
	return w.dataBuf.Bytes()
}

// 如果当前还没有response 重置当前 response 数据
func (w *ResWriterWrap) ResetResponse() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.rendered {
		return false
	}

	w.status = defaultStatus
	w.dataSize = notWritAnyData
	w.dataBuf.Reset()
	return true
}

// 如果还没有render，强制返回服务器错误，中断其它返回。否则啥也不做。
func (w *ResWriterWrap) RenderHijack(status int, body string) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.rendered {
		return nil
	}
	w.rendered = true

	w.ResponseWriter.WriteHeader(status)
	_, err = w.ResponseWriter.Write(bytesconv.StringToBytes(body))
	return
}

// Render才会真的往返回通道写数据，Render只执行一次
func (w *ResWriterWrap) RenderNow() (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 不允许 double render
	if w.rendered {
		return 0, errAlreadyRendered
	}

	w.rendered = true
	w.ResponseWriter.WriteHeader(w.status)
	n, err = w.ResponseWriter.Write(w.dataBuf.Bytes())
	return
}
