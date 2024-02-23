package jde

import (
	"errors"
	"github.com/qinchende/gofast/core/cst"
	"math"
	"reflect"
	"unsafe"
)

const (
	ptrMemSize    = int(unsafe.Sizeof(uintptr(0))) // 本机器指针占用字节数
	maxJsonStrLen = math.MaxInt32 - 1              // 最大解析2GB JSON字符串
)

var (
	nullBytes = []byte("null")
)

type (
	errType int
)

const (
	noErr        errType = 0  // 没有错误
	scanEOF      errType = -1 // 扫描结束
	errNormal    errType = -2 // 没找到期望的字符
	errJson      errType = -3 // 非法JSON格式
	errChar      errType = -4 // 非预期的字符
	errEscape    errType = -5
	errUnicode   errType = -6
	errOverflow  errType = -7
	errNumberFmt errType = -8
	errExceedMax errType = -9
	errInfinity  errType = -10 // 超出限制
	errMismatch  errType = -11
	errUTF8      errType = -12
	errKey       errType = -13
	errValue     errType = -14
	errKV        errType = -15
	errNull      errType = -16
	errObject    errType = -17
	errList      errType = -18
	errBool      errType = -19
	errSupport   errType = -20
)

var errDescription = []string{
	noErr:           "ok",
	-(scanEOF):      "Error eof",
	-(errNormal):    "Error normal",
	-(errJson):      "Error json format",
	-(errChar):      "Error char",
	-(errEscape):    "Error escape",
	-(errUnicode):   "Error unicode",
	-(errOverflow):  "Error overflow",
	-(errNumberFmt): "Error number format",
	-(errExceedMax): "Error exceed max depth",
	-(errInfinity):  "Error infinity",
	-(errMismatch):  "Error mismatch",
	-(errUTF8):      "Error utf8",
	-(errKey):       "Error key",
	-(errValue):     "Error value",
	-(errKV):        "Error kv map",
	-(errNull):      "Error null",
	-(errObject):    "Error object",
	-(errList):      "Error list",
	-(errBool):      "Error bool",
	-(errSupport):   "Error support",
}

var (
	errJsonTooLarge    = errors.New("jde: string too large")
	errValueType       = errors.New("jde: target value type error")
	errValueMustPtr    = errors.New("jde: target value must pointer type")
	errValueMustSlice  = errors.New("jde: target value must slice type")
	errValueMustStruct = errors.New("jde: target value must struct type")
	errValueIsNil      = errors.New("jde: target value is nil")
	errEmptyJsonStr    = errors.New("jde: empty of json string")
	errJsonRowStr      = errors.New("jde: wrong of GsonRow string")
	errJsonRowsStr     = errors.New("jde: wrong of GsonRows string")
	errPtrLevel        = errors.New("jde: target value is more than 3 layers of pointer")
	errMapType         = errors.New("jde: can't support the map type")
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//const (
//	isSpaceMask = (1 << ' ') | (1 << '\n') | (1 << '\r') | (1 << '\t')
//)
//
//func isBlank(c byte) bool {
//	return isSpaceMask&(1<<c) != 0
//}
//
//// 综合来说，判断空字符的综合性能是数组索引还不错，单一空字符多的情况下，直接||连接比较最好
//func isBlank(c byte) bool {
//	return isBlankChar[c]
//}

var (
	isBlankChar = [256]bool{
		' ':  true,
		'\n': true,
		'\r': true,
		'\t': true,
	}

	// ++++++++++++++++++++++++++++++++++++++需要转义的字符
	// \\ 反斜杠
	// \" 双引号
	// \' 单引号 （没有这个）
	// \/ 正斜杠
	// \b 退格符
	// \f 换页符
	// \t 制表符
	// \n 换行符
	// \r 回车符
	// \u 后面跟十六进制字符 （比如笑脸表情 \u263A）
	// +++++++++++++++++++++++++++++++++++++++++++++++++++
	unescapeChar = [256]byte{
		'u':  'u',
		'"':  '"',
		'\\': '\\',
		'/':  '/',

		'b': '\b',
		'f': '\f',
		'n': '\n',
		'r': '\r',
		't': '\t',
	}

	// escape unicode
	hexToInt = [256]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'A': 10,
		'B': 11,
		'C': 12,
		'D': 13,
		'E': 14,
		'F': 15,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
	}

	zeroNumValue = 0
	numUPtrVal   = *(*unsafe.Pointer)(reflect.ValueOf(&zeroNumValue).UnsafePointer())

	zeroBolValue = false
	bolUPtrVal   = *(*unsafe.Pointer)(reflect.ValueOf(&zeroBolValue).UnsafePointer())

	zeroStrValue = new(string)
	strUPtrVal   = *(*unsafe.Pointer)(reflect.ValueOf(&zeroStrValue).UnsafePointer())

	zeroValues = [27]unsafe.Pointer{
		reflect.Int:     numUPtrVal,
		reflect.Int8:    numUPtrVal,
		reflect.Int16:   numUPtrVal,
		reflect.Int32:   numUPtrVal,
		reflect.Int64:   numUPtrVal,
		reflect.Uint8:   numUPtrVal,
		reflect.Uint16:  numUPtrVal,
		reflect.Uint32:  numUPtrVal,
		reflect.Uint64:  numUPtrVal,
		reflect.Float32: numUPtrVal,
		reflect.Float64: numUPtrVal,
		reflect.Bool:    bolUPtrVal,
		reflect.String:  strUPtrVal,
	}

	rfTypeOfKV   = reflect.TypeOf(new(cst.KV)).Elem()
	rfTypeOfList = reflect.TypeOf(new([]any)).Elem()
	//rfTypeOfBytes = reflect.TypeOf(new([]byte)).Elem()
)

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// type Kind uint
//const (
//	kindsCount       = 27
//	isBaseTypeMask   = 17956862 // 0001 0001 0001 1111 1111 1111 1110
//	isNumKindMask    = 131068   // 0000 0000 0001 1111 1111 1111 1100
//	isIntKindMask    = 124      // 0000 0000 0000 0000 0000 0111 1100
//	isUintKindMask   = 8064     // 0000 0000 0000 0001 1111 1000 0000
//	receiveNumMask   = 1179644  // 0000 0001 0001 1111 1111 1111 1100
//	receiveIntMask   = 1056764  // 0000 0001 0000 0001 1111 1111 1100
//	receiveFloatMask = 1073152  // 0000 0001 0000 0110 0000 0000 0000
//	receiveStrMask   = 17825792 // 0001 0001 0000 0000 0000 0000 0000
//	receiveBoolMask  = 1048578  // 0000 0001 0000 0000 0000 0000 0010
//)
//
////go:inline
//func isNumKind(k reflect.Kind) bool {
//	return (1<<k)&isNumKindMask != 0
//}
//
////go:inline
//func isIntKind(k reflect.Kind) bool {
//	return (1<<k)&isIntKindMask != 0
//}
//
////go:inline
//func isUintKind(k reflect.Kind) bool {
//	return (1<<k)&isUintKindMask != 0
//}
//
//// 变量是否接收对应的值类型 ++++++++++++
////
////go:inline
//func allowNum(k reflect.Kind) bool {
//	return (1<<k)&receiveNumMask != 0
//}
//
////go:inline
//func allowInt(k reflect.Kind) bool {
//	return (1<<k)&receiveIntMask != 0
//}
//
//// 下面三种直接比较性能更好
////
////go:inline
//func allowFloat(k reflect.Kind) bool {
//	return (1<<k)&receiveFloatMask != 0
//}
//
////go:inline
//func allowStr(k reflect.Kind) bool {
//	return (1<<k)&receiveStrMask != 0
//}
//
////go:inline
//func allowBool(k reflect.Kind) bool {
//	return (1<<k)&receiveBoolMask != 0
//}
