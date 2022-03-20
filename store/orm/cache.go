package orm

import (
	"reflect"
	"sync"
)

var cachedModels sync.Map

type cachedModel struct {
	fieldIndexMap map[string]int
	fields        []string
}

func cacheSetModel(typ reflect.Type, val *cachedModel) {
	cachedModels.Store(typ, val)
}

func cacheGetModel(typ reflect.Type) *cachedModel {
	if ret, ok := cachedModels.Load(typ); ok {
		return ret.(*cachedModel)
	}
	return nil
}
