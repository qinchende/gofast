package mapx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/skill/json"
	"github.com/qinchende/gofast/skill/stringx"
	"reflect"
	"strconv"
	"time"
)

// 返回错误的原则是转换时候发现格式错误，不能转换
func sdxSetValue(dst reflect.Value, src interface{}, opt *fieldOptions, useName bool, useDef bool) error {
	if src == nil {
		return nil
	}

	srcT := reflect.TypeOf(src)
	switch srcT.Kind() {
	case reflect.String:
		return sdxSetWithString(dst, src.(string))
	case reflect.Array, reflect.Slice:
		return applyList(dst, src, useName, useDef)
	case reflect.Map:
		// 如果Map不是 cst.KV 类型，这里会抛异常
		// TODO: 目前不支持这种情况
		return applyKVToStruct(dst, src.(map[string]interface{}), useName, useDef)
	}

	// 实体对象字段类型
	switch dst.Kind() {
	case reflect.Bool:
		bv, err := sdxAsBool(src)
		if err == nil {
			dst.SetBool(bv.(bool))
		}
		return err
	case reflect.Float32, reflect.Float64:
		fv, err := sdxAsFloat64(src)
		if err == nil {
			dst.SetFloat(fv.(float64))
		}
		return err
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv, err := sdxAsInt64(src)
		if err == nil {
			dst.SetInt(iv.(int64))
		}
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		uiv, err := sdxAsUInt64(src)
		if err == nil {
			dst.SetUint(uiv.(uint64))
		}
		return err
	case reflect.Interface:
		// 初始化零值
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	case reflect.Slice, reflect.Array:
		// 此时包装一下, 将单个Value包装成 slice
		newSrc := []interface{}{src}
		return applyList(dst, newSrc, useName, useDef)
	case reflect.Map:
		return errNotKVType
	case reflect.Struct:
		// 这个时候值可能是时间类型
		if _, ok := dst.Interface().(time.Time); ok {
			return sdxSetTimeDuration(dst, sdxAsString(src))
		}
	}
	return nil
}

func sdxAsString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(sv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(sv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(sv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(sv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(sv.Bool())
	}
	//return fmt.Sprint("%v", src)
	return fmt.Sprint(src)
}

func sdxAsInt64(src interface{}) (interface{}, error) {
	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return sv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(sv.Uint()), nil
	case reflect.Float32:
		return strconv.ParseInt(strconv.FormatFloat(sv.Float(), 'g', -1, 32), 10, 64)
	case reflect.Float64:
		return strconv.ParseInt(strconv.FormatFloat(sv.Float(), 'g', -1, 64), 10, 64)
	}
	return nil, fmt.Errorf("sdx: couldn't convert %v (%T) into type int64", src, src)
}

func sdxAsUInt64(src interface{}) (interface{}, error) {
	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(sv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return sv.Uint(), nil
	case reflect.Float32:
		return strconv.ParseUint(strconv.FormatFloat(sv.Float(), 'g', -1, 32), 10, 64)
	case reflect.Float64:
		return strconv.ParseUint(strconv.FormatFloat(sv.Float(), 'g', -1, 64), 10, 64)
	}
	return nil, fmt.Errorf("sdx: couldn't convert %v (%T) into type uint64", src, src)
}

func sdxAsFloat64(src interface{}) (interface{}, error) {
	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(sv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(sv.Uint()), nil
	case reflect.Float64, reflect.Float32:
		return sv.Float(), nil
	}
	return nil, fmt.Errorf("sdx: couldn't convert %v (%T) into type float64", src, src)
}

func sdxAsBool(src interface{}) (interface{}, error) {
	switch s := src.(type) {
	case bool:
		return s, nil
	case []byte:
		b, err := strconv.ParseBool(string(s))
		if err != nil {
			return nil, fmt.Errorf("sdx: couldn't convert %q into type bool", s)
		}
		return b, nil
	}

	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv := sv.Int()
		if iv == 1 || iv == 0 {
			return iv == 1, nil
		}
		return nil, fmt.Errorf("sdx: couldn't convert %d into type bool", iv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uv := sv.Uint()
		if uv == 1 || uv == 0 {
			return uv == 1, nil
		}
		return nil, fmt.Errorf("sdx: couldn't convert %d into type bool", uv)
	}

	return nil, fmt.Errorf("sdx: couldn't convert %v (%T) into type bool", src, src)
}

// utils
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func sdxSetWithString(dst reflect.Value, src string) error {
	//if src == "" {
	//	return nil
	//}

	switch dst.Kind() {
	case reflect.Int:
		return sdxSetInt(dst, src, 0)
	case reflect.Int8:
		return sdxSetInt(dst, src, 8)
	case reflect.Int16:
		return sdxSetInt(dst, src, 16)
	case reflect.Int32:
		return sdxSetInt(dst, src, 32)
	case reflect.Int64:
		switch dst.Interface().(type) {
		case time.Duration:
			return sdxSetTimeDuration(dst, src)
		}
		return sdxSetInt(dst, src, 64)
	case reflect.Uint:
		return sdxSetUint(dst, src, 0)
	case reflect.Uint8:
		return sdxSetUint(dst, src, 8)
	case reflect.Uint16:
		return sdxSetUint(dst, src, 16)
	case reflect.Uint32:
		return sdxSetUint(dst, src, 32)
	case reflect.Uint64:
		return sdxSetUint(dst, src, 64)
	case reflect.Bool:
		return sdxSetBool(dst, src)
	case reflect.Float32:
		return sdxSetFloat(dst, src, 32)
	case reflect.Float64:
		return sdxSetFloat(dst, src, 64)
	case reflect.String:
		dst.SetString(src)
	case reflect.Slice:
		vs := []string{src}
		return sdxSetStringSlice(dst, vs)
	case reflect.Array:
		vs := []string{src}
		if len(vs) != dst.Len() {
			return fmt.Errorf("%q is not valid value for %s", vs, dst.Type().String())
		}
		return sdxSetStringArray(dst, vs)
	case reflect.Map:
		return json.Unmarshal(stringx.StringToBytes(src), dst.Addr().Interface())
	case reflect.Struct:
		switch dst.Interface().(type) {
		case time.Time:
			return sdxSetTimeDuration(dst, src)
		}
		return json.Unmarshal(stringx.StringToBytes(src), dst.Addr().Interface())
	default:
		//return nil
		return errors.New("unknown type")
	}
	return nil
}

func sdxSetInt(dst reflect.Value, src string, bitSize int) error {
	intVal, err := strconv.ParseInt(src, 10, bitSize)
	if err == nil {
		dst.SetInt(intVal)
	}
	return err
}

func sdxSetUint(dst reflect.Value, src string, bitSize int) error {
	uintVal, err := strconv.ParseUint(src, 10, bitSize)
	if err == nil {
		dst.SetUint(uintVal)
	}
	return err
}

func sdxSetBool(dst reflect.Value, src string) error {
	boolVal, err := strconv.ParseBool(src)
	if err == nil {
		dst.SetBool(boolVal)
	}
	return err
}

func sdxSetFloat(dst reflect.Value, src string, bitSize int) error {
	floatVal, err := strconv.ParseFloat(src, bitSize)
	if err == nil {
		dst.SetFloat(floatVal)
	}
	return err
}

func sdxSetStringArray(dst reflect.Value, items []string) error {
	for i, s := range items {
		err := sdxSetWithString(dst.Index(i), s)
		if err != nil {
			return err
		}
	}
	return nil
}

func sdxSetStringSlice(dest reflect.Value, values []string) error {
	slice := reflect.MakeSlice(dest.Type(), len(values), len(values))
	err := sdxSetStringArray(slice, values)
	if err != nil {
		return err
	}
	dest.Set(slice)
	return nil
}

func sdxSetTimeDuration(dst reflect.Value, src string) error {
	d, err := time.ParseDuration(src)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(d))
	return nil
}

//func sdxSetTime(dst reflect.Value, src string) error {
//	//timeFormat := field.Tag.Get("time_format")
//	//if timeFormat == "" {
//	//	timeFormat = time.RFC3339
//	//}
//
//	timeFormat := time.RFC3339
//	switch tf := strings.ToLower(timeFormat); tf {
//	case "unix", "unixnano":
//		tv, err := strconv.ParseInt(src, 10, 0)
//		if err != nil {
//			return err
//		}
//
//		d := time.Duration(1)
//		if tf == "unixnano" {
//			d = time.Second
//		}
//
//		t := time.Unix(tv/int64(d), tv%int64(d))
//		dst.Set(reflect.ValueOf(t))
//		return nil
//
//	}
//
//	if src == "" {
//		dst.Set(reflect.ValueOf(time.Time{}))
//		return nil
//	}
//
//	l := time.Local
//	//if isUTC, _ := strconv.ParseBool(field.Tag.Get("time_utc")); isUTC {
//	//	l = time.UTC
//	//}
//
//	//if locTag := field.Tag.Get("time_location"); locTag != "" {
//	//	loc, err := time.LoadLocation(locTag)
//	//	if err != nil {
//	//		return err
//	//	}
//	//	l = loc
//	//}
//
//	t, err := time.ParseInLocation(timeFormat, src, l)
//	if err != nil {
//		return err
//	}
//
//	dst.Set(reflect.ValueOf(t))
//	return nil
//}
