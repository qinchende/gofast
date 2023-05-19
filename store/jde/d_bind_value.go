package jde

import (
	"math"
	"unsafe"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindInt(ptr uintptr, v int64) {
	*(*int)(unsafe.Pointer(ptr)) = int(v)
}

func bindInt8(ptr uintptr, v int64) {
	if v < math.MinInt8 || v > math.MaxInt8 {
		panic(errInfinity)
	}
	*(*int8)(unsafe.Pointer(ptr)) = int8(v)
}

func bindInt16(ptr uintptr, v int64) {
	if v < math.MinInt16 || v > math.MaxInt16 {
		panic(errInfinity)
	}
	*(*int16)(unsafe.Pointer(ptr)) = int16(v)
}

func bindInt32(ptr uintptr, v int64) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		panic(errInfinity)
	}
	*(*int32)(unsafe.Pointer(ptr)) = int32(v)
}

func bindInt64(ptr uintptr, v int64) {
	*(*int64)(unsafe.Pointer(ptr)) = v
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindUint(ptr uintptr, v uint64) {
	*(*uint)(unsafe.Pointer(ptr)) = uint(v)
}

func bindUint8(ptr uintptr, v uint64) {
	if v > math.MaxUint8 {
		panic(errInfinity)
	}
	*(*uint8)(unsafe.Pointer(ptr)) = uint8(v)
}

func bindUint16(ptr uintptr, v uint64) {
	if v > math.MaxUint16 {
		panic(errInfinity)
	}
	*(*uint16)(unsafe.Pointer(ptr)) = uint16(v)
}

func bindUint32(ptr uintptr, v uint64) {
	if v > math.MaxUint32 {
		panic(errInfinity)
	}
	*(*uint32)(unsafe.Pointer(ptr)) = uint32(v)
}

func bindUint64(ptr uintptr, v uint64) {
	*(*uint64)(unsafe.Pointer(ptr)) = v
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindFloat32(ptr uintptr, v float64) {
	if v < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
		panic(errInfinity)
	}
	*(*float32)(unsafe.Pointer(ptr)) = float32(v)
}

func bindFloat64(ptr uintptr, v float64) {
	*(*float64)(unsafe.Pointer(ptr)) = v
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindString(ptr uintptr, v string) {
	*(*string)(unsafe.Pointer(ptr)) = v
}

func bindBool(ptr uintptr, v bool) {
	*(*bool)(unsafe.Pointer(ptr)) = v
}

func bindAny(ptr uintptr, v any) {
	*(*any)(unsafe.Pointer(ptr)) = v
}

func (sd *subDecode) resetArrLeftItems() {
	dfValue := zeroValues[sd.dm.itemKind]
	for i := sd.arrIdx; i < sd.dm.arrLen; i++ {
		*(*unsafe.Pointer)(unsafe.Pointer(sd.dstPtr + uintptr(i*sd.dm.itemSize))) = dfValue
	}
}
