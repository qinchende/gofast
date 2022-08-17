package gmp

import (
	"github.com/qinchende/gofast/logx"
)

// 启动新协程跑函数
func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer func() {
		if p := recover(); p != nil {
			logx.Stacks(p)
		}
	}()

	fn()
}
