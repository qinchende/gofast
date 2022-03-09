package sysx

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/sysx/cpux"
	"github.com/shirou/gopsutil/v3/cpu"
	"runtime"
	"time"
)

const (
	cpuRefreshInterval   = time.Second * 3 // 查询CPU使用率的最小间隔
	beta                 = 0.7             // 估算最近CPU使用率的参数
	printRefreshInterval = time.Minute     // 打印系统资源的时间间隔
)

var (
	CpuChecked bool // CPU 资源利用率监控是否启用

	CpuCurUsage    float64 // CPU 最近3秒内的平均利用率
	CpuSmoothUsage float64 // CPU 最近一段时间利用率

	cpuStatStep  []cpu.TimesStat
	cpuStatPrint []cpu.TimesStat
)

// 启动CPU和内存统计
func StartCpuCheck() {
	CpuChecked = true
	//// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++三方包CPU测试
	//// 打印CPU基础信息
	////infos, _ := cpu.Info()
	////logx.Info(infos)
	//
	//// 利用第三方包计算CPU的占用情况
	//totalPercent, _ := cpu.Percent(3*time.Second, false)
	//perPercents, _ := cpu.Percent(3*time.Second, true)
	//
	//decimalPlace2(totalPercent)
	//decimalPlace2(perPercents)
	//logx.Infof("CPU-Usage -> total: %v, per: %v", totalPercent, perPercents)
	//// ==++ NED

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
				cpuStatPrintLast := cpuStatPrint
				cpuStatPrint, _ = cpu.Times(false)
				printUsage(cpux.BusyPercent(cpuStatPrintLast[0], cpuStatPrint[0]))
			}
		}
	}()
}

// 打印系统资源的使用情况统计
func printUsage(cpuUsage float64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logx.Statf("CPU: [%.2f,%.2f,%.2f], Mem: [%.1fMB, %.1fMB, %.1fMB], Gor: %d, GC: %d",
		CpuCurUsage, CpuSmoothUsage, cpuUsage, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), runtime.NumGoroutine(), m.NumGC)
}

//// CPU当前统计周期内使用率
//func CpuCurUsage() float32 {
//	return float32(atomic.LoadInt32(&cpuCurrentRate)) / 100
//}
//
//// 打印周期内，CPU平均利用率
//func CpuSmoothUsage() float32 {
//	return float32(atomic.LoadInt32(&cpuSmoothRate)) / 100
//}

// 字节 到 MB 的转换.
func bToMb(b uint64) float32 {
	return float32(b) / 1024 / 1024
}

//// 切片中每个值都保留2位小数
//func decimalPlace2(arr []float64) {
//	for idx := range arr {
//		arr[idx], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", arr[idx]), 64)
//	}
//}
