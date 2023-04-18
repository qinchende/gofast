package lang

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// NOTE：STB 和 BTS 这种黑魔法转换是不推荐使用的，特殊场景可能会出现意想不到的错误。
// go 1.20后期版本中会提供标准库，实现类似的功能
// StringToBytes converts string to byte slice without a memory allocation.
//func StringToBytes(s string) []byte {
//	return STB(s)
//}
func STB(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
	// 还有下面这种写法
	//l := len(s)
	//return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
	//	Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
	//	Len:  l,
	//	Cap:  l,
	//}))
}

// BytesToString converts byte slice to string without a memory allocation.
//func BytesToString(b []byte) string {
//	return BTS(b)
//}
func BTS(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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
		s = fmt.Sprintf("%+v", v)
	}
	return
}
