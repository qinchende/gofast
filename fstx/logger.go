package fstx

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/stat"
)

// 全局的初始化日志
func InitLogger(lgConfig *logx.LogConf) {
	logx.MustSetup(*lgConfig)

	// log初始化完毕，接下来解析
	if lgConfig.NeedCpuMem {
		stat.StartCpuMemCollect()
	}
}
