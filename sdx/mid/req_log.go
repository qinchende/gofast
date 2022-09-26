// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/timex"
)

func Logger(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := &logx.ReqLogEntity{
		RawReq: c.ReqRaw,
	}
	p.Pms = c.Pms
	p.ClientIP = c.ClientIP()
	p.StatusCode = c.ResWrap.Status()
	p.ResData = c.ResWrap.WrittenData()
	p.BodySize = len(p.ResData)

	// 内部错误信息一般不返回给调用者，但是需要打印日志信息
	p.MsgBaskets = c.MsgBaskets()

	// Stop timer
	p.TimeStamp = timex.Now()
	p.Latency = p.TimeStamp - c.EnterTime

	// 打印请求日志
	logx.PrintReqLog(p)
}

func LoggerMini(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := &logx.ReqLogEntity{
		RawReq: c.ReqRaw,
	}
	p.ClientIP = c.ClientIP()
	p.StatusCode = c.ResWrap.Status()
	p.ResData = c.ResWrap.WrittenData()
	p.BodySize = len(p.ResData)

	// Stop timer
	p.TimeStamp = timex.Now()
	p.Latency = p.TimeStamp - c.EnterTime

	// 打印请求日志
	logx.PrintReqLogMini(p)
}
