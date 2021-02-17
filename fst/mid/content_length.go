package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

//
//func MaxBytesHandler(n int64) func(http.Handler) http.Handler {
//	if n <= 0 {
//		return func(next http.Handler) http.Handler {
//			return next
//		}
//	}
//
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			if r.ContentLength > n {
//				//internal.Errorf(r, "request entity too large, limit is %d, but got %d, rejected with code %d",
//				//	n, r.ContentLength, http.StatusRequestEntityTooLarge)
//				w.WriteHeader(http.StatusRequestEntityTooLarge)
//			} else {
//				next.ServeHTTP(w, r)
//			}
//		})
//	}
//}

func MaxReqContentLen(n int64) fst.IncHandler {
	return func(w http.ResponseWriter, r *fst.Request) {
		if n <= 0 {
			return
		}
		if r.RawReq.ContentLength > n {
			//internal.Errorf(r, "Request entity too large, limit is %d, but got %d, rejected with code %d",
			//	n, r.RawReq.ContentLength, http.StatusRequestEntityTooLarge)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			r.Abort()
		}
	}
}
