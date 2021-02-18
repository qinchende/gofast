package fstx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/mid"
)

// GoFast提供默认的全套拦截器
// 请求按照先后顺序依次执行这些过滤器
func AddDefaultFits(gft *fst.GoFast) *fst.GoFast {
	gft.Fit(mid.ReqLogger(gft.FitLogType))
	gft.Fit(mid.MaxReqContentLength(gft.FitMaxReqContentLen))
	gft.Fit(mid.GunzipFit)
	return gft
}
