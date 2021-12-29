// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/executors"
	"github.com/qinchende/gofast/skill/load"
	"os"
	"strconv"
)

// 请求统计管理员，负责分析每个路由的请求压力和处理延时情况
type RequestKeeper struct {
	executor  *executors.IntervalExecutor
	container *reqContainer

	// 熔断器
	Breakers []breaker.Breaker

	// 降载组件
	Shedding     load.Shedder
	SheddingStat *SheddingStat
}

func CreateReqKeeper(name string, fp FuncGetPath) *RequestKeeper {
	counter := &reqContainer{
		name:    name,
		pid:     os.Getpid(),
		getPath: fp,
	}

	return &RequestKeeper{
		executor:  executors.NewIntervalExecutor(LogInterval, counter),
		container: counter,
	}
}

// 每项路由都有自己单独的熔断器，熔断器采用滑动窗口限流算法
func (rk *RequestKeeper) InitKeeper(rtLen uint16) {
	// 初始化整个路由统计结构
	rk.container.sumRoutes = make([]routeSum, rtLen, rtLen)
	// rk.container.
	rk.Breakers = make([]breaker.Breaker, 0, rtLen)
	for i := 0; i < int(rtLen); i++ {
		rk.Breakers = append(rk.Breakers, breaker.NewBreaker(breaker.WithName(strconv.Itoa(i))))
	}

	// 初始化 降载 组件
	rk.SheddingStat = NewSheddingStat("cpu")
	rk.Shedding = load.NewAdaptiveShedder()
	//rk.Shedding = load.NewAdaptiveShedder(load.WithCpuThreshold(900))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 添加一次请求项目
func (rk *RequestKeeper) AddItem(item ReqItem) {
	rk.executor.Add(item)
}

// 添加一次被丢弃的请求，只需要标记本次
func (rk *RequestKeeper) AddDrop(idx uint16) {
	rk.executor.Add(ReqItem{
		RouteIdx: idx,
		Drop:     true,
	})
}
