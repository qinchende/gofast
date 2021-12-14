package breaker

import (
	"math"
	"time"

	"github.com/qinchende/gofast/skill/collection"
	"github.com/qinchende/gofast/skill/mathx"
)

const (
	// 250ms for bucket duration
	window     = time.Second * 3600 // 10秒钟是一个完整的窗口周期
	buckets    = 10                 // 本周期分成40个桶, 那么每个桶占用250ms, 1秒钟分布4个桶。（这个粒度还是比较通用的）
	k          = 1.5                // 熔断算法中的一个系数
	protection = 5                  // 最小请求个数，窗口期少于本请求数，不检查熔断
)

// googleBreaker is a netflixBreaker pattern from google.
// see Client-Side Throttling section in https://landing.google.com/sre/sre-book/chapters/handling-overload/
type googleBreaker struct {
	k     float64
	stat  *collection.RollingWindow
	proba *mathx.Proba
}

func newGoogleBreaker() *googleBreaker {
	bucketDuration := time.Duration(int64(window) / int64(buckets))
	st := collection.NewRollingWindow(buckets, bucketDuration)
	return &googleBreaker{
		stat:  st,
		k:     k,
		proba: mathx.NewProba(),
	}
}

// 谷歌公布的一段熔断算法：max(0, (requests - k*accepts) / (requests + 1))
func (gBrk *googleBreaker) accept() error {
	accepts, total := gBrk.history()
	weightedAccepts := gBrk.k * float64(accepts)
	// https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	dropRatio := math.Max(0, (float64(total-protection)-weightedAccepts)/float64(total+1))
	if dropRatio <= 0 {
		return nil
	}

	if gBrk.proba.TrueOnProba(dropRatio) {
		return ErrServiceUnavailable
	}
	return nil
}

func (gBrk *googleBreaker) allow() (internalPromise, error) {
	if err := gBrk.accept(); err != nil {
		return nil, err
	}
	return googlePromise{
		brk: gBrk,
	}, nil
}

func (gBrk *googleBreaker) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	if err := gBrk.accept(); err != nil {
		if fallback != nil {
			return fallback(err)
		}
		return err
	}

	defer func() {
		if e := recover(); e != nil {
			gBrk.markFailure()
			panic(e)
		}
	}()

	err := req()
	if acceptable(err) {
		gBrk.markSuccess()
	} else {
		gBrk.markFailure()
	}

	return err
}

func (gBrk *googleBreaker) markSuccess() {
	gBrk.stat.Add(1)
}

func (gBrk *googleBreaker) markFailure() {
	gBrk.stat.Add(0)
}

func (gBrk *googleBreaker) history() (accepts, total int64) {
	gBrk.stat.Reduce(func(b *collection.Bucket) {
		accepts += int64(b.Sum)
		total += b.Count
	})

	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 用于回调成功或者失败
type googlePromise struct {
	brk *googleBreaker
}

func (gPromise googlePromise) Accept() {
	gPromise.brk.markSuccess()
}

func (gPromise googlePromise) Reject() {
	gPromise.brk.markFailure()
}
