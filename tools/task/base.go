package task

import (
	"context"
	"time"
)

type TaskFunc func(ctx context.Context) bool // 返回执行成功或失败

const (
	LiteStoreKeyPrefix        = "GF.Lite."
	LiteStoreRunFlagExpireS   = 30 * 24 * 3600 // key的默认有效期秒
	LiteStoreRunFlagExpireTTL = LiteStoreRunFlagExpireS * time.Second

	DefLiteRunIntervalS = 60 // 循环运行任务的间隔时间second
	DefLitePetStartTime = "00:00"
	DefLitePetEndTime   = "23:59"
)
