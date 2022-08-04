package logx

import (
	"github.com/qinchende/gofast/cst"
	"net/http"
	"time"
)

// 日志参数实体
type ReqLogEntity struct {
	RawReq     *http.Request
	TimeStamp  time.Duration
	Latency    time.Duration
	ClientIP   string
	StatusCode int
	Pms        cst.KV
	BodySize   int
	ResData    []byte
	ErrorMsg   string
}

// 打印请求日志，可以指定不同的输出样式
func WriteReqLog(p *ReqLogEntity) {
	switch myCnf.logStyle {
	case LogStyleSdx:
		InfoDirect(genSdxReqLogString(p))
	case LogStyleSdxMini:
		InfoDirect(genSdxReqLogString(p))
	case LogStyleJson:
		InfoDirect(genSdxReqLogString(p))
	case LogStyleJsonMini:
		InfoDirect(genSdxReqLogString(p))
	default:
	}
}
