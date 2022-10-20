// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/exec"
	"github.com/qinchende/gofast/skill/load"
	"github.com/qinchende/gofast/skill/sysx"
	"os"
	"strconv"
)

// 请求统计管理员，负责分析每个路由的请求压力和处理延时情况
type RequestKeeper struct {
	// 访问量统计 Counter
	bucket  *reqContainer
	counter *exec.Interval

	// 熔断
	Breakers []breaker.Breaker

	// 降载
	Shedding     load.Shedder
	SheddingStat *sheddingStat
}

func NewReqKeeper(name string, paths []string) *RequestKeeper {
	reqC := &reqContainer{
		pid:   os.Getpid(),
		name:  name,
		paths: paths,
	}

	return &RequestKeeper{
		counter: exec.NewInterval(LogInterval, reqC),
		bucket:  reqC,
	}
}

// 每项路由都有自己单独的熔断器，熔断器采用滑动窗口限流算法
func (rk *RequestKeeper) InitKeeper(routesLen uint16) {
	// 初始化整个路由统计结构
	rk.bucket.sumRoutes = make([]routeSum, routesLen)

	// 初始化所有Breaker
	rk.Breakers = make([]breaker.Breaker, 0, routesLen)
	for i := 0; i < int(routesLen); i++ {
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
func (rk *RequestKeeper) CounterAdd(item OneReq) {
	rk.counter.Add(item)
}

func (rk *RequestKeeper) CounterDrop(idx uint16) {
	rk.counter.Add(OneReq{
		RouteIdx: idx,
		Drop:     true,
	})
}
