package orm

import (
	"fmt"
	"github.com/qinchende/gofast/skill/hashx"
	"reflect"
	"time"
)

func (ms *TableSchema) TableName() string {
	return ms.tAttrs.TableName
}

func (ms *TableSchema) CacheAll() bool {
	return ms.tAttrs.CacheAll
}

func (ms *TableSchema) CachePreFix() string {
	return ms.tAttrs.cacheKeyFmt
}

func (ms *TableSchema) CacheLineKey(dbName, id any) string {
	return fmt.Sprintf(ms.tAttrs.cacheKeyFmt, dbName, id)
}

func (ms *TableSchema) CacheSqlKey(sql string) string {
	return "Gf#Pet#" + hashx.Md5HexString(sql)
}

func (ms *TableSchema) ExpireS() uint32 {
	return ms.tAttrs.ExpireS
}

// 可以考虑加上随机 5% 左右的偏差，防止将来缓存统一过期导致缓存雪崩
func (ms *TableSchema) ExpireDuration() time.Duration {
	return time.Duration(ms.tAttrs.ExpireS) * time.Second
}

func (ms *TableSchema) Columns() []string {
	return ms.ss.Columns
}

func (ms *TableSchema) UpdatedIndex() int8 {
	return ms.updatedIndex
}

func (ms *TableSchema) PrimaryIndex() int8 {
	return ms.primaryIndex
}

func (ms *TableSchema) AutoIndex() int8 {
	return ms.autoIndex
}

func (ms *TableSchema) ColumnIndex(k string) int {
	return ms.ss.ColumnIndex(k)
}

func (ms *TableSchema) FieldIndex(k string) int {
	return ms.ss.FieldIndex(k)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (ms *TableSchema) InsertSQL(fn func(*TableSchema) string) string {
	if ms.insertSQL == "" {
		ms.insertSQL = fn(ms)
	}
	return ms.insertSQL
}

func (ms *TableSchema) UpdateSQL(fn func(*TableSchema) string) string {
	if ms.updateSQL == "" {
		ms.updateSQL = fn(ms)
	}
	return ms.updateSQL
}

func (ms *TableSchema) SelectSQL(fn func(*TableSchema) string) string {
	if ms.selectSQL == "" {
		ms.selectSQL = fn(ms)
	}
	return ms.selectSQL
}

func (ms *TableSchema) DeleteSQL(fn func(*TableSchema) string) string {
	if ms.deleteSQL == "" {
		ms.deleteSQL = fn(ms)
	}
	return ms.deleteSQL
}

func (ms *TableSchema) PrimaryValue(obj any) any {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	return rVal.FieldByIndex(ms.ss.FieldsAttr[ms.primaryIndex].RefIndex).Interface()
}

func (ms *TableSchema) ValueByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.ss.FieldsAttr[index].RefIndex).Interface()
}

func (ms *TableSchema) AddrByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.ss.FieldsAttr[index].RefIndex).Addr().Interface()
}
