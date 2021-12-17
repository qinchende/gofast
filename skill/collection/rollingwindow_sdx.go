package collection

import (
	"github.com/qinchende/gofast/skill/timex"
	"sync"
	"time"
)

type (
	// Bucket defines the bucket that holds sum and num of additions.
	BucketSdx struct {
		sum   float64 // 求和（因为是通用模块，其实可以用于一些计算值的场景）
		count int64   // 统计访问次数
	}

	// RollingWindowSdx defines a rolling windowSdx to calculate the events in buckets with time interval.
	RollingWindowSdx struct {
		interval      time.Duration
		baseTime      time.Duration // 滑动窗口创建的时间作为基准时间
		lock          sync.RWMutex
		size          int
		win           BucketSdx
		buckets       []BucketSdx
		lastReqOffset int
	}
)

// NewRollingWindow returns a RollingWindowSdx that with size buckets and time interval,
// use opts to customize the RollingWindowSdx.
func NewRollingWindowSdx(size int, interval time.Duration) *RollingWindowSdx {
	if size < 1 {
		panic("size must be greater than 0")
	}

	return &RollingWindowSdx{
		size:     size,
		buckets:  make([]BucketSdx, size, size),
		interval: interval,
		baseTime: timex.Now(),
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 返回当前滑动窗口值
func (rw *RollingWindowSdx) CurrWinValue() (float64, int64) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	rw.expireBuckets(rw.nowOffset())
	return rw.win.sum, rw.win.count
}

// Add adds value to current bucket.
// 加锁临界区的执行应该是越少越好
func (rw *RollingWindowSdx) Add(v float64) {
	currOffset := rw.nowOffset()
	bk := &rw.buckets[currOffset%rw.size]

	rw.lock.Lock()
	defer rw.lock.Unlock()

	rw.expireBuckets(currOffset)

	bk.sum += v // 每次的值可以任意指定。意味这个可以灵活应用于更多场景
	bk.count++  // 调用次数是递增的。可以用这个来记录调用次数

	rw.win.sum += v
	rw.win.count++

	rw.lastReqOffset = currOffset
}

// 剔除过期的桶
func (rw *RollingWindowSdx) expireBuckets(currOffset int) {
	diff := currOffset - rw.lastReqOffset
	if diff > rw.size {
		diff = rw.size
	}

	// 当前和最后一次写入记录时间，相差的桶个数，依次将最后一次写入后面的桶全部清空
	for i := 0; i < diff; i++ {
		rw.expireBucket(&rw.buckets[(rw.lastReqOffset+i+1)%rw.size])
	}
}

// 这个桶过期了，需要从全局变量中剔除
func (rw *RollingWindowSdx) expireBucket(b *BucketSdx) {
	rw.win.sum -= b.sum
	rw.win.count -= b.count

	// reset current bucket
	b.sum = 0
	b.count = 0
}

// 当前时间相对基准时间的偏移
func (rw *RollingWindowSdx) nowOffset() int {
	// TODO: 随着系统运行，如果不重启，dur会越来越大，是否需要重置一下时间呢？
	dur := timex.Since(rw.baseTime)

	// 如果系统时间被修改，比如导致时间倒退，需要重置基准时间
	if dur < 0 {
		dur = 0
		rw.baseTime = timex.Now()
	}
	return int(dur / rw.interval)
}
