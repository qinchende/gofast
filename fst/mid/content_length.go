package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

func MaxReqContentLength(n int64) fst.IncHandler {
	return func(w http.ResponseWriter, r *fst.Request) {
		if n <= 0 {
			return
		}
		if r.RawReq.ContentLength > n {
			r.Errorf("Request entity too large, limit is %d, but got %d, rejected with code %d",
				n, r.RawReq.ContentLength, http.StatusRequestEntityTooLarge)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			r.Abort()
		}
	}
}
