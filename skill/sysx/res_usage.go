package sysx

import (
	"fmt"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/shirou/gopsutil/v3/cpu"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	cpuRefreshInterval   = time.Second * 3 // 查询CPU使用率的最小间隔
	beta                 = 0.6             // 估算最近CPU使用率的参数
	printRefreshInterval = time.Minute     // 打印系统资源的时间间隔
)

var (
	CpuChecked bool // CPU 资源利用率监控是否启用

	cpuCurRate int32 // CPU 的利用率 * 100 （放大100倍，方便统计小数点后两位）最近3秒内的平均利用率
	cpuAveRate int32 // CPU 的利用率 * 100 （放大100倍，方便统计小数点后两位）最近一段时间利用率

	cpuUseSum float64 // 一段时间利用率求和
	cpuUseCt  int32   // 求和次数
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

		for {
			select {
			case <-cpuTicker.C: // 3秒统计一次
				gmp.RunSafe(func() {
					// curUsage := cpux.RefreshCpu()
					// TODO: 这里可能会有一个大BUG。如果别的地方也用下面这个方法统计CPU使用率
					// TODO: 下面这个方法统计的是上一次执行到本次执行这段时间内CPU使用率
					avg, _ := cpu.Percent(0, false)
					curUsage := avg[0] * 100
					cpuUseSum += curUsage
					cpuUseCt++

					prevUsage := atomic.LoadInt32(&cpuAveRate)
					// cpu = cpuᵗ⁻¹ * (1 - beta) + cpuᵗ * beta
					usage := int32(float64(prevUsage)*(1-beta) + curUsage*beta)
					atomic.StoreInt32(&cpuAveRate, usage)
					atomic.StoreInt32(&cpuCurRate, int32(curUsage))
				})
			case <-printTicker.C: // 60秒打印一次
				printUsage()
				cpuUseSum = 0.0
				cpuUseCt = 0
			}
		}
	}()
}

// 打印周期内，CPU平均利用率
func CpuSmoothUsage() float32 {
	return float32(atomic.LoadInt32(&cpuAveRate)) / 100
}

// CPU当前统计周期内使用率
func CpuCurUsage() float32 {
	return float32(atomic.LoadInt32(&cpuCurRate)) / 100
}

// 切片中每个值都保留2位小数
func decimalPlace2(arr []float64) {
	for idx := range arr {
		arr[idx], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", arr[idx]), 64)
	}
}

func bToMb(b uint64) float32 {
	return float32(b) / 1024 / 1024
}

// 打印系统资源的使用情况统计
func printUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logx.Statf("CPU: %.2f, Mem: [%.1fMB, %.1fMB, %.1fMB], Gor: %d, GC: %d",
		cpuUseSum/float64(cpuUseCt)/100, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), runtime.NumGoroutine(), m.NumGC)
}
