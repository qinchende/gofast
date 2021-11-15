package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

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
			ctx.AbortChain()
		}
	}
}
