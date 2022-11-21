// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/collect"
	"time"
)

const (
	shed_window  = time.Second * 10 // 窗口周期
	shed_buckets = 10               // 10个桶
)

type Limiter struct {
	sWin *collect.SlideWindowLimit
}

func NewLimiter() *Limiter {
	dur := time.Duration(int64(shed_window) / int64(shed_buckets))
	return &Limiter{
		sWin: collect.NewSlideWindowShed(shed_buckets, dur),
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (rk *RequestKeeper) LimiterIncome(idx uint16) {
	lt := rk.Limiters[idx]
	lt.sWin.MarkIncome()
}

func (rk *RequestKeeper) LimiterFinished(idx uint16, ms int32) {
	lt := rk.Limiters[idx]
	lt.sWin.MarkFinish(int64(ms))
}

func (rk *RequestKeeper) AccessLimit(idx uint16, timeout int32) bool {
	lt := rk.Limiters[idx]

	_, finish, totalTimeMS := lt.sWin.CurrWin()

	// 平均处理时间比超时时间都大，全部降载丢弃
	if float64(totalTimeMS)/float64(finish+1) > float64(timeout) {
		return true
	}
	return false
}
