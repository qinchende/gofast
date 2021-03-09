package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/stat"
	"github.com/qinchende/gofast/skill/timex"
	"net/http"
)

func Metric(metrics *stat.Metrics) fst.IncHandler {
	return func(w *fst.GFResponse, r *http.Request) {
		startTime := timex.Now()
		defer func() {
			metrics.Add(stat.Task{
				Duration: timex.Since(startTime),
			})
		}()
		w.NextFit(r)
	}
}
