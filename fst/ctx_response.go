// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/qinchende/gofast/logx"
	"net/http"
	"sync"
)

// 自定义接口 ResponseWriter
// 我们自己定义的 ResWriterWrap 结构需要实现这个接口
type ResponseWriter interface {
	http.ResponseWriter
	//http.Hijacker
	//http.Flusher
	//http.CloseNotifier

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// Returns true if the response body was already written.
	//WriteStarted() bool

	// Forces to write the http header (status code + headers).
	//WriteHeaderNow()
	WriteHeader(int)

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// get the http.Pusher for server push
	//Pusher() http.Pusher
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对标准 http.ResponseWriter 的包裹，加入对响应的状态管理
//const (
//	notWritAnyData = -1
//	defaultStatus  = http.StatusOK
//)

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

// 验证是否实现了接口所有的方法
var _ ResponseWriter = &ResWriterWrap{}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (w *ResWriterWrap) Reset(res http.ResponseWriter) {
	w.ResponseWriter = res
	w.status = -1   // 默认返回200 OK
	w.dataSize = -1 // 一定要初始化为-1，因为0代表已设置好返回状态
	w.dataBuf = new(bytes.Buffer)
	w.rendered = false
}

//// 只要不是初始化的-1，就代表已经开始写了，不管是不是只写了个返回状态
//func (w *ResWriterWrap) WriteStarted() bool {
//	return w.dataSize != notWritAnyData
//}

// 在没有调用 WriteHeaderNow() 之前，设置status code都是可以的，会对最终response起作用
// add by sdx 2021.08.25
// 否则：
// Gin: 只会改变这里的w.status值，而不会改变response给客户端的状态了。（这没有多大意义，GoFast做出改变）
// GoFast: 打印警告日志，不改变变量。
func (w *ResWriterWrap) WriteHeader(newStatus int) {
	checkWriteHeaderCode(newStatus)

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

// 重置当前缓存中写入的数据
func (w *ResWriterWrap) ResetData() {
	w.dataBuf.Reset()
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
	if w.status == -1 {
		w.status = http.StatusOK
	}
	w.ResponseWriter.WriteHeader(w.status)
	n, err = w.ResponseWriter.Write(w.dataBuf.Bytes())
	return
}

//// 第一次调用起作用，后面再调用不会改变response的状态了。
//func (w *ResWriterWrap) writeHeaderNow() {
//	if w.status == -1 {
//		w.status = http.StatusOK
//	}
//	w.ResponseWriter.WriteHeader(w.status) // 这是标准库的 WriteHeader，不是上面我们自定义的方法
//	//
//	//// 还没有任何写动作就可以设置返回状态，否则啥也不做，意味着返回状态只能被设置一次
//	//if w.dataSize == -1 {
//	//	// dataSize == 0 表示写已经准备开始写数据了
//	//	w.dataSize = 0
//	//	// 往底层写返回状态码，写入便不可改变了。
//	//	w.ResponseWriter.WriteHeader(w.status) // 这是标准库的 WriteHeader，不是上面我们自定义的方法
//	//}
//}

// copy from ./src/net/http/server.go
func checkWriteHeaderCode(code int) {
	// Issue 22880: require valid WriteHeader status codes.
	// For now we only enforce that it's three digits.
	// In the future we might block things over 599 (600 and above aren't defined
	// at https://httpwg.org/specs/rfc7231.html#status.codes)
	// and we might block under 200 (once we have more mature 1xx support).
	// But for now any three digits.
	//
	// We used to send "HTTP/1.1 000 0" on the wire in responses but there's
	// no equivalent bogus thing we can realistically send in HTTP/2,
	// so we'll consistently panic instead and help people find their bugs
	// early. (We can't return an error from WriteHeader even if we wanted to.)
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Hijack implements the http.Hijacker interface.
//func (w *ResWriterWrap) Hijack() (net.Conn, *bufio.ReadWriter, error) {
//	if w.dataSize < 0 {
//		w.dataSize = 0
//	}
//	return w.ResponseWriter.(http.Hijacker).Hijack()
//}

//// CloseNotify implements the http.CloseNotify interface.
//func (w *ResWriterWrap) CloseNotify() <-chan bool {
//	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
//}

//// Flush implements the http.Flush interface.
//func (w *ResWriterWrap) Flush() {
//	w.WriteHeaderNow()
//	w.ResponseWriter.(http.Flusher).Flush()
//}

//func (w *ResWriterWrap) Pusher() (pusher http.Pusher) {
//	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
//		return pusher
//	}
//	return nil
//}
