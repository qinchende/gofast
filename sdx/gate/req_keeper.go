// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/executors"
	"github.com/qinchende/gofast/skill/load"
	"github.com/qinchende/gofast/skill/sysx"
	"os"
	"strconv"
)

func CreateReqKeeper(name string, paths []string) *RequestKeeper {
	reqC := &reqContainer{
		pid:      os.Getpid(),
		name:     name,
		allPaths: paths,
	}

	return &RequestKeeper{
		counter: executors.NewIntervalExecutor(LogInterval, reqC),
		bucket:  reqC,
	}
}

// 每项路由都有自己单独的熔断器，熔断器采用滑动窗口限流算法
func (rk *RequestKeeper) InitKeeper(rtLen uint16) {
	// 初始化整个路由统计结构
	rk.bucket.sumRoutes = make([]routeSum, rtLen, rtLen)

	// rk.container.
	rk.Breakers = make([]breaker.Breaker, 0, rtLen)
	for i := 0; i < int(rtLen); i++ {
		rk.Breakers = append(rk.Breakers, breaker.NewBreaker(breaker.WithName(strconv.Itoa(i))))
	}

	// 有个前提是 CPU 的监控必须启动
	if sysx.CpuMonitor {
		// 初始化 降载 组件
		rk.SheddingStat = createSheddingStat()  // 降载信息打印
		rk.Shedding = load.NewAdaptiveShedder() // 降载统计分析
		// rk.Shedding = load.NewAdaptiveShedder(load.WithCpuThreshold(900))
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 添加一次请求项目
func (rk *RequestKeeper) CounterAdd(item ReqItem) {
	rk.counter.Add(item)
}

func (rk *RequestKeeper) CounterDrop(idx uint16) {
	rk.counter.Add(ReqItem{
		RouteIdx: idx,
		Drop:     true,
	})
}
