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
	Pms        map[string]string
	BodySize   int
	WriteBytes *[]byte
	// Keys       map[string]interface{}
	ErrorMsg string
}

var GenReqLogString = func(p *ReqLogParams) string {
	formatStr := `
[%s] %s (%s/%s) %d/%d [%d]
  P: %s
  R: %s
  E: %s
`
	// 最长打印出 1024个字节的结果
	tLen := len(*p.WriteBytes)
	if tLen > 1024 {
		tLen = 1024
	}

	// 请求参数
	var reqParams []byte
	if p.Pms != nil {
		reqParams, _ = json.Marshal(p.Pms)
	} else if p.Request.Form != nil {
		reqParams, _ = json.Marshal(p.Request.Form)
	}

	return fmt.Sprintf(formatStr,
		p.Method,
		p.Path,
		p.ClientIP,
		p.TimeStamp.Format("01-02 15:04:05"),
		p.StatusCode,
		p.BodySize,
		p.Latency/time.Millisecond,
		reqParams,
		(*p.WriteBytes)[:tLen],
		p.ErrorMsg,
	)
}

func WriteReqLog(p *ReqLogParams) {
	writeStringNow(GenReqLogString(p))
}

//func outputJson(writer io.Writer, info interface{}) {
//	if content, err := json.Marshal(info); err != nil {
//		log.Println(err.Error())
//	} else if atomic.LoadUint32(&initialized) == 0 || writer == nil {
//		log.Println(string(content))
//	} else {
//		writer.Write(append(content, '\n'))
//	}
//}
