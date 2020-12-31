package proto

import (
	"bufio"
	"gofast/skill"
	"io"
	"net"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

// ResponseWriter ...
type FIResWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (Status code + headers).
	WriteHeaderNow()

	// get the http.Pusher for server push
	Pusher() http.Pusher
}

type FResWriter struct {
	http.ResponseWriter
	Size   int
	Status int
}

var _ FIResWriter = &FResWriter{}

func (w *FResWriter) Reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.Size = noWritten
	w.Status = defaultStatus
}

func (w *FResWriter) WriteHeader(code int) {
	if code > 0 && w.Status != code {
		if w.Written() {
			skill.DebugPrint("[WARNING] Headers were already written. Wanted to override Status code %d with %d", w.Status, code)
		}
		w.Status = code
	}
}

func (w *FResWriter) WriteHeaderNow() {
	if !w.Written() {
		w.Size = 0
		w.ResponseWriter.WriteHeader(w.Status)
	}
}

func (w *FResWriter) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.Size += n
	return
}

func (w *FResWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.Size += n
	return
}

func (w *FResWriter) Written() bool {
	return w.Size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *FResWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.Size < 0 {
		w.Size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// Flush implements the http.Flush interface.
func (w *FResWriter) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *FResWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
