package task

import (
	"context"
	"time"
)

type TaskFunc func(ctx context.Context) bool // 返回执行成功或失败

const (
	liteStoreKeyPrefix        = "GF.Lite."
	liteStoreRunFlagExpireS   = 30 * 24 * 3600 // key的默认有效期秒
	liteStoreRunFlagExpireTTL = liteStoreRunFlagExpireS * time.Second
)
