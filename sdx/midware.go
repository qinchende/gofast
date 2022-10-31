package sdx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/sdx/mid"
	"github.com/qinchende/gofast/skill/sysx"
)

func SuperHandlers(app *fst.GoFast) *fst.GoFast {
	cnf := app.SdxConfig

	// 初始化一个全局的 请求管理器（记录访问数据，分析统计，限流熔断，定时日志）
	reqKeeper := gate.NewReqKeeper(app.ProjectName())
	app.OnBeforeBuildRoutes(func(app *fst.GoFast) {
		// 因为Routes的数量只能在加载完所有路由之后才知道,所以这里选择延时构造所有Breakers
		mid.AllAttrs.Rebuild(app.RoutesLen(), &cnf)        // 所有路由配置
		sysx.OpenSysMonitor(cnf.SysStatePrint)             // 系统资源监控
		reqKeeper.StartWorking(app.RoutePathsWithMethod()) // 警卫上岗
	})

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 第一级：HttpHandlers
	// NOTE：Fit系列是全局的，针对所有请求起作用，而且不区分路由，这个时候还根本没有开始匹配路由。
	// GoFast提供默认的全套拦截器，开启微服务治理
	// 请求按照先后顺序依次执行这些拦截器，顺序不可随意改变
	app.UseHttpHandler(mid.HttpAccessCount1(reqKeeper))                // 访问计数1
	app.UseHttpHandler(mid.HttpMaxConnections(cnf.MaxConnections))     // 最大同时接收请求数量
	app.UseHttpHandler(mid.HttpMaxContentLength(cnf.MaxContentLength)) // 请求头最大限制
	app.UseHttpHandler(mid.HttpLoadShedding(reqKeeper))                // 资源使用统计，超限就降载
	app.UseHttpHandler(mid.HttpAccessCount2(reqKeeper))                // 访问计数2

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 第二级：ContextHandlers 带上下文 fst.Context 的执行链
	app.Before(mid.AccessCount3(reqKeeper))    // 正确匹配路由的请求数
	app.Before(mid.Tracing(cnf.EnableTrack))   // 链路追踪
	app.Before(mid.Logger)                     // 请求日志
	app.Before(mid.Breaker(reqKeeper))         // 自适应熔断：不同route，不同熔断阈值
	app.Before(mid.Timeout(cnf.EnableTimeout)) // 超时自动返回，默认3000毫秒（后台处理继续）
	app.Before(mid.Recovery)                   // @@@ 截获所有异常，避免服务进程崩溃 @@@
	app.Before(mid.TimeMetric(reqKeeper))      // 耗时统计
	app.Before(mid.Prometheus)                 // 适合 prometheus 的统计信息
	app.Before(mid.ContentLength)              // 分路由判断请求长度
	app.Before(mid.Gunzip)                     // 自动 gunzip 解压缩

	// 下面的这些特性恐怕都需要用到 fork 时间模式添加监控。
	// app.Fit(mid.JwtAuthorize(app.FitJwtSecret))
	// app.Before(mid.LoadShedding(reqKeeper))        // 自适应降载：（判断CPU和最大并发数）（主要保护自己不跑爆）

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 特殊路由的处理链
	// 正确匹配路由之外的情况，比如特殊的404,504等路由处理链
	if cnf.UseSpecialHandlers {
		app.SpecialBefore(mid.LoggerMini)            // 特殊路径的日志
		app.SpecialBefore(mid.TimeMetric(reqKeeper)) // 耗时统计
	}
	return app
}

// ++++++++++++++++++++++++++++++++++ go-zero default handler chains
//  chain := alice.New(
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
//  chain = s.appendAuthHandler(fr, chain, verifier)
//
//  for _, middleware := range s.middlewares {
//  chain = chain.Append(convertMiddleware(middleware))
//  }
//  handle := chain.ThenFunc(route.Handler)
