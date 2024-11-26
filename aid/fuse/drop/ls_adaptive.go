package drop

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/aid/collect"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/aid/syncx"
	"github.com/qinchende/gofast/aid/sysx"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/core/cst"
	"math"
	"sync/atomic"
	"time"
)

const (
	defaultBuckets      = 50
	defaultWindow       = time.Second * 5
	defaultCpuThreshold = 900 // using 1000m notation, 900m is like 80%, keep it as var for unit test
	defaultMinRt        = float64(time.Second / time.Millisecond)
	flyingBeta          = 0.9         // moving average hyperparameter beta for calculating requests on the fly
	coolOffDuration     = time.Second // 冷却期默认为1秒
)

var (
	// ErrServiceOverloaded is returned by Shedder.Allow when the service is overloaded.
	ErrServiceOverloaded = errors.New("service overloaded")

	// default to be enabled
	enabled = syncx.ForAtomicBool(true)
	// make it a variable for unit test
	systemOverloadChecker = func(cpuThreshold float64) bool {
		return sysx.CpuSmoothUsage >= cpuThreshold
	}
)

// Disable lets callers disable load shedding.
func Disable() {
	enabled.Set(false)
}

// NewAdaptiveShedder returns an adaptive shedder.
// opts can be used to customize the Shedder.
func NewAdaptiveShedder(opts ...ShedderOption) Shedder {
	//if !enabled.True() {
	//	return newNopShedder()
	//}

	options := shedderOptions{
		window:       defaultWindow,
		buckets:      defaultBuckets,
		cpuThreshold: defaultCpuThreshold,
	}
	for _, opt := range opts {
		opt(&options)
	}
	bucketDuration := options.window / time.Duration(options.buckets)
	return &adaptiveShedder{
		cpuThreshold:    options.cpuThreshold,
		windows:         int64(time.Second / bucketDuration),
		dropTime:        syncx.NewAtomicDuration(),
		droppedRecently: syncx.NewAtomicBool(),
		passCounter: collect.NewRollingWindow(options.buckets, bucketDuration,
			collect.IgnoreCurrentBucket()),
		rtCounter: collect.NewRollingWindow(options.buckets, bucketDuration,
			collect.IgnoreCurrentBucket()),
	}
}

// Allow implements Shedder.Allow.
func (as *adaptiveShedder) Allow() (Promise, error) {
	if as.shouldDrop() {
		as.dropTime.Set(timex.NowDur())
		as.droppedRecently.Set(true)

		return nil, ErrServiceOverloaded
	}

	as.addFlying(1)

	return &promise{
		start:   timex.NowDur(),
		shedder: as,
	}, nil
}

func (as *adaptiveShedder) addFlying(delta int64) {
	flying := atomic.AddInt64(&as.flying, delta)
	// update avgFlying when the request is finished.
	// this strategy makes avgFlying have a little bit lag against flying, and smoother.
	// when the flying requests increase rapidly, avgFlying increase slower, accept more requests.
	// when the flying requests drop rapidly, avgFlying drop slower, accept less requests.
	// it makes the service to serve as more requests as possible.
	if delta < 0 {
		as.avgFlyingLock.Lock()
		as.avgFlying = as.avgFlying*flyingBeta + float64(flying)*(1-flyingBeta)
		as.avgFlyingLock.Unlock()
	}
}

func (as *adaptiveShedder) highThru() bool {
	as.avgFlyingLock.Lock()
	avgFlying := as.avgFlying
	as.avgFlyingLock.Unlock()
	maxFlight := as.maxFlight()
	return int64(avgFlying) > maxFlight && atomic.LoadInt64(&as.flying) > maxFlight
}

func (as *adaptiveShedder) maxFlight() int64 {
	// windows = buckets per second
	// maxQPS = maxPASS * windows
	// minRT = min average response time in milliseconds
	// maxQPS * minRT / milliseconds_per_second
	return int64(math.Max(1, float64(as.maxPass()*as.windows)*(as.minRt()/1e3)))
}

func (as *adaptiveShedder) maxPass() int64 {
	var result float64 = 1

	as.passCounter.Reduce(func(b *collect.Bucket) {
		if b.Sum > result {
			result = b.Sum
		}
	})

	return int64(result)
}

func (as *adaptiveShedder) minRt() float64 {
	result := defaultMinRt

	as.rtCounter.Reduce(func(b *collect.Bucket) {
		if b.Count <= 0 {
			return
		}

		avg := math.Round(b.Sum / float64(b.Count))
		if avg < result {
			result = avg
		}
	})

	return result
}

func (as *adaptiveShedder) shouldDrop() bool {
	if as.systemOverloaded() || as.stillHot() {
		if as.highThru() {
			flying := atomic.LoadInt64(&as.flying)
			as.avgFlyingLock.Lock()
			avgFlying := as.avgFlying
			as.avgFlyingLock.Unlock()
			msg := fmt.Sprintf(
				"dropreq, cpu: %d, maxPass: %d, minRt: %.2f, hot: %t, flying: %d, avgFlying: %.2f",
				sysx.CpuSmoothUsage, as.maxPass(), as.minRt(), as.stillHot(), flying, avgFlying)
			logx.Err().SendMsg(msg)
			logx.InfoReport(cst.KV{msg: msg})
			return true
		}
	}

	return false
}

func (as *adaptiveShedder) stillHot() bool {
	if !as.droppedRecently.True() {
		return false
	}

	dropTime := as.dropTime.Load()
	if dropTime == 0 {
		return false
	}

	hot := timex.NowDiffDur(dropTime) < coolOffDuration
	if !hot {
		as.droppedRecently.Set(false)
	}

	return hot
}

func (as *adaptiveShedder) systemOverloaded() bool {
	return systemOverloadChecker(as.cpuThreshold)
}

// WithBuckets customizes the Shedder with given number of buckets.
func WithBuckets(buckets int) ShedderOption {
	return func(opts *shedderOptions) {
		opts.buckets = buckets
	}
}

// WithCpuThreshold customizes the Shedder with given cpu threshold.
func WithCpuThreshold(threshold float64) ShedderOption {
	return func(opts *shedderOptions) {
		opts.cpuThreshold = threshold
	}
}

// WithWindow customizes the Shedder with given
func WithWindow(window time.Duration) ShedderOption {
	return func(opts *shedderOptions) {
		opts.window = window
	}
}

type promise struct {
	start   time.Duration
	shedder *adaptiveShedder
}

func (p *promise) Fail() {
	p.shedder.addFlying(-1)
}

func (p *promise) Pass() {
	rt := float64(timex.NowDiffDur(p.start)) / float64(time.Millisecond)
	p.shedder.addFlying(-1)
	p.shedder.rtCounter.Add(math.Ceil(rt))
	p.shedder.passCounter.Add(1)
}
