package mapx

import (
	"github.com/qinchende/gofast/cst"
	"reflect"
)

//func applyKV(dest interface{}, kvs cst.KV, useName bool, useDef bool) error {
//
//}

// 只用传入的值赋值对象
func applyKV(dest interface{}, kvs cst.KV, useName bool, useDef bool) error {
	if kvs == nil || len(kvs) == 0 {
		return nil
	}

	sm := Schema(dest)
	var fls []string
	if useName {
		fls = sm.fields
	} else {
		fls = sm.columns
	}
	flsOpts := sm.fieldsOpts

	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	for i := 0; i < len(fls); i++ {
		opt := flsOpts[i]
		fk := fls[i]
		sv, ok := kvs[fk]
		// 只要找到相应的Key，不管Value是啥，不能使用默认值，不能转换就抛异常，说明数据源错误
		if ok {
		} else if useDef && opt != nil && opt.DefExist {
			sv = opt.Default
		} else {
			continue
		}

		fv := sm.RefValueByIndex(&dstVal, int8(i))
		err := sdxSetValue(fv, sv, opt)
		errPanic(err)
	}
	return nil
}

//// 用传入的值赋值对象，没有的字段设置默认值
//func applyKVApplyDefault(dest interface{}, kvs cst.KV, useName bool, useDef bool) error {
//	return applyKV(dest, kvs, useName, useDef)
//}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//type setOptions struct {
//	isDefaultExists bool
//	defaultValue    string
//}
//
//// setter tries to set value on a walking by fields of a struct
//type setter interface {
//	TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSetted bool, err error)
//}
//
//var emptyField = reflect.StructField{}
//
//type pmsType map[string]interface{}
//
//func mapPmsByTag(ptr interface{}, pms cst.KV, tag string) error {
//	return mappingPmsByPtr(ptr, pmsType(pms), tag)
//}
//
//func mappingPmsByPtr(ptr interface{}, setter setter, tag string) error {
//	_, err := mappingPms(reflect.ValueOf(ptr), emptyField, setter, tag)
//	return err
//}
//
//func mappingPms(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
//	if field.Tag.Get(tag) == "-" { // just ignoring this field
//		return false, nil
//	}
//
//	var vKind = value.Kind()
//
//	if vKind == reflect.Ptr {
//		var isNew bool
//		vPtr := value
//		if value.IsNil() {
//			isNew = true
//			vPtr = reflect.New(value.Type().Elem())
//		}
//		isSetted, err := mappingPms(vPtr.Elem(), field, setter, tag)
//		if err != nil {
//			return false, err
//		}
//		if isNew && isSetted {
//			value.Set(vPtr)
//		}
//		return isSetted, nil
//	}
//
//	if vKind != reflect.Struct || !field.Anonymous {
//		ok, err := tryToSetValuePms(value, field, setter, tag)
//		if err != nil {
//			return false, err
//		}
//		if ok {
//			return true, nil
//		}
//	}
//
//	if vKind == reflect.Struct {
//		tValue := value.Type()
//
//		var isSetted bool
//		for i := 0; i < value.NumField(); i++ {
//			sf := tValue.Field(i)
//			if sf.PkgPath != "" && !sf.Anonymous { // unexported
//				continue
//			}
//			ok, err := mappingPms(value.Field(i), tValue.Field(i), setter, tag)
//			if err != nil {
//				return false, err
//			}
//			isSetted = isSetted || ok
//		}
//		return isSetted, nil
//	}
//	return false, nil
//}
//
//// TrySet tries to set a value by request's form source (like map[string][]string)
//func (pms pmsType) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
//	return setByPms(value, field, pms, tagValue, opt)
//}
//
//func tryToSetValuePms(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
//	var tagValue string
//	var setOpt setOptions
//
//	tagValue = field.Tag.Get(tag)
//	tagValue, opts := tagHead(tagValue, ",")
//
//	if tagValue == "" { // default value is FieldName
//		tagValue = field.Name
//		// modify by cd.net on 20220414 取字段的小写
//		//tagValue = stringx.Camel2Snake(field.Name)
//	}
//	if tagValue == "" { // when field is "emptyField" variable
//		return false, nil
//	}
//
//	var opt string
//	for len(opts) > 0 {
//		opt, opts = tagHead(opts, ",")
//
//		if k, v := tagHead(opt, "="); k == "default" {
//			setOpt.isDefaultExists = true
//			setOpt.defaultValue = v
//		}
//	}
//
//	return setter.TrySet(value, field, tagValue, setOpt)
//}
