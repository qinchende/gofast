package exec

import (
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/proc"
	"github.com/qinchende/gofast/skill/timex"
	"reflect"
	"sync"
	"time"
)

const idleRound = 10

type (
	// 外部可以自定义装载各种任务的容器，实现这些方法之后，就能赋予周期执行的特性
	TaskContainer interface {
		AddItem(item any) bool
		Execute(items any)
		RemoveAll() any
	}

	// 管理周期执行的实体对象
	Interval struct {
		messenger   chan any
		confirmChan chan lang.PlaceholderType
		interval    time.Duration
		container   TaskContainer
		newTicker   func(d time.Duration) timex.Ticker

		waitGroup sync.WaitGroup
		wgLock    sync.Mutex
		runLock   sync.Mutex
		isRunning bool
	}
)

// 初始化一个周期执行的实体对象
func NewInterval(dur time.Duration, box TaskContainer) *Interval {
	run := &Interval{
		messenger:   make(chan any, 1),               // 长度为1的有缓冲通道(buffer 1 to let the caller go quickly)
		confirmChan: make(chan lang.PlaceholderType), // 无缓冲通道
		interval:    dur,
		container:   box,
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(d)
		},
	}

	// 程序退出时要执行一次
	proc.AddShutdownListener(func() {
		run.Flush()
	})
	return run
}

func (run *Interval) Add(item any) {
	if values, ok := run.addAndCheck(item); ok {
		// ok == true -> 主动执行任务
		run.messenger <- values
		<-run.confirmChan
	}
}

// 执行一次定时任务。
func (run *Interval) Flush() bool {
	run.enterExecution()
	return run.executeTasks(func() any {
		run.runLock.Lock()
		defer run.runLock.Unlock()
		return run.container.RemoveAll()
	}())
}

func (run *Interval) Sync(fn func()) {
	run.runLock.Lock()
	defer run.runLock.Unlock()
	fn()
}

func (run *Interval) Wait() {
	run.Flush()

	run.wgLock.Lock()
	defer run.wgLock.Unlock()
	run.waitGroup.Wait()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (run *Interval) addAndCheck(item any) (any, bool) {
	run.runLock.Lock()
	defer func() {
		if run.isRunning == false {
			run.isRunning = true
			defer run.raiseLoopFlush()
		}
		run.runLock.Unlock()
	}()
	// 外部容器可以试图返回 true，这样能立刻执行统计请求。
	if run.container.AddItem(item) {
		return run.container.RemoveAll(), true
	}
	return nil, false
}

// 后台启动新的协程运行周期任务
func (run *Interval) raiseLoopFlush() {
	gmp.GoSafe(func() {
		// flush before quit goroutine to avoid missing items
		defer run.Flush()

		// 新建一个计时器，固定周期触发
		ticker := run.newTicker(run.interval)
		defer ticker.Stop()

		// 外部命令立即输出
		var active bool
		last := timex.Now()
		// 开启死循环循环检测。手动指令，或者定时任务 都可以输出统计结果
		for {
			select {
			case values := <-run.messenger: // 主动触发执行
				active = true
				run.enterExecution()
				run.confirmChan <- lang.Placeholder
				run.executeTasks(values)
				last = timex.Now()
			case <-ticker.Chan(): // 定时执行
				// 如果上面主动输出一次，那么本次自动输出将轮空
				if active {
					active = false
					continue
				}
				if run.Flush() {
					last = timex.Now()
				} else if run.quitLoop(last) {
					return
				}
			}
		}
	})
}

// 一定循环次数之后发现没有新任务，自动退出定时循环。下次自动循环由添加新任务时再唤起
func (run *Interval) quitLoop(last time.Duration) (stop bool) {
	if timex.Since(last) <= run.interval*idleRound {
		return false
	}

	run.runLock.Lock()
	run.isRunning = false
	run.runLock.Unlock()

	return true
}

//// 下面这个是go-zero的写法，我认为有Bug。当自动执行任务的协程执行这个准备退出时，突然来了一项任务，inflight的值将不再是0，
//// 而这里也不会把isRunning设置成false，程序永远也无法启动loopFlush()函数了。
//func (run *Interval) quitLoop(last time.Duration) (stop bool) {
//	if timex.Since(last) <= run.interval*idleRound {
//		return false
//	}
//	// checking run.inflight and setting run.guarded should be locked together
//	run.runLock.Lock()
//	if atomic.LoadInt32(&run.inflight) == 0 {
//		run.isRunning = false
//	}
//	run.runLock.Unlock()
//	return true
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (run *Interval) doneExecution() {
	run.waitGroup.Done()
}

func (run *Interval) enterExecution() {
	run.wgLock.Lock()
	defer run.wgLock.Unlock()
	run.waitGroup.Add(1)
}

func (run *Interval) executeTasks(items any) bool {
	defer run.doneExecution()
	ok := run.hasTasks(items)
	if ok {
		run.container.Execute(items)
	}
	return ok
}

func (run *Interval) hasTasks(items any) bool {
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
