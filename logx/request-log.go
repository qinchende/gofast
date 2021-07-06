package logx

import (
	"net/http"
	"time"
)

// 日志参数实体
type ReqLogParams struct {
	RawReq     *http.Request
	TimeStamp  time.Time
	Latency    time.Duration
	StatusCode int
	ClientIP   string
	// isTerm     bool
	Pms        map[string]interface{}
	BodySize   int
	WriteBytes *[]byte
	// Keys       map[string]interface{}
	ErrorMsg string
}

// 打印请求日志，可以指定不同的输出样式
func WriteReqLog(p *ReqLogParams) {
	switch currConfig.style {
	case styleSdx:
		writeSdxReqLog(p)
	case styleSdxMini:
		writeSdxReqLog(p)
	case styleJson:
		writeSdxReqLog(p)
	case styleJsonMini:
		writeSdxReqLog(p)
	default:
	}
}
