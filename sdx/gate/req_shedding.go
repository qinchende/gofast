// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/collect"
	"time"
)

type LimiterBucket struct {
	income      int64   // 新进请求数
	finish      int64   // 处理完返回的请求数
	totalTimeMS float64 // 处理总共耗时
}

func (bk *LimiterBucket) Reset() {
	bk.income = 0
	bk.finish = 0
	bk.totalTimeMS = 0
}

func (bk *LimiterBucket) Add(v float64) {
	bk.income++
}

func (bk *LimiterBucket) AddByFlag(v float64, flag int8) {
	bk.finish++
	bk.totalTimeMS += v
}

func (bk *LimiterBucket) Discount(past collect.SlideWinBucket) {
	pbk := past.(*LimiterBucket)
	bk.income -= pbk.income
	bk.finish -= pbk.finish
	bk.totalTimeMS -= pbk.totalTimeMS
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// TODO：注意这个熔断算法的关键参数
const (
	window  = time.Second * 10 // 10秒钟是一个完整的窗口周期
	buckets = 10               // 本周期分成40个桶, 那么每个桶占用250ms, 1秒钟分布4个桶。（这个粒度还是比较通用的）
)

type Limiter struct {
	sWin *collect.SlideWindow
}

func NewLimiter() *Limiter {

	win := new(LimiterBucket)
	bks := make([]collect.SlideWinBucket, buckets)
	for i := 0; i < buckets; i++ {
		bks[i] = new(LimiterBucket)
	}
	dur := time.Duration(int64(window) / int64(buckets))

	return &Limiter{
		sWin: collect.NewSlideWindow(win, bks, dur),
	}
}

func (rk *RequestKeeper) LimiterIncome(idx uint16) {
	lt := rk.Limiters[idx]
	lt.sWin.Add(1)
}

func (rk *RequestKeeper) LimiterFinished(idx uint16, ms int32) {
	lt := rk.Limiters[idx]
	lt.sWin.AddByFlag(float64(ms), 1)
}

func (rk *RequestKeeper) Shedding(idx uint16, timeout int32) bool {
	lt := rk.Limiters[idx]

	lbk := lt.sWin.CurrWin().(*LimiterBucket)

	// 平均处理时间比超时时间都大，全部降载丢弃
	if lbk.totalTimeMS/float64(lbk.finish+1) > float64(timeout) {
		return true
	}
	return false
}
