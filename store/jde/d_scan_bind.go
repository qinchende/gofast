package jde

import (
	"math"
	"unsafe"
)

// int
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindInt(p unsafe.Pointer, v int64) {
	*(*int)(p) = int(v)
}

func bindInt8(p unsafe.Pointer, v int64) {
	if v < math.MinInt8 || v > math.MaxInt8 {
		panic(errInfinity)
	}
	*(*int8)(p) = int8(v)
}

func bindInt16(p unsafe.Pointer, v int64) {
	if v < math.MinInt16 || v > math.MaxInt16 {
		panic(errInfinity)
	}
	*(*int16)(p) = int16(v)
}

func bindInt32(p unsafe.Pointer, v int64) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		panic(errInfinity)
	}
	*(*int32)(p) = int32(v)
}

func bindInt64(p unsafe.Pointer, v int64) {
	*(*int64)(p) = v
}

// uint
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindUint(p unsafe.Pointer, v uint64) {
	*(*uint)(p) = uint(v)
}

func bindUint8(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint8 {
		panic(errInfinity)
	}
	*(*uint8)(p) = uint8(v)
}

func bindUint16(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint16 {
		panic(errInfinity)
	}
	*(*uint16)(p) = uint16(v)
}

func bindUint32(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint32 {
		panic(errInfinity)
	}
	*(*uint32)(p) = uint32(v)
}

func bindUint64(p unsafe.Pointer, v uint64) {
	*(*uint64)(p) = v
}

// float
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindFloat32(p unsafe.Pointer, v float64) {
	if v < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
		panic(errInfinity)
	}
	*(*float32)(p) = float32(v)
}

func bindFloat64(p unsafe.Pointer, v float64) {
	*(*float64)(p) = v
}

// string & bool & any
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindString(p unsafe.Pointer, v string) {
	*(*string)(p) = v
}

func bindBool(p unsafe.Pointer, v bool) {
	*(*bool)(p) = v
}

func bindAny(p unsafe.Pointer, v any) {
	*(*any)(p) = v
}

// reset left array item
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) resetArrLeftItems() {
	var dfValue unsafe.Pointer
	if !sd.dm.isPtr {
		dfValue = zeroValues[sd.dm.itemBaseKind]
	}
	for i := sd.arrIdx; i < sd.dm.arrLen; i++ {
		*(*unsafe.Pointer)(unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(i*sd.dm.itemBytes))) = dfValue
	}
}
