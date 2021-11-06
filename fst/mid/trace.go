package mid

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

//// 启动链路追踪
//func Tracing(ctx *fst.Context) {
//	// 先禁用这个功能
//	if ctx != nil {
//		return
//	}
//
//	carrier, err := trace.Extract(trace.HttpFormat, ctx.ReqRaw.Header)
//	// ErrInvalidCarrier means no trace id was set in http header
//	if err != nil && err != trace.ErrInvalidCarrier {
//		logx.Error(err)
//	}
//
//	ctx, span := trace.StartServerSpan(ctx.ReqRaw.Context(), carrier, sysx.Hostname(), ctx.ReqRaw.RequestURI)
//	defer span.Finish()
//	ctx.ReqRaw = ctx.ReqRaw.WithContext(ctx)
//
//	ctx.Next()
//}
