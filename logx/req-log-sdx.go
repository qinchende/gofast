// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"encoding/json"
	"fmt"
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

// 通过模板构造字符串可能性能更好。
var genSdxReqLogString = func(p *ReqLogEntity) string {
	formatStr := `
[%s] %s (%s/%s) %d/%d [%d]
  B: %s
  P: %s
  R: %s%s
`
	// 最长打印出 1024个字节的结果
	tLen := len(p.ResData)
	if tLen > 1024 {
		tLen = 1024
	}

	// 这个时候可以随意改变 p.Pms ，这是请求最后一个执行的地方了
	var basePms = make(map[string]interface{})
	if p.Pms["tok"] != nil {
		basePms["tok"] = p.Pms["tok"]
		delete(p.Pms, "tok")
	}

	// 请求参数
	var reqParams []byte
	if p.Pms != nil {
		reqParams, _ = json.Marshal(p.Pms)
	} else if p.RawReq.Form != nil {
		reqParams, _ = json.Marshal(p.RawReq.Form)
	}
	// 请求 核心参数
	reqBaseParams, _ := json.Marshal(basePms)

	return fmt.Sprintf(formatStr,
		p.RawReq.Method,
		p.RawReq.URL.Path,
		p.ClientIP,
		timex.ToTime(p.TimeStamp).Format(timeFormatMini),
		p.StatusCode,
		p.BodySize,
		p.Latency/time.Millisecond,
		reqBaseParams,
		reqParams,
		(p.ResData)[:tLen],
		p.ErrorMsg,
	)
}
