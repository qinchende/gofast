// 下面封装一些常用的命令函数，不常用的自己用标准的调用方法。
package gfrds

import (
	"time"
)

func (rdx *GfRedis) Ping() (string, error) {
	return rdx.Cli.Ping(rdx.Ctx).Result()
}

func (rdx *GfRedis) Get(key string) (string, error) {
	return rdx.Cli.Get(rdx.Ctx, key).Result()
}

func (rdx *GfRedis) Set(key string, value interface{}, seconds time.Duration) (string, error) {
	return rdx.Cli.Set(rdx.Ctx, key, value, seconds).Result()
}

func (rdx *GfRedis) Expire(key string, seconds time.Duration) (bool, error) {
	return rdx.Cli.Expire(rdx.Ctx, key, seconds).Result()
}

func (rdx *GfRedis) Del(key string) (int64, error) {
	return rdx.Cli.Del(rdx.Ctx, key).Result()
}
