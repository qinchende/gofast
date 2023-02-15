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

type LiteTaskGroup struct {
	AppName   string
	ServerNo  string
	GroupName string

	Redis *gfrds.GfRedis // Note：最好用一个实时持久化的Redis数据库，否则脚本自行解决重入问题
	Tasks []*TaskPet

	key       string
	birthTime time.Duration

	stopRun    chan lang.PlaceholderType
	lock       sync.RWMutex
	scanTimes  int64
	lostTimes  int64
	isRunning  bool // 是否正在运行
	isStopping bool // 是否正在停止任务
}

func NewLiteTaskGroup(appName, serverNo, gpName string, rds *gfrds.GfRedis) *LiteTaskGroup {
	lite := &LiteTaskGroup{
		AppName:   appName,
		ServerNo:  serverNo,
		GroupName: gpName,
		Redis:     rds,
		Tasks:     make([]*TaskPet, 0),
		key:       LiteStoreKeyPrefix + appName + "." + gpName,
		birthTime: timex.Now(),
		stopRun:   make(chan lang.PlaceholderType, 1),
	}
	return lite
}

func (lite *LiteTaskGroup) AddTask(pet *TaskPet) {
	if pet.StartTime == "" {
		pet.StartTime = DefTaskPetStartTime
	}
	if pet.EndTime == "" {
		pet.EndTime = DefTaskPetEndTime
	}

	pet.group = lite
	pet.key = LiteStoreKeyPrefix + lite.AppName + "." + lang.NameOfFunc(pet.Task) + "." + pet.StartTime

	lite.Tasks = append(lite.Tasks, pet)
}

func (lite *LiteTaskGroup) Running() {
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
			time.Sleep(3 * time.Second)
			for {
				lite.scanController()
				time.Sleep(DefLoopIntervalS * time.Second)
			}
		}()
		wg.Wait()
		goto keepAlive
	}()
}

// 任务组开始扫描运行 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (lite *LiteTaskGroup) scanController() {
	lite.scanTimes++

	// 检查争夺运行权
	if str, err := lite.Redis.Get(lite.key); err == nil && str != "" {
		if val, err2 := jsonx.UnmarshalStringToKV(str); err2 == nil {
			lite.lostTimes = 0
			if val["ServerNo"] == lite.ServerNo {
				lite.scanTimes = 0
				if val["Status"] == "1" {
					lite.keepRunning()
				}
				lite.flushStatus("1")
			} else {
				lite.killMyself()

				lastTime, err3 := time.Parse(LiteStoreTimeFormat, lang.ToString(val["Time"]))
				diff := time.Since(lastTime) / time.Second
				if err3 == nil && diff > DefLoopIntervalS*2 && lite.scanTimes >= 3 {
					lite.flushStatus("0")
				} else {
					logx.TimerF("I wait. %d", lite.scanTimes)
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
	lite.killMyself()
	if lite.lostTimes >= 2 {
		logx.Timer("I am going to run......O(∩_∩)O")
		lite.flushStatus("0")
	}
}

func (lite *LiteTaskGroup) killMyself() {
	logx.Timer("killMyself")
	lite.lock.RLock()
	if lite.isRunning && lite.isStopping == false {
		lite.stopRun <- lang.ShareVal
		lite.isStopping = true
	}
	lite.lock.RUnlock()
}

// 有些任务耗时较长，可能一直运行而很久不退出，此时其它服务器也不能抢夺运行权利
func (lite *LiteTaskGroup) keepRunning() {
	logx.Timer("keepRunning")
	lite.lock.RLock()
	if lite.isRunning {
		lite.lock.RUnlock()
		return
	}
	lite.lock.RUnlock()

	gmp.GoSafe(func() {
		defer func() {
			lite.lock.Lock()
			lite.isRunning = false
			lite.isStopping = false
			lite.lock.Unlock()
		}()

		loopTicker := time.NewTicker(DefLoopIntervalS * time.Second)
		defer loopTicker.Stop()

		lite.lock.Lock()
		lite.isRunning = true
		lite.lock.Unlock()
		for {
			select {
			case <-lite.stopRun:
				return
			case <-loopTicker.C:
				logx.Timer("begin execute tasks")
				now := timex.Now()
				for _, task := range lite.Tasks {
					task.runTask(now)
				}
				logx.Timer("execute tasks finished")
			}
		}
	})
}

func (lite *LiteTaskGroup) flushStatus(status string) {
	kvData := cst.KV{
		"Time":     timex.Time().Format(LiteStoreTimeFormat),
		"ServerNo": lite.ServerNo,
		"Status":   status,
	}
	logx.Infos(kvData)
	data, _ := jsonx.Marshal(kvData)
	str, err := lite.Redis.Set(lite.key, data, LiteStoreRunFlagExpireTTL)
	logx.Infos(str)
	logx.Infos(err)
}
