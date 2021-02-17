package mid

import (
	"compress/gzip"
	"github.com/qinchende/gofast/fst"
	"net/http"
	"strings"

	"github.com/qinchende/gofast/skill/httpx"
)

func GunzipFit(w http.ResponseWriter, r *fst.Request) {
	if strings.Contains(r.RawReq.Header.Get(httpx.ContentEncoding), "gzip") {
		reader, err := gzip.NewReader(r.RawReq.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			r.Abort()
		}

		r.RawReq.Body = reader
	}
}
