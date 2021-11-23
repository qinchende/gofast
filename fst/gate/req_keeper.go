// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/executors"
	"os"
	"strconv"
)

// 请求统计管理员，负责分析每个路由的请求压力和处理延时情况
type RequestKeeper struct {
	executor  *executors.IntervalExecutor
	container *reqContainer
	Breakers  []breaker.Breaker
}

func CreateReqKeeper(name string, fp FuncGetPath) *RequestKeeper {
	container := &reqContainer{
		name:    name,
		pid:     os.Getpid(),
		getPath: fp,
	}

	return &RequestKeeper{
		executor:  executors.NewIntervalExecutor(LogInterval, container),
		container: container,
	}
}

// 每项路由都有自己单独的熔断器，熔断器采用滑动窗口限流算法
func (rk *RequestKeeper) SetBreakers(length uint16) {
	rk.Breakers = make([]breaker.Breaker, 0, length)
	for i := 0; i < int(length); i++ {
		rk.Breakers = append(rk.Breakers, breaker.NewBreaker(breaker.WithName(strconv.Itoa(i))))
	}
}

// 添加一次请求项目
func (rk *RequestKeeper) AddItem(item ReqItem) {
	rk.executor.Add(item)
}

// 添加一次被丢弃的请求，只需要标记本次
func (rk *RequestKeeper) AddDrop() {
	rk.executor.Add(ReqItem{
		Drop: true,
	})
}
