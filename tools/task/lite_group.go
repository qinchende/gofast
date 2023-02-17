package task

import (
	"context"
	"github.com/qinchende/gofast/connx/gfrds"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/timex"
	"sync"
	"time"
)

// Note: LiteGroup中的任务是支持分布式部署的。在应用多机部署的时候能满足高可用(这一点需要Redis数据库的保证)
type LiteGroup struct {
	appName   string
	serverNo  string
	groupName string
	tasks     []*LitePet

	rds *gfrds.GfRedis // Note：最好用一个实时持久化的Redis数据库
	key string

	createdTime time.Duration
	stopRun     chan lang.PlaceholderType
	lock        sync.RWMutex

	lastMark   string // 上次的运行标记
	waitTimes  int64  // 等待循环次数
	lostTimes  int64  // 无法正确获取标记数据的次数
	isRunning  bool   // 是否正在运行
	isStopping bool   // 是否正在停止任务
}

func NewLiteGroup(appName, serverNo, gpName string, rds *gfrds.GfRedis) *LiteGroup {
	return &LiteGroup{
		appName:     appName,
		serverNo:    serverNo,
		groupName:   gpName,
		tasks:       make([]*LitePet, 0),
		rds:         rds,
		key:         LiteStoreKeyPrefix + "Group." + appName + "." + gpName,
		createdTime: timex.Now(),
		stopRun:     make(chan lang.PlaceholderType, 1),
	}
}

func (lite *LiteGroup) AddTask(pet *LitePet) {
	if pet.StartTime == "" {
		pet.StartTime = DefLitePetStartTime
	}
	if pet.EndTime == "" {
		pet.EndTime = DefLitePetEndTime
	}
	if pet.IntervalS == 0 {
		pet.IntervalS = DefLiteRunIntervalS
	}
	if pet.EndTime < pet.StartTime {
		pet.crossDay = true
	}

	pet.group = lite
	pet.key = LiteStoreKeyPrefix + "Task." + lite.appName + "." + lang.FuncName(pet.Task) + "." + pet.StartTime

	lite.tasks = append(lite.tasks, pet)
}

func (lite *LiteGroup) StartRun() {
	// 因为在主程序中启动的协程运行，主程序不能异常退出，脚本必须安全运行
	go func() {
		wg := new(sync.WaitGroup)

	keepAlive:
		wg.Add(1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					logx.TimerError(lang.ToString(p))
				}
				wg.Done()
			}()
			time.Sleep(3 * time.Second) // 启动3秒之后再检查
			for {
				lite.scanController()
				time.Sleep(DefLiteRunIntervalS / 2 * time.Second)
			}
		}()
		wg.Wait()

		goto keepAlive
	}()
}

// 任务组开始扫描运行 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (lite *LiteGroup) scanController() {
	// 检查争夺运行权
	if str, err := lite.rds.Get(lite.key); err == nil && str != "" {
		if kvs, err2 := jsonx.UnmarshalStringToKV(str); err2 == nil {
			lite.lostTimes = 0

			if kvs["ServerNo"] == lite.serverNo {
				lite.waitTimes = 0

				if kvs["Status"] == "1" {
					lite.keepRunning()
					lite.flushTime(kvs)
				} else {
					lite.flushStatus(kvs, "1")
				}
			} else {
				lite.killMyself()

				if lite.lastMark == str {
					if lite.waitTimes > 0 {
						lite.waitTimes = 0
					}
					lite.waitTimes--
				} else {
					lite.lastMark = str
					lite.waitTimes++
				}

				// NOTE：要避免多个服务器竞争
				if lite.waitTimes < -4 {
					lite.flushStatus(kvs, "0")
				} else {
					logx.TimerF("Run by %s, wait %d", kvs["ServerNo"], lite.waitTimes)
				}
			}
		} else {
			goto lostFlag
		}
		return
	}

lostFlag:
	lite.lostTimes++

	// Note: 查不到redis数据或者数据解析错误，需要先关闭自己，然后试图夺取控制权
	if lite.lostTimes <= 1 {
		logx.TimerF("%s. Can't check status. %d", lite.key, lite.lostTimes)
		return
	} else {
		lite.killMyself()
	}
	if lite.lostTimes > 4 {
		lite.flushStatus(nil, "0")
	} else {
		logx.TimerF("Oh, Maybe it's my turn. %d", lite.lostTimes)
	}
}

func (lite *LiteGroup) killMyself() {
	lite.lock.RLock()
	if lite.isRunning && lite.isStopping == false {
		logx.Timer("Send stop sign to kill myself.")
		lite.stopRun <- lang.ShareVal
		lite.isStopping = true
	}
	lite.lock.RUnlock()
}

// 有些任务耗时较长，可能一直运行而很久不退出，此时其它服务器也不能抢夺运行权利
func (lite *LiteGroup) keepRunning() {
	lite.lock.RLock()
	if lite.isRunning {
		lite.lock.RUnlock()
		return
	}
	lite.lock.RUnlock()

	logx.Timer("Now start to run tasks.")
	gmp.GoSafe(func() {
		defer func() {
			lite.lock.Lock()
			lite.isRunning = false
			lite.isStopping = false
			lite.lock.Unlock()
		}()

		// 防止底层执行任务时，启动协程，而这里退出时协程泄露
		gorCtx, cancel := context.WithCancel(context.Background())
		defer cancel() // 结束协程链

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		lite.lock.Lock()
		lite.isRunning = true
		lite.lock.Unlock()
		for {
			select {
			case <-lite.stopRun:
				for _, task := range lite.tasks {
					task.resetTask()
				}
				return
			case <-ticker.C:
				now := timex.Now()
				for _, task := range lite.tasks {
					task.runTask(gorCtx, now)
				}
			}
		}
	})
}

// 设置分布式锁 +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (lite *LiteGroup) flushStatus(kvs cst.KV, status string) {
	logx.Timer("I am try to run. set status: " + status)

	if kvs == nil {
		kvs = cst.KV{}
	}
	kvs["ServerNo"] = lite.serverNo
	kvs["Status"] = status

	lite.flushTime(kvs)
}

func (lite *LiteGroup) flushTime(kvs cst.KV) {
	kvs["Time"] = time.Now().Format(cst.TimeFmtSaveReload)

	jsonStr, _ := jsonx.Marshal(kvs)
	_, _ = lite.rds.Set(lite.key, jsonStr, LiteStoreRunFlagExpireTTL)
}
