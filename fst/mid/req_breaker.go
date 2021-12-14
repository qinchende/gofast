// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"fmt"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/gate"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/httpx"
	"net/http"
)

//const breakerSeparator = "://"

//// 熔断器，针对不同路由统计
//func Breaker(method, path string, metrics *stat.Metrics) func(http.Handler) http.Handler {
//	brk := breaker.NewBreaker(breaker.WithName(strings.Join([]string{method, path}, breakerSeparator)))
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			promise, err := brk.Allow()
//			if err != nil && metrics != nil {
//				metrics.AddDrop()
//				logx.Errorf("[http] dropped, %s - %s - %s",
//					r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
//				w.WriteHeader(http.StatusServiceUnavailable)
//				return
//			}
//
//			cw := &security.WithCodeResponseWriter{Writer: w}
//			defer func() {
//				if cw.Code < http.StatusInternalServerError {
//					promise.Accept()
//				} else {
//					promise.Reject(fmt.Sprintf("%d %s", cw.Code, http.StatusText(cw.Code)))
//				}
//			}()
//			next.ServeHTTP(cw, r)
//		})
//	}
//}

// 熔断器，针对不同路由统计
//func Breaker(method, path string, metrics *stat.Metrics) func(http.Handler) http.Handler {
//	brk := breaker.NewBreaker(breaker.WithName(strings.Join([]string{method, path}, breakerSeparator)))
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		})
//	}
//}

// 请求分析，针对不同路由分别执行熔断策略
func Breaker(kp *gate.RequestKeeper) fst.CtxHandler {
	if kp == nil {
		return nil
	}
	return func(c *fst.Context) {
		brk := kp.Breakers[c.RouteIdx]

		// 检查是否允许本次访问通过，主要是滑动窗口判断是否达到熔断条件
		promise, err := brk.Allow()
		if err != nil {
			kp.AddDrop(c.RouteIdx)
			logx.Errorf("[http] dropped, %s - %s - %s", c.ReqRaw.RequestURI, httpx.GetRemoteAddr(c.ReqRaw), c.ReqRaw.UserAgent())
			c.ResWrap.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		defer func() {
			code := c.ResWrap.Status()
			// 5xx 以下的错误被认为是正常返回。否认就是服务器错误，被认定是拒绝服务
			if code < http.StatusInternalServerError {
				promise.Accept()
			} else {
				promise.Reject(fmt.Sprintf("%d %s", code, http.StatusText(code)))
			}
		}()

		// 执行后面的处理函数
		c.Next()
	}
}
