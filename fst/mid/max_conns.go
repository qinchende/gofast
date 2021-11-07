package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

// 限制最大并发连接数，相当于做一个请求资源数量连接池
func MaxConnections(limit int32) fst.IncHandler {
	// 并发数不做限制
	if limit <= 0 {
		return nil
	}

	//latch := syncx.Counter{Max: limit}
	return func(w http.ResponseWriter, r *http.Request) {
		//if latch.TryBorrow() {
		//	defer func() {
		//		if err := latch.Return(); err != nil {
		//			w.ErrorN(err)
		//			w.AbortFit()
		//		}
		//	}()
		//	w.NextFit(r)
		//} else {
		//	w.ErrorF("curr request %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
		//	w.ResWrap.WriteHeader(http.StatusServiceUnavailable)
		//	w.AbortFit()
		//}
	}
}
