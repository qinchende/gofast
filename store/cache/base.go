package cache

import "time"

type (
	// Cache interface is used to define the cache implementation.
	Cache interface {
		Del(keys ...string) error
		Get(key string, v any) error
		Set(key string, v any) error
		SetExpire(key string, v any, expire time.Duration) error

		// come from beego ++++
		// Get(key string) interface{}
		// GetMulti(keys []string) []interface{}
		// Put(key string, val interface{}, timeout int64) error
		// Delete(key string) error
		// Incr(key string) error
		// Decr(key string) error
		// IsExist(key string) bool
		// ClearAll() error
		// StartAndGC(config string) error
	}
)
