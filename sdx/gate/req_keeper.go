// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/aid/exec"
	"os"
	"time"
)

// 固定分钟作为统计周期
const countInterval = 1 * time.Minute

type KeeperCnf struct {
	ProjName   string
	PrintRoute bool
}

// 请求统计管理员，负责分析每个路由的请求压力和处理延时情况
type RequestKeeper struct {
	KeeperCnf
	// 目前支持以下几种功能
	*reqCounter            // 分路由访问计数器
	Breakers    []*Breaker // 分路由熔断器
	Limiters    []*Limiter // 分路由限流器
}

func NewReqKeeper(cnf KeeperCnf) *RequestKeeper {
	kp := &RequestKeeper{
		KeeperCnf: cnf,
		reqCounter: &reqCounter{
			pid:       os.Getpid(),
			needPrint: cnf.PrintRoute,
		},
	}
	kp.reqCounter.counter = exec.NewIntervalUnsafe(countInterval, kp.reqCounter)
	return kp
}

// 开启监控统计
func (rk *RequestKeeper) InitAndRun(routePaths, extraPaths []string) {
	rLen := len(routePaths)

	// 1. 初始化 reqCounter
	rk.reqCounter.paths = routePaths                          // 初始化整个路由统计结构
	rk.reqCounter.routes = make([]routeData, rLen)            // 初始化整个路由统计结构
	rk.reqCounter.extraPaths = extraPaths                     // 其它统计
	rk.reqCounter.extras = make([]extraData, len(extraPaths)) // 其它统计

	// 2. 初始化所有Breaker，每个路由都有自己单独的熔断计数器
	rk.Breakers = make([]*Breaker, rLen)
	for i := 0; i < rLen; i++ {
		rk.Breakers[i] = NewBreaker(rk.ProjName + "#" + routePaths[i])
	}

	// 3. 初始化降载信息收集器
	rk.Limiters = make([]*Limiter, rLen)
	for i := 0; i < rLen; i++ {
		rk.Limiters[i] = NewLimiter()
	}
}
