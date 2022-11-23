// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/collect"
	"time"
)

const (
	limit_window  = time.Second * 5 // 窗口周期
	limit_buckets = 10              // 窗口中桶的数量
)

type Limiter struct {
	sWin *collect.SlideWindowLimit
}

func NewLimiter() *Limiter {
	dur := time.Duration(int64(limit_window) / int64(limit_buckets))
	return &Limiter{
		sWin: collect.NewSlideWindowLimit(limit_buckets, dur),
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 新增请求数
func (rk *RequestKeeper) LimiterIncome(idx uint16) {
	rk.Limiters[idx].sWin.MarkIncome()
}

// 返回请求数以及耗时
func (rk *RequestKeeper) LimiterFinished(idx uint16, ms int32) {
	rk.Limiters[idx].sWin.MarkFinish(int64(ms))
}

// 是否允许本次请求通过
func (rk *RequestKeeper) LimiterAllow(idx uint16, defTimeoutMS int32) bool {
	lt := rk.Limiters[idx]

	income, finish, totalTimeMS := lt.sWin.CurrWin()

	// 平均处理时间比超时时间都大，全部降载丢弃
	//if finish >= 5 && income > finish+3 && float64(totalTimeMS)/float64(finish) > 1.2*float64(defTimeoutMS) {
	if finish >= 5 && income > 1 && float64(totalTimeMS)/float64(finish) > 1.2*float64(defTimeoutMS) {
		return true
	}
	return false
}
