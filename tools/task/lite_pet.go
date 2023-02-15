package task

import (
	"github.com/qinchende/gofast/skill/timex"
	"strconv"
	"time"
)

// 单个任务描述 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type TaskPet struct {
	group    *LiteTaskGroup
	key      string
	lastTime time.Duration // 上次运行时间

	Task TaskFunc

	JustOnce   bool  // 是否只运行一次
	JustDelayS int32 // 启动之后延时多少秒执行

	StartTime string // "00:00"
	EndTime   string // "23:59"
	IntervalS int32  // 循环执行间隔秒
}

func (pet *TaskPet) runTask(now time.Duration) {
	// 启动只执行一次的任务
	if pet.JustOnce {
		if pet.lastTime > 0 || int32(now-pet.group.birthTime) <= pet.JustDelayS {
			return
		}
		if pet.Task() {
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

	// 当前时间转换成 HH:MM 格式
	strNow := timex.ToTime(now).Format("15:04")
	diff := int32((now - pet.lastTime) / time.Second)
	if diff >= pet.IntervalS && strNow >= pet.StartTime && strNow <= pet.EndTime {
		if pet.Task() {
			pet.lastTime = now
		}
		_, _ = pet.group.Redis.Set(pet.key, pet.lastTime, time.Duration(pet.IntervalS)*time.Second)
	}
}
