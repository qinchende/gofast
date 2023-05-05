package jde

import (
	"reflect"
	"unsafe"
)

func (sd *subDecode) bindBoolArr(v bool) {
	p := unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))
	if sd.isAny {
		*(*any)(p) = v
	} else {
		*(*bool)(p) = v
	}
	sd.arrIdx++
}

func (sd *subDecode) bindStringArr(v string) {
	p := unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))
	if sd.isAny {
		*(*any)(p) = v
	} else {
		*(*string)(p) = v
	}
	sd.arrIdx++
}

func (sd *subDecode) bindIntArr(v int64) {
	p := unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))
	switch sd.dm.itemKind {
	case reflect.Int:
		*(*int)(p) = int(v)
	case reflect.Interface:
		*(*any)(p) = v
	case reflect.Int8:
		*(*int8)(p) = int8(v)
	case reflect.Int16:
		*(*int16)(p) = int16(v)
	case reflect.Int32:
		*(*int32)(p) = int32(v)
	case reflect.Int64:
		*(*int64)(p) = v
	}
	sd.arrIdx++
}

func (sd *subDecode) bindUintArr(v uint64) {
	p := unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))
	switch sd.dm.itemKind {
	case reflect.Uint:
		*(*uint)(p) = uint(v)
	case reflect.Interface:
		*(*any)(p) = v
	case reflect.Uint8:
		*(*uint8)(p) = uint8(v)
	case reflect.Uint16:
		*(*uint16)(p) = uint16(v)
	case reflect.Uint32:
		*(*uint32)(p) = uint32(v)
	case reflect.Uint64:
		*(*uint64)(p) = v
	}
	sd.arrIdx++
}

func (sd *subDecode) bindFloatArr(v float64) {
	p := unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))
	switch sd.dm.itemKind {
	case reflect.Float32:
		*(*float32)(p) = float32(v)
	case reflect.Float64:
		*(*float64)(p) = v
	case reflect.Interface:
		*(*any)(p) = v
	}
	sd.arrIdx++
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (sd *subDecode) bindAnyArr(v any) {
//	*(*any)(unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))) = v
//	sd.arrIdx++
//}
//
////func bindArrValue[T any | string | bool](sd *subDecode, v T) {
////	*(*T)(unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))) = v
////	sd.arrIdx++
////}
//
//func bindArrNumValue[T constraints.Integer | constraints.Float, T2 int64 | float64](sd *subDecode, v T2) {
//	*(*T)(unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))) = T(v)
//	sd.arrIdx++
//}
//
//type arrIntFunc func(int64)
//type arrFloatFunc func(float64)
//
//func setIntFunc(k reflect.Kind) arrIntFunc {
//	switch k {
//	case reflect.Int:
//		return bindArrNumValue[int, int64]
//	case reflect.Int8:
//		return bindArrNumValue[int8, int64]
//	case reflect.Int16:
//		return bindArrNumValue[int16, int64]
//	case reflect.Int32:
//		return bindArrNumValue[int32, int64]
//	case reflect.Int64:
//		return bindArrNumValue[int64, int64]
//
//	case reflect.Uint:
//		return bindArrNumValue[uint, int64]
//	case reflect.Uint8:
//		return bindArrNumValue[uint8, int64]
//	case reflect.Uint16:
//		return bindArrNumValue[uint16, int64]
//	case reflect.Uint32:
//		return bindArrNumValue[uint32, int64]
//	case reflect.Uint64:
//		return bindArrNumValue[uint64, int64]
//	}
//	return nil
//}
//
//func setFloatFunc(k reflect.Kind) arrFloatFunc {
//	switch k {
//	case reflect.Float32:
//		return bindArrNumValue[float32, float64]
//	case reflect.Float64:
//		return bindArrNumValue[float64, float64]
//	}
//	return nil
//}
