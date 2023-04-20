package jde

import (
	"errors"
	"math"
	"reflect"
	"unsafe"
)

const (
	maxJsonLength     = math.MaxInt32 - 1 // 最大解析2GB JSON字符串
	tempByteStackSize = 128               // 栈上分配一定空间，方便放临时字符串（不能太大，防止协程栈伸缩）| 或者单独申请内存并管理
)

const (
	bytesNull  = "null"
	bytesTrue  = "true"
	bytesFalse = "false"
)

const (
	noErr        int = 0  // 没有错误
	scanEOF      int = -1 // 扫描结束
	errNormal    int = -2 // 没找到期望的字符
	errJson      int = -3 // 非法JSON格式
	errChar      int = -4 // 非预期的字符
	errEscape    int = -5
	errUnicode   int = -6
	errOverflow  int = -7
	errNumberFmt int = -8
	errExceedMax int = -9
	errInfinity  int = -10
	errMismatch  int = -11
	errUTF8      int = -12
	errKey       int = -13
	errValue     int = -14
	errKV        int = -15
	errNull      int = -16
	errObject    int = -17
	errArray     int = -18
	errTrue      int = -19
	errFalse     int = -20

	//errNotSupportType int = -13
)

type Type struct{}
type emptyInterface struct {
	typ *Type
	ptr unsafe.Pointer
}

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
//	isSpaceMask = (1 << ' ') | (1 << '\t') | (1 << '\r') | (1 << '\n')
//)

//go:nosplit
//go:inline
//func isSpace(c byte) bool {
//	return isSpaceMask&(1<<c) != 0
//}

//go:inline
func isSpace(c byte) bool {
	return isBlankChar[c]
}

var (
	isBlankChar = [256]bool{}

	numChars = [256]bool{
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
	}
)

func init() {
	isBlankChar[' '] = true
	isBlankChar['\n'] = true
	isBlankChar['\t'] = true
	isBlankChar['\r'] = true
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//type Kind uint
const (
	kindsCount     = 27
	isBaseTypeMask = 17956862 // 0001 0001 0001 1111 1111 1111 1110
	isNumKindMask  = 131068   // 0000 0000 0001 1111 1111 1111 1100
	allowNumMask   = 1179644  // 0000 0001 0001 1111 1111 1111 1100
	allowIntMask   = 1056764  // 0000 0001 0000 0001 1111 1111 1100
	allowFloatMask = 1073152  // 0000 0001 0000 0110 0000 0000 0000
	allowStrMask   = 17825792 // 0001 0001 0000 0000 0000 0000 0000
	allowBoolMask  = 1048578  // 0000 0001 0000 0000 0000 0000 0010
)

//go:inline
func isNumKind(k reflect.Kind) bool {
	return (1<<k)&isNumKindMask != 0
}

//go:inline
func allowNum(k reflect.Kind) bool {
	return (1<<k)&allowNumMask != 0
}

//go:inline
func allowInt(k reflect.Kind) bool {
	return (1<<k)&allowIntMask != 0
}

//go:inline
func allowFloat(k reflect.Kind) bool {
	return (1<<k)&allowFloatMask != 0
}

//go:inline
func allowStr(k reflect.Kind) bool {
	return (1<<k)&allowStrMask != 0
}

//go:inline
func allowBool(k reflect.Kind) bool {
	return (1<<k)&allowBoolMask != 0
}
