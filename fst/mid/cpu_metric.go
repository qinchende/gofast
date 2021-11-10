package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/stat"
	"net/http"
	"time"
)

// ++++++++++++++++++++++ add by cd.net 2021.10.14
// 总说：定时统计（间隔60秒）系统资源利用情况 | 请求处理相应性能 | 请求量 等
//

//func CpuMetric(metrics *stat.Metrics) http.HandlerFunc {
//	if metrics == nil {
//		return nil
//	}
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		start := time.Now()
//		defer func() {
//			metrics.AddItem(stat.ReqItem{
//				Duration: time.Now().Sub(start),
//			})
//		}()
//
//	}
//}

func CpuMetric(metrics *stat.Metrics) fst.FitFunc {
	if metrics == nil {
		return nil
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				metrics.AddItem(stat.ReqItem{
					Duration: time.Now().Sub(start),
				})
			}()

			next(w, r)
		}
	}
}
