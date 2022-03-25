package orm

import (
	"reflect"
	"sync"
)

// 表结构体Schema, 限制表最多127列（用int8计数）
type ModelSchema struct {
	tableName    string          // 数据库表名称
	fields       map[string]int8 // Struct的字段
	fieldsIndex  [][]int         // Struct标准字段索引
	columns      []string        // 对应数据库的列名称
	primaryIndex int8            // 主键字段原始索引位置
	updatedIndex int8            // 更新字段原始索引位置，没有则为-1

	insertSQL string // 全字段insert（将来会建立通用缓存中心，这里暂时这样用）
	updateSQL string // 全字段update
	deleteSQL string // delete
	selectSQL string // select
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

func (ms *ModelSchema) PrimaryIndex() int8 {
	return ms.primaryIndex
}

func (ms *ModelSchema) InsertSQL(fn func(*ModelSchema) string) string {
	if ms.insertSQL == "" {
		ms.insertSQL = fn(ms)
	}
	return ms.insertSQL
}

func (ms *ModelSchema) UpdateSQL(fn func(*ModelSchema) string) string {
	if ms.updateSQL == "" {
		ms.updateSQL = fn(ms)
	}
	return ms.updateSQL
}

func (ms *ModelSchema) SelectSQL(fn func(*ModelSchema) string) string {
	if ms.selectSQL == "" {
		ms.selectSQL = fn(ms)
	}
	return ms.selectSQL
}

func (ms *ModelSchema) DeleteSQL(fn func(*ModelSchema) string) string {
	if ms.deleteSQL == "" {
		ms.deleteSQL = fn(ms)
	}
	return ms.deleteSQL
}

func (ms *ModelSchema) PrimaryValue(obj ApplyOrmStruct) interface{} {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	return rVal.FieldByIndex(ms.fieldsIndex[ms.primaryIndex]).Interface()
}

func (ms *ModelSchema) ValueByIndex(rVal *reflect.Value, index int8) interface{} {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Interface()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 缓存数据表的Schema
var cachedSchemas sync.Map

func cacheSetSchema(typ reflect.Type, val *ModelSchema) {
	cachedSchemas.Store(typ, val)
}

func cacheGetSchema(typ reflect.Type) *ModelSchema {
	if ret, ok := cachedSchemas.Load(typ); ok {
		return ret.(*ModelSchema)
	}
	return nil
}
