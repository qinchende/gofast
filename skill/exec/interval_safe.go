package exec

import (
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/proc"
	"github.com/qinchende/gofast/skill/timex"
	"sync"
	"sync/atomic"
	"time"
)

type (
	IntervalSafe struct {
		Interval
		taskLock sync.Mutex // 主动加锁
	}
)

func NewIntervalSafe(dur time.Duration, box TaskContainer) *IntervalSafe {
	run := &IntervalSafe{
		Interval: Interval{
			messenger:   make(chan any, 1),
			confirmChan: make(chan lang.PlaceholderType),
			interval:    dur,
			container:   box,
			newTicker: func(d time.Duration) timex.Ticker {
				return timex.NewTicker(d)
			},
		},
	}

	// 程序退出时要执行一次
	proc.AddShutdownListener(func() {
		run.Flush()
	})
	return run
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 执行一次定时任务。
// 返回是否真的有任务被执行过了
// 请调用者自己保证 RemoveAll 函数并发安全
func (run *IntervalSafe) Flush() (hasTasks bool) {
	run.enterExecution()
	run.taskLock.Lock()
	items := run.container.RemoveAll()
	run.taskLock.Unlock()
	return run.executeTasks(items)
}

// 请调用者自己保证AddItem函数并发安全
func (run *IntervalSafe) Add(item any) {
	run.checkLoop()
	run.taskLock.Lock()
	ok := run.container.AddItem(item)
	run.taskLock.Unlock()
	run.checkLoop()

	// ok == true -> 立即主动执行任务
	if ok {
		atomic.AddInt32(&run.inflight, 1)
		run.messenger <- run.container.RemoveAll()
		<-run.confirmChan
	}
}

// 请调用者自己保证fc函数并发安全
func (run *IntervalSafe) AddByFunc(fc AddFunc, item any) {
	run.checkLoop()
	run.taskLock.Lock()
	values, ok := fc(item)
	run.taskLock.Unlock()
	run.checkLoop()

	// ok == true -> 立即主动执行任务
	if ok {
		atomic.AddInt32(&run.inflight, 1)
		run.messenger <- values
		<-run.confirmChan
	}
}
