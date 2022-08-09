// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
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
		InfoStrDirect(genSdxReqLogString(p))
	case LogStyleSdxMini:
		InfoStrDirect(genSdxReqLogString(p))
	case LogStyleJson:
		InfoStrDirect(genSdxReqLogString(p))
	case LogStyleJsonMini:
		InfoStrDirect(genSdxReqLogString(p))
	default:
	}
}
