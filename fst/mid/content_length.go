package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

//func MaxContentLength(limit int64) fst.IncHandler {
//	// limit <= 0 意味着根本不检查ContentLength的限制
//	if limit <= 0 {
//		return nil
//	}
//
//	return func(w *fst.GFResponse, r *http.Request) {
//		// request body length
//		if r.ContentLength > limit {
//			w.ErrorF("Request body limit is %d, but got %d, rejected with code %d", limit,
//				r.ContentLength, http.StatusRequestEntityTooLarge)
//			w.ResWrap.WriteHeader(http.StatusRequestEntityTooLarge)
//			w.AbortFit()
//		}
//	}
//}

func MaxContentLength(limit int64) fst.CtxHandler {
	// limit <= 0 意味着根本不检查ContentLength的限制
	if limit <= 0 {
		return nil
	}

	return func(ctx *fst.Context) {
		// request body length
		if ctx.ReqRaw.ContentLength > limit {
			ctx.ErrorF("Request body limit is %d, but got %d, rejected with code %d", limit,
				ctx.ReqRaw.ContentLength, http.StatusRequestEntityTooLarge)
			ctx.ResWrap.WriteHeader(http.StatusRequestEntityTooLarge)
			ctx.Abort()
		}
	}
}
