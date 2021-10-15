package executors

import (
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/proc"
	"github.com/qinchende/gofast/skill/syncx"
	"github.com/qinchende/gofast/skill/timex"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

const idleRound = 10

type (
	// A type that satisfies executors.TaskContainer can be used as the underlying
	// container that used to do periodical executions.
	TaskContainer interface {
		// AddTask adds the task into the container.
		// Returns true if the container needs to be flushed after the addition.
		AddTask(task interface{}) bool
		// Execute handles the collected tasks by the container when flushing.
		Execute(tasks interface{})
		// RemoveAll removes the contained tasks, and return them.
		RemoveAll() interface{}
	}

	// 管理周期执行的实体对象
	IntervalExecutor struct {
		commander chan interface{}
		interval  time.Duration
		container TaskContainer
		waitGroup sync.WaitGroup
		// avoid race condition on waitGroup when calling wg.Add/Done/Wait(...)
		wgBarrier   syncx.Barrier
		confirmChan chan lang.PlaceholderType
		inflight    int32
		isRunning   bool
		newTicker   func(duration time.Duration) timex.Ticker
		lock        sync.Mutex
	}
)

// 初始化一个周期执行的实体对象
func NewPeriodicalExecutor(interval time.Duration, container TaskContainer) *IntervalExecutor {
	executor := &IntervalExecutor{
		// buffer 1 to let the caller go quickly
		commander:   make(chan interface{}, 1),
		interval:    interval,
		container:   container,
		confirmChan: make(chan struct{}),
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(d)
		},
	}
	proc.AddShutdownListener(func() {
		executor.Flush()
	})

	return executor
}

func (pe *IntervalExecutor) Add(task interface{}) {
	if values, ok := pe.addAndCheck(task); ok {
		pe.commander <- values
		<-pe.confirmChan
	}
}

// 执行一次定时任务。
func (pe *IntervalExecutor) Flush() bool {
	pe.enterExecution()
	return pe.executeTasks(func() interface{} {
		pe.lock.Lock()
		defer pe.lock.Unlock()
		return pe.container.RemoveAll()
	}())
}

func (pe *IntervalExecutor) Sync(fn func()) {
	pe.lock.Lock()
	defer pe.lock.Unlock()
	fn()
}

func (pe *IntervalExecutor) Wait() {
	pe.Flush()
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Wait()
	})
}

func (pe *IntervalExecutor) addAndCheck(task interface{}) (interface{}, bool) {
	pe.lock.Lock()
	defer func() {
		if !pe.isRunning {
			pe.isRunning = true
			// defer to unlock quickly
			defer pe.loopFlush()
		}
		pe.lock.Unlock()
	}()

	if pe.container.AddTask(task) {
		atomic.AddInt32(&pe.inflight, 1)
		return pe.container.RemoveAll(), true
	}

	return nil, false
}

// 后台启动新的协程运行周期任务
func (pe *IntervalExecutor) loopFlush() {
	gmp.GoSafe(func() {
		// flush before quit goroutine to avoid missing tasks
		defer pe.Flush()

		// 新建一个计时器，固定周期触发
		ticker := pe.newTicker(pe.interval)
		defer ticker.Stop()

		// 开启循环检测
		var commanded bool
		last := timex.Now()
		for {
			select {
			case values := <-pe.commander:
				commanded = true
				atomic.AddInt32(&pe.inflight, -1)
				pe.enterExecution()
				pe.confirmChan <- lang.Placeholder
				pe.executeTasks(values)
				last = timex.Now()
			case <-ticker.Chan():
				if commanded {
					commanded = false
				} else if pe.Flush() {
					last = timex.Now()
				} else if pe.quitLoop(last) {
					return
				}
			}
		}
	})
}

func (pe *IntervalExecutor) doneExecution() {
	pe.waitGroup.Done()
}

func (pe *IntervalExecutor) enterExecution() {
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Add(1)
	})
}

func (pe *IntervalExecutor) executeTasks(tasks interface{}) bool {
	defer pe.doneExecution()
	ok := pe.hasTasks(tasks)
	if ok {
		pe.container.Execute(tasks)
	}
	return ok
}

func (pe *IntervalExecutor) hasTasks(tasks interface{}) bool {
	if tasks == nil {
		return false
	}
	val := reflect.ValueOf(tasks)
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() > 0
	default:
		// unknown type, let caller execute it
		return true
	}
}

func (pe *IntervalExecutor) quitLoop(last time.Duration) (stop bool) {
	if timex.Since(last) <= pe.interval*idleRound {
		return
	}

	// checking pe.inflight and setting pe.guarded should be locked together
	pe.lock.Lock()
	if atomic.LoadInt32(&pe.inflight) == 0 {
		pe.isRunning = false
		stop = true
	}
	pe.lock.Unlock()

	return
}
