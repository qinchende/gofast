// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"fmt"
	"sync"
)

type (
	// 每个route消耗 28 字节（4字长）
	routeCounter struct {
		totalTimeMS int64  // 正常请求总共耗时
		maxTimeMS   int32  // 正常请求最长耗时
		accepts     uint32 // 正常处理的请求数
		drops       uint32 // 丢弃的请求数
	}

	reqBucket struct {
		pid  int
		name string

		lock  sync.Mutex
		times uint64

		// API访问统计
		paths  []string
		routes []routeCounter // 每个占4字长
		// 其它计数器
		extraPaths []string
		extras     []uint64 // 每个占1字长
	}

	printData struct {
		extras []uint64
		routes []routeCounter
	}
)

func (rb *reqBucket) initCounters() {
	rb.times = 0
	rb.routes = make([]routeCounter, len(rb.paths)) // 初始化整个路由统计结构
	rb.extras = make([]uint64, len(rb.extraPaths))  // 其它统计
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 实现 exec.Interval 接口方法，方便对所有请求进行定时统计
func (rb *reqBucket) AddItem(v any) bool {
	return false
}

// 返回当前容器中的所有数据，同时重置容器
func (rb *reqBucket) RemoveAll() any {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	fmt.Println(rb.times)

	// 没有数据需要处理，直接返回nil
	if rb.times == 0 {
		return nil
	}

	pack := &printData{
		extras: rb.extras,
		routes: rb.routes,
	}
	rb.initCounters()
	return pack
}

// 执行统计输出
func (rb *reqBucket) Execute(items any) {
	data := items.(*printData)
	rb.logPrintReqCounter(data)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 统计一个通过的请求
func (rk *RequestKeeper) CountRoutePass2(idx uint16, ms int32) {
	rk.counter.AddByFunc(func(any) (any, bool) {
		rk.bucket.lock.Lock()
		ct := &rk.bucket.routes[idx]
		ct.accepts++
		ct.totalTimeMS += int64(ms)
		if ct.maxTimeMS < ms {
			ct.maxTimeMS = ms
		}
		rk.bucket.times++
		rk.bucket.lock.Unlock()
		return nil, false
	}, nil)
}

func (rk *RequestKeeper) CountRoutePass(idx uint16) {
	rk.counter.AddByFunc(func(any) (any, bool) {
		rk.bucket.lock.Lock()
		rk.bucket.routes[idx].accepts++
		rk.bucket.times++
		rk.bucket.lock.Unlock()
		return nil, false
	}, nil)
}

// 统计一个被丢弃的请求
func (rk *RequestKeeper) CountRouteDrop(idx uint16) {
	rk.counter.AddByFunc(func(any) (any, bool) {
		rk.bucket.lock.Lock()
		rk.bucket.routes[idx].drops++
		rk.bucket.times++
		rk.bucket.lock.Unlock()
		return nil, false
	}, nil)
}

// 添加其它统计项
func (rk *RequestKeeper) CountExtras(pos uint16) {
	rk.counter.AddByFunc(func(any) (any, bool) {
		rk.bucket.lock.Lock()
		rk.bucket.extras[pos]++
		rk.bucket.times++
		rk.bucket.lock.Unlock()
		return nil, false
	}, nil)
}
