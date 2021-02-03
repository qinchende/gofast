package redis

import (
	"encoding/json"
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

func (rdx *RedigoX) Set(key string, data interface{}, time int) error {
	conn := rdx.conn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func (rdx *RedigoX) GetCommon(key string) (string, error) {
	conn := rdx.conn.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	//reply, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	return reply, nil
}

func (rdx *RedigoX) Get(key string) ([]byte, error) {
	conn := rdx.conn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	//reply, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (rdx *RedigoX) Exists(key string) bool {
	conn := rdx.conn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func (rdx *RedigoX) Del(key string) (bool, error) {
	conn := rdx.conn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func (rdx *RedigoX) LikeDeletes(key string) error {
	conn := rdx.conn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, v := range keys {
		_, err := rdx.Del(v)
		if err != nil {
			return err
		}
	}
	return nil
}
