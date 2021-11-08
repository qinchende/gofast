package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/timex"
)

//func ReqLogger(w *fst.GFResponse, r *http.Request) {
//	// 执行完后面的请求，再打印日志
//	w.NextFit(r)
//
//	// 请求处理完，并成功返回了，接下来就是打印请求日志
//	p := &logx.ReqLogEntity{
//		RawReq: r,
//		// isTerm:  isTerm,
//	}
//	if w.Ctx != nil {
//		p.Pms = w.Ctx.Pms
//	}
//	p.ClientIP = w.ClientIP(r)
//	p.StatusCode = w.ResWrap.Status()
//	p.WriteBytes = &w.ResWrap.WriteBytes
//	p.BodySize = w.ResWrap.Size()
//
//	// TODO: 内部错误信息一般不返回给调用者，而是打印日志
//	p.ErrorMsg = w.Errors.String(logx.Style())
//
//	// Stop timer
//	p.TimeStamp = timex.Now()
//	p.Latency = p.TimeStamp - w.EnterTime
//
//	// 打印请求日志
//	logx.WriteReqLog(p)
//}

func ReqLogger(ctx *fst.Context) {
	//if ctx == nil {
	//	return
	//}

	// 执行完后面的请求，再打印日志
	ctx.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := &logx.ReqLogEntity{
		RawReq: ctx.ReqRaw,
		// isTerm:  isTerm,
	}
	p.Pms = ctx.Pms
	p.ClientIP = ctx.ClientIP()
	p.StatusCode = ctx.ResWrap.Status()
	p.WriteBytes = &ctx.ResWrap.WriteBytes
	p.BodySize = ctx.ResWrap.Size()

	// TODO: 内部错误信息一般不返回给调用者，而是打印日志
	p.ErrorMsg = ctx.Errors.String(logx.Style())

	// Stop timer
	p.TimeStamp = timex.Now()
	p.Latency = p.TimeStamp - ctx.EnterTime

	// 打印请求日志
	logx.WriteReqLog(p)
}
