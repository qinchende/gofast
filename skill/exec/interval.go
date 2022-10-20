package exec

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
	// A type that satisfies executors.ItemContainer can be used as the underlying
	// container that used to do periodical executions.
	ItemContainer interface {
		// AddItem adds the task into the container.
		// Returns true if the container needs to be flushed after the addition.
		AddItem(item any) bool
		// Execute handles the collected items by the container when flushing.
		Execute(items any)
		// RemoveAll removes the contained items, and return them.
		RemoveAll() any
	}

	// 管理周期执行的实体对象
	Interval struct {
		commander   chan any
		confirmChan chan lang.PlaceholderType
		interval    time.Duration
		container   ItemContainer
		waitGroup   sync.WaitGroup
		wgBarrier   syncx.Barrier // avoid race condition on waitGroup when calling wg.Add/Done/Wait(...)
		inflight    int32
		isRunning   bool
		newTicker   func(duration time.Duration) timex.Ticker
		lock        sync.Mutex
	}
)

// 初始化一个周期执行的实体对象
func NewInterval(interval time.Duration, container ItemContainer) *Interval {
	executor := &Interval{
		// buffer 1 to let the caller go quickly
		commander:   make(chan any, 1),   // 长度为1的有缓冲通道
		confirmChan: make(chan struct{}), // 无缓冲通道
		interval:    interval,
		container:   container,
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(d)
		},
	}
	proc.AddShutdownListener(func() {
		executor.Flush()
	})

	return executor
}

func (pe *Interval) Add(item any) {
	// 外界可以强制刷新日志，并将values传入可执行函数
	if values, ok := pe.addAndCheck(item); ok {
		pe.commander <- values
		<-pe.confirmChan
	}
}

// 执行一次定时任务。
func (pe *Interval) Flush() bool {
	pe.enterExecution()
	return pe.executeTasks(func() any {
		pe.lock.Lock()
		defer pe.lock.Unlock()
		return pe.container.RemoveAll()
	}())
}

func (pe *Interval) Sync(fn func()) {
	pe.lock.Lock()
	defer pe.lock.Unlock()
	fn()
}

func (pe *Interval) Wait() {
	pe.Flush()
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Wait()
	})
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (pe *Interval) addAndCheck(item any) (any, bool) {
	pe.lock.Lock()
	defer func() {
		if !pe.isRunning {
			pe.isRunning = true
			// defer to unlock quickly
			defer pe.loopFlush()
		}
		pe.lock.Unlock()
	}()

	// 外部容器可以试图返回 true，这样能立刻执行统计请求。
	if pe.container.AddItem(item) {
		atomic.AddInt32(&pe.inflight, 1)
		return pe.container.RemoveAll(), true
	}

	return nil, false
}

// 后台启动新的协程运行周期任务
func (pe *Interval) loopFlush() {
	gmp.GoSafe(func() {
		// flush before quit goroutine to avoid missing items
		defer pe.Flush()

		// 新建一个计时器，固定周期触发
		ticker := pe.newTicker(pe.interval)
		defer ticker.Stop()

		// 外部命令立即输出
		var commanded bool
		last := timex.Now()
		// 开启死循环循环检测。手动指令，或者定时任务 都可以输出统计结果
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
				// 如果上面手动输出一次，那么本次自动输出将轮空
				if commanded {
					commanded = false
				} else if pe.Flush() {
					last = timex.Now()
				} else if pe.quitLoop(last) {
					return
				}
			}
		}
		// +++++++++++++++++++++++++++++++++++++++++++++++++++++
	})
}

// 一定次数循环发现没有新任务，自动退出定时循环
func (pe *Interval) quitLoop(last time.Duration) (stop bool) {
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (pe *Interval) doneExecution() {
	pe.waitGroup.Done()
}

func (pe *Interval) enterExecution() {
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Add(1)
	})
}

func (pe *Interval) executeTasks(items any) bool {
	defer pe.doneExecution()
	ok := pe.hasTasks(items)
	if ok {
		pe.container.Execute(items)
	}
	return ok
}

func (pe *Interval) hasTasks(items any) bool {
	if items == nil {
		return false
	}
	val := reflect.ValueOf(items)
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() > 0
	default:
		// unknown type, let caller execute it
		return true
	}
}
