package task

import (
	"github.com/qinchende/gofast/connx/gfrds"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/gmp"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/timex"
	"strconv"
	"time"
)

type LiteTaskGroup struct {
	AppName   string
	Redis     *gfrds.GfRedis // Note：最好用一个实时持久化的Redis数据库，否则脚本自行解决重入问题
	Tasks     []*TaskPet
	IntervalS int32

	createTime time.Duration
}

func NewLiteTaskGroup(name string, rds *gfrds.GfRedis) *LiteTaskGroup {
	lite := &LiteTaskGroup{
		AppName:    name,
		Redis:      rds,
		Tasks:      make([]*TaskPet, 0),
		IntervalS:  DefGroupIntervalS,
		createTime: timex.Now(),
	}
	return lite
}

func (lite *LiteTaskGroup) AddTask(pet *TaskPet) {
	pet.group = lite
	pet.key = LiteRedisKeyPrefix + lite.AppName + "." + lang.NameOfFunc(pet.Func) + "." + pet.StartTime
	if pet.StartTime == "" {
		pet.StartTime = DefTaskPetStartTime
	}
	if pet.EndTime == "" {
		pet.EndTime = DefTaskPetEndTime
	}
	lite.Tasks = append(lite.Tasks, pet)
}

func (lite *LiteTaskGroup) Run() {
	// 因为在主程序中启动的协程运行，主程序不能异常退出，脚本必须安全运行
	gmp.GoSafe(func() {
		checkTicker := time.NewTicker(time.Duration(lite.IntervalS) * time.Second)
		defer checkTicker.Stop()

		for {
			// 检查执行权
			if false {
				time.Sleep(30 * time.Second)
				continue
			}

			select {
			case <-checkTicker.C:
				logx.Info("checkTicker coming")
				lite.checkExecute()
			}
		}
	})
}

func (lite *LiteTaskGroup) checkExecute() {
	now := timex.Now()
	for _, task := range lite.Tasks {
		task.runTask(now)
	}
}

// 单个任务描述 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type TaskPet struct {
	group *LiteTaskGroup
	key   string
	Func  TaskFunc

	JustOnce   bool  // 是否只运行一次
	JustDelayS int32 // 启动之后延时多少秒执行

	StartTime string // "00:00"
	EndTime   string // "23:59"
	IntervalS int32  // 循环执行间隔秒

	lastTime time.Duration // 上次运行时间
}

func (pet *TaskPet) runTask(now time.Duration) {
	// 启动只执行一次的任务
	if pet.JustOnce {
		if pet.lastTime > 0 || int32(now-pet.group.createTime) <= pet.JustDelayS {
			return
		}
		if pet.Func() {
			pet.lastTime = now
		}
		return
	}

	// 可能需要反复执行的任务
	// 读取redis信息
	str, err := pet.group.Redis.Get(pet.key)
	if err != nil && str != "" {
		if val, err2 := strconv.Atoi(str); err2 != nil {
			pet.lastTime = time.Duration(val)
		}
	}

	strNow := timex.ToTime(now).Format("15:04")
	diff := int32((now - pet.lastTime) / time.Second)
	if diff >= pet.IntervalS && strNow >= pet.StartTime && strNow <= pet.EndTime {
		if pet.Func() {
			pet.lastTime = now
		}
		_, _ = pet.group.Redis.Set(pet.key, pet.lastTime, time.Duration(pet.IntervalS)*time.Second)
	}
}
