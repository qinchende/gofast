package fstx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/mid"
	"time"
)

// 开启微服务治理
// GoFast提供默认的全套拦截器
// 请求按照先后顺序依次执行这些过滤器
func AddDefaultFits(gft *fst.GoFast) *fst.GoFast {
	gft.Fit(mid.ReqLogger(gft.FitLogType)) // 所有请求写日志，放第一
	gft.Fit(mid.Recovery())                // 截获所有异常
	gft.Fit(mid.ReqTimeout(time.Duration(gft.FitReqTimeout) * time.Millisecond))
	gft.Fit(mid.MaxReqCounts(gft.FitMaxReqCount))
	gft.Fit(mid.MaxReqContentLength(gft.FitMaxReqContentLen))
	gft.Fit(mid.Gunzip)
	//gft.Fit(mid.JwtAuthorize(gft.FitJwtSecret))
	return gft
}
