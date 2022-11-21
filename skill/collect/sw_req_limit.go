// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package collect

import (
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

type (
	// 滑动窗口的数据结构，统计请求信息
	bucketLimit struct {
		income      int64 // 新进请求数
		finish      int64 // 处理完返回的请求数
		totalTimeMS int64 // 处理总共耗时
	}

	// 通用滑动窗口，实现了接口SlideWinBucket的不同数据类型
	SlideWindowLimit struct {
		*slideWindowBase

		win     bucketLimit   // 滑动窗口汇总后的桶数据
		buckets []bucketLimit // 所有桶数据
	}
)

func NewSlideWindowShed(size int, dur time.Duration) *SlideWindowLimit {
	if size < 1 {
		panic("size must be greater than 0")
	}

	return &SlideWindowLimit{
		slideWindowBase: &slideWindowBase{
			size:     size,
			interval: dur,
			baseTime: timex.Now(),
		},
		buckets: make([]bucketLimit, size),
	}
}

// 返回当前滑动窗口汇总数据
func (rw *SlideWindowLimit) CurrWin() (int64, int64, int64) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	rw.expireBuckets(rw.nowOffset())
	return rw.win.income, rw.win.finish, rw.win.totalTimeMS
}

func (rw *SlideWindowLimit) MarkIncome() {
	rw.addValue(0, 0)
}

func (rw *SlideWindowLimit) MarkFinish(ms int64) {
	rw.addValue(1, ms)
}

func (rw *SlideWindowLimit) addValue(flag int8, ms int64) {
	currOffset := rw.nowOffset()
	bk := &rw.buckets[currOffset%rw.size]

	rw.lock.Lock()
	defer rw.lock.Unlock()

	rw.expireBuckets(currOffset)

	if flag == 0 {
		bk.income++
		rw.win.income++
	} else {
		bk.finish++
		bk.totalTimeMS += ms
		rw.win.finish++
		rw.win.totalTimeMS += ms
	}

	rw.lastOffset = currOffset
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 剔除过期的桶
func (rw *SlideWindowLimit) expireBuckets(currOffset int) {
	diff := currOffset - rw.lastOffset
	if diff > rw.size {
		diff = rw.size
	}

	// 当前和最后一次写入记录时间，相差的桶个数，依次将最后一次写入后面的桶全部清空
	for i := 0; i < diff; i++ {
		rw.expireBucket(&rw.buckets[(rw.lastOffset+i+1)%rw.size])
	}
}

// 这个桶过期了，需要从全局变量中剔除
func (rw *SlideWindowLimit) expireBucket(b *bucketLimit) {
	rw.win.income -= b.income
	rw.win.finish -= b.finish
	rw.win.totalTimeMS -= b.totalTimeMS

	b.income = 0
	b.finish = 0
	b.totalTimeMS = 0
}
