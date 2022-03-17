package orm

import (
	"github.com/qinchende/gofast/logx"
	"reflect"
	"time"
)

const (
	dbTag      = "dbf"
	dbTagOther = "pms"
)

func DbFieldValues2(ins ...interface{}) ([]string, []interface{}) {
	fls := make([]string, 0)
	values := make([]interface{}, 0)

	for i := 0; i < len(ins); i++ {
		tF, tV := DbFieldValues(ins[i])

		fls = append(fls, tF...)
		values = append(values, tV...)
	}

	return fls, values
}

// 结构体中属性的数据库字段名称合集
func DbFieldValues(in interface{}, includes ...bool) ([]string, []interface{}) {
	fls := make([]string, 0)
	values := make([]interface{}, 0)

	sVal := reflect.ValueOf(in)
	if sVal.Kind() == reflect.Ptr {
		sVal = sVal.Elem()
	}

	// we only accept structs
	if sVal.Kind() != reflect.Struct {
		//panic(fmt.Errorf("ToMap only accepts structs; got %T", sVal))
		s, v := singleField(sVal.Interface())
		fls = append(fls, s)
		values = append(values, v)

		return fls, values
	}

	if _, ok := sVal.Interface().(time.Time); ok {
		s, v := singleField(sVal.Interface())
		fls = append(fls, s)
		values = append(values, v)
		return fls, values
	}

	typ := sVal.Type()
	for i := 0; i < sVal.NumField(); i++ {
		// value: get value
		va := sVal.Field(i)
		if va.Kind() == reflect.Struct {
			vaInter := va.Interface()
			if _, ok := vaInter.(time.Time); !ok {
				subFields, subValues := DbFieldValues(vaInter)
				fls = append(fls, subFields...)
				values = append(values, subValues...)
				continue
			}
		}
		values = append(values, va.Interface())

		// type: gets us a StructField
		fi := typ.Field(i)
		dbf := fi.Tag.Get(dbTag)
		if dbf == "" {
			dbf = fi.Tag.Get(dbTagOther)
		}
		if dbf != "" {
			fls = append(fls, dbf)
		} else {
			fls = append(fls, fi.Name)
		}
	}
	//var includeId bool
	//if len(includes) > 0 {
	//	includeId = includes[0]
	//}

	return fls, values
}

func singleField(in interface{}) (string, interface{}) {
	typ := reflect.TypeOf(in)
	sVal := reflect.ValueOf(in)

	logx.Info(typ)
	logx.Info(sVal)

	//sVal
	//
	//values = append(values, va.Interface())
	//
	//typ := sVal.Type()
	//
	//// type: gets us a StructField
	//fi := typ.Field(i)
	//dbf := fi.Tag.Get(dbTag)
	//if dbf == "" {
	//	dbf = fi.Tag.Get(dbTagOther)
	//}
	//if dbf != "" {
	//	fls = append(fls, dbf)
	//} else {
	//	fls = append(fls, fi.Name)
	//}

	return "", ""
}
