// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/skill/timex"
)

// ++++++++++++++++++++++ add by cd.net 2021.10.14
// 总说：定时统计（间隔60秒）系统资源利用情况 | 请求处理相应性能 | 请求量 等
func TimeMetric(kp *gate.RequestKeeper) fst.CtxHandler {
	if kp == nil {
		return nil
	}

	return func(c *fst.Context) {
		defer func() {
			// 统计耗时
			kp.CountRoutePass(c.RouteIdx, int32(timex.SinceMS(c.EnterTime)))
		}()

		c.Next()
	}
}
