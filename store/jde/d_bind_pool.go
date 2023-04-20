package jde

import (
	"reflect"
	"sync"
	"unsafe"
)

var jdePool = sync.Pool{New: func() any { return &fastPool{} }}

type fastPool struct {
	arrStr []string
	//arrStrPtr  []*string
	arrBool []bool
	//arrBoolPtr []*bool
	arrAny []any
	//arrAnyPtr  []*any
	arrI64 []int64
	//arrI64Ptr  []*int64
	arrF64 []float64
	//arrF64Ptr  []*float64

	// +++
	arr listMeta
	obj structMeta
}

func (pl *fastPool) initMem() {
}

func (sd *subDecode) startListPool() {
	if isNumKind(sd.arr.itemKind) {
		sd.pl.arrI64 = sd.pl.arrI64[0:0]
		sd.pl.arrF64 = sd.pl.arrF64[0:0]
	} else {
		sd.pl.arrStr = sd.pl.arrStr[0:0]
		sd.pl.arrBool = sd.pl.arrBool[0:0]
		sd.pl.arrAny = sd.pl.arrAny[0:0]
	}
}

func (sd *subDecode) endListPool() {
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

	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func sliceSetNum[T int8 | int16 | int32 | int | int64 | uint8 | uint16 | uint32 | uint | uint64 | float32 | float64,
	T2 int64 | float64](val []T2, arr *listMeta) {

	size := len(val)

	// 如果是数组 +++++++++++++++++++++++++
	if arr.isPtr == false && arr.arrSize > 0 {
		tmpSlice := []T{}
		bh := (*reflect.SliceHeader)(unsafe.Pointer(&tmpSlice))
		bh.Data, bh.Len, bh.Cap = uintptr((*emptyInterface)(unsafe.Pointer(&arr.dst)).ptr), size, size

		for i := 0; i < size; i++ {
			tmpSlice[i] = T(val[i])
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
	newArrPtr1 := make([]*T, size)
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

func sliceSetString(val []string, arr *listMeta) {
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
