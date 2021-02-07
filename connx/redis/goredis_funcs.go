// 下面封装一些常用的命令函数，不常用的自己用标准的调用方法。
package redis

import "time"

func (rdx *GoRedisX) Ping() (string, error) {
	return rdx.Cli.Ping(rdx.Ctx).Result()
}

func (rdx *GoRedisX) Get(key string) (string, error) {
	return rdx.Cli.Get(rdx.Ctx, key).Result()
}

func (rdx *GoRedisX) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	return rdx.Cli.Set(rdx.Ctx, key, value, expiration).Result()
}
