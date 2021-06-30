// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ReqLogParams struct {
	Request    *http.Request
	Method     string
	Path       string
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

var genSdxReqLogString = func(p *ReqLogParams) string {
	formatStr := `
[%s] %s (%s/%s) %d/%d [%d]
  B: %s
  P: %s
  R: %s
  E: %s
`
	// 最长打印出 1024个字节的结果
	tLen := len(*p.WriteBytes)
	if tLen > 1024 {
		tLen = 1024
	}

	// 这个时候可以随意改变 p.Pms ，这是请求最后一个执行的地方了
	var basePms = make(map[string]interface{})
	if p.Pms["tok"] != "" {
		basePms["tok"] = p.Pms["tok"]
		delete(p.Pms, "tok")
	}

	// 请求参数
	var reqParams []byte
	if p.Pms != nil {
		reqParams, _ = json.Marshal(p.Pms)
	} else if p.Request.Form != nil {
		reqParams, _ = json.Marshal(p.Request.Form)
	}
	// 请求 核心参数
	reqBaseParams, _ := json.Marshal(basePms)

	return fmt.Sprintf(formatStr,
		p.Method,
		p.Path,
		p.ClientIP,
		p.TimeStamp.Format(timeFormatMini),
		p.StatusCode,
		p.BodySize,
		p.Latency/time.Millisecond,
		reqBaseParams,
		reqParams,
		(*p.WriteBytes)[:tLen],
		p.ErrorMsg,
	)
}

func WriteSdxReqLog(p *ReqLogParams) {
	infoSync(genSdxReqLogString(p))
}
