// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/logx"
	"time"
)

// 固定分钟作为统计周期
var LogInterval = time.Minute

type (
	// 每个请求需要消耗 2个字长 16字节的空间
	OneReq struct {
		LossTime time.Duration // 单次请求耗时
		RouteIdx uint16        // 当前请求对应路由树节点的index，这用来单独统计不同route
		Drop     bool          // 是否是一个被丢弃的请求（熔断或者资源超限拒绝处理）
	}

	// 包裹一层，用于在定时任务器中传递对象
	deliverItems struct {
		reqs []OneReq
	}

	routeSum struct {
		sumTime time.Duration
		maxLoss time.Duration
		accepts uint32
		drops   uint32
	}

	// 存放所有请求的处理时间，作为统计的容器
	reqContainer struct {
		paths []string
		name  string
		pid   int

		currReqs  []OneReq
		sumRoutes []routeSum
	}
)

// 重置统计相关参数
func (rc *reqContainer) resetSum() {
	for i := 0; i < len(rc.sumRoutes); i++ {
		rc.sumRoutes[i].sumTime = 0
		rc.sumRoutes[i].accepts = 0
		rc.sumRoutes[i].drops = 0
		rc.sumRoutes[i].maxLoss = 0
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 添加统计项目
// 如果这里返回true，意味着要立刻刷新当前所有统计数据，这个开关用户自定义输出日志
// 一般这里都应该返回 false
func (rc *reqContainer) AddItem(v any) bool {
	if item, ok := v.(OneReq); ok {
		rc.currReqs = append(rc.currReqs, item)
	}
	return false
}

// 返回当前容器中的所有数据，同时重置容器
func (rc *reqContainer) RemoveAll() any {
	ret := deliverItems{
		reqs: rc.currReqs,
	}

	//reqs := rc.currReqs
	rc.currReqs = nil
	return ret
}

// 执行统计输出
func (rc *reqContainer) Execute(items any) {
	// 这里不需要断言判断类型转换的真假，因为结果是上面 RemoveAll 返回的
	ret := items.([]OneReq)
	reqs := ret
	rc.resetSum()

	// 用一次循环，分别统计不同route的访问情况
	rtsAll := &rc.sumRoutes[0]
	for _, req := range reqs {
		rts := &rc.sumRoutes[req.RouteIdx]
		if req.Drop {
			rtsAll.drops++
			rts.drops++
		} else {
			rtsAll.accepts++
			rts.accepts++

			rtsAll.sumTime += req.LossTime
			rts.sumTime += req.LossTime

			// 记录其中最长的响应时间
			if req.LossTime > rts.maxLoss {
				rts.maxLoss = req.LossTime
			}
		}
	}

	rc.logPrint()
}

// 打印每个路由的请求数据。
// TODO: 其实每项路由分钟级的日志应该是收集起来，放入数据库，可视化展示和分析
func (rc *reqContainer) logPrint() {
	for idx, route := range rc.sumRoutes {
		if idx != 0 && route.accepts == 0 && route.drops == 0 {
			continue
		}

		qps := float32(route.accepts) / float32(LogInterval/time.Second)
		var aveTime float32
		if route.accepts > 0 {
			aveTime = float32(route.sumTime/time.Millisecond) / float32(route.accepts)
		}

		logx.StatF("%s | suc: %d, drop: %d, qps: %.1f/s ave: %.1fms, max: %.1fms",
			rc.paths[idx], route.accepts, route.drops, qps, aveTime, float32(route.maxLoss/time.Millisecond))
	}
}
