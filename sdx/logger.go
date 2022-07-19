package sdx

import (
	"github.com/qinchende/gofast/logx"
)

// 全局的初始化日志
func InitLogger(cfg *logx.LogConfig) {
	logx.MustSetup(*cfg)
}
