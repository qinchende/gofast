package fst

import (
	"github.com/qinchende/gofast/fst/door"
	"github.com/qinchende/gofast/skill/stat"
	"github.com/qinchende/gofast/skill/timex"
)

// 创建日志模板
func (gft *GoFast) NewMetricsProject() *stat.Metrics {
	name := gft.Name
	if len(name) <= 0 {
		name = gft.Addr
	}
	return stat.NewMetrics(name)
}

// 统计当前路径的执行时间
func (w *GFResponse) AddRouteMetric() {
	var nodeIdx int16 = -1
	if w.Ctx != nil && w.Ctx.match.ptrNode != nil {
		nodeIdx = w.Ctx.match.ptrNode.routerIdx
	}
	door.Keeper.AddItem(door.ReqItem{
		RouterIdx: nodeIdx,
		Duration:  timex.Since(w.EnterTime),
	})
}
