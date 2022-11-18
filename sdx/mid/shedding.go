// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"net/http"
)

// 自适应降载，前面熔断还是有一定比例的请求会通过。这里通过各种参数动态调整熔断的敏感度
// 可能参考指标：
// 1. cpu利用率 > 95%
// 2. 通过的请求中超时的比例
// 3. 处理超时的严重程度
// 4. 业务优先级
// 5. 用户优先级
func LoadShedding(kp *gate.RequestKeeper, idx uint16) fst.CtxHandler {
	return func(c *fst.Context) {
		rt := AllAttrs[c.RouteIdx]

		if kp.Shedding(c.RouteIdx, rt.TimeoutMS) {
			kp.CountRouteDrop(c.RouteIdx)
			kp.CountExtras(idx) // Just for debug
			c.AbortDirect(http.StatusServiceUnavailable, midSheddingBody)
			return
		}

		kp.LimiterIncome(c.RouteIdx)
		c.Next()
	}
}
