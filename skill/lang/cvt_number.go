package lang

import (
	"math"
	"strconv"
)

// 任意类型的值，转换成 Int64，只要能转，不丢失精度都转，否则给出错误
func ToInt64(v any) (i64 int64, err error) {
	if v == nil {
		return 0, errorNilValue
	}

	switch vt := v.(type) {
	case string:
		i64, err = strconv.ParseInt(vt, 10, 64)
	case int:
		i64 = int64(vt)
	case uint:
		if vt <= math.MaxInt64 {
			i64 = int64(vt)
		} else {
			err = errorNumOutOfRange
		}
	case int64:
		i64 = vt
	case uint64:
		if vt <= math.MaxInt64 {
			i64 = int64(vt)
		} else {
			err = errorNumOutOfRange
		}
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
		i64, err = strconv.ParseInt(string(vt), 10, 64)
	default:
		err = errorConvertValue
	}
	return
}

// 任意类型的值，转换成 Int64，只要能转，不丢失精度都转，否则给出错误
func ToUint64(v any) (ui64 uint64, err error) {
	if v == nil {
		return 0, errorNilValue
	}

	switch vt := v.(type) {
	case string:
		ui64, err = strconv.ParseUint(vt, 10, 64)
	case int:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errorNumOutOfRange
		}
	case uint:
		ui64 = uint64(vt)
	case int64:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errorNumOutOfRange
		}
	case uint64:
		ui64 = vt
	case int32:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errorNumOutOfRange
		}
	case uint32:
		ui64 = uint64(vt)
	case int16:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errorNumOutOfRange
		}
	case uint16:
		ui64 = uint64(vt)
	case int8:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errorNumOutOfRange
		}
	case uint8:
		ui64 = uint64(vt)
	case []byte:
		ui64, err = strconv.ParseUint(string(vt), 10, 64)
	default:
		err = errorConvertValue
	}
	return
}

func ToFloat64(v any) (f64 float64, err error) {
	if v == nil {
		return 0.0, errorNilValue
	}

	switch vt := v.(type) {
	case string:
		f64, err = strconv.ParseFloat(vt, 64)
	case float32:
		f64 = float64(vt)
	case float64:
		f64 = vt
	case int:
		f64, err = strconv.ParseFloat(strconv.FormatInt(int64(vt), 10), 64)
	case uint:
		f64, err = strconv.ParseFloat(strconv.FormatUint(uint64(vt), 10), 64)
	case int64:
		f64, err = strconv.ParseFloat(strconv.FormatInt(vt, 10), 64)
	case uint64:
		f64, err = strconv.ParseFloat(strconv.FormatUint(vt, 10), 64)
	case int32:
		f64 = float64(vt)
	case uint32:
		f64 = float64(vt)
	case int16:
		f64 = float64(vt)
	case uint16:
		f64 = float64(vt)
	case int8:
		f64 = float64(vt)
	case uint8:
		f64 = float64(vt)
	case []byte:
		f64, err = strconv.ParseFloat(string(vt), 64)
	default:
		err = errorConvertValue
	}
	return
}

func ToFloat32(v any) (f32 float32, err error) {
	if v2, err2 := ToFloat64(v); err2 != nil {
		return 0.0, err2
	} else if v2 <= math.MaxFloat32 && v2 >= math.SmallestNonzeroFloat32 {
		return float32(v2), nil
	} else {
		return 0.0, errorNumOutOfRange
	}
}

func ToInt(v any) (i int, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt && v2 >= math.MinInt {
		return int(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}

func ToUint(v any) (ui uint, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint {
		return uint(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}

func ToInt32(v any) (i32 int32, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt32 && v2 >= math.MinInt32 {
		return int32(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}

func ToUint32(v any) (ui32 uint32, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint32 {
		return uint32(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}

func ToInt16(v any) (i16 int16, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt16 && v2 >= math.MinInt16 {
		return int16(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}

func ToUint16(v any) (ui16 uint16, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint16 {
		return uint16(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}

func ToInt8(v any) (i8 int8, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt8 && v2 >= math.MinInt8 {
		return int8(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}

func ToUint8(v any) (ui8 uint8, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint8 {
		return uint8(v2), nil
	} else {
		return 0, errorNumOutOfRange
	}
}
