package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

func MaxReqContentLength(limit int64) fst.IncHandler {
	return func(w *fst.GFResponse, r *http.Request) {
		if limit <= 0 {
			return
		}
		// request body length
		if r.ContentLength > limit {
			w.ErrorF("Request body limit is %d, but got %d, rejected with code %d", limit, r.ContentLength, http.StatusRequestEntityTooLarge)
			w.ResW.WriteHeader(http.StatusRequestEntityTooLarge)
			w.AbortFit()
		}
	}
}
