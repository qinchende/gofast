package cache

import "time"

type (
	// Cache interface is used to define the cache implementation.
	Cache interface {
		Del(keys ...string) error
		Get(key string, v interface{}) error
		Set(key string, v interface{}) error
		SetExpire(key string, v interface{}, expire time.Duration) error
	}
)
