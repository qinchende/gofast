package task

import (
	"context"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

// 任务项  +++++++++++++++++++++++++++++++++++++++
type LitePet struct {
	Task TaskFunc

	StartTime string        `v:"def=00:00,len=[5:5]"` // "00:00"
	EndTime   string        `v:"def=23:59,len=[5:5]"` // "23:59"
	IntervalS time.Duration `v:"def=60s"`             // 循环执行间隔s
	crossDay  bool          ``                        // 定时任务是否可跨日运行

	JustOnce   bool          `v:""`        // 是否只运行一次
	JustDelayS time.Duration `v:"def=60s"` // 启动之后延时多少秒执行

	group    *LiteGroup    `v:""`            // 分组
	key      string        `v:"len=[1:250]"` // 任务运行标记数据对应的key
	lastTime time.Duration `v:""`            // 上次运行时间
}

// 判断条件，运行任务+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (pet *LitePet) runTask(gorCtx context.Context, now time.Duration) {
	// 1. 启动只执行一次的任务
	if pet.JustOnce {
		if pet.lastTime > 0 || (now-pet.group.createdTime) <= pet.JustDelayS {
			return
		}
		pet.execute(gorCtx, now)
		return
	}

	// 2. 可能需要反复执行的任务
	// 获取上一次执行的时间
	if pet.lastTime == 0 {
		if str, err := pet.group.rds.Get(pet.key); err == nil && str != "" {
			if lst, err2 := time.Parse(cst.TimeFmtRFC3339, str); err2 == nil {
				pet.lastTime = timex.ToDuration(&lst)
			}
		}
	}

	// 上次运行到现在的时间差
	diffDur := now - pet.lastTime
	if diffDur >= pet.IntervalS {
		pet.lastTime = now

		// 当前时间转换成 HH:MM 格式
		nowHM := timex.ToTime(now).Format("15:04")
		if (pet.crossDay && (nowHM >= pet.StartTime || nowHM <= pet.EndTime)) ||
			(!pet.crossDay && nowHM >= pet.StartTime && nowHM <= pet.EndTime) {
			pet.execute(gorCtx, now)
		}
	}
}

func (pet *LitePet) execute(gorCtx context.Context, now time.Duration) {
	if pet.Task(gorCtx) {
		pet.group.rds.Set(pet.key, timex.ToTime(now).Format(cst.TimeFmtRFC3339), liteStoreRunFlagExpireTTL)
	}
}

func (pet *LitePet) resetTask() {
	pet.lastTime = 0
}
