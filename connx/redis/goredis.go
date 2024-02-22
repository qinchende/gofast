package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/qinchende/gofast/core/logx"
)

// go-redis
type (
	ConnCnf struct {
		// single redis
		Addr string `v:"match=ipv4:port"`

		// sentinel
		SentinelAddrs []string `v:"match=ipv4:port"`
		MasterName    string   `v:""`
		SentinelPass  string   `v:""`
		SlaveOnly     bool     `v:""`

		// common
		Pass     string `v:"must"`
		DB       int    `v:""`
		PoolSize int    `v:""`
		MinIdle  int    `v:""`

		// 扩展：节点的权重
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
		// 连接到具体的某个Redis
		rds.Cli = redis.NewClient(&redis.Options{
			Addr:         cf.Addr,
			Password:     cf.Pass,
			DB:           cf.DB,
			PoolSize:     cf.PoolSize,
			MinIdleConns: cf.MinIdle,

			//OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			//	//logx.Info(fmt.Sprintf("%s connected.", cn.String()))
			//	return nil
			//},
		})
		logx.Info(fmt.Sprintf("Redis alone %s created.", cf.Addr))
		_, err := rds.Ping()
		if err != nil {
			logx.Error(fmt.Sprintf("Redis alone %s connection error: %s", cf.Addr, err))
		}
	} else if cf.SentinelAddrs != nil {
		// 通过sentinel连接 redis
		rds.Cli = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs:    cf.SentinelAddrs,
			MasterName:       cf.MasterName,
			SentinelPassword: cf.SentinelPass,
			SlaveOnly:        cf.SlaveOnly,

			Password:     cf.Pass,
			DB:           cf.DB,
			PoolSize:     cf.PoolSize,
			MinIdleConns: cf.MinIdle,

			//OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			//	//logx.Info(fmt.Sprintf("%s connected.", cn.String()))
			//	return nil
			//},
		})

		roleName := "master"
		if cf.SlaveOnly == true {
			roleName = "slave"
		}
		logx.Info(fmt.Sprintf("Redis sentinels %s for %s(%s) created.", cf.SentinelAddrs, cf.MasterName, roleName))
		_, err := rds.Ping()
		if err != nil {
			logx.Error(fmt.Sprintf("Redis %s(%s) connection error: %s", cf.MasterName, roleName, err))
		}
	}

	return &rds
}
