package logx

import (
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
	// isTerm     bool
	Pms      map[string]any
	BodySize int
	ResData  []byte
	// Keys       map[string]interface{}
	ErrorMsg string
}

// 打印请求日志，可以指定不同的输出样式
func WriteReqLog(p *ReqLogEntity) {
	switch currConfig.style {
	case StyleSdx:
		Print(genSdxReqLogString(p))
	case StyleSdxMini:
		Print(genSdxReqLogString(p))
	case StyleJson:
		Print(genSdxReqLogString(p))
	case StyleJsonMini:
		Print(genSdxReqLogString(p))
	default:
	}
}
