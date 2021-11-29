// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/gate"
	"github.com/qinchende/gofast/skill/stat"
	"github.com/qinchende/gofast/skill/timex"
	"net/http"
	"time"
)

// ++++++++++++++++++++++ add by cd.net 2021.10.14
// 总说：定时统计（间隔60秒）系统资源利用情况 | 请求处理相应性能 | 请求量 等
func CpuMetric(metrics *stat.Metrics) fst.FitFunc {
	if metrics == nil {
		return nil
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				metrics.AddItem(stat.ReqItem{
					Duration: time.Now().Sub(start),
				})
			}()

			next(w, r)
		}
	}
}

// 自适应降载，主要是CPU实用率和最大并发数超过一定阈值，主动断开请求，等待一段时间的冷却
func LoadShedding(kp *gate.RequestKeeper) fst.CtxHandler {
	if kp == nil {
		return nil
	}

	return func(c *fst.Context) {
		c.Next()

		kp.AddItem(gate.ReqItem{
			RouterIdx: c.RouteID,
			Duration:  timex.Since(c.EnterTime),
		})
	}
}
