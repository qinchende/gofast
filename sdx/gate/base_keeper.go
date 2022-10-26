// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/exec"
	"github.com/qinchende/gofast/skill/load"
	"os"
	"strconv"
	"time"
)

// 固定分钟作为统计周期
const CountInterval = time.Minute

// 请求统计管理员，负责分析每个路由的请求压力和处理延时情况
type RequestKeeper struct {
	started bool           // 是否开始运行了
	bucket  *reqBucket     //
	counter *exec.Interval // 循环计数器

	Breakers     []breaker.Breaker // 不同路径的熔断统计器
	Shedding     load.Shedder
	SheddingStat *sheddingStat
}

func NewReqKeeper(name string) *RequestKeeper {
	bkt := &reqBucket{pid: os.Getpid(), name: name}
	return &RequestKeeper{
		counter: exec.NewInterval(CountInterval, bkt),
		bucket:  bkt,
	}
}

// 开启监控
func (rk *RequestKeeper) StartWorking(routePaths []string) {
	if rk.started == true {
		return
	}
	rk.started = true

	// 初始化整个路由统计结构
	routesLen := len(routePaths)
	rk.bucket.routes = make([]routeCounter, routesLen)
	rk.bucket.paths = routePaths

	// 初始化所有Breaker，每个路由都有自己单独的熔断计数器
	rk.Breakers = make([]breaker.Breaker, 0, routesLen)
	for i := 0; i < int(routesLen); i++ {
		rk.Breakers = append(rk.Breakers, breaker.NewBreaker(breaker.WithName(strconv.Itoa(i))))
	}

	//if sysx.MonitorStarted {
	//	rk.SheddingStat = createSheddingStat()  // 降载信息打印
	//	rk.Shedding = load.NewAdaptiveShedder() // 降载统计分析
	//	// rk.Shedding = load.NewAdaptiveShedder(load.WithCpuThreshold(900))
	//}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 添加一次请求项目
// 只有熔断，或者降载发生的时候，才算是一个drop请求
//func (rk *RequestKeeper) AddOne(it oneReq) {
//	rk.counter.Add(it)
//}

// 统计一个通过的请求
func (rk *RequestKeeper) CountPass(idx uint16, ms int32) {
	rk.counter.Add(oneReq{routeIdx: idx, takeTimeMS: ms})
}

// 统计一个被丢弃的请求
func (rk *RequestKeeper) CountDrop(idx uint16) {
	rk.counter.Add(oneReq{routeIdx: idx, isDrop: true})
}
