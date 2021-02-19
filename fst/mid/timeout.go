package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
	"time"
)

const reason = "Request Timeout"

// 单位是毫秒
func ReqTimeout(dur time.Duration, handler http.Handler) fst.IncHandler {
	return func(w *fst.GFResponse, r *http.Request) {
		if dur > 0 {
			http.TimeoutHandler(handler, dur, reason)
		}
	}
}
