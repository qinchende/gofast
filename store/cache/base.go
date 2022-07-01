package cache

import "time"

type (
	// Cache interface is used to define the cache implementation.
	Cache interface {
		Del(keys ...string) error
		Get(key string, v any) error
		Set(key string, v any) error
		SetExpire(key string, v any, expire time.Duration) error
	}
)
