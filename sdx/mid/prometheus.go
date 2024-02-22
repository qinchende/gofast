// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/aid/metric"
	"github.com/qinchende/gofast/fst"
)

const serverNamespace = "http_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "code"},
	})
)

//func PrometheusHandler(path string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			startTime := timex.Now()
//			cw := &security.WithCodeResponseWriter{Writer: w}
//			defer func() {
//				metricServerReqDur.Observe(int64(timex.Since(startTime)/time.Millisecond), path)
//				metricServerReqCodeTotal.Inc(path, strconv.Itoa(cw.Code))
//			}()
//
//			next.ServeHTTP(cw, r)
//		})
//	}
//}

func Prometheus(c *fst.Context) {
	//startTime := timex.Now()
	//cw := &security.WithCodeResponseWriter{Writer: w}
	//defer func() {
	// metricServerReqDur.Observe(int64(timex.Since(startTime)/time.Millisecond), path)
	// metricServerReqCodeTotal.Inc(path, strconv.Itoa(cw.Code))
	//}()
}
