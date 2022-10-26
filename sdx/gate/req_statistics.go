// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

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

	reqBucket struct {
		paths []string
		name  string
		pid   int

		reqs   []oneReq
		routes []routeCounter
	}
)

// 重置统计相关参数
func (rb *reqBucket) reset() {
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
		rb.reqs = nil
	}()
	return rb.reqs
}

// 执行统计输出
func (rb *reqBucket) Execute(items any) {
	// 这里不需要断言判断类型转换的真假，因为结果是上面 RemoveAll 返回的
	//reqs := items.([]oneReq)

	//rtsAll := &rb.routes[0]
	for _, req := range rb.reqs {
		route := &rb.routes[req.routeIdx]
		if req.isDrop {
			//rtsAll.drops++
			route.drops++
			continue
		}

		//rtsAll.accepts++
		route.accepts++
		//rtsAll.totalTimeMS += int64(req.takeTimeMS)
		route.totalTimeMS += int64(req.takeTimeMS)
		if req.takeTimeMS > route.maxTimeMS {
			route.maxTimeMS = req.takeTimeMS
		}
	}

	rb.logPrint()
	rb.reset()
}
