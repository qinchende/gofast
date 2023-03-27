package fst

import (
	"net/http"
)

type RequestWrap struct {
	Raw *http.Request
}

func (rw *RequestWrap) Reset(r *http.Request) {
	rw.Raw = r
}
