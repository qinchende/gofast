// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package exec

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/qinchende/gofast/skill/syncx"
	"github.com/qinchende/gofast/skill/timex"
)

// 降频执行器（主要是减少对某一类任务的执行频率）
type Reduce struct {
	skipTimes int32 // 跳过多少次执行

	interval time.Duration
	lastTime *syncx.AtomicDuration
	lock     sync.Mutex
}

func NewReduce(step time.Duration) *Reduce {
	return &Reduce{
		interval: step,
		lastTime: syncx.NewAtomicDuration(),
	}
}

// 至少间隔一定时间才执行指定Task
func (rd *Reduce) DoInterval(flush bool, task func(skipTimes int32)) bool {
	rd.lock.Lock()
	defer rd.lock.Unlock()

	now := timex.NowDur()
	lastTime := rd.lastTime.Load()
	// 首次需要执行
	if flush || lastTime+rd.interval < now || lastTime == 0 {
		rd.lastTime.Set(now)
		task(atomic.SwapInt32(&rd.skipTimes, 0))
		return true
	}
	atomic.AddInt32(&rd.skipTimes, 1)
	return false
}
