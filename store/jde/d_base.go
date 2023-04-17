package jde

import (
	"errors"
	"math"
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
	errValueIsNil   = errors.New("jde: target value is nil")
	errJsonEmpty    = errors.New("jde: json content empty")
)
