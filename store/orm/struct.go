package orm

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	dbDefPrimaryKeyName = "ID"            // 默认主键的字段名
	dbDefUpdatedKeyName = "UpdatedAt"     // 默认主键的字段名
	dbConfigTag         = "dbc"           // 数据库字段配置tag头
	dbPrimaryKeyFlag    = "primary_field" // 数据库主键tag头中配置值
	dbCreatedKeyFlag    = "created_field" // 创建时间
	dbUpdatedKeyFlag    = "updated_field" // 更新时间

	dbColumnNameTag  = "dbf" // 数据库字段名称，对应的tag
	dbColumnNameTag2 = "pms" // 数据库字段名称，次优先级
)

// 结构体中属性的数据库字段名称合集
func SchemaValues(obj ApplyOrmStruct) (*ModelSchema, []interface{}) {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	if rVal.Kind() != reflect.Struct {
		panic(fmt.Errorf("must structs; got %T", rVal))
	}

	// 看类型，缓存有就直接用，否则计算一次并缓存
	rTyp := rVal.Type()
	cModel := cacheGetModel(rTyp)
	if cModel == nil {
		// primary, updated
		mFields := [2]string{}
		fDB, fStruct := structFields(obj, &mFields)
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
			panic(fmt.Errorf("%T, model has no primary key", rVal)) // 不能没有主键
		}
		// 2. updated 的索引位置
		var updatedIndex = -1
		for idx, f := range fStruct {
			if f == mFields[1] {
				updatedIndex = idx
				break
			}
		}

		// db column name
		fKeyName := fDB[priIndex]
		fDBNew := make([]string, 0, len(fDB))
		fDBNew = append(fDBNew, fKeyName)
		fDBNew = append(fDBNew, fDB[:priIndex]...)
		fDBNew = append(fDBNew, fDB[priIndex+1:]...)

		// struct filed name
		fStructNew := make([]string, 0, len(fStruct))
		fStructNew = append(fStructNew, mFields[0])
		fStructNew = append(fStructNew, fStruct[:priIndex]...)
		fStructNew = append(fStructNew, fStruct[priIndex+1:]...)
		// +++++++++++++++++++++++++

		//fields = fDBNew
		cModel = &ModelSchema{
			tableName:    obj.TableName(),
			fields:       make(map[string]int8, len(fDBNew)),
			columns:      fDBNew,
			primaryIndex: int8(priIndex),
			updatedIndex: int8(updatedIndex),
		}
		for idx, name := range fStructNew {
			cModel.fields[name] = int8(idx)
		}
		cacheSetModel(rTyp, cModel)
	}

	// 反射取值
	values := make([]interface{}, cModel.Length())
	var valIndex int8 = 0
	var priIndex = cModel.primaryIndex
	pValIndex := &valIndex
	pPriIndex := &priIndex
	structValues(&values, pValIndex, pPriIndex, obj)

	return cModel, values
}

// 反射提取结构体的字段（支持嵌套递归）
func structFields(obj interface{}, mFields *[2]string) ([]string, []string) {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	rTyp := rVal.Type()

	fColumns := make([]string, 0)
	fFields := make([]string, 0)

	for i := 0; i < rVal.NumField(); i++ {
		// 通过值类型来确定后面
		va := rVal.Field(i)
		if va.Kind() == reflect.Struct {
			vaI := va.Interface()
			if _, ok := vaI.(time.Time); !ok {
				c, f := structFields(vaI, mFields)
				fColumns = append(fColumns, c...)
				fFields = append(fFields, f...)
				continue
			}
		}

		// 通过字段类型，查找其中的标记
		fi := rTyp.Field(i)
		dbf := fi.Tag.Get(dbColumnNameTag)
		if dbf == "" {
			dbf = fi.Tag.Get(dbColumnNameTag2)
		}
		if dbf != "" {
			fColumns = append(fColumns, dbf)
		} else {
			fColumns = append(fColumns, strings.ToLower(fi.Name))
		}
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
	}
	return fColumns, fFields
}

// 反射提取结构体的值（支持嵌套递归）
func structValues(values *[]interface{}, nextIndex *int8, priIndex *int8, obj interface{}) {
	rVal := reflect.Indirect(reflect.ValueOf(obj))

	for i := 0; i < rVal.NumField(); i++ {
		va := rVal.Field(i)
		if va.Kind() == reflect.Struct {
			vaI := va.Interface()
			if _, ok := vaI.(time.Time); !ok {
				structValues(values, nextIndex, priIndex, vaI)
				continue
			}
		}
		if *nextIndex == *priIndex {
			(*values)[0] = va.Interface()
			*priIndex = -1
		} else {
			*nextIndex++
			(*values)[*nextIndex] = va.Interface()
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// utils
