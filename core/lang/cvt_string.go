package lang

import (
	"fmt"
	"strconv"
	"unsafe"
)

// NOTE：S2B 和 B2S 这种黑魔法转换是不推荐使用的，特殊场景可能会出现意想不到的错误。
// go 1.20后期版本中会提供标准库，实现类似的功能
func S2B(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
	//sh := *(*rt.StringHeader)(unsafe.Pointer(&s))
	//bh := (*rt.SliceHeader)(unsafe.Pointer(&b))
	//bh.DataPtr, bh.Len, bh.Cap = sh.DataPtr, sh.Len, sh.Len
	//return b
}

func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
	//return *(*string)(unsafe.Pointer(&b))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func ToString(v any) (s string) {
	s, _ = ToString2(v)
	return
}

// ToString 获取变量的字符串值
// 浮点型 3.0 将会转换成字符串3 -> "3", 非数值或字符类型的变量将会被转换成JSON格式字符串
func ToString2(v any) (s string, err error) {
	if v == nil {
		return "", errNilValue
	}

	switch vt := v.(type) {
	case string:
		s = vt
	case bool:
		s = strconv.FormatBool(vt)
	case error:
		s = vt.Error()
	case float32:
		s = strconv.FormatFloat(float64(vt), 'g', -1, 32)
	case float64:
		s = strconv.FormatFloat(vt, 'g', -1, 64)
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
		//s = fmt.Sprint(v)
		s = fmt.Sprintf("%+v", v)
	}
	return
}
