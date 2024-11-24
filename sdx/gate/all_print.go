package gate

import (
	"github.com/qinchende/gofast/aid/fuse"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/aid/proc"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/lang"
	"time"
)

const specialRouteMethod = "NA@"

func (rb *reqCounter) logPrintReqCounter(data *printData) {
	// 输出扩展统计
	for idx := 0; idx < len(data.extras); idx++ {
		extra := &data.extras[idx]
		if extra.total == 0 {
			continue
		}
		logx.InfoStat().Any("Status", cst.KV{
			"typ": logx.LogStatRouteReq.Type,
			"pth": specialRouteMethod + rb.extraPaths[idx],
			//"fls": []string{"accept", "timeout", "drop", "qps", "ave", "max"}
			"val": [6]any{extra.total, 0, 0, 0.00, 0.00, 0},
		}).Send()
	}

	// 输出路由统计
	for idx := 0; idx < len(data.routes); idx++ {
		rt := &data.routes[idx]
		if rt.accepts == 0 && rt.drops == 0 {
			continue
		}

		qps := float32(rt.accepts) / float32(countInterval/time.Second)
		var aveTimeMS int64
		if rt.accepts > 0 {
			aveTimeMS = rt.totalTimeMS / int64(rt.accepts)
		}

		logx.InfoStat().Any("Status", cst.KV{
			"typ": logx.LogStatRouteReq.Type,
			"pth": rb.paths[idx],
			//"fls": []string{"accept", "timeout", "drop", "qps", "ave", "max"}
			"val": [6]any{rt.accepts, rt.timeouts, rt.drops, lang.Round32(qps, 2), aveTimeMS, rt.maxTimeMS},
		}).Send()
	}
}

func (bk *Breaker) LogError(err error) {
	bk.reduceLog.DoInterval(false, func(skipTimes int32) {
		if err != fuse.ErrServiceUnavailable {
			return
		}
		logx.InfoStat().
			Uint("typ", uint(logx.LogStatBreakerOpen.Type)).
			Str("proc", proc.ProcessName()+"/"+lang.ToString(proc.Pid())).
			Str("callee", bk.name).
			Int("skip", int(skipTimes)).
			Msg(bk.Breaker.Errors(","))
	})
}
