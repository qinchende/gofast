package httpx

import (
	"errors"
	"mime/multipart"
	"net/http"
)

type RequestWrap struct {
	Raw *http.Request
	//parsed bool
}

func (rw *RequestWrap) Reset(r *http.Request) {
	rw.Raw = r
}

func (rw *RequestWrap) FormFile(name string) (*multipart.FileHeader, error) {
	mForm := rw.Raw.MultipartForm
	if mForm != nil && mForm != multipartByReader && len(mForm.File) > 0 {
		if fhs := mForm.File[name]; len(fhs) > 0 {
			return fhs[0], nil
		}
	}
	return nil, errors.New("http: no such file")
}
