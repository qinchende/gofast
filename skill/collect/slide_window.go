// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package collect

import (
	"github.com/qinchende/gofast/skill/timex"
	"sync"
	"time"
)

type (
	// SlideWindow bucket
	SlideWinBucket interface {
		Reset()
		Add(v float64)
		AddByFlag(v float64, flag int8)
		Discount(past SlideWinBucket)
	}

	// SlideWindow defines a rolling windowSdx to calculate the events in buckets with time interval.
	SlideWindow struct {
		interval      time.Duration
		baseTime      time.Duration // 滑动窗口创建的时间作为基准时间
		lock          sync.RWMutex
		size          int
		win           SlideWinBucket
		buckets       []SlideWinBucket
		lastReqOffset int
	}
)

// NewRollingWindow returns a SlideWindow that with size buckets and time interval,
// use opts to customize the SlideWindow.
func NewSlideWindow(win SlideWinBucket, buckets []SlideWinBucket, interval time.Duration) *SlideWindow {
	if len(buckets) < 1 {
		panic("size must be greater than 0")
	}

	return &SlideWindow{
		size:     len(buckets),
		win:      win,
		buckets:  buckets,
		interval: interval,
		baseTime: timex.Now(),
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 返回当前滑动窗口值
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

	rw.lastReqOffset = currOffset
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

	rw.lastReqOffset = currOffset
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 剔除过期的桶
func (rw *SlideWindow) expireBuckets(currOffset int) {
	diff := currOffset - rw.lastReqOffset
	if diff > rw.size {
		diff = rw.size
	}

	// 当前和最后一次写入记录时间，相差的桶个数，依次将最后一次写入后面的桶全部清空
	for i := 0; i < diff; i++ {
		rw.expireBucket(rw.buckets[(rw.lastReqOffset+i+1)%rw.size])
	}
}

// 这个桶过期了，需要从全局变量中剔除
func (rw *SlideWindow) expireBucket(b SlideWinBucket) {
	rw.win.Discount(b)
	b.Reset()
}

// 当前时间相对基准时间的偏移
func (rw *SlideWindow) nowOffset() int {
	// TODO: 随着系统运行，如果不重启，dur会越来越大，是否需要重置一下时间呢？
	dur := timex.Since(rw.baseTime)

	// 如果系统时间被修改，比如导致时间倒退，需要重置基准时间
	if dur < 0 {
		dur = 0
		rw.baseTime = timex.Now()
	}
	return int(dur / rw.interval)
}
