package task

import "time"

const (
	LiteStoreKeyPrefix        = "GF.LT."
	LiteStoreRunFlagExpireS   = 3600 * 24 * 30
	LiteStoreRunFlagExpireTTL = LiteStoreRunFlagExpireS * time.Second
	LiteStoreTimeFormat       = "2006-01-02 15:04:05"

	DefLoopIntervalS    = 60 // 循环运行任务的间隔时间
	DefTaskPetStartTime = "00:00"
	DefTaskPetEndTime   = "23:59"
)

type TaskFunc func() bool
