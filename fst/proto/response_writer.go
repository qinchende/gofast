// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package proto

import (
	"bufio"
	"github.com/qinchende/gofast/skill"
	"io"
	"net"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

// 自定义接口FResponseWriter
type FResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (Status code + headers).
	WriteHeaderNow()

	// get the http.Pusher for server push
	Pusher() http.Pusher
}

type InnerResWrite struct {
	http.ResponseWriter
	Size   int
	Status int
}
// 验证是否实现了接口所有的方法
var _ FResponseWriter = &InnerResWrite{}

func (w *InnerResWrite) Reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.Size = noWritten
	w.Status = defaultStatus
}

func (w *InnerResWrite) WriteHeader(code int) {
	if code > 0 && w.Status != code {
		if w.Written() {
			skill.DebugPrint("[WARNING] Headers were already written. Wanted to override Status code %d with %d", w.Status, code)
		}
		w.Status = code
	}
}

func (w *InnerResWrite) WriteHeaderNow() {
	if !w.Written() {
		w.Size = 0
		w.ResponseWriter.WriteHeader(w.Status)
	}
}

func (w *InnerResWrite) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.Size += n
	return
}

func (w *InnerResWrite) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.Size += n
	return
}

func (w *InnerResWrite) Written() bool {
	return w.Size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *InnerResWrite) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.Size < 0 {
		w.Size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *InnerResWrite) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flush interface.
func (w *InnerResWrite) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *InnerResWrite) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
