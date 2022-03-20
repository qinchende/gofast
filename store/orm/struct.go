package orm

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	dbTag      = "dbf"
	dbTagOther = "pms"
)

// 结构体中属性的数据库字段名称合集
func FieldValuesForDB(obj interface{}, includes ...bool) ([]string, []interface{}) {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	if rVal.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", rVal))
	}

	// 看类型，缓存有就直接用，否则计算一次并缓存
	var fields []string
	rTyp := rVal.Type()
	cModel := cacheGetModel(rTyp)
	if cModel == nil {
		f1, f2 := structFields(obj)
		fields = f1

		cModel = &cachedModel{
			fields:        fields,
			fieldIndexMap: make(map[string]int, len(fields)),
		}
		for idx, name := range f2 {
			cModel.fieldIndexMap[name] = idx
		}
		cacheSetModel(rTyp, cModel)
	} else {
		fields = cModel.fields
	}

	// 反射取值
	values := structValues(obj)
	return fields, values
}

func structFields(obj interface{}) ([]string, []string) {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	rTyp := rVal.Type()

	fls1 := make([]string, 0)
	fls2 := make([]string, 0)

	for i := 0; i < rVal.NumField(); i++ {
		// 通过值类型来确定后面
		va := rVal.Field(i)
		if va.Kind() == reflect.Struct {
			vaI := va.Interface()
			if _, ok := vaI.(time.Time); !ok {
				s1, s2 := structFields(vaI)
				fls1 = append(fls1, s1...)
				fls2 = append(fls2, s2...)
				continue
			}
		}

		// 通过字段类型，查找其中的标记
		fi := rTyp.Field(i)
		dbf := fi.Tag.Get(dbTag)
		if dbf == "" {
			dbf = fi.Tag.Get(dbTagOther)
		}
		if dbf != "" {
			fls1 = append(fls1, dbf)
		} else {
			fls1 = append(fls1, strings.ToLower(fi.Name))
		}
		fls2 = append(fls2, fi.Name)
	}
	return fls1, fls2
}

func structValues(obj interface{}) []interface{} {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	values := make([]interface{}, 0)

	for i := 0; i < rVal.NumField(); i++ {
		va := rVal.Field(i)
		if va.Kind() == reflect.Struct {
			vaI := va.Interface()
			if _, ok := vaI.(time.Time); !ok {
				subValues := structValues(vaI)
				values = append(values, subValues...)
				continue
			}
		}
		values = append(values, va.Interface())
	}
	return values
}
