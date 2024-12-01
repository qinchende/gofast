// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/aid/sysx"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/sdx/mid"
	"time"
)

// 闪电侠实现的中间件控制参数
type MidConfig struct {
	//SysStateMonitor  bool  `v:"def=true"`                  // 是否启动系统资源使用情况的定时检查工作
	PrintSysState    bool  `v:"def=true"`                  // 定时打印系统资源状态检查日志
	PrintRouteState  bool  `v:"def=true"`                  // 定时打印路由访问统计数据
	MaxContentLength int64 `v:"def=33554432"`              // 最大请求字节数，32MB（33554432），传0不限制
	MaxConnections   int32 `v:"def=0,range=[0:100000000]"` // 最大同时请求数，0不限制

	EnableSpecialHandlers bool          `v:"def=true"`   // 是否启用默认的特殊路由中间件
	EnableTrack           bool          `v:"def=false"`  // 启动链路追踪
	EnableGunzip          bool          `v:"def=false"`  // 启动gunzip
	EnableShedding        bool          `v:"def=true"`   // 启动降载限制访问
	EnableTimeout         bool          `v:"def=true"`   // 启动超时拦截
	ReqTimeout            time.Duration `v:"def=3000ms"` // 默认请求超时时间
}

var cnf MidConfig

func SetMidConfig(cf *MidConfig) {
	cnf = *cf
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GoFast default handlers chain
func SuperHandlers(app *fst.GoFast) *fst.GoFast {
	// 初始化一个全局的 请求管理器（记录访问数据，分析统计，限流降载熔断，定时日志）
	keeper := gate.NewReqKeeper(gate.KeeperCnf{
		ProjName:   app.ProjectName(),
		PrintRoute: cnf.PrintRouteState,
	})
	app.OnBeforeBuildRoutes(func(app *fst.GoFast) {
		// 因为Routes的数量只能在加载完所有路由之后才知道,所以这里选择延时构造所有Breakers
		mid.RoutesAttrs.Rebuild(app.RoutesLen(), cnf.ReqTimeout) // 所有路由配置
		sysx.OpenSysMonitor(cnf.PrintSysState)                   // 系统资源监控

		routePaths := app.RoutePathsWithMethod()
		extraPaths := []string{"AllRequest", "RouteMatched", "LoadShedding"}
		keeper.InitAndRun(routePaths, extraPaths) // 看守上岗
	})

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 第一级：HttpHandlers
	// NOTE：Fit系列是全局的，针对所有请求起作用，而且不区分路由，这个时候还根本没有开始匹配路由。
	// GoFast提供默认的全套拦截器，开启微服务治理
	// 请求按照先后顺序依次执行这些拦截器，顺序不可随意改变
	app.UseHttpHandler(mid.HttpReqCountPos(keeper, 0))                 // 访问计数1
	app.UseHttpHandler(mid.HttpMaxConnections(cnf.MaxConnections))     // 最大同时处理请求数量
	app.UseHttpHandler(mid.HttpMaxContentLength(cnf.MaxContentLength)) // 请求头最大限制

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 第二级：ContextHandlers 带上下文 fst.Context 的执行链
	app.Before(mid.ReqCountPos(keeper, 1))                      // 正确匹配路由的请求数
	app.Before(mid.Tracing(app.AppName, cnf.EnableTrack))       // 链路追踪
	app.Before(mid.Logger)                                      // 请求日志
	app.Before(mid.Breaker(keeper))                             // 自适应熔断
	app.Before(mid.LoadShedding(keeper, cnf.EnableShedding, 2)) // 过载保护
	app.Before(mid.Timeout(keeper, cnf.EnableTimeout))          // 超时自动返回（请求在后台任然继续执行）
	app.Before(mid.Recovery)                                    // @@@ 截获所有异常，避免服务进程崩溃 @@@
	app.Before(mid.TimeMetric(keeper))                          // 耗时统计
	app.Before(mid.MaxContentLength)                            // 分路由判断请求长度
	app.Before(mid.Gunzip(cnf.EnableGunzip))                    // 自动 gunzip 解压缩

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 特殊路由的处理链
	// 正确匹配路由之外的情况，比如特殊的404,504等路由处理链
	if cnf.EnableSpecialHandlers {
		app.SpecialBefore(mid.ReqCount(keeper)) // 数量统计
		app.SpecialBefore(mid.LoggerMini)       // 特殊路径的日志
	}
	return app
}

// ++++++++++++++++++++++++++++++++++ go-zero default handlers chain
//  hds := alice.New(
//  handler.TracingHandler,
//  s.getLogHandler(),
//  handler.MaxConns(s.conf.MaxConns),
//  handler.BreakerHandler(route.Method, route.Path, metrics),
//  handler.SheddingHandler(s.getShedder(fr.priority), metrics),
//  handler.TimeoutHandler(time.Duration(s.conf.Timeout)*time.Millisecond),
//  handler.RecoverHandler,
//  handler.MetricHandler(metrics),
//  handler.PromethousHandler(route.Path),
//  handler.MaxBytesHandler(s.conf.MaxBytes),
//  handler.GunzipHandler,
//  )
//  hds = s.appendAuthHandler(fr, hds, verifier)
//
//  for _, middleware := range s.middlewares {
//  hds = chain.Append(convertMiddleware(middleware))
//  }
//  handle := chain.ThenFunc(route.Handler)
