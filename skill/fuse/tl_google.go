package fuse

import (
	"fmt"
	"github.com/qinchende/gofast/skill/proc"

	"github.com/qinchende/gofast/skill/collect"
	"github.com/qinchende/gofast/skill/mathx"
	"math"
	"time"
)

const (
	window     = time.Second * 10 // 10秒钟是一个完整的窗口周期
	buckets    = 40               // 本周期分成40个桶, 那么每个桶占用250ms, 1秒钟分布4个桶。（这个粒度还是比较通用的）
	k          = 1.5              // 熔断算法中的敏感系数
	protection = 5                // 最小请求个数，窗口期请求总数<=本请求数，即使都出错也不熔断
)

// 用于回调成功或者失败
type googlePromise struct {
	gtl *googleThrottle
}

func (p googlePromise) Accept() {
	p.gtl.markSuc()
}

func (p googlePromise) Reject(reason string) {
	p.gtl.errWin.add(reason)
	p.gtl.markFai()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// googleThrottle is a netflixBreaker pattern from google.
// see Client-Side Throttling section in https://landing.google.com/sre/sre-book/chapters/handling-overload/
type googleThrottle struct {
	name   string
	k      float64
	rWin   *collect.RollingWindowSdx
	may    *mathx.Maybe
	errWin *errorWindow
}

func newGoogleThrottle(name string) *googleThrottle {
	dur := time.Duration(int64(window) / int64(buckets))
	rWin := collect.NewRollingWindowSdx(buckets, dur)

	return &googleThrottle{
		name: name,
		k:    k,
		rWin: rWin,
		may:  mathx.NewMaybe(),
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gtl *googleThrottle) allow() (Promise, error) {
	err := gtl.accept()
	if err != nil {
		return nil, gtl.logError(err)
	}
	return googlePromise{gtl: gtl}, nil
}

func (gtl *googleThrottle) doReq(req funcReq, fallback funcFallback, acceptable Acceptable) error {
	if err := gtl.accept(); err != nil {
		if fallback != nil {
			return gtl.logError(fallback(err))
		}
		return gtl.logError(err)
	}

	defer func() {
		if e := recover(); e != nil {
			gtl.markFai()
			panic(e)
		}
	}()

	err := req()
	if acceptable(err) {
		gtl.markSuc()
	} else {
		gtl.markFai()
	}

	return gtl.logError(err)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gtl *googleThrottle) markSuc() {
	gtl.rWin.Add(1)
}

func (gtl *googleThrottle) markFai() {
	gtl.rWin.Add(0)
}

// 是否接收本次请求
// 谷歌公布的一段熔断算法：max(0, (requests - k*accepts) / (requests + 1))
func (gtl *googleThrottle) accept() error {
	accepts, total := gtl.rWin.CurrWinValue()
	// https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	dropRatio := math.Max(0, (float64(total-protection)-gtl.k*accepts)/float64(total+1))
	// logx.Error("total: ", total, " accepts: ", accepts, " dropRatio: ", dropRatio)
	if dropRatio <= 0 {
		return nil
	}

	// 取一个 0-1 之间的随机数 和 失败比率 做比较。失败比率越大越容易触发熔断。
	// 这种算法也决定了，在窗口熔断期内 还是随机存在一定比例的请求会被放过，起到了在熔断窗口期试探的作用。
	if gtl.may.TrueOnMaybe(dropRatio) {
		return ErrServiceUnavailable
	}
	return nil
}

func (gtl *googleThrottle) logError(err error) error {
	if err == ErrServiceUnavailable {
		Report(fmt.Sprintf("proc(%s/%d), callee: %s, breaker is open and requests dropped\nlast errors:\n%s",
			proc.ProcessName(), proc.Pid(), gtl.name, gtl.errWin))
	}
	return err
}

//func (gtl *googleThrottle) history() (accepts, total int64) {
//	gtl.rWin.Reduce(func(b *collection.Bucket) {
//		accepts += int64(b.Sum)
//		total += b.Count
//	})
//
//	return
//}

//// 获取当前滑动窗口的数据
//func (gtl *googleThrottle) historySdx() (float64, int64) {
//	return gtl.rWin.TotalValue()
//}
