package gate

import (
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/executors"
	"github.com/qinchende/gofast/skill/load"
)

// 请求统计管理员，负责分析每个路由的请求压力和处理延时情况
type RequestKeeper struct {
	// 访问量统计 Counter
	bucket  *reqContainer
	counter *executors.IntervalExecutor

	// 熔断器
	Breakers []breaker.Breaker

	// 降载组件
	Shedding     load.Shedder
	SheddingStat *sheddingStat
}
