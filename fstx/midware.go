package fstx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/door"
	"github.com/qinchende/gofast/fst/mid"
	"time"
)

// 第一级：
// NOTE：Fit系列是全局的，针对所有请求起作用，而且不区分路由，这个时候还根本没有开始匹配路由。
// GoFast提供默认的全套拦截器，开启微服务治理
// 请求按照先后顺序依次执行这些拦截器，顺序不可随意改变
func AddDefaultFits(gft *fst.GoFast) *fst.GoFast {
	gft.Fit(mid.MaxConnections(gft.FitMaxConnections)) // 最大同时处理请求数量：100万
	gft.Fit(mid.CpuMetric(nil))                        // cpu 统计 | 熔断
	return gft
}

// 第二级：
// 带上下文 gofast.fst.Context 的执行链
func AddDefaultHandlers(gft *fst.GoFast) *fst.GoFast {
	// 初始化一个全局的路由统计器
	door.InitKeeper(gft.FullPath)

	//gft.Fit(mid.Tracing)                                                         // 加入调用链路追踪标记
	gft.Before(mid.ReqLogger)                                                              // 所有请求写日志，根据配置输出日志样式
	gft.Before(mid.ReqTimeout(time.Duration(gft.FitReqTimeout) * 1000 * time.Millisecond)) // 超时自动返回，后台处理继续，默认3000毫秒
	gft.Before(mid.Recovery())                                                             // 截获所有异常
	gft.Before(mid.RouteMetric)                                                            // path 访问统计
	gft.Before(mid.MaxContentLength(gft.FitMaxContentLength))                              // 最大的请求头限制，默认32MB
	gft.Before(mid.Gunzip)                                                                 // 自动 gunzip 解压缩

	// 下面的这些特性恐怕都需要用到 fork 时间模式添加监控。
	//gft.Fit(mid.BreakerDoor())
	//gft.Fit(mid.JwtAuthorize(gft.FitJwtSecret))
	return gft
}

// ++++++++++++++++++++++++++++++++++ go-zero default handler chains
//  chain := alice.New(
//  handler.TracingHandler,
//  s.getLogHandler(),
//  handler.MaxConns(s.conf.MaxConns),
//  handler.BreakerHandler(route.Method, route.Path, metrics),
//  handler.SheddingHandler(s.getShedder(fr.priority), metrics),.
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
