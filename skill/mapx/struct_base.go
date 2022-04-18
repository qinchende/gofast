package mapx

import (
	"github.com/qinchende/gofast/skill/stringx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	columnNameTag  = "pms" // 字段名称，对应的tag
	columnNameTag2 = "cnf" // 字段名称，次优先级
	columnNameTag3 = "dbf" // 字段名称，3优先级
)

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

// 表结构体Schema, 限制表最多127列（用int8计数）
type GfStruct struct {
	attrs       orm.ModelAttrs  // 实体类型的相关控制属性
	columnsKV   map[string]int8 // pms_name index
	fieldsKV    map[string]int8 // field_name index
	fieldsIndex [][]int         // reflect fields index
	fieldsOpts  []*fieldOptions // 字段的属性
}

func (ms *GfStruct) ColumnsKV() map[string]int8 {
	return ms.columnsKV
}

func (ms *GfStruct) FieldsKV() map[string]int8 {
	return ms.fieldsKV
}

//func (ms *GfStruct) Columns() []string {
//	return ms.columns
//}
//
func (ms *GfStruct) ValueByIndex(rVal *reflect.Value, index int8) interface{} {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Interface()
}

func (ms *GfStruct) AddrByIndex(rVal *reflect.Value, index int8) interface{} {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Addr().Interface()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Schema(obj interface{}) *GfStruct {
	return fetchSchema(reflect.TypeOf(obj))
}

func SchemaOfType(rTyp reflect.Type) *GfStruct {
	return fetchSchema(rTyp)
}

// 提取结构体变量的Schema元数据
func fetchSchema(rTyp reflect.Type) *GfStruct {
	for rTyp.Kind() == reflect.Ptr {
		rTyp = rTyp.Elem()
	}

	mSchema := cacheGetSchema(rTyp) // 看类型，缓存有就直接用，否则计算一次并缓存
	if mSchema == nil {
		rootIdx := make([]int, 0)
		fJson, fFields, fIndexes := structFields(rTyp, rootIdx)

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
		//fJsonNew := make([]string, len(fJson))
		//copy(fJsonNew, fJson)
		fIndexesNew := make([][]int, len(fIndexes))
		copy(fIndexesNew, fIndexes)
		// 构造ORM Model元数据
		mSchema = &GfStruct{
			attrs:       *mdAttrs,
			columnsKV:   make(map[string]int8, len(fJson)),
			fieldsKV:    make(map[string]int8, len(fFields)),
			fieldsIndex: fIndexesNew,
		}
		// 这样做的目的是收缩Map对象占用的空间
		for idx, name := range fJson {
			mSchema.columnsKV[name] = int8(idx)
		}
		for idx, name := range fFields {
			mSchema.fieldsKV[name] = int8(idx)
		}
		cacheSetSchema(rTyp, mSchema)
	}

	return mSchema
}

// 反射提取结构体的字段（支持嵌套递归）
func structFields(rTyp reflect.Type, parentIdx []int) ([]string, []string, [][]int) {
	fJson := make([]string, 0)
	fFields := make([]string, 0)
	fIndexes := make([][]int, 0)

	for i := 0; i < rTyp.NumField(); i++ {
		fi := rTyp.Field(i)

		// 通过值类型来确定后面
		fdType := fi.Type
		if fdType.Kind() == reflect.Struct {
			vaI := reflect.New(fdType).Interface()
			if _, ok := vaI.(*time.Time); !ok {
				newPIdx := make([]int, 0)
				newPIdx = append(newPIdx, parentIdx...)
				newPIdx = append(newPIdx, i)

				c, f, x := structFields(fdType, newPIdx)
				fJson = append(fJson, c...)
				fFields = append(fFields, f...)
				fIndexes = append(fIndexes, x...)
				continue
			}
		}

		// 1. 查找tag，确定列名称
		dbf, _ := tagHead(fi.Tag.Get(columnNameTag), ",")
		if dbf == "" {
			dbf, _ = tagHead(fi.Tag.Get(columnNameTag2), ",")
		}
		if dbf == "" {
			dbf, _ = tagHead(fi.Tag.Get(columnNameTag3), ",")
		}
		if dbf == "" {
			dbf = stringx.Camel2Snake(fi.Name)
		}
		fJson = append(fJson, dbf)

		// 2. 确定结构体字段名称
		fFields = append(fFields, fi.Name)

		// 3. index
		cIdx := make([]int, 0)
		cIdx = append(cIdx, parentIdx...)
		cIdx = append(cIdx, i)
		fIndexes = append(fIndexes, cIdx)
	}
	return fJson, fFields, fIndexes
}

func tagHead(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}
