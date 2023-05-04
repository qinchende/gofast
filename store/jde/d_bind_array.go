package jde

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

func bindArrValue[T any | string | bool](sd *subDecode, v T) {
	*(*T)(unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))) = v
	sd.arrIdx++
}

func bindArrNumValue[T constraints.Integer | constraints.Float, T2 int64 | float64](sd *subDecode, v T2) {
	*(*T)(unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))) = T(v)
	sd.arrIdx++
}

type arrIntFunc func(*subDecode, int64)
type arrFloatFunc func(*subDecode, float64)

func setIntFunc(k reflect.Kind) arrIntFunc {
	switch k {
	case reflect.Int:
		return bindArrNumValue[int, int64]
	case reflect.Int8:
		return bindArrNumValue[int8, int64]
	case reflect.Int16:
		return bindArrNumValue[int16, int64]
	case reflect.Int32:
		return bindArrNumValue[int32, int64]
	case reflect.Int64:
		return bindArrNumValue[int64, int64]

	case reflect.Uint:
		return bindArrNumValue[uint, int64]
	case reflect.Uint8:
		return bindArrNumValue[uint8, int64]
	case reflect.Uint16:
		return bindArrNumValue[uint16, int64]
	case reflect.Uint32:
		return bindArrNumValue[uint32, int64]
	case reflect.Uint64:
		return bindArrNumValue[uint64, int64]
	}
	return nil
}

func setFloatFunc(k reflect.Kind) arrFloatFunc {
	switch k {
	case reflect.Float32:
		return bindArrNumValue[float32, float64]
	case reflect.Float64:
		return bindArrNumValue[float64, float64]
	}
	return nil
}
