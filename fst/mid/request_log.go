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
		raw := r.URL.RawQuery

		// time.Sleep(1 * time.Second)
		// 执行完后面的请求，再打印日志
		w.NextFit(r)

		p := &logx.LogReqParams{
			Request: r,
			// Keys:    r.Keys,
			// isTerm:  isTerm,
		}

		// Stop timer
		p.TimeStamp = time.Now()
		p.Latency = p.TimeStamp.Sub(start)

		p.ClientIP = w.ClientIP(r)
		p.Method = r.Method
		p.StatusCode = w.ResW.Status()
		p.ErrorMessage = w.Errors.ByType(fst.ErrorTypePrivate).String()
		p.BodySize = w.ResW.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		p.Path = path
		logx.WriteReqLog(p)

		// 错误信息也写给客户端
		w.ResW.WriteString(p.ErrorMessage)
	}
}
