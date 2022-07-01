package mapx

import (
	"fmt"
	"github.com/qinchende/gofast/skill/stringx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"sync"
)

const (
	columnNameTag  = "pms" // 字段名称，对应的tag
	columnNameTag2 = "dbf" // 字段名称，次优先级
	columnOptTag   = "opt" // 字段附加选项
)

// 表结构体Schema, 限制表最多127列（用int8计数）
type GfStruct struct {
	attrs       orm.ModelAttrs  // 实体类型的相关控制属性
	columns     []string        // 安顺序存放的tag列名
	columnsKV   map[string]int8 // pms_name index
	fields      []string        // 安顺序存放的字段名
	fieldsKV    map[string]int8 // field_name index
	fieldsIndex [][]int         // reflect fields index
	fieldsOpts  []*fieldOptions // 字段的属性
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 缓存所有需要反序列化的实体结构的解析数据，防止反复不断的进行反射解析操作。
var cachedSchemas sync.Map

func cacheSetSchema(typ reflect.Type, val *GfStruct) {
	cachedSchemas.Store(typ, val)
}

func cacheGetSchema(typ reflect.Type) *GfStruct {
	if ret, ok := cachedSchemas.Load(typ); ok {
		return ret.(*GfStruct)
	}
	return nil
}

func (ms *GfStruct) ValueByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Interface()
}

func (ms *GfStruct) AddrByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Addr().Interface()
}

func (ms *GfStruct) RefValueByIndex(rVal *reflect.Value, index int8) reflect.Value {
	idxArr := ms.fieldsIndex[index]
	if len(idxArr) == 1 {
		return rVal.Field(idxArr[0])
	}
	tmpVal := *rVal
	for _, x := range idxArr {
		tmpVal = tmpVal.Field(x)
	}
	return tmpVal
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func SchemaNoCache(obj any) *GfStruct {
	return structSchema(reflect.TypeOf(obj))
}

func Schema(obj any) *GfStruct {
	return fetchSchemaCache(reflect.TypeOf(obj))
}

func SchemaOfType(rTyp reflect.Type) *GfStruct {
	return fetchSchemaCache(rTyp)
}

// 提取结构体变量的Schema元数据
func fetchSchemaCache(rTyp reflect.Type) *GfStruct {
	for rTyp.Kind() == reflect.Ptr {
		rTyp = rTyp.Elem()
	}
	// 看类型，缓存有就直接用，否则计算一次并缓存
	mSchema := cacheGetSchema(rTyp)
	if mSchema == nil {
		mSchema = structSchema(rTyp)
		cacheSetSchema(rTyp, mSchema)
	}
	return mSchema
}

func structSchema(rTyp reflect.Type) *GfStruct {
	rootIdx := make([]int, 0)
	fColumns, fFields, fIndexes, fOptions := structFields(rTyp, rootIdx)

	// 获取 Model的所有控制属性
	rTypVal := reflect.ValueOf(reflect.New(rTyp).Interface())
	mdAttrs := new(orm.ModelAttrs)
	attrsFunc := rTypVal.MethodByName("GfAttrs")
	if !attrsFunc.IsZero() {
		mdAttrs = attrsFunc.Call(nil)[0].Interface().(*orm.ModelAttrs)
	}
	if mdAttrs == nil {
		mdAttrs = &orm.ModelAttrs{}
	}
	//mdAttrs.cacheKeyFmt = "Gf#Line#%s#" + mdAttrs.TableName + "#%v"

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
	mSchema := GfStruct{
		attrs:       *mdAttrs,
		columns:     fColumnsNew,
		columnsKV:   make(map[string]int8, len(fColumns)),
		fields:      fFieldsNew,
		fieldsKV:    make(map[string]int8, len(fFields)),
		fieldsIndex: fIndexesNew,
		fieldsOpts:  fOptionsNew,
	}
	// 这样做的目的是收缩Map对象占用的空间
	for idx, name := range fColumns {
		mSchema.columnsKV[name] = int8(idx)
	}
	for idx, name := range fFields {
		mSchema.fieldsKV[name] = int8(idx)
	}
	return &mSchema
}

// 反射提取结构体的字段（支持嵌套递归）
func structFields(rTyp reflect.Type, parentIdx []int) ([]string, []string, [][]int, []*fieldOptions) {
	if rTyp.Kind() != reflect.Struct {
		panic(fmt.Errorf("%T is not like struct", rTyp))
	}

	fColumns := make([]string, 0)
	fFields := make([]string, 0)
	fIndexes := make([][]int, 0)
	fOptions := make([]*fieldOptions, 0)

	for i := 0; i < rTyp.NumField(); i++ {
		fi := rTyp.Field(i)

		// 结构体，需要递归提取其中的字段
		fiType := fi.Type
		if fiType.Kind() == reflect.Struct && fiType.String() != "time.Time" {
			newPIdx := make([]int, 0)
			newPIdx = append(newPIdx, parentIdx...)
			newPIdx = append(newPIdx, i)

			c, f, x, z := structFields(fiType, newPIdx)
			fColumns = append(fColumns, c...)
			fFields = append(fFields, f...)
			fIndexes = append(fIndexes, x...)
			fOptions = append(fOptions, z...)
			continue
		}

		// 1. 查找tag，确定列名称
		col := fi.Tag.Get(columnNameTag)
		if col == "" {
			col = fi.Tag.Get(columnNameTag2)
		}
		if col == "" {
			col = stringx.Camel2Snake(fi.Name)
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
		optStr := fi.Tag.Get(columnOptTag)
		fOpt, err := doParseKeyAndOptions(&fi, optStr)
		errPanic(err)
		fOptions = append(fOptions, fOpt)
	}
	return fColumns, fFields, fIndexes, fOptions
}
