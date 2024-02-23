// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package exec

import (
	"github.com/qinchende/gofast/aid/proc"
	"sync/atomic"
	"time"
)

type (
	IntervalUnsafe struct {
		*Interval
	}
)

func NewIntervalUnsafe(dur time.Duration, box TaskContainer) *IntervalUnsafe {
	run := &IntervalUnsafe{Interval: createInterval(dur, box)}
	// 程序退出时要执行一次
	proc.AddShutdownListener(func() {
		run.Flush()
	})
	return run
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 执行一次定时任务。
// 返回是否真的有任务被执行过了，请调用者自己保证 RemoveAll 函数并发安全
func (run *IntervalUnsafe) Flush() (hasTasks bool) {
	run.enterExecution()
	items := run.container.RemoveAll()
	return run.executeTasks(items)
}

// 请调用者自己保证AddItem函数并发安全
func (run *IntervalUnsafe) Add(item any) {
	run.checkLoop()
	ok := run.container.AddItem(item)
	run.checkLoop()

	// ok == true -> 立即主动执行任务
	if ok {
		atomic.AddInt32(&run.inflight, 1)
		run.messenger <- run.container.RemoveAll()
		<-run.confirmChan
	}
}

// 请调用者自己保证fc函数并发安全
func (run *IntervalUnsafe) AddByFunc(fc AddFunc, item any) {
	run.checkLoop()
	items, ok := fc(item)
	run.checkLoop()

	// ok == true -> 立即主动执行任务
	if ok {
		atomic.AddInt32(&run.inflight, 1)
		run.messenger <- items
		<-run.confirmChan
	}
}
