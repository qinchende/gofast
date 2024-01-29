// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package collect

import (
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

type (
	// 滑动窗口的数据结构，统计请求信息
	bucketBreak struct {
		total   int64   // 总访问次数
		accepts float64 // 处理成功的次数
	}

	// 通用滑动窗口，实现了接口SlideWinBucket的不同数据类型
	SlideWindowBreak struct {
		*slideWindowBase

		win     bucketBreak   // 滑动窗口汇总后的桶数据
		buckets []bucketBreak // 所有桶数据
	}
)

func NewSlideWindowBreak(size int, dur time.Duration) *SlideWindowBreak {
	if size < 1 {
		panic("size must be greater than 0")
	}

	return &SlideWindowBreak{
		slideWindowBase: &slideWindowBase{
			size:     size,
			interval: dur,
			baseTime: timex.NowDur(),
		},
		buckets: make([]bucketBreak, size),
	}
}

// 返回当前滑动窗口汇总数据
func (rw *SlideWindowBreak) CurrWin() (float64, int64) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	rw.expireBuckets(rw.nowOffset())
	return rw.win.accepts, rw.win.total
}

// Add adds value to current bucket.
// 加锁临界区的执行应该是越少越好
func (rw *SlideWindowBreak) Add(v float64) {
	currOffset := rw.nowOffset()
	bk := &rw.buckets[currOffset%rw.size]

	rw.lock.Lock()
	defer rw.lock.Unlock()

	rw.expireBuckets(currOffset)

	bk.total++
	bk.accepts += v
	rw.win.total++
	rw.win.accepts += v

	rw.lastOffset = currOffset
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 剔除过期的桶
func (rw *SlideWindowBreak) expireBuckets(currOffset int) {
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
func (rw *SlideWindowBreak) expireBucket(b *bucketBreak) {
	// discount expire bucket
	rw.win.total -= b.total
	rw.win.accepts -= b.accepts
	// reset current bucket
	b.total = 0
	b.accepts = 0
}
