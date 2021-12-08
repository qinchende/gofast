package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"net/http"
)

// 全局判断所有请求类型的最大长度
func MaxContentData(limit int64) fst.FitFunc {
	// limit <= 0 意味着根本不检查ContentLength的限制
	if limit <= 0 {
		return nil
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// request body length
			if r.ContentLength > limit {
				logx.Errorf("Request body limit is %d, but got %d, rejected with code %d", limit,
					r.ContentLength, http.StatusRequestEntityTooLarge)
				w.WriteHeader(http.StatusRequestEntityTooLarge)
			} else {
				next(w, r)
			}
		}
	}
}

// 限制当前路径的请求最大数据长度
func MaxContentLength(ctx *fst.Context) {
	rt := RConfigs[ctx.RouteIdx]
	if rt.MaxContentLen <= 0 {
		return
	}

	// request body length
	if ctx.ReqRaw.ContentLength > rt.MaxContentLen {
		ctx.ErrorF("Request body limit is %d, but got %d, rejected with code %d", rt.MaxContentLen,
			ctx.ReqRaw.ContentLength, http.StatusRequestEntityTooLarge)
		ctx.ResWrap.WriteHeader(http.StatusRequestEntityTooLarge)
		ctx.AbortChain()
	}
}
