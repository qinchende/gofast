package orm

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/hash"
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
	"strings"
	"sync"
	"time"
)

func Schema(obj any) *ModelSchema {
	return fetchSchema(reflect.TypeOf(obj))
}

func SchemaOfType(rTyp reflect.Type) *ModelSchema {
	return fetchSchema(rTyp)
}

// 结构体中属性的数据库字段名称合集
func SchemaValues(obj any) (*ModelSchema, []any) {
	ms := Schema(obj)

	var vIndex int8 = 0 // 反射取值索引
	values := make([]any, len(ms.columns))
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
func fetchSchema(rTyp reflect.Type) *ModelSchema {
	eTyp := rTyp.Elem()
	if eTyp.Kind() == reflect.Slice {
		rTyp = eTyp.Elem()
	}

	for rTyp.Kind() == reflect.Ptr {
		rTyp = rTyp.Elem()
	}

	mSchema := cacheGetSchema(rTyp) // 看类型，缓存有就直接用，否则计算一次并缓存
	if mSchema == nil {
		mPath := rTyp.PkgPath()
		mFullName := rTyp.String()
		mName := rTyp.Name()

		if rTyp.Kind() != reflect.Struct {
			// 如果是 KV map 类型的。统一
			if mName == "KV" {
				mSchema = &ModelSchema{}
				cacheSetSchema(rTyp, mSchema)
				return mSchema
			}
			cst.PanicString(fmt.Sprintf("Target object must be structs; but got %T", rTyp))
		}

		// primary, updated
		mFields := [3]string{}
		rootIdx := make([]int, 0)
		fDB, fStruct, fIndexes := structFields(rTyp, rootIdx, &mFields)
		if mFields[0] == "" {
			mFields[0] = dbDefAutoIncKeyName
		}
		if mFields[1] == "" {
			mFields[1] = dbDefPrimaryKeyName
		}
		if mFields[2] == "" {
			mFields[2] = dbDefUpdatedKeyName
		}

		// 0. 自增的索引位置 ++++++++++
		var autoIndex = -1
		for idx, f := range fStruct {
			if f == mFields[0] {
				autoIndex = idx
				break
			}
		}
		// 1. 主键的索引位置 ++++++++++
		var priIndex = -1
		for idx, f := range fStruct {
			if f == mFields[1] {
				priIndex = idx
				break
			}
		}
		if priIndex == -1 {
			cst.PanicString(fmt.Sprintf("%T, model has no primary key", rTyp)) // 不能没有主键
		}
		// 2. updated 的索引位置
		var updateIndex = -1
		for idx, f := range fStruct {
			if f == mFields[2] {
				updateIndex = idx
				break
			}
		}

		// 获取 Model的所有控制属性
		rTypVal := reflect.ValueOf(reflect.New(rTyp).Interface())
		attrsFunc := rTypVal.MethodByName("GfAttrs")
		var mdAttrs *ModelAttrs
		if attrsFunc.IsValid() {
			vls := []reflect.Value{rTypVal}
			mdAttrs = attrsFunc.Call(vls)[0].Interface().(*ModelAttrs)
		}
		if mdAttrs == nil {
			mdAttrs = &ModelAttrs{}
		}
		if mdAttrs.TableName == "" {
			mdAttrs.TableName = lang.Camel2Snake(rTyp.Name())
		}
		mdAttrs.hashNumber = hash.Hash(lang.StringToBytes(strings.Join(fDB, ",")))
		hashStr := lang.ToString(mdAttrs.hashNumber)
		mdAttrs.cacheKeyFmt = "Gf#Line#%v#" + mdAttrs.TableName + "#" + hashStr + "#" + mFields[1] + "#%v"

		// 收缩切片
		fIndexesNew := make([][]int, len(fIndexes))
		copy(fIndexesNew, fIndexes)
		fDBNew := make([]string, len(fDB))
		copy(fDBNew, fDB)
		// 构造ORM Model元数据
		mSchema = &ModelSchema{
			pkgPath:  mPath,
			fullName: mFullName,
			name:     mName,
			attrs:    *mdAttrs,

			columns:      fDBNew,
			fieldsKV:     make(map[string]int8, len(fStruct)),
			columnsKV:    make(map[string]int8, len(fStruct)),
			fieldsIndex:  fIndexesNew,
			autoIndex:    int8(autoIndex),
			primaryIndex: int8(priIndex),
			updatedIndex: int8(updateIndex),
		}
		for idx, name := range fStruct {
			mSchema.fieldsKV[name] = int8(idx)
		}
		for idx, name := range fDBNew {
			mSchema.columnsKV[name] = int8(idx)
		}
		cacheSetSchema(rTyp, mSchema)
	}

	return mSchema
}

// 反射提取结构体的字段（支持嵌套递归）
func structFields(rTyp reflect.Type, parentIdx []int, mFields *[3]string) ([]string, []string, [][]int) {
	if rTyp.Kind() != reflect.Struct {
		cst.PanicString(fmt.Sprintf("%T is not like struct", rTyp))
	}

	fColumns := make([]string, 0)
	fFields := make([]string, 0)
	fIndexes := make([][]int, 0)

	for i := 0; i < rTyp.NumField(); i++ {
		fi := rTyp.Field(i)

		// NOTE: 结构体，需要递归提取其中的字段
		// 这里有个疑问，应该是匿名结构体才需要开箱取其中的字段
		fiType := fi.Type
		// if fiType.Kind() == reflect.Struct && fiType.String() != "time.Time" {
		if fi.Anonymous && fiType.Kind() == reflect.Struct && fiType.String() != "time.Time" {
			newPIdx := make([]int, 0)
			newPIdx = append(newPIdx, parentIdx...)
			newPIdx = append(newPIdx, i)

			c, f, x := structFields(fiType, newPIdx, mFields)
			fColumns = append(fColumns, c...)
			fFields = append(fFields, f...)
			fIndexes = append(fIndexes, x...)
			continue
		}

		// 1. 查找tag，确定数据库列名称
		dbf := fi.Tag.Get(cst.FieldTagDB)
		if dbf == "" {
			dbf = fi.Tag.Get(cst.FieldTag)
		}
		if dbf == "" {
			dbf = lang.Camel2Snake(fi.Name)
		}
		fColumns = append(fColumns, dbf)

		// 2. 确定结构体字段名称
		fFields = append(fFields, fi.Name)
		// 查找 auto
		if mFields[0] == "" {
			dbc := fi.Tag.Get(dbConfigTag)
			if strings.HasSuffix(dbc, dbAutoIncKeyFlag) {
				mFields[0] = fi.Name
			}
		}
		// 查找 primary
		if mFields[1] == "" {
			dbc := fi.Tag.Get(dbConfigTag)
			if strings.HasSuffix(dbc, dbPrimaryKeyFlag) {
				mFields[1] = fi.Name
			}
		}
		// 查找 updated
		if mFields[2] == "" {
			dbc := fi.Tag.Get(dbConfigTag)
			if strings.HasSuffix(dbc, dbUpdatedKeyFlag) {
				mFields[2] = fi.Name
			}
		}

		// 3. index
		cIdx := make([]int, 0)
		cIdx = append(cIdx, parentIdx...)
		cIdx = append(cIdx, i)
		fIndexes = append(fIndexes, cIdx)
	}
	return fColumns, fFields, fIndexes // 结构体的字段别名，字段名，顺序索引位置
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
