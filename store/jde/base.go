package jde

import (
	"errors"
	"math"
	"reflect"
)

const (
	maxJsonLength = math.MaxInt32 - 1 // 最大解析2GB JSON字符串
)

type errType int

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
	errInfinity  errType = -10
	errMismatch  errType = -11
	errUTF8      errType = -12
	errKey       errType = -13
	errValue     errType = -14
	errKV        errType = -15
	errNull      errType = -16
	errObject    errType = -17
	errList      errType = -18
	errBool      errType = -19
)

//var errorStrings = []string{
//	0:                      "ok",
//	-(scanEOF):              "eof",
//	ERR_INVALID_CHAR:       "invalid char",
//	ERR_INVALID_ESCAPE:     "invalid escape char",
//	ERR_INVALID_UNICODE:    "invalid unicode escape",
//	ERR_INTEGER_OVERFLOW:   "integer overflow",
//	ERR_INVALID_NUMBER_FMT: "invalid number format",
//	ERR_RECURSE_EXCEED_MAX: "recursion exceeded max depth",
//	ERR_FLOAT_INFINITY:     "float number is infinity",
//	ERR_MISMATCH:           "mismatched type with value",
//	ERR_INVALID_UTF8:       "invalid UTF8",
//}

var (
	//sErr            = errors.New("jsonx: json syntax error.")
	errJsonTooLarge = errors.New("jde: string too large")
	errValueType    = errors.New("jde: target value type error")
	errValueMustPtr = errors.New("jde: target value must pointer type")
	errValueIsNil   = errors.New("jde: target value is nil")
	errJsonEmpty    = errors.New("jde: json content empty")
	errPtrLevel     = errors.New("jde: target value is more than 3 layers of pointer")
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//const (
//	isSpaceMask = (1 << ' ') | (1 << '\n') | (1 << '\r') | (1 << '\t')
//)

////go:nosplit
////go:inline
//func isBlank(c byte) bool {
//	return isSpaceMask&(1<<c) != 0
//}

//// 综合来说，判断空字符的综合性能是数组索引还不错，单一空字符多的情况下，直接||连接比较最好
////go:inline
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

	unescapeChar = [256]byte{
		'"':  '"',
		'\\': '\\',
		'/':  '/',
		'b':  '\b',
		'f':  '\f',
		'n':  '\n',
		'r':  '\r',
		't':  '\t',
	}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//type Kind uint
const (
	kindsCount       = 27
	isBaseTypeMask   = 17956862 // 0001 0001 0001 1111 1111 1111 1110
	isNumKindMask    = 131068   // 0000 0000 0001 1111 1111 1111 1100
	receiveNumMask   = 1179644  // 0000 0001 0001 1111 1111 1111 1100
	receiveIntMask   = 1056764  // 0000 0001 0000 0001 1111 1111 1100
	receiveFloatMask = 1073152  // 0000 0001 0000 0110 0000 0000 0000
	receiveStrMask   = 17825792 // 0001 0001 0000 0000 0000 0000 0000
	receiveBoolMask  = 1048578  // 0000 0001 0000 0000 0000 0000 0010
)

//go:inline
func isNumKind(k reflect.Kind) bool {
	return (1<<k)&isNumKindMask != 0
}

// 变量是否接收对应的值类型 ++++++++++++
//go:inline
func allowNum(k reflect.Kind) bool {
	return (1<<k)&receiveNumMask != 0
}

//go:inline
func allowInt(k reflect.Kind) bool {
	return (1<<k)&receiveIntMask != 0
}

// 下面三种直接比较性能更好
//go:inline
func allowFloat(k reflect.Kind) bool {
	return (1<<k)&receiveFloatMask != 0
}

//go:inline
func allowStr(k reflect.Kind) bool {
	return (1<<k)&receiveStrMask != 0
}

//go:inline
func allowBool(k reflect.Kind) bool {
	return (1<<k)&receiveBoolMask != 0
}
