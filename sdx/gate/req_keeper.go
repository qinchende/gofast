// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/exec"
	"github.com/qinchende/gofast/skill/fuse"
	"github.com/qinchende/gofast/skill/sysx"
	"os"
	"time"
)

// 固定分钟作为统计周期
const CountInterval = time.Minute

// 请求统计管理员，负责分析每个路由的请求压力和处理延时情况
type RequestKeeper struct {
	counter *reqCounter          // 请求统计器
	execute *exec.IntervalUnsafe // 定时打印统计数据

	Breakers     []fuse.Breaker // 不同路径的熔断统计器
	Shedding     fuse.Shedder   // 降载，服务器资源超限
	SheddingStat *sheddingStat
}

func NewReqKeeper(name string) *RequestKeeper {
	ct := &reqCounter{pid: os.Getpid(), name: name}
	return &RequestKeeper{
		execute: exec.NewIntervalUnsafe(CountInterval, ct),
		counter: ct,
	}
}

// 开启监控统计
func (rk *RequestKeeper) InitAndRun(routePaths, extraPaths []string) {
	rk.counter.paths = routePaths                             // 初始化整个路由统计结构
	rk.counter.routes = make([]routeCounter, len(routePaths)) // 初始化整个路由统计结构
	rk.counter.extraPaths = extraPaths                        // 其它统计
	rk.counter.extras = make([]extraCounter, len(extraPaths)) // 其它统计

	routesLen := len(routePaths)
	// 初始化所有Breaker，每个路由都有自己单独的熔断计数器
	rk.Breakers = make([]fuse.Breaker, routesLen)
	for i := 0; i < routesLen; i++ {
		rk.Breakers[i] = fuse.NewGBreaker(rk.counter.name+"#"+routePaths[i], true)
	}

	if sysx.MonitorStarted {
		rk.SheddingStat = createSheddingStat()  // 降载信息打印
		rk.Shedding = fuse.NewAdaptiveShedder() // 降载统计分析
		// rk.Shedding = load.NewAdaptiveShedder(load.WithCpuThreshold(900))
	}
}
