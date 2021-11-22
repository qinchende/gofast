// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"bufio"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/bytesconv"
	"net"
	"net/http"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对标准 http.ResponseWriter 的包裹，加入对响应的状态管理
const (
	notWritAnyData = -1
	defaultStatus  = http.StatusOK
)

// 自定义 ResponseWriter, 对标准库的一层包裹处理
// 实现接口 ResponseWriter
type ResWriterWrap struct {
	http.ResponseWriter
	size       int
	status     int
	WriteBytes []byte // 记录响应的数据
}

// 自定义接口 ResponseWriter
// 我们自己定义的 GFResponse 结构需要实现这个接口
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// Returns true if the response body was already written.
	WriteStarted() bool

	// Forces to write the http header (status code + headers).
	WriteHeaderNow()

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// get the http.Pusher for server push
	Pusher() http.Pusher
}

// 验证是否实现了接口所有的方法
var _ ResponseWriter = &ResWriterWrap{}

func (w *ResWriterWrap) Reset(res http.ResponseWriter) {
	w.ResponseWriter = res
	w.size = notWritAnyData  // 一定要初始化为-1，因为0代表已设置好返回状态
	w.status = defaultStatus // 默认返回200 OK
}

func (w *ResWriterWrap) WriteStarted() bool {
	// 只要不是初始化的-1，就代表已经开始写了，不管是不是只写了个返回状态
	return w.size != notWritAnyData
}

// 在没有调用 WriteHeaderNow() 之前，设置status code都是可以的，会对最终response起作用
// add by sdx 2021.08.25
// 否则：
// Gin: 只会改变这里的w.status值，而不会改变response给客户端的状态了。（这没有多大意义，GoFast做出改变）
// GoFast: 打印警告日志，不改变变量。
func (w *ResWriterWrap) WriteHeader(newStatus int) {
	if newStatus > 0 && w.status != newStatus {
		if w.WriteStarted() {
			logx.DebugPrint("[WARNING] HTTP status %d rendered, so status %d is useless.", w.status, newStatus)
		} else {
			logx.DebugPrint("[WARNING] HTTP status %d, now change to %d.", w.status, newStatus)
			w.status = newStatus
		}
	}
}

// 第一次调用起作用，后面再调用不会改变response的状态了。
func (w *ResWriterWrap) WriteHeaderNow() {
	// 还没有任何写动作就可以设置返回状态，否则啥也不做，意味着返回状态只能被设置一次
	if !w.WriteStarted() {
		// size == 0 表示写已经准备开始写数据了
		w.size = 0
		// 往底层写返回状态码，写入便不可改变了。
		w.ResponseWriter.WriteHeader(w.status) // 这是标准库的 WriteHeader，不是上面我们自定义的方法
	}
}

// 最后都要通过这个函数Render所有数据
// 问题1: 是否要避免 double render?
// 答：目前不需要管这个事，调用多少次Write，就往返回流写入多少数据。double render是前段业务逻辑的问题，开发应该主动避免。
func (w *ResWriterWrap) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.WriteBytes = data[:n] // 多次Write的情况下，暂时只记录最后一次输出给客户端的数据
	w.size += n
	return
}

// bytesconv.StringToBytes 高性能字符串和字节切片的转换
func (w *ResWriterWrap) WriteString(s string) (n int, err error) {
	return w.Write(bytesconv.StringToBytes(s))
}

func (w *ResWriterWrap) Status() int {
	return w.status
}

func (w *ResWriterWrap) Size() int {
	return w.size
}

// Hijack implements the http.Hijacker interface.
func (w *ResWriterWrap) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *ResWriterWrap) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flush interface.
func (w *ResWriterWrap) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *ResWriterWrap) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
