package sdx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/sdx/mid"
)

func SuperHandlers(app *fst.GoFast) *fst.GoFast {
	// 初始化一个全局的 请求管理器（记录访问数据，分析统计，限流熔断，定时日志）
	reqKeeper := gate.CreateReqKeeper(app.ProjectName(), app.RoutePaths())
	// 因为Routes的数量只能在加载完所有路由之后才知道,所以这里选择延时构造所有Breakers
	app.OnBeforeBuildRoutes(func(app *fst.GoFast) {
		reqKeeper.InitKeeper(app.RouteLen())
		mid.RAttrsList.Reordering(app, app.RouteLen())
	})

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 第一级：HttpHandlers
	// NOTE：Fit系列是全局的，针对所有请求起作用，而且不区分路由，这个时候还根本没有开始匹配路由。
	// GoFast提供默认的全套拦截器，开启微服务治理
	// 请求按照先后顺序依次执行这些拦截器，顺序不可随意改变
	app.UseHttpHandler(mid.HttpMaxConnections(app.SdxConfig.MaxConnections))     // 最大同时接收请求数量
	app.UseHttpHandler(mid.HttpMaxContentLength(app.SdxConfig.MaxContentLength)) // 请求头限制，最大32MB（但这是对所有请求的限制）
	// gft.UseGlobalFit(mid.CpuMetric(nil))                        // cpu 统计 | 熔断

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 第二级：ContextHandlers 带上下文 fst.Context 的执行链
	//app.Before(mid.Tracing)                        // 链路追踪，在日志打印之前执行，日志才会体现出标记
	app.Before(mid.Logger)             // 所有请求写日志，根据配置输出日志样式
	app.Before(mid.Breaker(reqKeeper)) // 自适应熔断：针对不同route，启动熔断机制（主要保护下游资源不被挤兑）
	//app.Before(mid.LoadShedding(reqKeeper))        // 自适应降载：（判断CPU和最大并发数）（主要保护自己不跑爆）
	app.Before(mid.Timeout(app.SdxConfig.EnableTimeout)) // 超时自动返回，后台处理继续，默认3000毫秒
	app.Before(mid.Recovery)                             // @@@@@ 截获所有异常 @@@@@
	app.Before(mid.HandlerTime(reqKeeper))               // 请求处理耗时统计
	app.Before(mid.Prometheus)                           // 适合 prometheus 的统计信息
	app.Before(mid.MaxContentLength)                     // 最大的请求头限制，默认32MB（这个可以单独限制不同的路径）
	app.Before(mid.Gunzip)                               // 自动 gunzip 解压缩（前面的处理都完成了再解压缩）

	// 下面的这些特性恐怕都需要用到 fork 时间模式添加监控。
	//app.Fit(mid.JwtAuthorize(app.FitJwtSecret))

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// 特殊路由的处理链（中间件）
	// 正确匹配路由之外的情况，比如特殊的404,504等路由处理链
	if app.SdxConfig.UseSpecialHandlers {
		app.SpecialBefore(mid.LoggerMini)
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
