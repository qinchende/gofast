package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/syncx"
	"net/http"
)

//// 限制最大并发连接数，相当于做一个请求资源数量连接池
//func MaxConnections(gft *fst.GoFast, limit int32) http.HandlerFunc {
//	// 并发数不做限制
//	if limit <= 0 {
//		return nil
//	}
//
//	latch := syncx.Counter{Max: limit}
//	return func(w http.ResponseWriter, r *http.Request) {
//		if latch.TryBorrow() {
//			defer func() {
//				if err := latch.Return(); err != nil {
//					//w.ErrorN(err)
//					logx.Error(err)
//					gft.AbortFit()
//				}
//			}()
//			gft.NextFit(w, r)
//		} else {
//			//w.ErrorF("curr request %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
//			logx.Errorf("curr request %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
//			w.WriteHeader(http.StatusServiceUnavailable)
//			gft.AbortFit()
//		}
//	}
//}

// 限制最大并发连接数，相当于做一个请求资源数量连接池
func MaxConnections(limit int32) fst.FitFunc {
	// 并发数不做限制
	if limit <= 0 {
		return nil
	}

	latch := syncx.Counter{Max: limit}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if latch.TryBorrow() {
				defer func() {
					if err := latch.Return(); err != nil {
						//w.ErrorN(err)
						logx.Error(err)
					}
				}()
				next(w, r)
			} else {
				//w.ErrorF("curr request %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
				logx.Errorf("curr request %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		}
	}
}
