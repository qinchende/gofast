package orm

import (
	"fmt"
	"github.com/qinchende/gofast/aid/hashx"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/dts"
	lang2 "github.com/qinchende/gofast/core/lang"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Orm中的TableSchema
// 只解决单条记录绑定到单个object，或者多条记录绑定到 list object 类型的变量。
func Schema(obj any) *TableSchema {
	return SchemaByType(reflect.TypeOf(obj))
}

func SchemaNil(obj any) *TableSchema {
	if obj == nil {
		return nilTableSchema
	}
	return SchemaByType(reflect.TypeOf(obj))
}

func SchemaByType(typ reflect.Type) *TableSchema {
	// typ 可能是：*Struct，Struct，*[]Struct，[]Struct, *[]*Struct, []*Struct, *[]cst.KV 等类型
	kd := typ.Kind()
	if kd == reflect.Pointer {
		typ = typ.Elem() // 只剥离一层 pointer
		kd = typ.Kind()
	}
	if kd == reflect.Slice {
		typ = typ.Elem() // 只剥离一层 slice
		kd = typ.Kind()
	}
	if kd == reflect.Pointer {
		typ = typ.Elem() // 再剥离一层 pointer
		kd = typ.Kind()
	}
	// 此时必须是 struct
	if kd != reflect.Struct {
		// 如果是 cst.KV 类型也默认支持，意味着此时想到得到JSON数据
		// 非 struct 类型中，只支持 cst.KV
		if typ == cst.TypeCstKV {
			ts := kvTableSchema
			cacheSetSchema(typ, ts)
			return ts
		}
		cst.PanicString(fmt.Sprintf("Target must be struct, but got %T", typ.Kind()))
	}

	// @@@ 开始拆解结构体并缓存
	return fetchSchema(typ)
}

// 对传入的类型检查已经做过，直接取TableSchema即可
func SchemaByTypeDirect(typ reflect.Type) *TableSchema {
	return fetchSchema(typ)
}

// 结构体中属性的数据库字段名称合集
func SchemaValues(obj any) (*TableSchema, []any) {
	ms := Schema(obj)

	var vIndex int8 = 0 // 反射取值索引
	values := make([]any, len(ms.SS.Columns))
	structValues(&values, &vIndex, obj)

	return ms, values
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
func fetchSchema(typ reflect.Type) *TableSchema {
	ts := cacheGetSchema(typ) // 看类型，缓存有就直接用，否则计算一次并缓存
	if ts != nil {
		return ts
	}

	// 如果是 Struct 类型
	ss := dts.SchemaAsDBByType(typ)
	// 构造ORM Model元数据
	ts = &TableSchema{
		SS:           *ss, // NOTE：这里是一个赋值操作，而不是指针引用
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
		cst.PanicString(fmt.Sprintf("%T, model has no primary key", typ)) // 不能没有主键
	}

	// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	// 获取 Model的所有控制属性
	rTypVal := reflect.ValueOf(reflect.New(typ).Interface())
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
		mdAttrs.TableName = lang2.Camel2Snake(typ.Name())
	}

	// Important Note:
	// 表字段的hash值，决定了数据存贮对应的字段以及顺序。
	// 这个特性一定程度能解决，表结构在重构过程中字段发生变化的问题，此表缓存的数据也将失效。
	mdAttrs.columnsHash = hashx.Sum64(lang2.STB(strings.Join(ss.Columns, ",")))
	hashStr := lang2.ToString(mdAttrs.columnsHash)
	priKeyName := ss.FieldsAttr[ts.primaryIndex].RefField.Name
	// 行记录缓存 Key format
	mdAttrs.cacheKeyFmt = "Gf#Line#%v#" + mdAttrs.TableName + "#" + hashStr + "#" + priKeyName + "#%v"
	// 默认 行记录 缓存 3600 second
	if mdAttrs.CacheAll == true && mdAttrs.ExpireS <= 0 {
		mdAttrs.ExpireS = 3600
	}

	ts.tAttrs = *mdAttrs

	cacheSetSchema(typ, ts)
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
