package mapx

import (
	"reflect"
)

func setBySdxApplyDefault(dst reflect.Value, src interface{}, opt *fieldOptions) error {
	switch reflect.TypeOf(src).Kind() {
	case reflect.String:
		return setWithProperType(dst, src.(string))
	case reflect.Bool:
		dst.SetBool(src.(bool))
	case reflect.Float32, reflect.Float64:
		val := src.(float64)
		dst.SetFloat(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := int64(src.(float64))
		dst.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val := uint64(src.(float64))
		dst.SetUint(val)
	default:
		// 实体对象字段类型
		switch dst.Kind() {
		//case reflect.Slice:
		//	if !ok {
		//		vs = []string{opt.Default}
		//	}
		//	return setSlice(vs, dst, field)
		//	//return true, nil
		//case reflect.Array:
		//	if !ok {
		//		vs = []string{opt.Default}
		//	}
		//	if len(vs) != dst.Len() {
		//		return fmt.Errorf("%q is not valid dst for %s", vs, dst.Type().String())
		//	}
		//	return setArray(vs, dst, field)
		default:
			return nil
		}
	}
	return nil
}
