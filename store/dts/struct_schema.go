// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/validx"
	"reflect"
	"sync"
)

type kvMix struct {
	items []string
	idxes []int
}

// 表结构体Schema, 限制表最多127列（用int8计数）
type StructSchema struct {
	// attrs       orm.ModelAttrs  // 实体类型的相关控制属性
	columns     []string        // 按顺序存放的tag列名
	columnsKV   kvMix           // pms_name index
	fields      []string        // 按顺序存放的字段名
	fieldsKV    kvMix           // field_name index
	fieldsIndex [][]int         // reflect fields index
	fieldsOpts  []*fieldOptions // 字段的属性
	addrOff     []int
}

//type StructAttrs struct {
//	TableName string // 数据库表名称
//	CacheAll  bool   // 是否缓存所有记录
//	ExpireS   uint32 // 过期时间（秒）默认7天
//
//	// 内部状态标记
//	hashNumber  uint64 // 本结构体的哈希值
//	cacheKeyFmt string // 行记录缓存的Key前缀
//}

type fieldOptions struct {
	valid  *validx.ValidOptions // 验证
	sField *reflect.StructField // 原始值，方便后期自定义验证特殊Tag
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

	// 获取 Model的所有控制属性
	//rTypVal := reflect.ValueOf(reflect.New(rTyp).Interface())
	//mdAttrs := new(orm.ModelAttrs)
	//attrsFunc := rTypVal.MethodByName("GfAttrs")
	//if attrsFunc.IsValid() {
	//	vls := []reflect.Value{rTypVal}
	//	mdAttrs = attrsFunc.Call(vls)[0].Interface().(*orm.ModelAttrs)
	//}
	//if mdAttrs == nil {
	//	mdAttrs = &orm.ModelAttrs{}
	//}

	// 收缩切片
	fColumnsNew := make([]string, len(fColumns))
	copy(fColumnsNew, fColumns)
	fFieldsNew := make([]string, len(fFields))
	copy(fFieldsNew, fFields)
	fIndexesNew := make([][]int, len(fIndexes))
	copy(fIndexesNew, fIndexes)
	fOptionsNew := make([]*fieldOptions, len(fOptions))
	copy(fOptionsNew, fOptions)
	// 构造ORM Model元数据
	ss := StructSchema{
		//attrs:       *mdAttrs,
		columns:     fColumnsNew,
		fields:      fFieldsNew,
		fieldsIndex: fIndexesNew,
		fieldsOpts:  fOptionsNew,
	}
	ss.columnsKV.items = make([]string, len(fColumns))
	ss.columnsKV.idxes = make([]int, len(fColumns))
	ss.fieldsKV.items = make([]string, len(fFields))
	ss.fieldsKV.idxes = make([]int, len(fFields))

	// 这样做的目的是收缩Map对象占用的空间
	copy(ss.columnsKV.items, ss.columns)
	lang.SortByLen(ss.columnsKV.items)
	for idx, name := range ss.columnsKV.items {
		for sIdx, sName := range ss.columns {
			if name == sName {
				ss.columnsKV.idxes[idx] = sIdx
			}
		}
	}
	copy(ss.fieldsKV.items, ss.fields)
	lang.SortByLen(ss.fieldsKV.items)
	for idx, name := range ss.fieldsKV.items {
		for sIdx, sName := range ss.fields {
			if name == sName {
				ss.fieldsKV.idxes[idx] = sIdx
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
		//if col == "" {
		//	col = fi.Tag.Get(cst.FieldTagDB)
		//}
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
