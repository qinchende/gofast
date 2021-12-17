package test

import (
	"github.com/qinchende/gofast/skill/collection"
	"runtime"
	"testing"
	"time"
)

func init() {
	runtime.GOMAXPROCS(4)
}

//PS gofast\skill\collection\test> go test -bench=Go* -benchmem -benchtime=10s
//goos: windows
//goarch: amd64
//pkg: github.com/qinchende/gofast/skill/collection/test
//cpu: Intel(R) Core(TM) i7-10700 CPU @ 2.90GHz
//BenchmarkGozeroRW-4      8464778              1282 ns/op               2 B/op          0 allocs/op
//BenchmarkGofastRW-4     59552853               201.4 ns/op             0 B/op          0 allocs/op
//PASS
//ok      github.com/qinchende/gofast/skill/collection/test       34.894s
// 看看上面的性能差异，桶越多 go-zero的性能越差，和单个桶的时间间隔差别不大

const duration = time.Millisecond * 50
const winSize = 100
const concurrencyNum = 100000
const loopTimes = 1

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// go-zero 滑动窗口的性能
func BenchmarkGozeroRW(b *testing.B) {
	rWin := collection.NewRollingWindow(winSize, duration)

	b.ReportAllocs()
	b.ResetTimer()

	// 并发测试模式
	b.SetParallelism(concurrencyNum)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			goZeroRollingWindow(rWin)
		}
	})
}

func goZeroRollingWindow(rw *collection.RollingWindow) {
	for i := 0; i < loopTimes; i++ {
		var accepts, total int64
		rw.Reduce(func(b *collection.Bucket) {
			accepts += int64(b.Sum)
			total += b.Count
		})
		rw.Add(1)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// gofast 滑动窗口的性能
func BenchmarkGofastRW(b *testing.B) {
	rWin := collection.NewRollingWindowSdx(winSize, duration)

	b.ReportAllocs()
	b.ResetTimer()

	// 并发测试模式
	b.SetParallelism(concurrencyNum)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gofastRollingWindow(rWin)
		}
	})
}

func gofastRollingWindow(rw *collection.RollingWindowSdx) {
	for i := 0; i < loopTimes; i++ {
		rw.CurrWinValue()
		rw.Add(1)
	}
}
