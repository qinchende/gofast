package lang

import (
	"fmt"
	"strconv"
)

// ToString 获取变量的字符串值
// 浮点型 3.0 将会转换成字符串3 -> "3", 非数值或字符类型的变量将会被转换成JSON格式字符串
func ToString(v any) (s string, err error) {
	if v == nil {
		return "", errorNilValue
	}

	switch vt := v.(type) {
	case string:
		s = vt
	case bool:
		s = strconv.FormatBool(vt)
	case error:
		s = vt.Error()
	case float32:
		s = strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		s = strconv.FormatFloat(vt, 'f', -1, 64)
	case int:
		s = strconv.Itoa(vt)
	case int8:
		s = strconv.Itoa(int(vt))
	case int16:
		s = strconv.Itoa(int(vt))
	case int32:
		s = strconv.Itoa(int(vt))
	case int64:
		s = strconv.FormatInt(vt, 10)
	case uint:
		s = strconv.FormatUint(uint64(vt), 10)
	case uint8:
		s = strconv.FormatUint(uint64(vt), 10)
	case uint16:
		s = strconv.FormatUint(uint64(vt), 10)
	case uint32:
		s = strconv.FormatUint(uint64(vt), 10)
	case uint64:
		s = strconv.FormatUint(vt, 10)
	case []byte:
		s = string(vt)
	case fmt.Stringer:
		s = vt.String()
	default:
		s = fmt.Sprint(v)
	}
	return
}
