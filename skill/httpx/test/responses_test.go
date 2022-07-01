package test

import (
	"errors"
	"github.com/qinchende/gofast/skill/httpx"
	"net/http"
	"strings"
	"testing"

	"github.com/qinchende/gofast/logx"
	"github.com/stretchr/testify/assert"
)

type message struct {
	Name string `json:"name"`
}

func init() {
	logx.Disable()
}

func TestError(t *testing.T) {
	const (
		body        = "foo"
		wrappedBody = `"foo"`
	)

	tests := []struct {
		name         string
		input        string
		errorHandler func(error) (int, any)
		expectBody   string
		expectCode   int
	}{
		{
			name:       "default error handler",
			input:      body,
			expectBody: body,
			expectCode: http.StatusBadRequest,
		},
		{
			name:  "customized error handler return string",
			input: body,
			errorHandler: func(err error) (int, any) {
				return http.StatusForbidden, err.Error()
			},
			expectBody: wrappedBody,
			expectCode: http.StatusForbidden,
		},
		{
			name:  "customized error handler return error",
			input: body,
			errorHandler: func(err error) (int, any) {
				return http.StatusForbidden, err
			},
			expectBody: body,
			expectCode: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := tracedResponseWriter{
				headers: make(map[string][]string),
			}
			if test.errorHandler != nil {
				//httpx.lock.RLock()
				//prev := httpx.errorHandler
				//httpx.lock.RUnlock()
				//httpx.SetErrorHandler(test.errorHandler)
				//defer func() {
				//	httpx.lock.Lock()
				//	httpx.errorHandler = prev
				//	httpx.lock.Unlock()
				//}()
			}
			httpx.Error(&w, errors.New(test.input))
			assert.Equal(t, test.expectCode, w.code)
			assert.Equal(t, test.expectBody, strings.TrimSpace(w.builder.String()))
		})
	}
}

func TestOk(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	httpx.Ok(&w)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestOkJson(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := message{Name: "anyone"}
	httpx.OkJson(&w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "{\"name\":\"anyone\"}", w.builder.String())
}

func TestWriteJsonTimeout(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		timeout: true,
	}
	msg := message{Name: "anyone"}
	httpx.WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteJsonLessWritten(t *testing.T) {
	w := tracedResponseWriter{
		headers:     make(map[string][]string),
		lessWritten: true,
	}
	msg := message{Name: "anyone"}
	httpx.WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

type tracedResponseWriter struct {
	headers     map[string][]string
	builder     strings.Builder
	code        int
	lessWritten bool
	timeout     bool
}

func (w *tracedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *tracedResponseWriter) Write(bytes []byte) (n int, err error) {
	if w.timeout {
		return 0, http.ErrHandlerTimeout
	}

	n, err = w.builder.Write(bytes)
	if w.lessWritten {
		n -= 1
	}
	return
}

func (w *tracedResponseWriter) WriteHeader(code int) {
	w.code = code
}
