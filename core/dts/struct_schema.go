// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

// Core Package
// @@@ Data transfer system @@@ -> dts

import (
	"fmt"
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/aid/validx"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/rt"
	"math"
	"reflect"
	"sync"
	"unsafe"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 表结构体Schema, 限制表最多127列（用int8计数）
type (
	StructSchema struct {
		structAttrs             // 结构体元数据
		FieldsAttr  []fieldAttr // 字段元数据
		Fields      []string    // 按顺序存放的字段名
		Columns     []string    // 按顺序存放的tag列名

		fTips stringsTips // field_name index
		cTips stringsTips // pms_name index
	}

	// 基本信息
	structAttrs struct {
		Type        reflect.Type
		MemSize     int
		HasPtrField bool // 是否保护有指针类型的字段
	}

	// 所有字段按照长度从小到大排序，用于快速索引
	stringsTips struct {
		items  []string
		idxes  []uint8
		lenOff []uint8
	}

	// 给字段绑定值
	kvBinderFunc func(fPtr unsafe.Pointer, v any)
	SqlValueFunc func(sPtr unsafe.Pointer) any

	// 方便字段数据处理
	fieldAttr struct {
		RefIndex []int                // 字段定位（反射用到）
		RefField *reflect.StructField // 原始值，方便后期自定义验证特殊Tag
		Valid    *validx.ValidOptions // 验证

		Type      reflect.Type   // 字段最终的类型，剥开指针(Pointer)之后的类型
		TypeAbi   unsafe.Pointer // 模拟*runtime.abiType
		Kind      reflect.Kind   // 字段最终类型的Kind类型，
		Offset    uintptr        // 字段在结构体中的地址偏移量
		PtrLevel  uint8          // 字段指针层级
		IsMixType bool           // 是否为混合数据类型（非基础数据类型之外的类型，比如Struct,Map,Array,Slice）

		KVBinder kvBinderFunc // 绑定函数
		SqlValue SqlValueFunc
	}

	fieldOptions struct {
		valid  *validx.ValidOptions // 验证
		sField *reflect.StructField // 原始值，方便后期自定义验证特殊Tag
	}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 提取结构体变量的Schema元数据
func fetchSchemaCache(typ reflect.Type, opts *BindOptions) *StructSchema {
	// 看类型，缓存有就直接用，否则计算一次并缓存
	ss := cacheGetSchema(typ)
	if ss == nil {
		ss = buildStructSchema(typ, opts)
		cacheSetSchema(typ, ss)
	}
	return ss
}

func buildStructSchema(typ reflect.Type, opts *BindOptions) *StructSchema {
	rootIdx := make([]int, 0)
	fColumns, fFields, fIndexes, fOptions := structFields(typ, rootIdx, opts)

	if len(fColumns) <= 0 {
		cst.PanicString("Struct not contain any fields")
	}
	if len(fColumns) > math.MaxUint8 {
		cst.PanicString("Struct field items large the 256")
	}

	// 构造ORM Model元数据
	ss := StructSchema{}
	ss.Type = typ
	ss.MemSize = int(typ.Size())

	// 收缩切片占用的空间，因为原slice可能有多余的cap
	ss.Columns = make([]string, len(fColumns))
	copy(ss.Columns, fColumns)
	ss.Fields = make([]string, len(fFields))
	copy(ss.Fields, fFields)

	// 缓存字段的属性和实用方法 ++++++++++
	ss.FieldsAttr = make([]fieldAttr, len(fOptions))
	for i := range fOptions {
		fa := &ss.FieldsAttr[i]

		fa.RefIndex = fIndexes[i]
		fa.Valid = fOptions[i].valid
		fa.RefField = fOptions[i].sField

		fa.Offset = fOptions[i].sField.Offset
		fa.Type = fOptions[i].sField.Type
		// 如果struct字段中有指针类型，需要穿透到非指针类型，同时记录指针层级数
		for fa.Type.Kind() == reflect.Pointer {
			fa.PtrLevel++
			fa.Type = fa.Type.Elem()
			ss.HasPtrField = true // 结构体保护指针字段
		}
		fa.Kind = fa.Type.Kind()
		fa.TypeAbi = (*rt.AFace)(unsafe.Pointer(&fa.Type)).DataPtr
		switch fa.Kind {
		case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
			fa.IsMixType = true
		default:
		}

		// +++ set KVBinder function and other attributes
		// Note: if opts.model == AsConfig
		// 不同的模式，解析函数可能是不一样的，当前仅支持 AsConfig 模式
		switch fa.Kind {
		case reflect.Int:
			fa.KVBinder = setInt
			fa.SqlValue = fa.intValue
		case reflect.Int8:
			fa.KVBinder = setInt8
			fa.SqlValue = fa.int8Value
		case reflect.Int16:
			fa.KVBinder = setInt16
			fa.SqlValue = fa.int16Value
		case reflect.Int32:
			fa.KVBinder = setInt32
			fa.SqlValue = fa.int32Value
		case reflect.Int64:
			if fa.Type == cst.TypeDuration {
				fa.KVBinder = setDuration
				fa.SqlValue = fa.durationValue
			} else {
				fa.KVBinder = setInt64
				fa.SqlValue = fa.int64Value
			}

		case reflect.Uint:
			fa.KVBinder = setUint
			fa.SqlValue = fa.uintValue
		case reflect.Uint8:
			fa.KVBinder = setUint8
			fa.SqlValue = fa.uint8Value
		case reflect.Uint16:
			fa.KVBinder = setUint16
			fa.SqlValue = fa.uint16Value
		case reflect.Uint32:
			fa.KVBinder = setUint32
			fa.SqlValue = fa.uint32Value
		case reflect.Uint64:
			fa.KVBinder = setUint64
			fa.SqlValue = fa.uint64Value

		case reflect.Float32:
			fa.KVBinder = setFloat32
			fa.SqlValue = fa.float32Value
		case reflect.Float64:
			fa.KVBinder = setFloat64
			fa.SqlValue = fa.float64Value

		case reflect.String:
			fa.KVBinder = setString
			fa.SqlValue = fa.stringValue
		case reflect.Bool:
			fa.KVBinder = setBool
			fa.SqlValue = fa.boolValue
		case reflect.Interface:
			fa.KVBinder = setAny
			fa.SqlValue = fa.anyValue

		case reflect.Struct:
			if fa.Type == cst.TypeTime {
				fa.KVBinder = setTime
				fa.SqlValue = fa.timeValue
			} else {

			}
		case reflect.Pointer:
		case reflect.Map, reflect.Array, reflect.Slice:
		default:
		}
	}

	// cTips
	// 方便检索字符串项，这里做一些数据冗余的优化处理
	ss.cTips.items = make([]string, len(fColumns))
	ss.cTips.idxes = make([]uint8, len(fColumns))

	copy(ss.cTips.items, ss.Columns)
	lang.SortByLen(ss.cTips.items)
	lastLen := len(ss.cTips.items[len(ss.cTips.items)-1]) // 按长度排序后，最后一个字符串就是最长的
	if lastLen > math.MaxUint8 {
		cst.PanicString("Struct has field large the 256 chars")
	}
	ss.cTips.lenOff = make([]uint8, lastLen+1) // 将来字符串长度直接作为下标索引，所以这里必须+1
	lastLen = 0

	for idx, item := range ss.cTips.items {
		if lastLen != len(item) {
			ss.cTips.lenOff[len(item)] = uint8(idx)
			lastLen = len(item)
		}
		for fIdx := range ss.Columns {
			if item == ss.Columns[fIdx] {
				// 排序后的字段对应 struct 中原始 field 的索引
				ss.cTips.idxes[idx] = uint8(fIdx)
				break
			}
		}
	}

	// fTips
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++
	ss.fTips.items = make([]string, len(fFields))
	ss.fTips.idxes = make([]uint8, len(fFields))

	copy(ss.fTips.items, ss.Fields)
	lang.SortByLen(ss.fTips.items)
	lastLen = len(ss.fTips.items[len(ss.fTips.items)-1])
	if lastLen > math.MaxUint8 {
		cst.PanicString("Struct has field large the 256 chars")
	}
	ss.fTips.lenOff = make([]uint8, lastLen+1)
	lastLen = 0

	for idx, item := range ss.fTips.items {
		if lastLen != len(item) {
			ss.fTips.lenOff[len(item)] = uint8(idx)
			lastLen = len(item)
		}
		for fIdx := range ss.Fields {
			if item == ss.Fields[fIdx] {
				ss.fTips.idxes[idx] = uint8(fIdx)
				break
			}
		}
	}

	return &ss
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 反射提取结构体的字段（支持嵌套递归）
func structFields(typ reflect.Type, parentIdx []int, opts *BindOptions) ([]string, []string, [][]int, []fieldOptions) {
	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%T is not like struct", typ))
	}

	// 需要解析 struct 的几种数据 ++++
	fLen := typ.NumField()
	fColumns := make([]string, 0, fLen)
	fFields := make([]string, 0, fLen)
	fIndexes := make([][]int, 0, fLen)
	fOptions := make([]fieldOptions, 0, fLen)

	for i := 0; i < fLen; i++ {
		fi := typ.Field(i)

		// 排除掉非导出字段
		if !fi.IsExported() {
			continue
		}

		// TODO: 匿名的结构体（非time.Time），需要递归提取其中的字段
		fiType := fi.Type
		if fi.Anonymous && fiType.Kind() == reflect.Struct && fiType != cst.TypeTime {
			newPIdx := make([]int, 0)
			newPIdx = append(newPIdx, parentIdx...)
			newPIdx = append(newPIdx, i)

			c, f, x, z := structFields(fiType, newPIdx, opts)
			fColumns = append(fColumns, c...)
			fFields = append(fFields, f...)
			fIndexes = append(fIndexes, x...)
			fOptions = append(fOptions, z...)
			continue
		}

		// 1. 查找tag，确定列名称
		col := fi.Tag.Get(opts.FieldTag)
		// 如果指定的tag找不到，就试这去找通用tag
		if col == "" && opts.FieldTag != cst.FieldTag {
			col = fi.Tag.Get(cst.FieldTag)
		}
		if col == "" {
			col = lang.Camel2Snake(fi.Name)
		}
		fColumns = append(fColumns, col)

		// 2. 确定结构体字段名称
		fFields = append(fFields, fi.Name)

		// 3. index
		cIdx := make([]int, 0)
		cIdx = append(cIdx, parentIdx...)
		cIdx = append(cIdx, i)
		fIndexes = append(fIndexes, cIdx)

		// 4. options
		optStr := fi.Tag.Get(opts.ValidTag)
		vOpt, err := validx.ParseOptions(&fi, optStr)
		cst.PanicIfErr(err) // 解析不对，直接抛异常
		fOptions = append(fOptions, fieldOptions{valid: vOpt, sField: &fi})
	}
	return fColumns, fFields, fIndexes, fOptions
}

// cache
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
