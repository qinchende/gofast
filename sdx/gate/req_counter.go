// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"sync"
)

type (
	// 每个route消耗 28 字节（4字长）
	routeCounter struct {
		rtLock      sync.Mutex
		totalTimeMS int64  // 处理完成总共耗时
		maxTimeMS   int32  // 处理完成最长耗时
		accepts     uint32 // 进入处理的请求数 或者 处理完成的请求数
		timeouts    uint32 // 处理超时请求数
		drops       uint32 // 熔断丢弃请求数
	}

	// 额外的计数，需要2个字长
	extraCounter struct {
		extLock sync.Mutex
		total   uint64 // 只记调用次数
	}

	reqCounter struct {
		pid  int
		name string

		rmLock sync.Mutex

		// API访问统计
		paths  []string
		routes []routeCounter
		// 其它计数器
		extraPaths []string
		extras     []extraCounter
	}

	printData struct {
		extras []extraCounter
		routes []routeCounter
	}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 实现 exec.Interval 接口方法，方便对所有请求进行定时统计
func (rb *reqCounter) AddItem(v any) bool {
	return false
}

// 返回当前容器中的所有数据，同时重置容器
func (rb *reqCounter) RemoveAll() any {
	times := uint64(0)
	tpExtras := make([]extraCounter, len(rb.extras))
	tpRoutes := make([]routeCounter, len(rb.routes))

	rb.rmLock.Lock()
	defer rb.rmLock.Unlock()

	for i := 0; i < len(rb.extras); i++ {
		ct := &rb.extras[i]

		ct.extLock.Lock()
		tpExtras[i].total = ct.total
		ct.total = 0
		ct.extLock.Unlock()

		times += tpExtras[i].total
	}

	for i := 0; i < len(rb.routes); i++ {
		ct := &rb.routes[i]

		ct.rtLock.Lock()
		tpRoutes[i].maxTimeMS = ct.maxTimeMS
		tpRoutes[i].totalTimeMS = ct.totalTimeMS
		tpRoutes[i].accepts = ct.accepts
		tpRoutes[i].timeouts = ct.timeouts
		tpRoutes[i].drops = ct.drops
		ct.maxTimeMS = 0
		ct.totalTimeMS = 0
		ct.accepts = 0
		ct.timeouts = 0
		ct.drops = 0
		ct.rtLock.Unlock()

		times += uint64(tpRoutes[i].accepts)
		times += uint64(tpRoutes[i].drops)
	}

	// 没有数据需要处理，直接返回nil
	if times == 0 {
		return nil
	}

	pack := &printData{
		extras: tpExtras,
		routes: tpRoutes,
	}
	return pack
}

// 执行统计输出。这里的输入参数来自于上面 RemoveAll 的返回值
func (rb *reqCounter) Execute(items any) {
	data := items.(*printData)
	rb.logPrintReqCounter(data)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 统计一个通过的请求
func (rk *RequestKeeper) CountRoutePass2(idx uint16, ms int32) {
	rk.execute.AddByFunc(func(any) (any, bool) {
		ct := &rk.counter.routes[idx]

		ct.rtLock.Lock()
		ct.accepts++
		ct.totalTimeMS += int64(ms)
		if ct.maxTimeMS < ms {
			ct.maxTimeMS = ms
		}
		ct.rtLock.Unlock()

		return nil, false
	}, nil)
}

func (rk *RequestKeeper) CountRoutePass(idx uint16) {
	rk.execute.AddByFunc(func(any) (any, bool) {
		ct := &rk.counter.routes[idx]

		ct.rtLock.Lock()
		ct.accepts++
		ct.rtLock.Unlock()

		return nil, false
	}, nil)
}

// 统计一个处理超时的请求
func (rk *RequestKeeper) CountRouteTimeout(idx uint16) {
	rk.execute.AddByFunc(func(any) (any, bool) {
		ct := &rk.counter.routes[idx]

		ct.rtLock.Lock()
		ct.timeouts++
		ct.rtLock.Unlock()

		return nil, false
	}, nil)
}

// 统计一个被丢弃的请求
func (rk *RequestKeeper) CountRouteDrop(idx uint16) {
	rk.execute.AddByFunc(func(any) (any, bool) {
		ct := &rk.counter.routes[idx]

		ct.rtLock.Lock()
		ct.drops++
		ct.rtLock.Unlock()

		return nil, false
	}, nil)
}

// 添加其它统计项
func (rk *RequestKeeper) CountExtras(pos uint16) {
	rk.execute.AddByFunc(func(any) (any, bool) {
		ct := &rk.counter.extras[pos]

		ct.extLock.Lock()
		ct.total++
		ct.extLock.Unlock()

		return nil, false
	}, nil)
}
