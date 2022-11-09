package exec

import (
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/proc"
	"github.com/qinchende/gofast/skill/timex"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

const idleRound = 10

// 间隔定时执行器（间隔执行某些任务，长时间没任务就退出执行器。有新任务自动启动）
type (
	AddFunc func(item any) (any, bool)

	// 外部可以自定义装载各种任务的容器，实现这些方法之后，就能赋予周期执行的特性
	TaskContainer interface {
		AddItem(item any) bool // 返回true，立即触发任务执行
		RemoveAll() any
		Execute(items any)
	}

	// 管理周期执行的实体对象
	Interval struct {
		container   TaskContainer
		interval    time.Duration
		newTicker   func(d time.Duration) timex.Ticker
		messenger   chan any
		confirmChan chan lang.PlaceholderType
		inflight    int32

		taskLock sync.Mutex
		execWG   sync.WaitGroup

		isRunning bool
		runLock   sync.Mutex
	}
)

func createInterval(dur time.Duration, box TaskContainer) *Interval {
	return &Interval{
		messenger:   make(chan any, 1),               // 长度为1的有缓冲通道(buffer 1 to let the caller go quickly)
		confirmChan: make(chan lang.PlaceholderType), // 无缓冲通道
		interval:    dur,
		container:   box,
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(d)
		},
	}
}

// 初始化一个周期执行的实体对象
func NewInterval(dur time.Duration, box TaskContainer) *Interval {
	run := createInterval(dur, box)
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
func (run *Interval) Flush() (hasTasks bool) {
	run.enterExecution()

	run.taskLock.Lock()
	items := run.container.RemoveAll()
	run.taskLock.Unlock()

	return run.executeTasks(items)
}

// 请调用者自己保证AddItem函数并发安全
func (run *Interval) Add(item any) {
	run.checkLoop()

	run.taskLock.Lock()
	ok := run.container.AddItem(item)
	run.taskLock.Unlock()

	run.checkLoop()

	// ok == true -> 立即主动执行任务
	if ok {
		run.taskLock.Lock()
		items := run.container.RemoveAll()
		run.taskLock.Unlock()

		atomic.AddInt32(&run.inflight, 1)
		run.messenger <- items
		<-run.confirmChan
	}
}

// 请调用者自己保证fc函数并发安全
func (run *Interval) AddByFunc(fc AddFunc, item any) {
	run.checkLoop()

	run.taskLock.Lock()
	items, ok := fc(item)
	run.taskLock.Unlock()

	run.checkLoop()

	// ok == true -> 立即主动执行任务
	if ok {
		atomic.AddInt32(&run.inflight, 1)
		run.messenger <- items
		<-run.confirmChan
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 有增加任务的动作，就要想办法激活循环检测
// 因为这个没有加锁运行，所以在添加任务的时候要前后都检测一次。
func (run *Interval) checkLoop() {
	if run.isRunning == false {
		run.raiseLoop()
	}
}

// 后台启动新的协程运行周期任务
func (run *Interval) raiseLoop() {
	run.runLock.Lock()
	if run.isRunning {
		run.runLock.Unlock()
		return
	}
	run.isRunning = true
	run.runLock.Unlock()

	gmp.GoSafe(func() {
		defer run.Flush()                     // flush before quit goroutine to avoid missing items
		ticker := run.newTicker(run.interval) // 定时器
		defer ticker.Stop()                   // 退出定时器

		var active bool // 外部命令立即输出
		last := timex.Now()
		// 开启死循环循环检测。手动指令，或者定时任务 都可以输出统计结果
		for {
			select {
			case items := <-run.messenger: // 主动触发执行
				active = true
				atomic.AddInt32(&run.inflight, -1)
				run.enterExecution()
				run.confirmChan <- lang.Placeholder
				run.executeTasks(items)
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

// 而这里也不会把 guarded 设置成false，程序永远也无法启动 backgroundFlush()函数了。
func (run *Interval) quitLoop(last time.Duration) (stop bool) {
	if timex.Since(last) <= run.interval*idleRound {
		return
	}

	run.runLock.Lock()
	if atomic.LoadInt32(&run.inflight) == 0 {
		run.isRunning = false
		stop = true
	}
	run.runLock.Unlock()

	return
}

// 执行体的安全控制
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (run *Interval) Wait() {
	run.Flush()
	run.execWG.Wait()
}

func (run *Interval) enterExecution() {
	run.execWG.Add(1)
}

func (run *Interval) executeTasks(items any) (hasTasks bool) {
	defer run.execWG.Done()
	hasTasks = run.hasTasks(items)
	if hasTasks {
		run.container.Execute(items)
	}
	return
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
