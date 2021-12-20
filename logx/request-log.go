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
	Pms      map[string]interface{}
	BodySize int
	ResData  []byte
	// Keys       map[string]interface{}
	ErrorMsg string
}

// 打印请求日志，可以指定不同的输出样式
func WriteReqLog(p *ReqLogEntity) {
	switch currConfig.style {
	case StyleSdx:
		writeSdxReqLog(p)
	case StyleSdxMini:
		writeSdxReqLog(p)
	case StyleJson:
		writeSdxReqLog(p)
	case StyleJsonMini:
		writeSdxReqLog(p)
	default:
	}
}
