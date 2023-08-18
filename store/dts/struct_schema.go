// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/validx"
	"math"
	"reflect"
	"unsafe"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 表结构体Schema, 限制表最多127列（用int8计数）
type (
	StructSchema struct {
		Attrs      structAttrs // 结构体元数据
		FieldsAttr []fieldAttr // 字段元数据

		Columns []string    // 按顺序存放的tag列名
		Fields  []string    // 按顺序存放的字段名
		cTips   stringsTips // pms_name index
		fTips   stringsTips // field_name index
	}

	// 基本信息
	structAttrs struct {
		Type    reflect.Type
		MemSize int
	}

	// 所有字段按照长度从小到大排序，用于快速索引
	stringsTips struct {
		items  []string
		idxes  []uint8
		lenOff []uint8
	}

	// 给字段绑定值
	valueBinder func(p unsafe.Pointer, v any)

	// 方便字段数据处理
	fieldAttr struct {
		rIndex []int                // 字段定位（反射用到）
		Valid  *validx.ValidOptions // 验证
		sField *reflect.StructField // 原始值，方便后期自定义验证特殊Tag

		Type      reflect.Type // 字段最终的类型，剥开指针(Pointer)之后的类型
		Kind      reflect.Kind // 字段最终类型的Kind类型
		Offset    uintptr      // 字段在结构体中的地址偏移量
		PtrLevel  uint8        // 字段指针层级
		IsMixType bool         // 是否为混合数据类型（非基础数据类型之外的类型，比如Struct,Map,Array,Slice）
		KVBinder  valueBinder  // 绑定函数
	}

	fieldOptions struct {
		valid  *validx.ValidOptions // 验证
		sField *reflect.StructField // 原始值，方便后期自定义验证特殊Tag
	}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 提取结构体变量的Schema元数据
func fetchSchemaCache(rTyp reflect.Type, opts *BindOptions) *StructSchema {
	for rTyp.Kind() == reflect.Pointer {
		rTyp = rTyp.Elem()
	}
	// 看类型，缓存有就直接用，否则计算一次并缓存
	ss := cacheGetSchema(rTyp)
	if ss == nil {
		ss = buildStructSchema(rTyp, opts)
		cacheSetSchema(rTyp, ss)
	}
	return ss
}

func buildStructSchema(rTyp reflect.Type, opts *BindOptions) *StructSchema {
	rootIdx := make([]int, 0)
	fColumns, fFields, fIndexes, fOptions := structFields(rTyp, rootIdx, opts)

	if len(fColumns) <= 0 {
		panic("Struct not contain any fields")
	}
	if len(fColumns) > math.MaxUint8 {
		panic("Struct field items large the 256")
	}

	// 构造ORM Model元数据
	ss := StructSchema{}
	ss.Attrs.Type = rTyp
	ss.Attrs.MemSize = int(rTyp.Size())

	// 收缩切片占用的空间，因为原slice可能有多余的cap
	ss.Columns = make([]string, len(fColumns))
	copy(ss.Columns, fColumns)
	ss.Fields = make([]string, len(fFields))
	copy(ss.Fields, fFields)

	// 缓存字段的属性和实用方法 ++++++++++
	ss.FieldsAttr = make([]fieldAttr, len(fOptions))
	for i := range fOptions {
		fa := &ss.FieldsAttr[i]

		fa.rIndex = fIndexes[i]
		fa.Valid = fOptions[i].valid
		fa.sField = fOptions[i].sField

		fa.Offset = fOptions[i].sField.Offset
		fa.Type = fOptions[i].sField.Type
		for fa.Type.Kind() == reflect.Pointer {
			fa.PtrLevel++
			fa.Type = fa.Type.Elem()
		}
		fa.Kind = fa.Type.Kind()
		switch fa.Kind {
		case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
			fa.IsMixType = true
		}

		// +++ set KVBinder function and other attributes
		// Note: if opts.model == AsConfig
		// 不同的模式，解析函数可能是不一样的，当前仅支持 AsConfig 模式
		switch fa.Kind {
		case reflect.Int:
			fa.KVBinder = setInt
		case reflect.Int8:
			fa.KVBinder = setInt8
		case reflect.Int16:
			fa.KVBinder = setInt16
		case reflect.Int32:
			fa.KVBinder = setInt32
		case reflect.Int64:
			fa.KVBinder = setInt64

		case reflect.Uint:
			fa.KVBinder = setUint
		case reflect.Uint8:
			fa.KVBinder = setUint8
		case reflect.Uint16:
			fa.KVBinder = setUint16
		case reflect.Uint32:
			fa.KVBinder = setUint32
		case reflect.Uint64:
			fa.KVBinder = setUint64

		case reflect.Float32:
			fa.KVBinder = setFloat32
		case reflect.Float64:
			fa.KVBinder = setFloat64

		case reflect.String:
			fa.KVBinder = setString
		case reflect.Bool:
			fa.KVBinder = setBool
		case reflect.Interface:
			fa.KVBinder = setAny

		case reflect.Pointer:

		case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		}
	}

	// cTips
	// 方便检索字符串项，这里做一些数据冗余的优化处理
	ss.cTips.items = make([]string, len(fColumns))
	ss.cTips.idxes = make([]uint8, len(fColumns))

	copy(ss.cTips.items, ss.Columns)
	lang.SortByLen(ss.cTips.items)
	lastLen := len(ss.cTips.items[len(ss.cTips.items)-1]) // 最长string长度（最后一个就是最长的）
	if lastLen > math.MaxUint8 {
		panic("Struct has field large the 256 chars")
	}
	ss.cTips.lenOff = make([]uint8, lastLen+1)
	lastLen = 0

	for idx, item := range ss.cTips.items {
		if lastLen != len(item) {
			ss.cTips.lenOff[len(item)] = uint8(idx)
			lastLen = len(item)
		}
		for sIdx := range ss.Columns {
			if item == ss.Columns[sIdx] {
				ss.cTips.idxes[idx] = uint8(sIdx)
				break
			}
		}
	}

	// fTips
	// +++++++++++++++
	ss.fTips.items = make([]string, len(fFields))
	ss.fTips.idxes = make([]uint8, len(fFields))

	copy(ss.fTips.items, ss.Fields)
	lang.SortByLen(ss.fTips.items)
	lastLen = len(ss.fTips.items[len(ss.fTips.items)-1])
	if lastLen > math.MaxUint8 {
		panic("Struct has field large the 256 chars")
	}
	ss.fTips.lenOff = make([]uint8, lastLen+1)
	lastLen = 0

	for idx, item := range ss.fTips.items {
		if lastLen != len(item) {
			ss.fTips.lenOff[len(item)] = uint8(idx)
			lastLen = len(item)
		}
		for sIdx := range ss.Fields {
			if item == ss.Fields[sIdx] {
				ss.fTips.idxes[idx] = uint8(sIdx)
				break
			}
		}
	}

	return &ss
}

// 反射提取结构体的字段（支持嵌套递归）
func structFields(rTyp reflect.Type, parentIdx []int, opts *BindOptions) ([]string, []string, [][]int, []fieldOptions) {
	if rTyp.Kind() != reflect.Struct {
		cst.PanicString(fmt.Sprintf("%T is not like struct", rTyp))
	}

	fColumns := make([]string, 0)
	fFields := make([]string, 0)
	fIndexes := make([][]int, 0)
	fOptions := make([]fieldOptions, 0)

	for i := 0; i < rTyp.NumField(); i++ {
		fi := rTyp.Field(i)

		// 结构体，需要递归提取其中的字段
		fiType := fi.Type
		if fi.Anonymous && fiType.Kind() == reflect.Struct && fiType.String() != "time.Time" {
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
