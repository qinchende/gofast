// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/timex"
	"net/http"
	"time"
)

func Logger(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := &ReqLogEntity{
		RawReq: c.ReqRaw,
		// isTerm:  isTerm,
	}
	p.Pms = c.Pms
	p.ClientIP = c.ClientIP()
	p.StatusCode = c.ResWrap.Status()
	p.ResData = c.ResWrap.WrittenData()
	p.BodySize = len(p.ResData)

	// TODO: 内部错误信息一般不返回给调用者，而是打印日志
	p.ErrorMsg = c.Errors.ToString(logx.LogStyleType())

	// Stop timer
	p.TimeStamp = timex.Now()
	p.Latency = p.TimeStamp - c.EnterTime

	// 打印请求日志
	WriteReqLog(p)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

const timeFormatMini = "01-02 15:04:05"

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
	switch logx.LogStyleType() {
	case logx.LogStyleSdx:
		logx.InfoDirect(genSdxReqLogString(p))
	case logx.LogStyleSdxMini:
		logx.InfoDirect(genSdxReqLogString(p))
	case logx.LogStyleJson:
		logx.InfoDirect(genSdxReqLogString(p))
	case logx.LogStyleJsonMini:
		logx.InfoDirect(genSdxReqLogString(p))
	default:
	}
}

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
	var basePms = make(map[string]any)
	if p.Pms["tok"] != nil {
		basePms["tok"] = p.Pms["tok"]
		delete(p.Pms, "tok")
	}

	// 请求参数
	var reqParams []byte
	if p.Pms != nil {
		reqParams, _ = jsonx.Marshal(p.Pms)
	} else if p.RawReq.Form != nil {
		reqParams, _ = jsonx.Marshal(p.RawReq.Form)
	}
	// 请求 核心参数
	reqBaseParams, _ := jsonx.Marshal(basePms)

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
