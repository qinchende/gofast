package jde

import (
	"math"
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
		//if v < math.MinInt || v > math.MaxInt {
		//	goto errPanic
		//}
		*(*int)(p) = int(v)
	case reflect.Int8:
		if v < math.MinInt8 || v > math.MaxInt8 {
			goto errPanic
		}
		*(*int8)(p) = int8(v)
	case reflect.Int16:
		if v < math.MinInt16 || v > math.MaxInt16 {
			goto errPanic
		}
		*(*int16)(p) = int16(v)
	case reflect.Int32:
		if v < math.MinInt32 || v > math.MaxInt32 {
			goto errPanic
		}
		*(*int32)(p) = int32(v)
	case reflect.Int64:
		*(*int64)(p) = v
		//case reflect.Interface:
		//	*(*any)(p) = v
	}
	sd.arrIdx++
	return

errPanic:
	panic(errInfinity)
}

func (sd *subDecode) bindUintArr(v uint64) {
	p := unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))
	switch sd.dm.itemKind {
	case reflect.Uint:
		//if v > math.MaxUint {
		//	goto errPanic
		//}
		*(*uint)(p) = uint(v)
	case reflect.Uint8:
		if v > math.MaxUint8 {
			goto errPanic
		}
		*(*uint8)(p) = uint8(v)
	case reflect.Uint16:
		if v > math.MaxUint16 {
			goto errPanic
		}
		*(*uint16)(p) = uint16(v)
	case reflect.Uint32:
		if v > math.MaxUint32 {
			goto errPanic
		}
		*(*uint32)(p) = uint32(v)
	case reflect.Uint64:
		*(*uint64)(p) = v
		//case reflect.Interface:
		//	*(*any)(p) = v
	}
	sd.arrIdx++
	return

errPanic:
	panic(errInfinity)
}

func (sd *subDecode) bindFloatArr(v float64) {
	p := unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize))
	switch sd.dm.itemKind {
	case reflect.Float32:
		if v < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
			goto errPanic
		}
		*(*float32)(p) = float32(v)
	case reflect.Float64:
		*(*float64)(p) = v
		//case reflect.Interface:
		//	*(*any)(p) = v
	}
	sd.arrIdx++
	return

errPanic:
	panic(errInfinity)
}

func (sd *subDecode) resetArrLeftItems() {
	dfValue := zeroValues[sd.dm.itemKind]
	for i := sd.arrIdx; i < sd.dm.arrLen; i++ {
		*(*unsafe.Pointer)(unsafe.Pointer(sd.dstPtr + uintptr(i*sd.dm.itemSize))) = dfValue
	}
}
