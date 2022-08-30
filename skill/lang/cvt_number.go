package lang

import (
	"errors"
	"strconv"
)

var (
	errorNilValue     = errors.New("Value is nil.")
	errorConvertValue = errors.New("Data convert error.")
)

func ToInt64(v any) (i64 int64, err error) {
	if v == nil {
		return 0, errorNilValue
	}

	switch vt := v.(type) {
	case string:
		i64, err = strconv.ParseInt(v.(string), 10, 64)
	case int:
		i64 = int64(vt)
	case uint:
		i64 = int64(vt)
	case int64:
		i64 = int64(vt)
	case int32:
		i64 = int64(vt)
	case uint32:
		i64 = int64(vt)
	case int16:
		i64 = int64(vt)
	case uint16:
		i64 = int64(vt)
	case int8:
		i64 = int64(vt)
	case uint8:
		i64 = int64(vt)
	case []byte:
		srcStr := string(v.([]byte))
		i64, err = strconv.ParseInt(srcStr, 10, 64)
	default:
		err = errorConvertValue
	}
	return
}

func ToUInt64(src any) (ui64 uint64, err error) {
	if src == nil {
		return 0, errorNilValue
	}

	switch src.(type) {
	case string:
		ui64, err = strconv.ParseUint(src.(string), 10, 64)
	case uint:
		ui64 = uint64(src.(uint))
	case uint64:
		ui64 = src.(uint64)
	case uint32:
		ui64 = uint64(src.(uint32))
	case uint16:
		ui64 = uint64(src.(uint16))
	case uint8:
		ui64 = uint64(src.(uint8))
	case []byte:
		srcStr := string(src.([]byte))
		ui64, err = strconv.ParseUint(srcStr, 10, 64)
	default:
		err = errorConvertValue
	}
	return
}

//func ToInt(v any) (i int, err error) {
//	if v == nil {
//		return 0, errorNilValue
//	}
//
//	switch src := v.(type) {
//	case string:
//		i64, errT := strconv.ParseInt(v.(string), 10, 0)
//		i = int(i64)
//		err = errT
//	case int:
//		i64 = int64(src.(int))
//	case uint:
//		i64 = int64(src.(uint))
//	case int64:
//		i64 = src.(int64)
//	case uint64:
//		i64 = int64(src.(uint64))
//	case int32:
//		i64 = int64(src.(int32))
//	case uint32:
//		i64 = int64(src.(uint32))
//	case int16:
//		i64 = int64(src.(int16))
//	case uint16:
//		i64 = int64(src.(uint16))
//	case int8:
//		i64 = int64(src.(int8))
//	case uint8:
//		i64 = int64(src)
//	case []byte:
//		srcStr := string(src.([]byte))
//		i64, err = strconv.ParseInt(srcStr, 10, 64)
//	default:
//		err = errorConvertValue
//	}
//	return
//}
