package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
	"time"
)

//TODO: 这种用法在这里有些问题，有待进一步研究
const timeoutReason = "[GoFast]Request Timeout."

// 单位是毫秒
func ReqTimeout(dur time.Duration) fst.IncHandler {
	return func(wGF *fst.GFResponse, r *http.Request) {
		if dur <= 0 {
			return
		}

		hp := &FitHelper{}
		hp.nextHandler = func(w http.ResponseWriter, r *http.Request) {
			wGF.NextFit(r)
		}

		handler := http.TimeoutHandler(hp, dur, timeoutReason)
		handler.ServeHTTP(wGF.ResW, r)
	}
}
