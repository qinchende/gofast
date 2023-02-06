// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// 启动链路追踪
func Tracing(appName string, useTracing bool) fst.CtxHandler {
	if useTracing == false {
		return nil
	}

	return func(c *fst.Context) {
		propagator := otel.GetTextMapPropagator()
		tracer := otel.GetTracerProvider().Tracer(trace.TraceName)

		ctx := propagator.Extract(c.ReqRaw.Context(), propagation.HeaderCarrier(c.ReqRaw.Header))
		spanName := c.FullPath()
		if len(spanName) == 0 {
			spanName = c.ReqRaw.URL.Path
		}
		spanCtx, span := tracer.Start(
			ctx,
			spanName,
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(appName, spanName, c.ReqRaw)...),
		)
		defer span.End()

		// convenient for tracking error messages
		propagator.Inject(spanCtx, propagation.HeaderCarrier(c.ResWrap.Header()))
		c.ReqRaw = c.ReqRaw.WithContext(spanCtx)

		c.Next()
	}
}

//
//// TracingHandler return a middleware that process the opentelemetry.
//func TracingHandler(serviceName, path string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		propagator := otel.GetTextMapPropagator()
//		tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
//
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
//			spanName := path
//			if len(spanName) == 0 {
//				spanName = r.URL.Path
//			}
//			spanCtx, span := tracer.Start(
//				ctx,
//				spanName,
//				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
//				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(
//					serviceName, spanName, r)...),
//			)
//			defer span.End()
//
//			// convenient for tracking error messages
//			propagator.Inject(spanCtx, propagation.HeaderCarrier(w.Header()))
//			next.ServeHTTP(w, r.WithContext(spanCtx))
//		})
//	}
//}
