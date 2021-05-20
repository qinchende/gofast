package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/syncx"
	"net/http"
)

func MaxReqCounts(limit int32) fst.IncHandler {
	latch := syncx.Counter{Max: limit}

	return func(w *fst.GFResponse, r *http.Request) {
		if limit <= 0 {
			return
		}

		//log.Printf("curr %d", latch.Curr)
		if latch.TryBorrow() {
			defer func() {
				if err := latch.Return(); err != nil {
					w.ErrorN(err)
					w.AbortFit()
				}
			}()
			w.NextFit(r)
		} else {
			w.ErrorF("curr request %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
			w.ResWrap.WriteHeader(http.StatusServiceUnavailable)
			w.AbortFit()
		}
	}
}
