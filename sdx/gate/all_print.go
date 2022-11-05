package gate

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/sysx"
	"time"
)

const specialRouteMethod = "NA@"

func (rb *reqBucket) logPrintReqCounter(data *printData) {
	// 输出扩展统计
	for idx := 0; idx < len(data.extras); idx++ {
		total := data.extras[idx]
		if total == 0 {
			continue
		}
		logx.StatKV(cst.KV{
			"typ": logx.LogStatRouteReq.Type,
			"pth": specialRouteMethod + rb.extraPaths[idx],
			//"fls": []string{"suc", "drop", "qps", "ave", "max"}
			"val": [5]any{total, 0, 0.00, 0.00, 0},
		})
	}

	// 输出路由统计
	for idx := 0; idx < len(data.routes); idx++ {
		route := &data.routes[idx]
		if route.accepts == 0 && route.drops == 0 {
			continue
		}

		qps := float32(route.accepts) / float32(CountInterval/time.Second)
		var aveTimeMS float32
		if route.accepts > 0 {
			aveTimeMS = float32(route.totalTimeMS) / float32(route.accepts)
		}

		logx.StatKV(cst.KV{
			"typ": logx.LogStatRouteReq.Type,
			"pth": rb.paths[idx],
			//"fls": []string{"suc", "drop", "qps", "ave", "max"}
			"val": [5]any{route.accepts, route.drops, lang.Round32(qps, 2), lang.Round32(aveTimeMS, 2), route.maxTimeMS},
		})
	}
}

// 单独的协程运行这个定时任务。启动定时日志输出
func (s *sheddingStat) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// 定时器，每分钟执行一次，死循环
	for range ticker.C {
		st := s.reset()
		if st.total == 0 && st.pass == 0 && st.drop == 0 {
			continue
		}

		logx.StatKV(cst.KV{
			"typ": logx.LogStatCpuUsage.Type,
			//"fls": [5]string{"cpu", "total", "pass", "drop"},
			"val": []any{sysx.CpuSmoothUsage, st.total, st.pass, st.drop},
		})
	}
}
