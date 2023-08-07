package dts

import (
	"errors"
	"math"
	"reflect"
	"sync"
	"unsafe"
)

var (
	errNumOutOfRange  = errors.New("dts: number out of range")
	errNotSupportType = errors.New("dts: can't support the value type")
)

// 缓存所有需要反序列化的实体结构的解析数据，防止反复不断的进行反射解析操作。
var cachedStructSchemas sync.Map

func cacheSetSchema(typ reflect.Type, val *StructSchema) {
	cachedStructSchemas.Store(typ, val)
}

func cacheGetSchema(typ reflect.Type) *StructSchema {
	if ret, ok := cachedStructSchemas.Load(typ); ok {
		return ret.(*StructSchema)
	}
	return nil
}

// 通用的绑定函数，将给定值写入指定的地址内存
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// int
func BindInt(p unsafe.Pointer, v int64) {
	*(*int)(p) = int(v)
}

func BindInt8(p unsafe.Pointer, v int64) {
	if v < math.MinInt8 || v > math.MaxInt8 {
		panic(errNumOutOfRange)
	}
	*(*int8)(p) = int8(v)
}

func BindInt16(p unsafe.Pointer, v int64) {
	if v < math.MinInt16 || v > math.MaxInt16 {
		panic(errNumOutOfRange)
	}
	*(*int16)(p) = int16(v)
}

func BindInt32(p unsafe.Pointer, v int64) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		panic(errNumOutOfRange)
	}
	*(*int32)(p) = int32(v)
}

func BindInt64(p unsafe.Pointer, v int64) {
	*(*int64)(p) = v
}

// uint
func BindUint(p unsafe.Pointer, v uint64) {
	*(*uint)(p) = uint(v)
}

func BindUint8(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint8 {
		panic(errNumOutOfRange)
	}
	*(*uint8)(p) = uint8(v)
}

func BindUint16(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint16 {
		panic(errNumOutOfRange)
	}
	*(*uint16)(p) = uint16(v)
}

func BindUint32(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint32 {
		panic(errNumOutOfRange)
	}
	*(*uint32)(p) = uint32(v)
}

func BindUint64(p unsafe.Pointer, v uint64) {
	*(*uint64)(p) = v
}

// float
func BindFloat32(p unsafe.Pointer, v float64) {
	if v < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
		panic(errNumOutOfRange)
	}
	*(*float32)(p) = float32(v)
}

func BindFloat64(p unsafe.Pointer, v float64) {
	*(*float64)(p) = v
}

// string & bool & any
func BindString(p unsafe.Pointer, v string) {
	*(*string)(p) = v
}

func BindBool(p unsafe.Pointer, v bool) {
	*(*bool)(p) = v
}

func BindAny(p unsafe.Pointer, v any) {
	*(*any)(p) = v
}
