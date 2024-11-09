// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/fst"
)

func Logger(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := logx.InfoReq()

	p.RawReq = c.Req.Raw
	p.Pms = c.Pms
	p.RemoteAddr = c.ClientIP()
	p.StatusCode = c.Res.Status()
	p.ResData = c.Res.WrittenData()
	p.BodySize = len(p.ResData)

	// 内部错误信息一般不返回给调用者，但是需要打印日志信息
	p.CarryItems = c.CarryMsgItems()

	// Stop timer
	p.TimeStamp = timex.SdxNowDur()
	p.Latency = p.TimeStamp - c.EnterTime

	p.Send()
}

func LoggerMini(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := logx.InfoReq()

	p.RawReq = c.Req.Raw
	p.RemoteAddr = c.ClientIP()
	p.StatusCode = c.Res.Status()
	p.ResData = c.Res.WrittenData()
	p.BodySize = len(p.ResData)

	// Stop timer
	p.TimeStamp = timex.SdxNowDur()
	p.Latency = p.TimeStamp - c.EnterTime

	p.Send()
}
