package orm

import (
	"reflect"
	"sync"
)

type ModelSchema struct {
	tableName    string
	fields       map[string]int8
	columns      []string
	primaryIndex int8 // 主键字段原始索引位置
	updatedIndex int8 // 更新字段调整之后的索引位置
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

func (ms *ModelSchema) UpdatedIndex() int8 {
	return ms.updatedIndex
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
