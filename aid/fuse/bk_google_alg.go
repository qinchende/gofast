// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fuse

import (
	"github.com/qinchende/gofast/aid/collect"
	"github.com/qinchende/gofast/aid/mathx"
	"math"
	"time"
)

// TODO：注意这个熔断算法的关键参数
const (
	bk_window  = time.Second * 10 // 10秒钟是一个完整的窗口周期
	bk_buckets = 40               // 本周期分成40个桶, 那么每个桶占用250ms, 1秒钟分布4个桶。（这个粒度还是比较通用的）
	k          = 1.5              // 熔断算法中的敏感系数
	protection = 5                // 最小请求个数，窗口期请求总数<=本请求数，即使都出错也不熔断
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// googleThrottle is a netflixBreaker pattern from google.
// see Client-Side Throttling section in https://landing.google.com/sre/sre-book/chapters/handling-overload/
type googleThrottle struct {
	sWin *collect.SlideWindowBreak
	k    float64
	may  *mathx.Maybe
}

func newGoogleThrottle() *googleThrottle {
	dur := time.Duration(int64(bk_window) / int64(bk_buckets))

	return &googleThrottle{
		sWin: collect.NewSlideWindowBreak(bk_buckets, dur),
		k:    k,
		may:  mathx.NewMaybe(),
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gtl *googleThrottle) allow() error {
	if err := gtl.accept(); err != nil {
		return err
	}
	return nil
}

// 降级逻辑
func (gtl *googleThrottle) doReq(req funcReq, fb funcFallback, cpt funcAcceptable) error {
	if err := gtl.accept(); err != nil {
		if fb != nil {
			return fb(err)
		}
		return err
	}

	defer func() {
		if e := recover(); e != nil {
			gtl.markValue(0)
			panic(e)
		}
	}()

	err := req()
	if cpt(err) {
		gtl.markValue(1)
	} else {
		gtl.markValue(0)
	}

	return err
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gtl *googleThrottle) markValue(val float64) {
	gtl.sWin.Add(val)
}

// 是否接收本次请求
// 谷歌公布的一段熔断算法：max(0, (requests - k*accepts) / (requests + 1))
func (gtl *googleThrottle) accept() error {
	accepts, total := gtl.sWin.CurrWin()

	// https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	// 比例(k-1)/k的请求出现错误，才会进入熔断的判断。如k=1.5时，失败达到33.3%以上可能熔断
	dropRatio := math.Max(0, (float64(total-protection)-gtl.k*accepts)/float64(total+1)) // 出错概率值[0,1)之间
	if dropRatio <= 0 {
		return nil
	}

	// 巧妙算法：取一个 0-1 之间的随机数 和 失败比率 做比较。失败比率越大越容易触发熔断。
	// 这种算法也决定了，在窗口熔断期内 还是随机存在一定比例的请求会被放过，起到了在熔断窗口期试探的作用。
	if gtl.may.TrueOnMaybe(dropRatio) {
		return ErrServiceUnavailable
	}
	return nil
}
