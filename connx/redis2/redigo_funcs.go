package redis2

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

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
