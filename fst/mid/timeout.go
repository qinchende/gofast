package mid

const reason = "Request Timeout"

//
//func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		if duration > 0 {
//			return http.TimeoutHandler(next, duration, reason)
//		} else {
//			return next
//		}
//	}
//}

//func TimeoutFit(dur time.Duration) fst.IncHandler {
//	return func(w http.ResponseWriter, r *fst.Request) {
//		if dur > 0 {
//			http.TimeoutHandler(next, dur, reason)
//		}
//	}
//}
