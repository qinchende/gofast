// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import "sync"

// 请求数据增长20%
const reqsGrowsRate = 0.2

type (
	// 每个请求需要消耗 7 字节（一段时间内请求一般都很多，要省内存）
	oneReq struct {
		takeTimeMS int32  // 单次请求耗时
		routeIdx   uint16 // 当前请求对应路由树节点的index，这用来单独统计不同route
		isDrop     bool   // 是否是一个被丢弃的请求（熔断或者资源超限拒绝处理）
	}

	// 每个route消耗 20 字节
	routeCounter struct {
		totalTimeMS int64
		maxTimeMS   int32
		accepts     uint32
		drops       uint32
	}

	otherCounter struct {
		lock  sync.Mutex
		total uint64
	}

	reqBucket struct {
		name string
		pid  int

		lastReqsLen int
		reqs        []oneReq

		// API访问统计
		paths  []string
		routes []routeCounter

		// 其它计数器
		otherPaths []string
		others     []otherCounter
	}
)

// 重置统计相关参数
func (rb *reqBucket) resetRouteCounters() {
	for i := 0; i < len(rb.routes); i++ {
		rb.routes[i].totalTimeMS = 0
		rb.routes[i].accepts = 0
		rb.routes[i].drops = 0
		rb.routes[i].maxTimeMS = 0
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 实现 exec.Interval 接口方法，方便这个集装箱进行定时统计

// 如果这里返回true，意味着要立刻刷新当前所有统计数据，这个开关用户自定义输出日志
// 一般这里都应该返回 false
func (rb *reqBucket) AddItem(v any) bool {
	if req, ok := v.(oneReq); ok {
		rb.reqs = append(rb.reqs, req)
	}
	return false
}

// 返回当前容器中的所有数据，同时重置容器
func (rb *reqBucket) RemoveAll() any {
	defer func() {
		// 初始化长度为上次数量的倍数 reqsGrowsRate，防止中途频繁扩容的开销
		//rb.reqs = nil
		if len(rb.reqs) > 0 || cap(rb.reqs) == 0 {
			rb.reqs = make([]oneReq, 0, int(float64(rb.lastReqsLen+1)*(1+reqsGrowsRate)))
		}
	}()
	rb.lastReqsLen = len(rb.reqs)

	// 特殊情况也需要统计
	if len(rb.reqs) == 0 && len(rb.others) > 0 && rb.others[0].total > 0 {
		rb.reqs = append(rb.reqs, oneReq{routeIdx: 0, isDrop: true, takeTimeMS: -1})
	}
	return rb.reqs
}

// 执行统计输出
func (rb *reqBucket) Execute(items any) {
	defer func() {
		rb.logPrintOthers()
		rb.logPrintRoutes()
		rb.resetRouteCounters()
	}()

	// 这里不需要断言判断类型转换的真假，因为结果是上面 RemoveAll 返回的
	// 只能取 items 的值，不能取rb.reqs，因为这个已经清空了
	reqs := items.([]oneReq)
	//first := rb.reqs[0] // 至少有一个，否则不会触发

	for _, req := range reqs {
		route := &rb.routes[req.routeIdx]
		if req.isDrop {
			route.drops++
			continue
		}

		route.accepts++
		route.totalTimeMS += int64(req.takeTimeMS)
		if req.takeTimeMS > route.maxTimeMS {
			route.maxTimeMS = req.takeTimeMS
		}
	}
}
