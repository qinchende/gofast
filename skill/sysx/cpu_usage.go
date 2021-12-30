package sysx

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/sysx/cpu"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	// 1s and 0.80 as beta will count the average cpu load for past 5 seconds
	// cpuRefreshInterval   = time.Millisecond * 250
	cpuRefreshInterval = time.Second * 1
	// moving average beta hyperparameter
	beta = 0.8

	// 日志输出的时间间隔
	printRefreshInterval = time.Minute
)

var (
	cpuUsage   int64 // CPU的利用率
	CpuChecked bool  // CPU 资源利用率监控是否启用
)

// 启动CPU和内存统计
func StartCpuCheck() {
	CpuChecked = true

	go func() {
		cpuTicker := time.NewTicker(cpuRefreshInterval)
		defer cpuTicker.Stop()
		printTicker := time.NewTicker(printRefreshInterval)
		defer printTicker.Stop()

		for {
			select {
			case <-cpuTicker.C:
				gmp.RunSafe(func() {
					curUsage := cpu.RefreshCpu()
					prevUsage := atomic.LoadInt64(&cpuUsage)
					// cpu = cpuᵗ⁻¹ * beta + cpuᵗ * (1 - beta)
					usage := int64(float64(prevUsage)*beta + float64(curUsage)*(1-beta))
					atomic.StoreInt64(&cpuUsage, usage)
				})
			case <-printTicker.C:
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

	logx.Statf("CPU: %d, Mem: [%.1fMB, %.1fMB, %.1fMB], NumGo: %d, NumGC: %d",
		CpuUsage(), bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), runtime.NumGoroutine(), m.NumGC)
}
