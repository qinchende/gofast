// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"reflect"
)

type pmsType map[string]interface{}

func mapPms(ptr interface{}, pms map[string]interface{}) error {
	return mapPmsByTag(ptr, pms, "pms")
}

func mapPmsByTag(ptr interface{}, pms map[string]interface{}, tag string) error {
	return mappingPmsByPtr(ptr, pmsType(pms), tag)
}

func mappingPmsByPtr(ptr interface{}, setter setter, tag string) error {
	_, err := mappingPms(reflect.ValueOf(ptr), emptyField, setter, tag)
	return err
}

func mappingPms(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
	if field.Tag.Get(tag) == "-" { // just ignoring this field
		return false, nil
	}

	var vKind = value.Kind()

	if vKind == reflect.Ptr {
		var isNew bool
		vPtr := value
		if value.IsNil() {
			isNew = true
			vPtr = reflect.New(value.Type().Elem())
		}
		isSetted, err := mappingPms(vPtr.Elem(), field, setter, tag)
		if err != nil {
			return false, err
		}
		if isNew && isSetted {
			value.Set(vPtr)
		}
		return isSetted, nil
	}

	if vKind != reflect.Struct || !field.Anonymous {
		ok, err := tryToSetValuePms(value, field, setter, tag)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	if vKind == reflect.Struct {
		tValue := value.Type()

		var isSetted bool
		for i := 0; i < value.NumField(); i++ {
			sf := tValue.Field(i)
			if sf.PkgPath != "" && !sf.Anonymous { // unexported
				continue
			}
			ok, err := mappingPms(value.Field(i), tValue.Field(i), setter, tag)
			if err != nil {
				return false, err
			}
			isSetted = isSetted || ok
		}
		return isSetted, nil
	}
	return false, nil
}

// TrySet tries to set a value by request's form source (like map[string][]string)
func (pms pmsType) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByPms(value, field, pms, tagValue, opt)
}

func tryToSetValuePms(value reflect.Value, field reflect.StructField, setter setter, tag string) (bool, error) {
	var tagValue string
	var setOpt setOptions

	tagValue = field.Tag.Get(tag)
	tagValue, opts := head(tagValue, ",")

	if tagValue == "" { // default value is FieldName
		tagValue = field.Name
	}
	if tagValue == "" { // when field is "emptyField" variable
		return false, nil
	}

	var opt string
	for len(opts) > 0 {
		opt, opts = head(opts, ",")

		if k, v := head(opt, "="); k == "default" {
			setOpt.isDefaultExists = true
			setOpt.defaultValue = v
		}
	}

	return setter.TrySet(value, field, tagValue, setOpt)
}

func setByPms(value reflect.Value, field reflect.StructField, pms pmsType, tagValue string, opt setOptions) (isSetted bool, err error) {
	vs, ok := pms[tagValue]
	if !ok && !opt.isDefaultExists {
		return false, nil
	}

	switch value.Kind() {
	case reflect.Slice:
		//if !ok {
		//	vs = []string{opt.defaultValue}
		//}
		//return true, setPmsSlice(vs, value, field)
		return true, nil
	case reflect.Array:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		//if len(vs) != value.Len() {
		//	return false, fmt.Errorf("%q is not valid value for %s", vs, value.Type().String())
		//}
		//return true, setArray(vs, value, field)
		return true, nil
	default:
		var val string
		if !ok {
			val = opt.defaultValue
		}

		//if len(vs) > 0 {
		//	val = vs[0]
		//}
		val = vs.(string)
		return true, setWithProperType(val, value, field)
	}
}

//func setPmsSlice(vals interface{}, value reflect.Value, field reflect.StructField) error {
//	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
//	err := setArray(vals, slice, field)
//	if err != nil {
//		return err
//	}
//	value.Set(slice)
//	return nil
//}
