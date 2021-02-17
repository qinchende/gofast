// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"bufio"
	"github.com/qinchende/gofast/logx"
	"io"
	"net"
	"net/http"
)

type Response struct {
}

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

type ResWriteWrap struct {
	http.ResponseWriter
	size   int
	status int
}

// 自定义接口FResponseWriter
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
	w.size = noWritten
	w.status = defaultStatus
}

func (w *ResWriteWrap) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			logx.DebugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

func (w *ResWriteWrap) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *ResWriteWrap) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

func (w *ResWriteWrap) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}

func (w *ResWriteWrap) Status() int {
	return w.status
}

func (w *ResWriteWrap) Size() int {
	return w.size
}

func (w *ResWriteWrap) Written() bool {
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
