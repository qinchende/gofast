// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/collect"
	"time"
)

const (
	limit_win_seconds  = 3
	limit_win_duration = limit_win_seconds * time.Second // 窗口周期(秒)
	limit_buckets      = limit_win_seconds * 4           // 窗口中桶的数量(250ms一个桶)
	limit_timeout_rate = 1.1                             // 超时倍数
)

type Limiter struct {
	sWin *collect.SlideWindowLimit
}

func NewLimiter() *Limiter {
	dur := time.Duration(int64(limit_win_duration) / int64(limit_buckets))
	return &Limiter{
		sWin: collect.NewSlideWindowLimit(limit_buckets, dur),
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 新增请求数
func (rk *RequestKeeper) LimiterIncome(idx uint16) {
	rk.Limiters[idx].sWin.MarkIncome()
}

// 记录请求耗时
func (rk *RequestKeeper) LimiterFinished(idx uint16, ms int32, defMS int32) {
	fixMS := int32(float64(defMS) * limit_timeout_rate)
	if ms > fixMS {
		ms = fixMS
	}
	rk.Limiters[idx].sWin.MarkFinish(int64(ms))
}

// 是否允许本次请求通过
func (rk *RequestKeeper) LimiterAllow(idx uint16, defMS int32) bool {
	lt := rk.Limiters[idx]

	_, finish, totalTimeMS := lt.sWin.CurrWin()

	// 过去3秒，处理完成至少3个请求，而且所有请求几乎都是超时，这个时候要降载了
	if finish >= limit_win_seconds && float64(totalTimeMS)/float64(finish) > float64(defMS) {
		return true
	}
	return false
}
