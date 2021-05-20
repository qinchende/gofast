package mid

import (
	"compress/gzip"
	"github.com/qinchende/gofast/fst"
	"net/http"
	"strings"

	"github.com/qinchende/gofast/skill/httpx"
)

func Gunzip(w *fst.GFResponse, r *http.Request) {
	if strings.Contains(r.Header.Get(httpx.ContentEncoding), "gzip") {
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			w.ResWrap.WriteHeader(http.StatusBadRequest)
			w.AbortFit()
		}
		r.Body = reader
	}
}
