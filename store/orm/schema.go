package orm

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/hashx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"strings"
	"sync"
	"time"
)

func Schema(obj any) *TableSchema {
	return fetchSchema(reflect.TypeOf(obj))
}

func SchemaOfType(rTyp reflect.Type) *TableSchema {
	return fetchSchema(rTyp)
}

// 结构体中属性的数据库字段名称合集
func SchemaValues(obj any) (*TableSchema, []any) {
	ms := Schema(obj)

	var vIndex int8 = 0 // 反射取值索引
	values := make([]any, len(ms.ss.Columns))
	structValues(&values, &vIndex, obj)

	return ms, values
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func CheckDestType(rTyp reflect.Type) reflect.Type {
	eTyp := rTyp.Elem()
	if eTyp.Kind() == reflect.Slice {
		rTyp = eTyp.Elem()
	}
	for rTyp.Kind() == reflect.Pointer {
		rTyp = rTyp.Elem()
	}
	// @@@@@@@@@@ 开始拆解结构体并缓存
	if rTyp.Kind() != reflect.Struct {
		//// 如果是 KV map 类型的。统一
		//if rTyp.Name() == "KV" {
		//	ts = &TableSchema{}
		//	cacheSetSchema(rTyp, ts)
		//	return ts
		//}
		cst.PanicString(fmt.Sprintf("Target object must be structs; but got %T", rTyp))
	}
	return rTyp
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 反射提取结构体的值（支持内联递归）
func structValues(values *[]any, nextIndex *int8, obj any) {
	rVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < rVal.NumField(); i++ {
		va := rVal.Field(i)
		vaI := va.Interface()

		if va.Kind() == reflect.Struct {
			if _, ok := vaI.(time.Time); !ok {
				structValues(values, nextIndex, vaI)
				continue
			}
		}
		(*values)[*nextIndex] = vaI
		*nextIndex++
	}
}

// 提取结构体变量的ORM Schema元数据
func fetchSchema(rTyp reflect.Type) *TableSchema {
	//eTyp := rTyp.Elem()
	//if eTyp.Kind() == reflect.Slice {
	//	rTyp = eTyp.Elem()
	//}
	//for rTyp.Kind() == reflect.Pointer {
	//	rTyp = rTyp.Elem()
	//}
	//// @@@@@@@@@@ 开始拆解结构体并缓存
	//if rTyp.Kind() != reflect.Struct {
	//	//// 如果是 KV map 类型的。统一
	//	//if rTyp.Name() == "KV" {
	//	//	ts = &TableSchema{}
	//	//	cacheSetSchema(rTyp, ts)
	//	//	return ts
	//	//}
	//	cst.PanicString(fmt.Sprintf("Target object must be structs; but got %T", rTyp))
	//}

	ts := cacheGetSchema(rTyp) // 看类型，缓存有就直接用，否则计算一次并缓存
	if ts != nil {
		return ts
	}

	// 如果是 Struct 类型
	ss := dts.SchemaAsDBByType(rTyp)
	// 构造ORM Model元数据
	ts = &TableSchema{
		ss:           *ss, // NOTE：这里是一个赋值操作，而不是指针引用
		autoIndex:    -1,
		primaryIndex: -1,
		updatedIndex: -1,
	}

	// 数据库表的特殊字段索引 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	// auto_inc, primary, updated
	for idx := 0; idx < len(ss.FieldsAttr); idx++ {
		dbc := ss.FieldsAttr[idx].RefField.Tag.Get(dbConfigTag)
		if dbc == "" {
			continue
		}

		// 查找 auto
		if ts.autoIndex == -1 {
			if strings.HasSuffix(dbc, dbAutoIncKeyFlag) {
				ts.autoIndex = int8(idx)
			}
		}
		// 查找 primary
		if ts.primaryIndex == -1 {
			if strings.HasSuffix(dbc, dbPrimaryKeyFlag) {
				ts.primaryIndex = int8(idx)
			}
		}
		// 查找 updated
		if ts.updatedIndex == -1 {
			if strings.HasSuffix(dbc, dbUpdatedKeyFlag) {
				ts.updatedIndex = int8(idx)
			}
		}
	}

	for idx, f := range ss.Fields {
		if ts.autoIndex == -1 {
			if f == dbDefAutoIncKeyName {
				ts.autoIndex = int8(idx)
			}
		}
		if ts.primaryIndex == -1 {
			if f == dbDefPrimaryKeyName {
				ts.primaryIndex = int8(idx)
			}
		}
		if ts.updatedIndex == -1 {
			if f == dbDefUpdatedKeyName {
				ts.updatedIndex = int8(idx)
			}
		}
	}
	if ts.primaryIndex == -1 {
		cst.PanicString(fmt.Sprintf("%T, model has no primary key", rTyp)) // 不能没有主键
	}

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	// 获取 Model的所有控制属性
	rTypVal := reflect.ValueOf(reflect.New(rTyp).Interface())
	attrsFunc := rTypVal.MethodByName("GfAttrs")
	var mdAttrs *TableAttrs
	if attrsFunc.IsValid() {
		vls := []reflect.Value{rTypVal}
		mdAttrs = attrsFunc.Call(vls)[0].Interface().(*TableAttrs)
	}
	if mdAttrs == nil {
		mdAttrs = &TableAttrs{}
	}
	if mdAttrs.TableName == "" {
		mdAttrs.TableName = lang.Camel2Snake(rTyp.Name())
	}

	// Important Note:
	// 表字段的hash值，决定了数据存贮对应的字段以及顺序。
	// 这个特性一定程度能解决，表结构在重构过程中字段发生变化的问题，此表缓存的数据也将失效。
	mdAttrs.columnsHash = hashx.Sum64(lang.STB(strings.Join(ss.Columns, ",")))
	hashStr := lang.ToString(mdAttrs.columnsHash)
	priKeyName := ss.FieldsAttr[ts.primaryIndex].RefField.Name
	mdAttrs.cacheKeyFmt = "Gf#Line#%v#" + mdAttrs.TableName + "#" + hashStr + "#" + priKeyName + "#%v"
	ts.tAttrs = *mdAttrs

	cacheSetSchema(rTyp, ts)
	return ts
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 缓存数据表的Schema
var cachedTableSchemas sync.Map

func cacheSetSchema(typ reflect.Type, val *TableSchema) {
	cachedTableSchemas.Store(typ, val)
}

func cacheGetSchema(typ reflect.Type) *TableSchema {
	if ret, ok := cachedTableSchemas.Load(typ); ok {
		return ret.(*TableSchema)
	}
	return nil
}
