// Copyright 2024 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package cdo

import (
	"errors"
	"math"
	"reflect"
	"unsafe"
)

// cdo (Compact data of objects)
// All type encoded format
// 无特殊说明，本编码方案都遵从小端序原则
const (
	TypeMask        byte = 0b11000000
	TypeValMask     byte = 0b00111111
	TypeListMask    byte = 0b11100000
	TypeListValMask byte = 0b00011111

	// TypeFixed & TypeList are TypeMixed, both start with 00
	TypeFixed     byte = 0b00000000 // 000
	TypeList      byte = 0b00100000 // 001
	TypeVarIntPos byte = 0b01000000 // 01 >=0
	TypeVarIntNeg byte = 0b10000000 // 10 < 0
	TypeStr       byte = 0b11000000 // 11

	// TypeFixed subtypes +++++++++++++++++++++++
	FixFalse    byte = 0x00 // 0
	FixTrue     byte = 0x01 // 1
	FixF32      byte = 0x02 // 2
	FixF64      byte = 0x03 // 3
	FixTime     byte = 0x04 // 4 unix纳秒数，UTC时间
	FixNil      byte = 0x0E // 14
	FixNilMixed byte = 0x0F // 15
	FixMax      byte = 0x1F // 31

	// TypeList subtypes ++++++++++++++++++++++++
	// 00 | 000000
	ListMask      byte = 0b11000000
	ListValMask   byte = 0b00111111
	ListVarIntPos byte = 0b00000000
	ListVarIntNeg byte = 0b10000000

	ListVarInt    byte = 0x00
	ListF32       byte = 0x01
	ListF64       byte = 0x02
	ListBool      byte = 0x03
	ListStr       byte = 0x04
	ListKV        byte = 0x05
	ListAny       byte = 0x06
	ListTime      byte = 0x07
	ListFixInt08  byte = 0x10 // 固定长度的数值类型
	ListFixInt16  byte = 0x11
	ListFixInt32  byte = 0x12
	ListFixInt64  byte = 0x13
	ListFixUint08 byte = 0x14
	ListFixUint16 byte = 0x15
	ListFixUint32 byte = 0x16
	ListFixUint64 byte = 0x17

	// 01
	ListObjFields byte = 0b01000000 // 都是object，提供所有的Fields
	// 10
	ListObjIndex byte = 0b10000000 // 都是object，提供前面出现的索引号
	// 11
	ListExt byte = 0b11000000 // 预留
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

	MaxUint   uint64 = math.MaxUint
	OverInt   uint64 = -math.MinInt // 此程序只支持64位机器
	OverInt08 uint64 = -math.MinInt8
	OverInt16 uint64 = -math.MinInt16
	OverInt32 uint64 = -math.MinInt32
	OverInt64 uint64 = -math.MinInt64
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

const (
	ptrMemSize      = int(unsafe.Sizeof(uintptr(0))) // 本机器指针占用字节数
	maxCdoStrLen    = math.MaxUint32 - 1             // 最大解析 4GB Cdo 字符串
	maxNestingDepth = 1000                           // 编解码对象最大嵌套层级
)

type (
	errType int
)

const (
	noErr        errType = 0  // 没有错误
	errEOF       errType = -1 // 文件结尾错误
	errNormal    errType = -2 // 没找到期望的字符
	errCdoChar   errType = -3 // 非法Cdo格式
	errUnicode   errType = -6
	errOverflow  errType = -7
	errNumberFmt errType = -8
	errExceedMax errType = -9
	errInfinity  errType = -10 // 数值超出类型值范围
	errMismatch  errType = -11
	errUTF8      errType = -12
	errKey       errType = -13
	errValue     errType = -14
	errMap       errType = -15
	errNull      errType = -16
	errObject    errType = -17
	errList      errType = -18
	errListType  errType = -18
	errBool      errType = -19
	errSupport   errType = -20
	errOutRange  errType = -21
	errListSize  errType = -22
	errValType   errType = -23
)

var errDescription = []string{
	noErr:           "ok",
	-(errEOF):       "Error eof",
	-(errNormal):    "Error normal",
	-(errCdoChar):   "Error cdo format",
	-(errUnicode):   "Error unicode",
	-(errOverflow):  "Error overflow",
	-(errNumberFmt): "Error number format",
	-(errExceedMax): "Error exceed max depth",
	-(errInfinity):  "Error infinity",
	-(errMismatch):  "Error mismatch",
	-(errUTF8):      "Error utf8",
	-(errKey):       "Error key",
	-(errValue):     "Error value",
	-(errMap):       "Error kv map",
	-(errNull):      "Error null",
	-(errObject):    "Error object",
	-(errList):      "Error list",
	-(errBool):      "Error bool",
	-(errSupport):   "Error support",
	-(errOutRange):  "Error out of range",
	-(errListSize):  "Error wrong size of list",
	-(errValType):   "Error value type",
}

var (
	errValueIsNil      = errors.New("cdo: target value is nil")
	errCdoTooLarge     = errors.New("cdo: string length out of limit")
	errOutOfRange      = errors.New("cdo: out of range")
	errValueType       = errors.New("cdo: target value type error")
	errValueMustPtr    = errors.New("cdo: target value must pointer type")
	errValueMustSlice  = errors.New("cdo: target value must slice type")
	errValueMustStruct = errors.New("cdo: target value must struct type")
	errEmptyCdoStr     = errors.New("cdo: empty of cdo string")
	errCdoRowStr       = errors.New("cdo: wrong of GsonRow string")
	errCdoRowsStr      = errors.New("cdo: wrong of GsonRows string")
	errPtrLevel        = errors.New("cdo: target value is more than 3 layers of pointer")
	errMapType         = errors.New("cdo: can't support the map type")
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var (
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
)
