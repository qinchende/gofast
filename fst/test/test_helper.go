package test

import (
	"net/http"
	"net/http/httptest"
)

type header struct {
	Key   string
	Value string
}

func ExecRequest(app http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	return w
}

type CustomerResWriter struct{}

func (m *CustomerResWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *CustomerResWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *CustomerResWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *CustomerResWriter) WriteHeader(int) {}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
