package task

import (
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
	waitTimes   int64
	lostTimes   int64
	isRunning   bool // 是否正在运行
	isStopping  bool // 是否正在停止任务
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

	pet.group = lite
	pet.key = LiteStoreKeyPrefix + "Task." + lite.appName + "." + lang.FuncName(pet.Task) + "." + pet.StartTime

	lite.tasks = append(lite.tasks, pet)
}

func (lite *LiteGroup) Running() {
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

// 任务组开始扫描运行 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (lite *LiteGroup) scanController() {
	// 检查争夺运行权
	if str, err := lite.rds.Get(lite.key); err == nil && str != "" {
		if val, err2 := jsonx.UnmarshalStringToKV(str); err2 == nil {
			lite.lostTimes = 0
			if val["ServerNo"] == lite.serverNo {
				lite.waitTimes = 0
				if val["Status"] == "1" {
					lite.flushTime(val)
					lite.keepRunning()
				} else {
					lite.flushStatus("1")
				}
			} else {
				lite.waitTimes++
				lite.killMyself()

				lst, _ := time.Parse(time.RFC3339, lang.ToString(val["Time"]))
				diff := timex.SinceS(timex.ToDuration(&lst))

				// NOTE：要避免多个服务器竞争
				if diff > DefLiteRunIntervalS*2 && lite.waitTimes > 3 {
					lite.flushStatus("0")
				} else {
					logx.TimerF("I wait. %d", lite.waitTimes)
				}
			}
		} else {
			goto unknownFlag
		}
		return
	}

unknownFlag:
	lite.lostTimes++
	// TODO: 查不到redis数据或者数据解析错误，需要先关闭自己，然后试图夺取控制权
	if lite.lostTimes <= 1 {
		logx.TimerF("%s. Can't check status.", lite.key)
		return
	} else {
		lite.killMyself()
	}
	if lite.lostTimes >= 3 {
		lite.flushStatus("0")
	} else {
		logx.Timer("Oh, I am ready to run...")
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

		// TODO: 如果定时没有处理完，又产生到时信号，会如何
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
					task.runTask(now)
				}
			}
		}
	})
}

func (lite *LiteGroup) flushStatus(status string) {
	logx.Timer("I am try to run. set status: " + status)

	kvData := cst.KV{
		"Time":     time.Now().Format(time.RFC3339),
		"ServerNo": lite.serverNo,
		"Status":   status,
	}
	data, _ := jsonx.Marshal(kvData)
	_, _ = lite.rds.Set(lite.key, data, LiteStoreRunFlagExpireTTL)
}

func (lite *LiteGroup) flushTime(kvData cst.KV) {
	kvData["Time"] = time.Now().Format(time.RFC3339)
	data, _ := jsonx.Marshal(kvData)
	_, _ = lite.rds.Set(lite.key, data, LiteStoreRunFlagExpireTTL)
}
