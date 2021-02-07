package redis2

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// Redigo初始化的配置参数
type (
	RedigoConfig struct {
		Type        string `json:",default=node,options=node|cluster"`
		Addr        string `json:",optional"`
		Pass        string `json:",optional"`
		Key         string `json:",optional"`
		MaxIdleConn int    `json:",optional"`
		MaxOpenConn int    `json:",optional"`
	}
	RedigoX struct {
		conn *redis.Pool
	}
)

func NewRedigo(cf *RedigoConfig) *RedigoX {
	rdgX := RedigoX{}
	rdgX.conn = &redis.Pool{
		MaxIdle:     cf.MaxIdleConn,
		MaxActive:   cf.MaxOpenConn,
		IdleTimeout: 200,
		Dial: func() (redis.Conn, error) {
			// addr := fmt.Sprintf("%s:%s", cf.Host, cf.Port)
			c, err := redis.Dial("tcp", cf.Addr)
			if err != nil {
				fmt.Printf("Redigo error: %v", err)
				//os.Exit(1)
				return nil, err
			}

			if cf.Pass != "" {
				if _, err := c.Do("AUTH", cf.Pass); err != nil {
					_ = c.Close()
					fmt.Printf("Redigo error: %v", err)
					//os.Exit(1)
					return nil, err
				}
			}
			println("Redis " + cf.Addr + " conn success.")
			return c, err
		},
		//TestOnBorrow: func(c redis.Conn, t time.Time) error {
		//	_, err := c.Do("PING")
		//	return err
		//},
	}
	return &rdgX
}
