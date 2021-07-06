package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"net/http"
	"time"
)

func ReqLogger() fst.IncHandler {
	return func(w *fst.GFResponse, r *http.Request) {
		// Start timer
		start := time.Now()

		// 执行完后面的请求，再打印日志
		w.NextFit(r)

		// 请求处理完，并成功返回了，接下来就是打印请求日志
		p := &logx.ReqLogParams{
			RawReq: r,
			// isTerm:  isTerm,
		}
		if w.Ctx != nil {
			p.Pms = w.Ctx.Pms
		}
		p.ClientIP = w.ClientIP(r)
		p.StatusCode = w.ResWrap.Status()
		p.ErrorMsg = w.Errors.ByType(fst.ErrorTypePrivate).String()
		p.WriteBytes = &w.ResWrap.WriteBytes
		p.BodySize = w.ResWrap.Size()

		// Stop timer
		p.TimeStamp = time.Now()
		p.Latency = p.TimeStamp.Sub(start)

		// 打印请求日志
		logx.WriteReqLog(p)

		// TODO: 错误信息返回给调用端，这个地方是否要打开？
		//if p.ErrorMsg != "" {
		//	w.ResWrap.WriteString(p.ErrorMsg)
		//}
	}
}
