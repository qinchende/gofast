package httpx

import "net/http"

// A WithCodeResponseWriter is a helper to delay sealing a http.ResponseWriter on writing code.
type ResponseWriterWrapCode struct {
	Writer http.ResponseWriter
	Code   int
}

// Flush flushes the response writer.
func (w *ResponseWriterWrapCode) Flush() {
	if flusher, ok := w.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Header returns the http header.
func (w *ResponseWriterWrapCode) Header() http.Header {
	return w.Writer.Header()
}

// Write writes bytes into w.
func (w *ResponseWriterWrapCode) Write(bytes []byte) (int, error) {
	return w.Writer.Write(bytes)
}

// WriteHeader writes code into w, and not sealing the writer.
func (w *ResponseWriterWrapCode) WriteHeader(code int) {
	w.Writer.WriteHeader(code)
	w.Code = code
}
