// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/gate"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/httpx"
	"github.com/qinchende/gofast/skill/sysx"
	"net/http"
)

// 自适应降载，主要是CPU使用率和请求大量超时，主动断开请求，等待一段时间的冷却
// 判断高负载主要取决于两个指标（必须同时满足才能降载）：
// 1. cpu 是否过载。利用率 > 95%
// 2. 请求大量熔断，最好不分路由的普遍发生熔断。
func LoadShedding(kp *gate.RequestKeeper) fst.CtxHandler {
	if kp == nil || sysx.CpuChecked == false {
		return nil
	}

	return func(c *fst.Context) {
		kp.SheddingStat.Total()

		shedding, err := kp.Shedding.Allow()
		if err != nil {
			kp.CounterAddDrop(c.RouteIdx)
			kp.SheddingStat.Drop()

			logx.Errorf("[http] load shedding, %s - %s - %s", c.ReqRaw.RequestURI, httpx.GetRemoteAddr(c.ReqRaw), c.ReqRaw.UserAgent())
			c.AbortAndHijack(http.StatusServiceUnavailable, "LoadShedding!!!")
			return
		}

		defer func() {
			if c.ResWrap.Status() == http.StatusServiceUnavailable {
				shedding.Fail()
			} else {
				kp.SheddingStat.Pass()
				shedding.Pass()
			}
		}()

		// 执行后面的处理函数
		c.Next()
	}
}

// ++++++++++++++++++++++ add by cd.net 2021.10.14
// 总说：定时统计（间隔60秒）系统资源利用情况 | 请求处理相应性能 | 请求量 等
//func CpuMetric(metrics *stat.Metrics) fst.FitFunc {
//	if metrics == nil {
//		return nil
//	}
//
//	return func(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			start := time.Now()
//			defer func() {
//				metrics.AddItem(stat.ReqItem{
//					Duration: time.Now().Sub(start),
//				})
//			}()
//
//			next(w, r)
//		}
//	}
//}

//// SheddingHandler returns a middleware that does load shedding.
//func SheddingHandler(shedder load.Shedder, metrics *stat.Metrics) func(http.Handler) http.Handler {
//	if shedder == nil {
//		return func(next http.Handler) http.Handler {
//			return next
//		}
//	}
//
//	ensureSheddingStat()
//
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			sheddingStat.IncrementTotal()
//			promise, err := shedder.Allow()
//			if err != nil {
//				metrics.AddDrop()
//				sheddingStat.IncrementDrop()
//				logx.Errorf("[http] dropped, %s - %s - %s",
//					r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
//				w.WriteHeader(http.StatusServiceUnavailable)
//				return
//			}
//
//			cw := &security.WithCodeResponseWriter{Writer: w}
//			defer func() {
//				if cw.Code == http.StatusServiceUnavailable {
//					promise.Fail()
//				} else {
//					sheddingStat.IncrementPass()
//					promise.Pass()
//				}
//			}()
//			next.ServeHTTP(cw, r)
//		})
//	}
//}
