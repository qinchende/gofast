// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
)

// ++++++++++++++++++++++ add by cd.net 2021.10.14
// 总说：定时统计（间隔60秒）系统资源利用情况 | 请求处理相应性能 | 请求量 等
func TimeMetric(kp *gate.RequestKeeper) fst.CtxHandler {
	return func(c *fst.Context) {
		defer func() {
			rt := RoutesAttrs[c.RouteIdx]

			// 无论是否panic，在统计访问量的模块，本次都算一次正常触达请求，并统计耗时
			diffMS := int32(timex.SdxNowDiffMS(c.EnterTime))
			kp.LimiterFinished(c.RouteIdx, diffMS, rt.TimeoutMS)
			kp.CountRoutePass2(c.RouteIdx, diffMS)
		}()

		c.Next()
	}
}
