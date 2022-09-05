package orm

import (
	"fmt"
	"github.com/qinchende/gofast/skill/stringx"
	"reflect"
	"strings"
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
	sm := Schema(obj)

	var vIndex int8 = 0 // 反射取值索引
	values := make([]any, len(sm.columns))
	structValues(&values, &vIndex, obj)

	return sm, values
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
			panic(fmt.Errorf("target item type must be structs; but got %T", rTyp))
		}

		// primary, updated
		mFields := [2]string{}
		rootIdx := make([]int, 0)
		fDB, fStruct, fIndexes := structFields(rTyp, rootIdx, &mFields)
		if mFields[0] == "" {
			mFields[0] = dbDefPrimaryKeyName
		}
		if mFields[1] == "" {
			mFields[1] = dbDefUpdatedKeyName
		}

		// 1. 主键的索引位置 ++++++++++
		var priIndex = -1
		for idx, f := range fStruct {
			if f == mFields[0] {
				priIndex = idx
				break
			}
		}
		if priIndex == -1 {
			panic(fmt.Errorf("%T, model has no primary key", rTyp)) // 不能没有主键
		}
		//if rVal.FieldByName(mFields[0]).Kind() != reflect.Uint {
		//	panic("primary key must uint") // 主键必须是 uint 类型
		//}
		// 2. updated 的索引位置
		var updateIndex = -1
		for idx, f := range fStruct {
			if f == mFields[1] {
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
			mdAttrs.TableName = stringx.Camel2Snake(rTyp.Name())
		}
		mdAttrs.cacheKeyFmt = "Gf#Line#%v#" + mdAttrs.TableName + "#%v"

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
func structFields(rTyp reflect.Type, parentIdx []int, mFields *[2]string) ([]string, []string, [][]int) {
	if rTyp.Kind() != reflect.Struct {
		panic(fmt.Errorf("%T is not like struct", rTyp))
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
		dbf := fi.Tag.Get(dbColumnNameTag)
		if dbf == "" {
			dbf = fi.Tag.Get(dbColumnNameTag2)
		}
		if dbf == "" {
			dbf = stringx.Camel2Snake(fi.Name)
		}
		fColumns = append(fColumns, dbf)

		// 2. 确定结构体字段名称
		fFields = append(fFields, fi.Name)
		// 查找 primary
		if mFields[0] == "" {
			dbc := fi.Tag.Get(dbConfigTag)
			if strings.HasSuffix(dbc, dbPrimaryKeyFlag) {
				mFields[0] = fi.Name
			}
		}
		// 查找 updated
		if mFields[1] == "" {
			dbc := fi.Tag.Get(dbConfigTag)
			if strings.HasSuffix(dbc, dbUpdatedKeyFlag) {
				mFields[1] = fi.Name
			}
		}

		// 3. index
		cIdx := make([]int, 0)
		cIdx = append(cIdx, parentIdx...)
		cIdx = append(cIdx, i)
		fIndexes = append(fIndexes, cIdx)
	}
	return fColumns, fFields, fIndexes
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func SchemaValuesByFields(obj ApplyOrmStruct, fields []string) (*ModelSchema, []interface{}) {
//	mSchema := fetchSchema(obj)
//	clsLen := mSchema.Length()
//
//	// 反射取值
//	values := make([]interface{}, clsLen, clsLen+1)
//	var valIndex int8 = 0               // 反射取值索引
//	var priIndex = mSchema.primaryIndex // 主键索引位置
//	pValIndex := &valIndex
//	pPriIndex := &priIndex
//	structValues2(&values, pValIndex, pPriIndex, obj)
//
//	return mSchema, values
//}
//
//// 反射提取结构体的值（支持嵌套递归）
//func structValues2(values *[]interface{}, nextIndex *int8, priIndex *int8, obj interface{}) {
//	rVal := reflect.Indirect(reflect.ValueOf(obj))
//
//	for i := 0; i < rVal.NumField(); i++ {
//		va := rVal.Field(i)
//		vaI := va.Interface()
//
//		if va.Kind() == reflect.Struct {
//			if _, ok := vaI.(time.Time); !ok {
//				structValues2(values, nextIndex, priIndex, vaI)
//				continue
//			}
//		}
//		if *nextIndex == *priIndex {
//			(*values)[0] = vaI
//			*priIndex = -1
//		} else {
//			*nextIndex++
//			(*values)[*nextIndex] = vaI
//		}
//	}
//}
