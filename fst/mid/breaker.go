package mid

import (
	"fmt"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/httpx"
	"github.com/qinchende/gofast/skill/security"
	"github.com/qinchende/gofast/skill/stat"
	"net/http"
	"strings"
)

const breakerSeparator = "://"

//func BreakerDoor() http.HandlerFunc {
//	return func(w *fst.GFResponse, r *http.Request) {
//
//		brk := breaker.NewBreaker(breaker.WithName(strings.Join([]string{r.Method, r.URL.Path}, breakerSeparator)))
//
//		promise, _ := brk.Allow()
//		//promise, err := brk.Allow()
//		//if err != nil && metrics != nil {
//		//	metrics.AddDrop()
//		//	logx.Errorf("[http] dropped, %s - %s - %s",
//		//		r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
//		//	w.WriteHeader(http.StatusServiceUnavailable)
//		//	return
//		//}
//
//		cw := &security.WithCodeResponseWriter{Writer: w.ResWrap}
//		defer func() {
//			if cw.Code < http.StatusInternalServerError {
//				promise.Accept()
//			} else {
//				promise.Reject(fmt.Sprintf("%d %s", cw.Code, http.StatusText(cw.Code)))
//			}
//		}()
//		//next.ServeHTTP(cw, r)
//	}
//}

func BreakerHandler(method, path string, metrics *stat.Metrics) func(http.Handler) http.Handler {
	brk := breaker.NewBreaker(breaker.WithName(strings.Join([]string{method, path}, breakerSeparator)))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			promise, err := brk.Allow()
			if err != nil && metrics != nil {
				metrics.AddDrop()
				logx.Errorf("[http] dropped, %s - %s - %s",
					r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			cw := &security.WithCodeResponseWriter{Writer: w}
			defer func() {
				if cw.Code < http.StatusInternalServerError {
					promise.Accept()
				} else {
					promise.Reject(fmt.Sprintf("%d %s", cw.Code, http.StatusText(cw.Code)))
				}
			}()
			next.ServeHTTP(cw, r)
		})
	}
}
