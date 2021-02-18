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

	return func(w http.ResponseWriter, r *fst.Request) {
		// Start timer
		start := time.Now()
		path := r.RawReq.URL.Path
		raw := r.RawReq.URL.RawQuery

		r.NextFit(w)

		p := &logx.LogReqParams{
			Request: r.RawReq,
			// Keys:    r.Keys,
			// isTerm:  isTerm,
		}

		// Stop timer
		p.TimeStamp = time.Now()
		p.Latency = p.TimeStamp.Sub(start)

		p.ClientIP = r.ClientIP()
		p.Method = r.RawReq.Method
		//p.StatusCode = .Status()
		p.ErrorMessage = r.Errors.ByType(fst.ErrorTypePrivate).String()
		//p.BodySize = r.Reply.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		p.Path = path
		logx.WriteReqLog(p)
	}
}
