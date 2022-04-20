package mapx

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/qinchende/gofast/skill/json"
	"github.com/qinchende/gofast/skill/stringx"
)

func setWithProperType(dst reflect.Value, src string) error {
	switch dst.Kind() {
	case reflect.Int:
		return setIntField(src, 0, dst)
	case reflect.Int8:
		return setIntField(src, 8, dst)
	case reflect.Int16:
		return setIntField(src, 16, dst)
	case reflect.Int32:
		return setIntField(src, 32, dst)
	case reflect.Int64:
		switch dst.Interface().(type) {
		case time.Duration:
			return setTimeDuration(src, dst)
		}
		return setIntField(src, 64, dst)
	case reflect.Uint:
		return setUintField(src, 0, dst)
	case reflect.Uint8:
		return setUintField(src, 8, dst)
	case reflect.Uint16:
		return setUintField(src, 16, dst)
	case reflect.Uint32:
		return setUintField(src, 32, dst)
	case reflect.Uint64:
		return setUintField(src, 64, dst)
	case reflect.Bool:
		return setBoolField(src, dst)
	case reflect.Float32:
		return setFloatField(src, 32, dst)
	case reflect.Float64:
		return setFloatField(src, 64, dst)
	case reflect.String:
		dst.SetString(src)
	case reflect.Struct:
		switch dst.Interface().(type) {
		case time.Time:
			return setTimeField(src, dst)
		}
		return json.Unmarshal(stringx.StringToBytes(src), dst.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal(stringx.StringToBytes(src), dst.Addr().Interface())
	default:
		return errors.New("unknown type")
	}
	return nil
}

func setIntField(src string, bitSize int, dst reflect.Value) error {
	if src == "" {
		src = "0"
	}
	intVal, err := strconv.ParseInt(src, 10, bitSize)
	if err == nil {
		dst.SetInt(intVal)
	}
	return err
}

func setUintField(src string, bitSize int, dst reflect.Value) error {
	if src == "" {
		src = "0"
	}
	uintVal, err := strconv.ParseUint(src, 10, bitSize)
	if err == nil {
		dst.SetUint(uintVal)
	}
	return err
}

func setBoolField(src string, dst reflect.Value) error {
	if src == "" {
		src = "false"
	}
	boolVal, err := strconv.ParseBool(src)
	if err == nil {
		dst.SetBool(boolVal)
	}
	return err
}

func setFloatField(src string, bitSize int, dst reflect.Value) error {
	if src == "" {
		src = "0.0"
	}
	floatVal, err := strconv.ParseFloat(src, bitSize)
	if err == nil {
		dst.SetFloat(floatVal)
	}
	return err
}

func setTimeField(src string, dst reflect.Value) error {
	//timeFormat := field.Tag.Get("time_format")
	//if timeFormat == "" {
	//	timeFormat = time.RFC3339
	//}

	timeFormat := time.RFC3339
	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(src, 10, 0)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		dst.Set(reflect.ValueOf(t))
		return nil

	}

	if src == "" {
		dst.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	//if isUTC, _ := strconv.ParseBool(field.Tag.Get("time_utc")); isUTC {
	//	l = time.UTC
	//}

	//if locTag := field.Tag.Get("time_location"); locTag != "" {
	//	loc, err := time.LoadLocation(locTag)
	//	if err != nil {
	//		return err
	//	}
	//	l = loc
	//}

	t, err := time.ParseInLocation(timeFormat, src, l)
	if err != nil {
		return err
	}

	dst.Set(reflect.ValueOf(t))
	return nil
}

func setArray(items []string, dst reflect.Value) error {
	for i, s := range items {
		err := setWithProperType(dst.Index(i), s)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(values []string, dest reflect.Value) error {
	slice := reflect.MakeSlice(dest.Type(), len(values), len(values))
	err := setArray(values, slice)
	if err != nil {
		return err
	}
	dest.Set(slice)
	return nil
}

func setTimeDuration(src string, dst reflect.Value) error {
	d, err := time.ParseDuration(src)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(d))
	return nil
}
