package task

import (
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

// 单个任务描述 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type LitePet struct {
	Task TaskFunc

	StartTime string // "00:00"
	EndTime   string // "23:59"
	IntervalS int32  // 循环执行间隔s

	// Note: 这种情况几乎不会用到，有被删除的可能
	JustOnce   bool  // 是否只运行一次
	JustDelayS int32 // 启动之后延时多少秒执行

	group    *LiteGroup
	key      string
	lastTime time.Duration // 上次运行时间
}

func (pet *LitePet) runTask(now time.Duration) {
	// 1. 启动只执行一次的任务
	if pet.JustOnce {
		if pet.lastTime > 0 || int32(timex.DiffS(now, pet.group.createdTime)) <= pet.JustDelayS {
			return
		}
		pet.execute(now)
		return
	}

	// 2. 可能需要反复执行的任务
	// 获取上一次执行的时间
	if pet.lastTime == 0 {
		if str, err := pet.group.rds.Get(pet.key); err == nil && str != "" {
			if lst, err2 := time.Parse(time.RFC3339, str); err2 == nil {
				pet.lastTime = timex.ToDuration(&lst)
			}
		}
	}
	// 当前时间转换成 HH:MM 格式
	diff := int32(timex.DiffS(now, pet.lastTime))
	if diff >= pet.IntervalS {
		pet.lastTime = now

		nowHM := timex.ToTime(now).Format("15:04")
		if nowHM >= pet.StartTime && nowHM <= pet.EndTime {
			pet.execute(now)
		}
	}
}

func (pet *LitePet) execute(now time.Duration) {
	if pet.Task() {
		pet.group.rds.Set(pet.key, timex.ToTime(now).Format(time.RFC3339), LiteStoreRunFlagExpireTTL)
	}
}

func (pet *LitePet) resetTask() {
	pet.lastTime = 0
}
