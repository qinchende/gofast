package gate

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/sysx"
	"time"
)

// 打印每个路由的请求数据。
// TODO: 其实每项路由分钟级的日志应该是收集起来，放入数据库，可视化展示和分析
func (rb *reqBucket) logPrint() {
	for idx, route := range rb.routes {
		if idx != 0 && route.accepts == 0 && route.drops == 0 {
			continue
		}

		qps := float32(route.accepts) / float32(CountInterval/time.Second)
		var aveTime float32
		if route.accepts > 0 {
			aveTime = float32(route.totalTimeMS) / float32(route.accepts)
		}

		logx.StatF("%s | suc: %d, drop: %d, qps: %.1f/s ave: %.1fms, max: %.1fms",
			rb.paths[idx], route.accepts, route.drops, qps, aveTime, float32(route.maxTimeMS))
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
		cpu := sysx.CpuSmoothUsage
		logx.StatF("down-1m | cpu: %d, total: %d, pass: %d, drop: %d", cpu, st.total, st.pass, st.drop)

		//if st.drop == 0 {
		//	logx.Statf("down-1m | cpu: %d, total: %d, pass: %d, drop: %d", cpu, st.total, st.pass, st.drop)
		//} else {
		//	logx.ErrorF("down-1m | cpu: %d, total: %d, pass: %d, drop: %d", cpu, st.total, st.pass, st.drop)
		//}
	}
}
