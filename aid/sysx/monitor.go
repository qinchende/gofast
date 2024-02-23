package sysx

import (
	"github.com/qinchende/gofast/aid/gmp"
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/aid/sysx/cpux"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/logx"
	"github.com/shirou/gopsutil/v3/cpu"
	"runtime"
	"sync"
	"time"
)

const (
	cpuRefreshInterval   = time.Second * 3 // 查询CPU使用率的最小间隔
	beta                 = 0.75            // 估算最近CPU使用率的参数
	printRefreshInterval = time.Minute     // 打印系统资源的时间间隔
)

var (
	MonitorStarted bool    // 资源利用率监控是否启用
	CpuCurUsage    float64 // CPU 最近3秒内的平均利用率
	CpuSmoothUsage float64 // CPU 最近一段时间平滑利用率

	ckLook       sync.Mutex
	cpuStatStep  []cpu.TimesStat
	cpuStatPrint []cpu.TimesStat
)

// 启动系统资源统计
func OpenSysMonitor(print bool) {
	ckLook.Lock()
	defer ckLook.Unlock()

	if MonitorStarted {
		return
	}
	MonitorStarted = true

	go func() {
		cpuTicker := time.NewTicker(cpuRefreshInterval)
		defer cpuTicker.Stop()
		printTicker := time.NewTicker(printRefreshInterval)
		defer printTicker.Stop()

		cpuStatStep, _ = cpu.Times(false)
		cpuStatPrint = cpuStatStep
		for {
			select {
			case <-cpuTicker.C: // 3秒统计一次
				gmp.RunSafe(func() {
					// 当前小周期内的使用率
					cpuStatStepLast := cpuStatStep
					cpuStatStep, _ = cpu.Times(false)
					CpuCurUsage = cpux.BusyPercent(cpuStatStepLast[0], cpuStatStep[0])

					// 平滑使用率
					// cpu = cpuᵗ⁻¹ * (1 - beta) + cpuᵗ * beta
					CpuSmoothUsage = CpuSmoothUsage*(1-beta) + CpuCurUsage*beta
				})
			case <-printTicker.C: // 60秒打印一次
				// 启用CPU利用率的统计，但并不意味要打印状态信息
				if print {
					cpuStatPrintLast := cpuStatPrint
					cpuStatPrint, _ = cpu.Times(false)
					printSysResourceStatus(cpux.BusyPercent(cpuStatPrintLast[0], cpuStatPrint[0]))
				}
			}
		}
	}()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 打印系统资源的使用情况统计
// const resJsonStr = `{"CPU":[%.2f,%.2f],"Mem":[%.1f,%.1f,%.1f],"Gor":%d,"GC":%d}`
func printSysResourceStatus(cpuAvaUsage float64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logx.StatKV(cst.KV{
		"typ": logx.LogStatSysMonitor.Type,
		//"fls": []string{"cpu", "mem", "gor", "gc"},
		"val": []any{
			[2]float64{lang.Round64(CpuSmoothUsage, 2), lang.Round64(cpuAvaUsage, 2)},
			[3]float32{bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys)},
			runtime.NumGoroutine(),
			m.NumGC,
		},
	})
}

// 字节 到 MB 的转换.
func bToMb(b uint64) float32 {
	return lang.Round32(float32(b)/1024/1024, 2)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func testPsutil() {
//	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++三方包CPU测试
//	// 打印CPU基础信息
//	//infos, _ := cpu.Info()
//	//logx.Info(infos)
//
//	// 利用第三方包计算CPU的占用情况
//	totalPercent, _ := cpu.Percent(3*time.Second, false)
//	perPercents, _ := cpu.Percent(3*time.Second, true)
//
//	decimalPlace2(totalPercent)
//	decimalPlace2(perPercents)
//	logx.InfoF("CPU-Usage -> total: %v, per: %v", totalPercent, perPercents)
//	// ==++ NED
//}
//
//// 切片中每个值都保留2位小数
//func decimalPlace2(arr []float64) {
//	for idx := range arr {
//		arr[idx], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", arr[idx]), 64)
//	}
//}

//// CPU当前统计周期内使用率
//func CpuCurUsage() float32 {
//	return float32(atomic.LoadInt32(&cpuCurrentRate)) / 100
//}
//
//// 打印周期内，CPU平均利用率
//func CpuSmoothUsage() float32 {
//	return float32(atomic.LoadInt32(&cpuSmoothRate)) / 100
//}
//
