// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/bytesconv"
	"net"
	"net/http"
	"strings"
)

// 自定义 Response
type GFResponse struct {
	ResW *ResWriteWrap
	PCtx *Context

	// 用于上下文
	gftApp *GoFast
	fitIdx int
	Errors errorMsgs // []*Error
}

func (w *GFResponse) ClientIP(r *http.Request) string {
	if w.gftApp.ForwardedByClientIP {
		clientIP := r.Header.Get("X-Forwarded-For")
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-Ip")
		}
		if clientIP != "" {
			return clientIP
		}
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func (w *GFResponse) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	w.Errors = append(w.Errors, parsedError)
	return parsedError
}

func (w *GFResponse) ErrorN(err error) {
	_ = w.Error(err)
}

func (w *GFResponse) ErrorF(format string, v ...interface{}) {
	_ = w.Error(errors.New(fmt.Sprintf(format, v...)))
}

func (w *GFResponse) AbortWithStatus(code int) {
	w.ResW.WriteHeader(code)
	w.ResW.WriteHeaderNow()
	w.AbortFit()
}

func (w *GFResponse) AbortWithError(code int, err error) *Error {
	w.AbortWithStatus(code)
	return w.Error(err)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

// 实现接口 ResponseWriter
type ResWriteWrap struct {
	http.ResponseWriter
	size       int
	status     int
	WriteBytes []byte // 最多保存两组返回结果
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

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (status code + headers).
	WriteHeaderNow()

	// get the http.Pusher for server push
	Pusher() http.Pusher
}

// 验证是否实现了接口所有的方法
var _ ResponseWriter = &ResWriteWrap{}

func (w *ResWriteWrap) Reset(res http.ResponseWriter) {
	w.ResponseWriter = res
	w.size = noWritten       // 一定要初始化为-1，因为0代表已设置好返回状态
	w.status = defaultStatus // 默认返回200 OK
}

// 在没有调用 WriteHeaderNow() 之前，设置status code都是可以的，会对最终response起作用
// 否则只会改变这里的w.status值，而不会改变response给客户端的状态了。切记。
func (w *ResWriteWrap) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			logx.DebugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

// 第一次调用起作用，后面再调用不会改变response的状态了。
func (w *ResWriteWrap) WriteHeaderNow() {
	// 还没有任何写动作就可以设置返回状态，否则啥也不做，意味着返回状态只能被设置一次
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

// 返回结果都是通过这里的两个函数处理的
func (w *ResWriteWrap) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.WriteBytes = data[:n] // 记录最后一次输出给客户端的数据
	w.size += n
	return
}

func (w *ResWriteWrap) WriteString(s string) (n int, err error) {
	return w.Write(bytesconv.StringToBytes(s))
	//w.WriteHeaderNow()
	//n, err = io.WriteString(w.ResponseWriter, s)
	//w.size += n
	//return
}

func (w *ResWriteWrap) Status() int {
	return w.status
}

func (w *ResWriteWrap) Size() int {
	return w.size
}

func (w *ResWriteWrap) Written() bool {
	// 只要不是初始化的-1，就代表已经开始写了，不管是不是只写了个返回状态
	return w.size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *ResWriteWrap) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *ResWriteWrap) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flush interface.
func (w *ResWriteWrap) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *ResWriteWrap) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
