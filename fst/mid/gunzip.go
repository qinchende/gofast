package mid

import (
	"compress/gzip"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/httpx"
	"net/http"
	"strings"
)

//func Gunzip(w *fst.GFResponse, r *http.Request) {
//	if strings.Contains(r.Header.Get(httpx.ContentEncoding), "gzip") {
//		reader, err := gzip.NewReader(r.Body)
//		if err != nil {
//			w.ResWrap.WriteHeader(http.StatusBadRequest)
//			w.AbortFit()
//		}
//		r.Body = reader
//	}
//}

func Gunzip(ctx *fst.Context) {
	if strings.Contains(ctx.ReqRaw.Header.Get(httpx.ContentEncoding), "gzip") {
		reader, err := gzip.NewReader(ctx.ReqRaw.Body)
		if err != nil {
			ctx.ResWrap.WriteHeader(http.StatusBadRequest)
			ctx.Abort()
		}
		ctx.ReqRaw.Body = reader
	}
}
