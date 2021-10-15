package stat

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/stat/tips"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	// 250ms and 0.95 as beta will count the average cpu load for past 5 seconds
	cpuRefreshInterval = time.Millisecond * 250
	allRefreshInterval = time.Minute
	// moving average beta hyperparameter
	beta = 0.95
)

var cpuUsage int64

func StartCpuMemCollect() {
	go func() {
		cpuTicker := time.NewTicker(cpuRefreshInterval)
		defer cpuTicker.Stop()
		allTicker := time.NewTicker(allRefreshInterval)
		defer allTicker.Stop()

		for {
			select {
			case <-cpuTicker.C:
				gmp.RunSafe(func() {
					curUsage := tips.RefreshCpu()
					prevUsage := atomic.LoadInt64(&cpuUsage)
					// cpu = cpuᵗ⁻¹ * beta + cpuᵗ * (1 - beta)
					usage := int64(float64(prevUsage)*beta + float64(curUsage)*(1-beta))
					atomic.StoreInt64(&cpuUsage, usage)
				})
			case <-allTicker.C:
				printUsage()
			}
		}
	}()
}

func CpuUsage() int64 {
	return atomic.LoadInt64(&cpuUsage)
}

func bToMb(b uint64) float32 {
	return float32(b) / 1024 / 1024
}

func printUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logx.Statf("CPU: %dm, MEMORY: Alloc=%.1fMi, TotalAlloc=%.1fMi, Sys=%.1fMi, NumGC=%d",
		CpuUsage(), bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}
