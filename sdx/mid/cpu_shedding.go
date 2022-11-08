// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/skill/httpx"
	"github.com/qinchende/gofast/skill/sysx"
	"net/http"
)

func HttpLoadShedding(kp *gate.RequestKeeper, pos uint16) fst.HttpHandler {
	if kp == nil || sysx.MonitorStarted == false {
		return nil
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			kp.SheddingStat.Total()
			shedding, err := kp.Shedding.Allow()
			if err != nil {
				kp.CountExtras(pos)
				kp.SheddingStat.Drop()
				logx.ErrorF("[http] load shedding, %s - %s - %s", r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			// 包裹ResponseWriter
			cw := &httpx.ResponseWriterWrapCode{Writer: w}
			defer func() {
				if cw.Code == http.StatusServiceUnavailable {
					shedding.Fail()
				} else {
					kp.SheddingStat.Pass()
					shedding.Pass()
				}
			}()

			// 执行后面的处理函数
			next(cw, r)
		}
	}
}

// 自适应降载，主要是CPU使用率和请求大量超时，主动断开请求，等待一段时间的冷却
// 判断高负载主要取决于两个指标（必须同时满足才能降载）：
// 1. cpu 是否过载。利用率 > 95%
// 2. 请求大量熔断，最好不分路由的普遍发生熔断。
//func LoadShedding(kp *gate.RequestKeeper) fst.CtxHandler {
//	if kp == nil || sysx.MonitorStarted == false {
//		return nil
//	}
//
//	return func(c *fst.Context) {
//		kp.SheddingStat.Total()
//
//		shedding, err := kp.Shedding.Allow()
//		if err != nil {
//			kp.CountRouteDrop(c.RouteIdx)
//			kp.SheddingStat.Drop()
//
//			r := c.ReqRaw
//			logx.ErrorF("[http] load shedding, %s - %s - %s", r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
//			c.AbortDirect(http.StatusServiceUnavailable, "LoadShedding!!!")
//			return
//		}
//
//		defer func() {
//			if c.ResWrap.Status() == http.StatusServiceUnavailable {
//				shedding.Fail()
//			} else {
//				kp.SheddingStat.Pass()
//				shedding.Pass()
//			}
//		}()
//
//		// 执行后面的处理函数
//		c.Next()
//	}
//}
