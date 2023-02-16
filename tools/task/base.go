package task

import "time"

const (
	LiteStoreKeyPrefix        = "GF.Lite."
	LiteStoreRunFlagExpireS   = 3600 * 24 * 30 // key的默认有效期
	LiteStoreRunFlagExpireTTL = LiteStoreRunFlagExpireS * time.Second
	//LiteStoreTimeFormat       = "2006-01-02 15:04:05"

	DefLiteRunIntervalS = 60 // 循环运行任务的间隔时间second
	DefLitePetStartTime = "00:00"
	DefLitePetEndTime   = "23:59"
)

type TaskFunc func() bool
