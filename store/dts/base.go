package dts

import (
	"errors"
	"reflect"
	"sync"
)

var (
	errNumOutOfRange  = errors.New("dts: number out of range")
	errNotSupportType = errors.New("dts: can't support the value type")
)

// 缓存所有需要反序列化的实体结构的解析数据，防止反复不断的进行反射解析操作。
var cachedStructSchemas sync.Map

func cacheSetSchema(typ reflect.Type, val *StructSchema) {
	cachedStructSchemas.Store(typ, val)
}

func cacheGetSchema(typ reflect.Type) *StructSchema {
	if ret, ok := cachedStructSchemas.Load(typ); ok {
		return ret.(*StructSchema)
	}
	return nil
}
