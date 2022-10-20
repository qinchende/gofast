// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/sysx/host"
	"github.com/qinchende/gofast/skill/trace"
)

//// 启动链路追踪
//func Tracing(w *fst.GFResponse, r *http.Request) {
//	// 先禁用这个功能
//	if w != nil {
//		return
//	}
//
//	carrier, err := trace.Extract(trace.HttpFormat, r.Header)
//	// ErrInvalidCarrier means no trace id was set in http header
//	if err != nil && err != trace.ErrInvalidCarrier {
//		logx.Error(err)
//	}
//
//	ctx, span := trace.StartServerSpan(r.Context(), carrier, sysx.Hostname(), r.RequestURI)
//	defer span.Finish()
//	r = r.WithContext(ctx)
//
//	w.NextFit(r)
//}

// 启动链路追踪

func Tracing(useTracing bool) fst.CtxHandler {
	if useTracing == false {
		return nil
	}

	return func(c *fst.Context) {
		carrier, err := trace.Extract(trace.HttpFormat, c.ReqRaw.Header)
		// ErrInvalidCarrier means no trace id was set in http header
		if err != nil && err != trace.ErrInvalidCarrier {
			logx.Error(err.Error())
		}

		newCtx, span := trace.StartServerSpan(c.ReqRaw.Context(), carrier, host.Hostname(), c.ReqRaw.RequestURI)
		defer span.Finish()
		c.ReqRaw = c.ReqRaw.WithContext(newCtx)

		// 有 defer ，这里的 ctx.Next() 有意义
		c.Next()
	}
}
