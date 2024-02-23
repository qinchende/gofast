package mid

import (
	"fmt"
	"github.com/qinchende/gofast/core/logx"
	"github.com/qinchende/gofast/fst"
	"net/http"
)

// 全局判断所有请求类型的最大长度
func HttpMaxContentLength(limit int64) fst.HttpHandler {
	// limit <= 0 意味着根本不检查ContentLength的限制
	if limit <= 0 {
		return nil
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// request body length
			if r.ContentLength > limit {
				logx.ErrorF("Request body limit is %d, but got %d, rejected with code %d", limit,
					r.ContentLength, http.StatusRequestEntityTooLarge)
				w.WriteHeader(http.StatusRequestEntityTooLarge)
			} else {
				next(w, r)
			}
		}
	}
}

// 限制当前路径的请求最大数据长度
func MaxContentLength(c *fst.Context) {
	rt := RoutesAttrs[c.RouteIdx]
	if rt.MaxLen <= 0 {
		return
	}

	// request body length
	if c.Req.Raw.ContentLength > rt.MaxLen {
		c.AbortDirect(http.StatusRequestEntityTooLarge, fmt.Sprintf("Request body large then %d", rt.MaxLen))
	}

	c.Next()
}
