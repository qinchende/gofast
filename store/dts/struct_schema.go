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
	"sync"
)

// 表结构体Schema, 限制表最多127列（用int8计数）
type StructSchema struct {
	columns []string    // 按顺序存放的tag列名
	fields  []string    // 按顺序存放的字段名
	cTips   stringsTips // pms_name index
	fTips   stringsTips // field_name index

	fieldsIndex [][]int         // reflect fields index
	fieldsOpts  []*fieldOptions // 字段的属性
	fieldsKind  []reflect.Kind
	fieldsPtr   []uintptr
	binds       []BindValue
}

type stringsTips struct {
	items  []string
	idxes  []uint8
	lenOff []uint8
}

type fieldOptions struct {
	valid  *validx.ValidOptions // 验证
	sField *reflect.StructField // 原始值，方便后期自定义验证特殊Tag
}

type BindValue interface {
	BindInt()
}

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

	// 收缩切片占用的空间，因为原slice可能有多余的cap
	ss.columns = make([]string, len(fColumns))
	copy(ss.columns, fColumns)
	ss.fields = make([]string, len(fFields))
	copy(ss.fields, fFields)
	ss.fieldsIndex = make([][]int, len(fIndexes))
	copy(ss.fieldsIndex, fIndexes)
	ss.fieldsOpts = make([]*fieldOptions, len(fOptions))
	copy(ss.fieldsOpts, fOptions)

	// 抽取出字段的类型和偏移地址
	ss.fieldsKind = make([]reflect.Kind, len(fOptions))
	ss.fieldsPtr = make([]uintptr, len(fOptions))
	for i := range fOptions {
		ss.fieldsKind[i] = fOptions[i].sField.Type.Kind()
		ss.fieldsPtr[i] = fOptions[i].sField.Offset
	}

	// 方便检索字符串项，这里做一些数据冗余的优化处理
	ss.cTips.items = make([]string, len(fColumns))
	ss.cTips.idxes = make([]uint8, len(fColumns))

	copy(ss.cTips.items, ss.columns)
	lang.SortByLen(ss.cTips.items)
	lastLen := len(ss.cTips.items[len(ss.cTips.items)-1])
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
		for sIdx := range ss.columns {
			if item == ss.columns[sIdx] {
				ss.cTips.idxes[idx] = uint8(sIdx)
				break
			}
		}
	}

	// +++++++++++++++
	ss.fTips.items = make([]string, len(fFields))
	ss.fTips.idxes = make([]uint8, len(fFields))

	copy(ss.fTips.items, ss.fields)
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
		for sIdx := range ss.fields {
			if item == ss.fields[sIdx] {
				ss.fTips.idxes[idx] = uint8(sIdx)
				break
			}
		}
	}

	return &ss
}

// 反射提取结构体的字段（支持嵌套递归）
func structFields(rTyp reflect.Type, parentIdx []int, opts *BindOptions) ([]string, []string, [][]int, []*fieldOptions) {
	if rTyp.Kind() != reflect.Struct {
		cst.PanicString(fmt.Sprintf("%T is not like struct", rTyp))
	}

	fColumns := make([]string, 0)
	fFields := make([]string, 0)
	fIndexes := make([][]int, 0)
	fOptions := make([]*fieldOptions, 0)

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
		fOptions = append(fOptions, &fieldOptions{valid: vOpt, sField: &fi})
	}
	return fColumns, fFields, fIndexes, fOptions
}

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
