package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"net/http"
	"time"
)

func ReqLogger(logType string) fst.IncHandler {
	if logType == "" {
		logType = fst.LogTypeConsole
	}

	return func(w *fst.GFResponse, r *http.Request) {
		// Start timer
		start := time.Now()
		path := r.URL.Path
		// raw := r.URL.RawQuery

		// time.Sleep(1 * time.Second)
		// 执行完后面的请求，再打印日志
		w.NextFit(r)

		p := &logx.ReqLogParams{
			Request: r,
			//isTerm:  isTerm,
		}
		//if w.PCtx != nil {
		//	p.Keys = w.PCtx.Keys
		//}

		p.ClientIP = w.ClientIP(r)
		p.Method = r.Method
		p.StatusCode = w.ResW.Status()
		p.ErrorMsg = w.Errors.ByType(fst.ErrorTypePrivate).String()
		p.WriteBytes = &w.ResW.WriteBytes
		p.BodySize = w.ResW.Size()
		//if raw != "" {
		//	path = path + "?" + raw
		//}
		p.Path = path

		// Stop timer
		p.TimeStamp = time.Now()
		p.Latency = p.TimeStamp.Sub(start)

		// 打印请求日志
		logx.WriteReqLog(p)
		// 错误信息返回给调用端
		//if p.ErrorMsg != "" {
		//	w.ResW.WriteString(p.ErrorMsg)
		//}
	}
}
