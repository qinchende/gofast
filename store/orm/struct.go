package orm

import (
	"fmt"
	"reflect"
)

const (
	dbTag      = "dbf"
	dbTagOther = "pms"
)

// 结构体中属性的数据库字段名称合集
func DbFieldNames(in interface{}, includes ...bool) []string {
	out := make([]string, 0)
	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// we only accept structs
	if val.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", val))
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		dbf := fi.Tag.Get(dbTag)
		if dbf == "" {
			dbf = fi.Tag.Get(dbTagOther)
		}
		if dbf != "" {
			out = append(out, dbf)
		} else {
			out = append(out, fi.Name)
		}
	}

	return out
}

//
//// ToMap converts interface into map
//func ToMap(in interface{}) map[string]interface{} {
//	out := make(map[string]interface{})
//	v := reflect.ValueOf(in)
//	if v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//
//	// we only accept structs
//	if v.Kind() != reflect.Struct {
//		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
//	}
//
//	typ := v.Type()
//	for i := 0; i < v.NumField(); i++ {
//		// gets us a StructField
//		fi := typ.Field(i)
//		if tagv := fi.Tag.Get(dbTag); tagv != "" {
//			// set key of map to value in struct field
//			val := v.Field(i)
//			zero := reflect.Zero(val.Type()).Interface()
//			current := val.Interface()
//
//			if reflect.DeepEqual(current, zero) {
//				continue
//			}
//			out[tagv] = current
//		}
//	}
//
//	return out
//}

//// RawFieldNames converts golang struct field into slice string
//func RawFieldNames(in interface{}, postgresSql ...bool) []string {
//	out := make([]string, 0)
//	v := reflect.ValueOf(in)
//	if v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//
//	var pg bool
//	if len(postgresSql) > 0 {
//		pg = postgresSql[0]
//	}
//
//	// we only accept structs
//	if v.Kind() != reflect.Struct {
//		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
//	}
//
//	typ := v.Type()
//	for i := 0; i < v.NumField(); i++ {
//		// gets us a StructField
//		fi := typ.Field(i)
//		if tagv := fi.Tag.Get(dbTag); tagv != "" {
//			if pg {
//				out = append(out, fmt.Sprintf("%s", tagv))
//			} else {
//				out = append(out, fmt.Sprintf("`%s`", tagv))
//			}
//		} else {
//			if pg {
//				out = append(out, fmt.Sprintf("%s", fi.Name))
//			} else {
//				out = append(out, fmt.Sprintf("`%s`", fi.Name))
//			}
//		}
//	}
//
//	return out
//}
