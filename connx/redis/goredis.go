package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// go-redis
type (
	ConnConfig struct {
		Addr     string `json:",optional"`
		Pass     string `json:",optional"`
		DB       int    `json:",optional"`
		PoolSize int    `json:",optional"`
		MinIdle  int    `json:",optional"`
	}
	ConnSentinelConfig struct {
		ConnConfig
		MasterName   string   `json:",optional"`
		SentinelAddr []string `json:",optional"`
		SentinelPass string   `json:",optional"`
	}
	GoRedisX struct {
		Cli *redis.Client
		Ctx context.Context
	}
)

// 直接连接redis
// go-redis 底层自带连接池功能，不需要你再管理了。
func NewGoRedis(cf *ConnConfig) *GoRedisX {
	rds := GoRedisX{Ctx: context.Background()}
	rds.Cli = redis.NewClient(&redis.Options{
		Addr:     cf.Addr,
		Password: cf.Pass,
		DB:       cf.DB,
	})
	return &rds
}

// 通过sentinel连接 redis
func NewGoRedisBySentinel(cf *ConnSentinelConfig) *GoRedisX {
	rds := GoRedisX{Ctx: context.Background()}
	rds.Cli = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       cf.MasterName,
		SentinelAddrs:    cf.SentinelAddr,
		SentinelPassword: cf.SentinelPass,
		Password:         cf.Pass,
	})
	return &rds
}
