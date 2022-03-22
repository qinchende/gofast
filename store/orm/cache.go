package orm

import (
	"reflect"
	"sync"
)

type ModelSchema struct {
	tableName    string
	fields       map[string]int8
	columns      []string
	primaryIndex int8
	updatedIndex int8
}

func (ms *ModelSchema) Length() int8 {
	return int8(len(ms.columns))
}

func (ms *ModelSchema) TableName() string {
	return ms.tableName
}

func (ms *ModelSchema) Fields() map[string]int8 {
	return ms.fields
}

func (ms *ModelSchema) Columns() []string {
	return ms.columns
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 缓存数据表的Schema
var cachedModels sync.Map

func cacheSetModel(typ reflect.Type, val *ModelSchema) {
	cachedModels.Store(typ, val)
}

func cacheGetModel(typ reflect.Type) *ModelSchema {
	if ret, ok := cachedModels.Load(typ); ok {
		return ret.(*ModelSchema)
	}
	return nil
}
