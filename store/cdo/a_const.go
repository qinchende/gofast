// Copyright 2024 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package cdo

import (
	"errors"
	"github.com/qinchende/gofast/core/cst"
	"math"
	"reflect"
	"unsafe"
)

// cdo (Compact data of object)
// All type encoded format
const (
	TypeMask        byte = 0b11000000
	TypeValMask     byte = 0b00111111
	TypeListMask    byte = 0b11100000
	TypeListValMask byte = 0b00011111

	TypePosInt byte = 0b00000000 // 00
	TypeNegInt byte = 0b01000000 // 01
	TypeStr    byte = 0b10000000 // 10
	TypeMixed  byte = 0b11000000 // 11
	TypeList   byte = 0b11000000 // 110
	TypeFixed  byte = 0b11100000 // 111

	// TypeFixed subtypes +++++++++++++++++++++++
	FixNil      byte = 0xE0 // 0
	FixMixedNil byte = 0xE1 // 1
	FixTrue     byte = 0xE2 // 2
	FixFalse    byte = 0xE3 // 3
	FixF32      byte = 0xE4 // 8
	FixF64      byte = 0xE5 // 9
	FixDateTime byte = 0xE6 // 10
	FixDate     byte = 0xE7 // 11
	FixDuration byte = 0xE8 // 12
	FixTime     byte = 0xE9 // 13
	FixMax      byte = 0xFF // 31

	// TypeList subtypes ++++++++++++++++++++++++
	ListMask      byte = 0b11000000
	ListValMask   byte = 0b00111111
	ListVarIntPos byte = 0b00000000
	ListVarIntNeg byte = 0b10000000

	ListVarInt    byte = 0x00
	ListFixInt8   byte = 0x01
	ListFixInt16  byte = 0x02
	ListFixInt32  byte = 0x03
	ListFixInt64  byte = 0x04
	ListFixUint8  byte = 0x05
	ListFixUint16 byte = 0x06
	ListFixUint32 byte = 0x07
	ListFixUint64 byte = 0x08
	ListF32       byte = 0x09
	ListF64       byte = 0x0A
	ListStr       byte = 0x0B
	ListBool      byte = 0x0C
	ListKV        byte = 0x0D
	ListAny       byte = 0x0E
	ListObjFields byte = 0b01000000 // 都是object，提供所有的Fields
	ListObjIndex  byte = 0b10000000 // 都是object，提供前面出现的索引号
	ListExt       byte = 0b11000000 // 预留
)

const (
	MaxUint05 uint64 = 0x000000000000001F // 5
	MaxUint08 uint64 = 0x00000000000000FF // 8
	MaxUint13 uint64 = 0x0000000000001FFF // 5 + 8
	MaxUint16 uint64 = 0x000000000000FFFF // 8 + 8
	MaxUint21 uint64 = 0x00000000001FFFFF // 5 + 8 + 8
	MaxUint24 uint64 = 0x0000000000FFFFFF
	MaxUint32 uint64 = 0x00000000FFFFFFFF
	MaxUint40 uint64 = 0x000000FFFFFFFFFF
	MaxUint48 uint64 = 0x0000FFFFFFFFFFFF
	MaxUint56 uint64 = 0x00FFFFFFFFFFFFFF
	MaxUint64 uint64 = 0xFFFFFFFFFFFFFFFF
)

//func typeValue(b byte) (uint8, uint8) {
//	return b | TypeMask, b | TypeValMask
//}
//
//func extBytes2(v uint8) uint8 {
//	return v - 2
//}
//
//func extBytes4(v uint8) uint8 {
//	return v - 4
//}
//
//func extBytes8(v uint8) uint8 {
//	return v - 8
//}

//const (
//	// Type: 0, Format: 000|XXXXXX
//	// some fixed type
//	TypeFixed byte = 0b00000000 // 0
//	// Type0 subtypes
//	FixNil      byte = 0b00000000 // 0
//	FixMixedNil byte = 0b00000001 // 1
//	FixTrue     byte = 0b00000010 // 2
//	FixFalse    byte = 0b00000011 // 3
//	FixFloat32  byte = 0b00001000 // 8
//	FixFloat64  byte = 0b00001001 // 9
//	FixDateTime byte = 0b00001010 // 10
//	FixDate     byte = 0b00001011 // 11
//	FixDuration byte = 0b00001100 // 12
//	FixTime     byte = 0b00001101 // 13
//	FixMax      byte = 0b00011111 // 31
//
//	// Type: 1, Format: 001|XXXXXX
//	// all int numbers  which >= 0
//	TypePosInt byte = 0b00100000
//
//	// Type: 2, Format: 010|XXXXXX
//	// all int numbers  which < 0
//	TypeNegInt byte = 0b01000000
//
//	// Type: 3, Format: 011|XXXXXX
//	// all bytes array such as string/bytes
//	TypeBytes byte = 0b01100000
//
//	// Type: 4, Format: 100|XXXXXX
//	// array data
//	TypeList byte = 0b10000000
//
//	//// Type: 5, Format: 101|XXXXXX
//	//// just kvs
//	//TypeArrSame byte = 0b10100000
//	//
//	//// Type: 6, Format: 110|XXXXXX
//	////
//	//TypeMap byte = 0b11000000
//	//
//	//// Type: 7, Format: 111|XXXXXX
//	////
//	//TypeExt byte = 0b11100000
//
//	// ++++++++++++++++++++++++++++++++++++++
//	TypeSizeOffset2 uint8 = 2
//	TypeSizeOffset4 uint8 = 4
//	TypeSizeOffset8 uint8 = 8
//
//	TypeMask      byte = 0b11100000
//	TypeValueMask byte = 0b00011111
//)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

const (
	ptrMemSize   = int(unsafe.Sizeof(uintptr(0))) // 本机器指针占用字节数
	maxCdoStrLen = math.MaxUint32 - 1             // 最大解析 4GB Cdo 字符串
)

type (
	errType int
)

const (
	noErr        errType = 0  // 没有错误
	scanEOF      errType = -1 // 扫描结束
	errNormal    errType = -2 // 没找到期望的字符湖北 天门 铁路
	errCdo       errType = -3 // 非法Cdo格式
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
	errListType  errType = -18
	errBool      errType = -19
	errSupport   errType = -20
)

var errDescription = []string{
	noErr:           "ok",
	-(scanEOF):      "Error eof",
	-(errNormal):    "Error normal",
	-(errCdo):       "Error cdo format",
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
	errCdoTooLarge     = errors.New("cdo: string too large")
	errOutOfRange      = errors.New("cdo: out of range")
	errValueType       = errors.New("cdo: target value type error")
	errValueMustPtr    = errors.New("cdo: target value must pointer type")
	errValueMustSlice  = errors.New("cdo: target value must slice type")
	errValueMustStruct = errors.New("cdo: target value must struct type")
	errValueIsNil      = errors.New("cdo: target value is nil")
	errEmptyCdoStr     = errors.New("cdo: empty of cdo string")
	errCdoRowStr       = errors.New("cdo: wrong of GsonRow string")
	errCdoRowsStr      = errors.New("cdo: wrong of GsonRows string")
	errPtrLevel        = errors.New("cdo: target value is more than 3 layers of pointer")
	errMapType         = errors.New("cdo: can't support the map type")
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
