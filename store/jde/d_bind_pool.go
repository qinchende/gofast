package jde

import (
	"reflect"
	"sync"
	"unsafe"
)

var jdePool = sync.Pool{New: func() any { return &fastPool{} }}

type fastPool struct {
	arrI64  []int64
	arrF64  []float64
	arrStr  []string
	arrBool []bool
	arrAny  []any

	// +++++
	arr listPost
	obj structPost
}

type bindIntFunc func(*listPost, int64)
type bindFloatFunc func(*listPost, int64)
type bindStrFunc func(*listPost, int64)
type bindBoolFunc func(*listPost, int64)

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

func bindArrValue[T string | bool | float32 | float64](a *listPost, v T) {
	*(*T)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.arrSize))) = v
	a.arrIdx++
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (pl *fastPool) initMem() {
}

func (sd *subDecode) startListPool() {
	if isNumKind(sd.arr.itemKind) {
		sd.pl.arrI64 = sd.pl.arrI64[0:0]
		sd.pl.arrF64 = sd.pl.arrF64[0:0]
		sd.pl.arrStr = sd.pl.arrStr[0:0]
	} else {
		sd.pl.arrStr = sd.pl.arrStr[0:0]
		sd.pl.arrBool = sd.pl.arrBool[0:0]
		sd.pl.arrAny = sd.pl.arrAny[0:0]
	}
}

func (sd *subDecode) endListPool() {
	if sd.isArray && !sd.arr.isPtr {
		return
	}

	switch sd.arr.itemKind {
	case reflect.Int8:
		sliceSetNum[int8, int64](sd.pl.arrI64, sd.arr)
	case reflect.Int16:
		sliceSetNum[int16, int64](sd.pl.arrI64, sd.arr)
	case reflect.Int32:
		sliceSetNum[int32, int64](sd.pl.arrI64, sd.arr)
	case reflect.Int64:
		sliceSetNum[int64, int64](sd.pl.arrI64, sd.arr)
	case reflect.Int:
		sliceSetNum[int, int64](sd.pl.arrI64, sd.arr)
	case reflect.Uint8:
		sliceSetNum[uint8, int64](sd.pl.arrI64, sd.arr)
	case reflect.Uint16:
		sliceSetNum[uint16, int64](sd.pl.arrI64, sd.arr)
	case reflect.Uint32:
		sliceSetNum[uint32, int64](sd.pl.arrI64, sd.arr)
	case reflect.Uint64:
		sliceSetNum[uint64, int64](sd.pl.arrI64, sd.arr)
	case reflect.Uint:
		sliceSetNum[int, int64](sd.pl.arrI64, sd.arr)
	case reflect.Float32:
		sliceSetNum[float32, float64](sd.pl.arrF64, sd.arr)
	case reflect.Float64:
		sliceSetNum[float64, float64](sd.pl.arrF64, sd.arr)
	case reflect.String:
		sliceSetString(sd.pl.arrStr, sd.arr)
	case reflect.Interface:
		sliceSetString(sd.pl.arrStr, sd.arr)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func sliceSetNum[T int8 | int16 | int32 | int | int64 | uint8 | uint16 | uint32 | uint | uint64 | float32 | float64,
	T2 int64 | float64](val []T2, arr *listPost) {

	size := len(val)

	// 如果是数组 +++++++++++++++++++++++++
	if !arr.isPtr && arr.arrSize > 0 {
		dstArr := []T{}
		bh := (*reflect.SliceHeader)(unsafe.Pointer(&dstArr))
		bh.Data, bh.Len, bh.Cap = uintptr((*emptyInterface)(unsafe.Pointer(&arr.dst)).ptr), size, size

		for i := 0; i < size; i++ {
			dstArr[i] = T(val[i])
		}
		return
	}

	// 如果是 Slice ++++++++++++++++++++++
	newArr := make([]T, size)
	for i := 0; i < len(newArr); i++ {
		newArr[i] = T(val[i])
	}
	if arr.ptrLevel <= 0 {
		*(arr.dst.(*[]T)) = newArr
		return
	}

	// 第一级指针
	var newArrPtr1 []*T
	if arr.ptrLevel == 1 && arr.isPtr {
		//newArrPtr1 =

		bh := (*reflect.SliceHeader)(unsafe.Pointer(&newArrPtr1))
		bh.Data, bh.Len, bh.Cap = uintptr((*emptyInterface)(unsafe.Pointer(&arr.dst)).ptr), size, size

		for i := 0; i < size; i++ {
			newArrPtr1[i] = &newArr[i]
		}

		arr.ptrLevel = 0
		return
	}
	newArrPtr1 = make([]*T, size)
	for i := 0; i < len(newArr); i++ {
		newArrPtr1[i] = &newArr[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		*(arr.dst.(*[]*T)) = newArrPtr1
		return
	}

	// 第二级指针
	newArrPtr2 := make([]**T, size)
	for i := 0; i < len(newArrPtr1); i++ {
		newArrPtr2[i] = &newArrPtr1[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		*(arr.dst.(*[]**T)) = newArrPtr2
		return
	}

	// 第三级指针
	newArrPtr3 := make([]***T, size)
	for i := 0; i < len(newArrPtr2); i++ {
		newArrPtr3[i] = &newArrPtr2[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		*(arr.dst.(*[]***T)) = newArrPtr3
	}
	return
}

func sliceSetString(val []string, arr *listPost) {
	size := len(val)

	newArr := make([]string, size)
	copy(newArr, val)
	if arr.ptrLevel <= 0 {
		// arr.refVal.Set(reflect.ValueOf(newArr))
		*(arr.dst.(*[]string)) = newArr
		return
	}

	// 第一级指针
	newArrPtr1 := make([]*string, size)
	for i := 0; i < len(newArr); i++ {
		newArrPtr1[i] = &newArr[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		*(arr.dst.(*[]*string)) = newArrPtr1
		return
	}

	// 第二级指针
	newArrPtr2 := make([]**string, size)
	for i := 0; i < len(newArrPtr1); i++ {
		newArrPtr2[i] = &newArrPtr1[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		*(arr.dst.(*[]**string)) = newArrPtr2
		return
	}

	// 第三级指针
	newArrPtr3 := make([]***string, size)
	for i := 0; i < len(newArrPtr2); i++ {
		newArrPtr3[i] = &newArrPtr2[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		*(arr.dst.(*[]***string)) = newArrPtr3
	}
	return
}
