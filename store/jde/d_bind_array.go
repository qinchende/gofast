package jde

import "unsafe"

type bindIntFunc func(*listPost, int64)
type bindFloatFunc func(*listPost, int64)
type bindStrFunc func(*listPost, int64)
type bindBoolFunc func(*listPost, int64)

func bindArrValue[T string | bool | float32 | float64](a *listPost, v T) {
	*(*T)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = v
	a.arrIdx++
}

var (
	kindIntFunc = [27]bindIntFunc{
		2: func(a *listPost, v int64) {
			*(*int)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = int(v)
			a.arrIdx++
		},
		3: func(a *listPost, v int64) {
			*(*int8)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = int8(v)
			a.arrIdx++
		},
		4: func(a *listPost, v int64) {
			*(*int16)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = int16(v)
			a.arrIdx++
		},
		5: func(a *listPost, v int64) {
			*(*int32)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = int32(v)
			a.arrIdx++
		},
		6: func(a *listPost, v int64) {
			*(*int64)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = v
			a.arrIdx++
		},

		7: func(a *listPost, v int64) {
			*(*uint)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = uint(v)
			a.arrIdx++
		},
		8: func(a *listPost, v int64) {
			*(*uint8)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = uint8(v)
			a.arrIdx++
		},
		9: func(a *listPost, v int64) {
			*(*uint16)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = uint16(v)
			a.arrIdx++
		},
		10: func(a *listPost, v int64) {
			*(*uint32)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = uint32(v)
			a.arrIdx++
		},
		11: func(a *listPost, v int64) {
			*(*uint64)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = uint64(v)
			a.arrIdx++
		},
	}
)
