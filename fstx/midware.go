package fstx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/mid"
	"time"
)

// 开启微服务治理
// GoFast提供默认的全套拦截器
// 请求按照先后顺序依次执行这些拦截器
func AddDefaultFits(gft *fst.GoFast) *fst.GoFast {
	//gft.Fit(mid.Tracing)                   // 加入调用链路追踪标记
	gft.Fit(mid.ReqLogger())                                                     // 所有请求写日志，放第一
	gft.Fit(mid.Recovery())                                                      // 截获所有异常
	gft.Fit(mid.ReqTimeout(time.Duration(gft.FitReqTimeout) * time.Millisecond)) // 超时自动返回，后台处理继续，默认3000毫秒
	gft.Fit(mid.MaxReqCounts(gft.FitMaxReqCount))                                // 最大处理请求数量限制 100万
	gft.Fit(mid.MaxReqContentLength(gft.FitMaxReqContentLen))                    // 最大的请求头限制，默认32MB
	//gft.Fit(mid.Metric(gft.CreateMetrics())) // 系统访问频率按响应时间分段统计
	gft.Fit(mid.Gunzip) // 自动 gunzip 解压缩
	//gft.Fit(mid.JwtAuthorize(gft.FitJwtSecret))
	return gft
}

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
