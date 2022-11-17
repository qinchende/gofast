// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/skill/sysx"
)

// 自适应降载，前面熔断还是有一定比例的请求会通过。这里通过各种参数动态调整熔断的敏感度
// 可能参考指标：
// 1. cpu利用率 > 95%
// 2. 通过的请求中超时的比例
// 3. 处理超时的严重程度
// 4. 业务优先级
// 5. 用户优先级
func LoadShedding(kp *gate.RequestKeeper) fst.CtxHandler {
	if kp == nil || sysx.MonitorStarted == false {
		return nil
	}

	return func(c *fst.Context) {
		//rt := AllAttrs[c.RouteIdx]

		c.Next()
	}
}

//
//func LoadShedding(kp *gate.RequestKeeper) fst.CtxHandler {
//	if kp == nil || sysx.MonitorStarted == false {
//		return nil
//	}
//
//	return func(c *fst.Context) {
//		//rt := AllAttrs[c.RouteIdx]
//
//		kp.SheddingStat.Total()
//
//		shedding, err := kp.Shedding.Allow()
//		if err != nil {
//			kp.CountRouteDrop(c.RouteIdx)
//			kp.SheddingStat.Drop()
//
//			//r := c.ReqRaw
//			//logx.ErrorF("[http] load shedding, %s - %s - %s", r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
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
