package gfrds

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/qinchende/gofast/logx"
)

// go-redis
type (
	ConnCnf struct {
		// single redis
		Addr string `v:"match=ipv4:port"`

		// sentinel
		SentinelAddr []string `v:"match=ipv4:port"`
		MasterName   string   `v:""`
		SentinelPass string   `v:""`
		SlaveOnly    bool     `v:""`

		// common
		Pass     string `v:"required"`
		DB       int    `v:""`
		PoolSize int    `v:""`
		MinIdle  int    `v:""`

		// 扩展
		Weight uint16 `v:""` // 权重
	}
	GfRedis struct {
		Cli    *redis.Client
		Ctx    context.Context
		Weight uint16
	}
)

// 直接连接redis
// go-redis 底层自带连接池功能，不需要你再管理了。
// 看看 go-redis/redis.go 中的代码：
// 第330行 func (c *baseClient) process(ctx context.Context, cmd Cmder)
// 第288行 func (c *baseClient) withConn(
// 第292行 cn, err := c.getConn(ctx)
func NewGoRedis(cf *ConnCnf) *GfRedis {
	rds := GfRedis{Ctx: context.Background(), Weight: cf.Weight}

	if cf.Addr != "" {
		rds.Cli = redis.NewClient(&redis.Options{
			Addr:         cf.Addr,
			Password:     cf.Pass,
			DB:           cf.DB,
			PoolSize:     cf.PoolSize,
			MinIdleConns: cf.MinIdle,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				logx.Info(fmt.Sprintf("%s connected.", cn.String()))
				return nil
			},
		})
		logx.Info(fmt.Sprintf("Redis %s created.", cf.Addr))
	} else if cf.SentinelAddr != nil {
		// 通过sentinel连接 redis
		rds.Cli = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       cf.MasterName,
			SentinelAddrs:    cf.SentinelAddr,
			SentinelPassword: cf.SentinelPass,
			Password:         cf.Pass,
			SlaveOnly:        cf.SlaveOnly,
			DB:               cf.DB,
			PoolSize:         cf.PoolSize,
			MinIdleConns:     cf.MinIdle,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				logx.Info(fmt.Sprintf("%s connected.", cn.String()))
				return nil
			},
		})
		logx.Info(fmt.Sprintf("Redis %s created.", cf.MasterName))
	}

	return &rds
}

//// 通过sentinel连接 redis
//func NewGoRedisBySentinel(cf *ConnConfig) *GfRedis {
//	rds := GfRedis{Ctx: context.Background()}
//	rds.Cli = redis.NewFailoverClient(&redis.FailoverOptions{
//		MasterName:       cf.MasterName,
//		SentinelAddrs:    cf.SentinelAddr,
//		SentinelPassword: cf.SentinelPass,
//		Password:         cf.Pass,
//		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
//			logx.Info(fmt.Sprintf("%s connected.", cn.String()))
//			return nil
//		},
//	})
//	logx.Info(fmt.Sprintf("Redis %s created.", cf.Addr))
//	return &rds
//}
