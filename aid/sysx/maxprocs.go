package sysx

import (
	"go.uber.org/automaxprocs/maxprocs"
	"log"
)

// 自动设置
// Automatically set GOMAXPROCS to match Linux container CPU quota.
func initClose() {
	info, err := maxprocs.Set(maxprocs.Logger(log.Printf))
	info()
	if err != nil {
		log.Printf("Auto set GOMAXPROCS err: %#v", err)
	}
}
