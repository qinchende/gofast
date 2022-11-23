// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package collect

import (
	"github.com/qinchende/gofast/skill/timex"
	"sync"
	"time"
)

// 这里通过接口的形式实现滑动窗口功能，不同场景下桶内数据可以自定义。
// 但为了追求更好的性能，建议每种数据结构独立实现
type (
	// SlideWindow bucket
	SlideWinBucket interface {
		Reset()
		Add(v float64)
		AddByFlag(v float64, flag int8)
		WeedOut(past SlideWinBucket)
	}

	slideWindowBase struct {
		interval time.Duration
		baseTime time.Duration // 滑动窗口创建的时间作为基准时间
		lock     sync.RWMutex

		size       int // 桶的数量
		lastOffset int // 上次处理偏移量
	}

	// 通用滑动窗口，实现了接口SlideWinBucket的不同数据类型
	SlideWindow struct {
		*slideWindowBase
		win     SlideWinBucket   // 滑动窗口汇总后的桶数据
		buckets []SlideWinBucket // 所有桶数据
	}
)

func NewSlideWindow(win SlideWinBucket, buckets []SlideWinBucket, dur time.Duration) *SlideWindow {
	if len(buckets) < 1 {
		panic("size must be greater than 0")
	}

	return &SlideWindow{
		slideWindowBase: &slideWindowBase{
			size:     len(buckets),
			interval: dur,
			baseTime: timex.Now(),
		},
		win:     win,
		buckets: buckets,
	}
}

// 当前时间相对基准时间的偏移
func (rw *slideWindowBase) nowOffset() int {
	// TODO: 随着系统运行，如果不重启，dur会越来越大，是否需要重置一下时间呢？
	dur := timex.Since(rw.baseTime)

	// 如果系统时间被修改，比如导致时间倒退，需要重置基准时间
	if dur < 0 {
		dur = 0
		rw.baseTime = timex.Now()
	}
	return int(dur / rw.interval)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 返回当前滑动窗口汇总数据
func (rw *SlideWindow) CurrWin() SlideWinBucket {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	rw.expireBuckets(rw.nowOffset())
	return rw.win
}

// Add adds value to current bucket.
// 加锁临界区的执行应该是越少越好
func (rw *SlideWindow) Add(v float64) {
	currOffset := rw.nowOffset()
	bk := rw.buckets[currOffset%rw.size]

	rw.lock.Lock()
	defer rw.lock.Unlock()

	rw.expireBuckets(currOffset)

	bk.Add(v)
	rw.win.Add(v)

	rw.lastOffset = currOffset
}

// 高级添加函数，带标签
func (rw *SlideWindow) AddByFlag(v float64, flag int8) {
	currOffset := rw.nowOffset()
	bk := rw.buckets[currOffset%rw.size]

	rw.lock.Lock()
	defer rw.lock.Unlock()

	rw.expireBuckets(currOffset)

	bk.AddByFlag(v, flag)
	rw.win.AddByFlag(v, flag)

	rw.lastOffset = currOffset
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 剔除过期的桶
func (rw *SlideWindow) expireBuckets(currOffset int) {
	diff := currOffset - rw.lastOffset
	if diff > rw.size {
		diff = rw.size
	}

	// 当前和最后一次写入记录时间，相差的桶个数，依次将最后一次写入后面的桶全部清空
	for i := 0; i < diff; i++ {
		rw.expireBucket(rw.buckets[(rw.lastOffset+i+1)%rw.size])
	}
}

// 这个桶过期了，需要从全局变量中剔除
func (rw *SlideWindow) expireBucket(b SlideWinBucket) {
	rw.win.WeedOut(b)
	b.Reset()
}

// Demo ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//type LimiterBucket struct {
//	income      int64   // 新进请求数
//	finish      int64   // 处理完返回的请求数
//	totalTimeMS float64 // 处理总共耗时
//}
//
//func (bk *LimiterBucket) Reset() {
//	bk.income = 0
//	bk.finish = 0
//	bk.totalTimeMS = 0
//}
//
//func (bk *LimiterBucket) Add(v float64) {
//	bk.income++
//}
//
//func (bk *LimiterBucket) AddByFlag(v float64, flag int8) {
//	bk.finish++
//	bk.totalTimeMS += v
//}
//
//func (bk *LimiterBucket) Discount(past collect.SlideWinBucket) {
//	pbk := past.(*LimiterBucket)
//	bk.income -= pbk.income
//	bk.finish -= pbk.finish
//	bk.totalTimeMS -= pbk.totalTimeMS
//}
