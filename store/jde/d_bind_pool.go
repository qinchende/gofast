package jde

import (
	"reflect"
	"sync"
	"unsafe"
)

var jdePool = sync.Pool{New: func() any { return &fastPool{} }}

// TODO: buffer pool 需要有个机制，释放那些某次偶发申请太大的buffer，而导致后面一致不释放的问题
type fastPool struct {
	bufI64 []int64
	bufF64 []float64
	bufStr []string
	bufBol []bool
	bufAny []any

	// ++++++++++++
	arr listPost
	obj structPost
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) resetListPool() {
	//if isNumKind(sd.arr.itemKind) {
	//	sd.pl.bufI64 = sd.pl.bufI64[0:0]
	//	sd.pl.bufF64 = sd.pl.bufF64[0:0]
	//	sd.pl.bufStr = sd.pl.bufStr[0:0]
	//} else {
	//	sd.pl.bufStr = sd.pl.bufStr[0:0]
	//	sd.pl.bufBol = sd.pl.bufBol[0:0]
	//	sd.pl.bufAny = sd.pl.bufAny[0:0]
	//}

	sd.pl.bufI64 = sd.pl.bufI64[0:0]
	sd.pl.bufF64 = sd.pl.bufF64[0:0]
	sd.pl.bufStr = sd.pl.bufStr[0:0]
	sd.pl.bufBol = sd.pl.bufBol[0:0]
	sd.pl.bufAny = sd.pl.bufAny[0:0]
}

func (sd *subDecode) flushListPool() {
	//if sd.isArray && !sd.arr.isPtr {
	//	return
	//}

	switch sd.arr.itemKind {
	case reflect.Int:
		sliceSetNum[int, int64](sd.pl.bufI64, sd.arr)
	case reflect.Int8:
		sliceSetNum[int8, int64](sd.pl.bufI64, sd.arr)
	case reflect.Int16:
		sliceSetNum[int16, int64](sd.pl.bufI64, sd.arr)
	case reflect.Int32:
		sliceSetNum[int32, int64](sd.pl.bufI64, sd.arr)
	case reflect.Int64:
		sliceSetNum[int64, int64](sd.pl.bufI64, sd.arr)

	case reflect.Uint:
		sliceSetNum[int, int64](sd.pl.bufI64, sd.arr)
	case reflect.Uint8:
		sliceSetNum[uint8, int64](sd.pl.bufI64, sd.arr)
	case reflect.Uint16:
		sliceSetNum[uint16, int64](sd.pl.bufI64, sd.arr)
	case reflect.Uint32:
		sliceSetNum[uint32, int64](sd.pl.bufI64, sd.arr)
	case reflect.Uint64:
		sliceSetNum[uint64, int64](sd.pl.bufI64, sd.arr)

	case reflect.Float32:
		sliceSetNum[float32, float64](sd.pl.bufF64, sd.arr)
	case reflect.Float64:
		sliceSetNum[float64, float64](sd.pl.bufF64, sd.arr)

	case reflect.String:
		sliceSetString(sd.pl.bufStr, sd.arr)
	case reflect.Interface:
		sliceSetString(sd.pl.bufStr, sd.arr)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func sliceSetNum[T int8 | int16 | int32 | int | int64 | uint8 | uint16 | uint32 | uint | uint64 | float32 | float64,
	T2 int64 | float64](val []T2, arr *listPost) {

	size := len(val)

	// 如果是数组 +++++++++++++++++++++++++
	if !arr.isPtr && arr.arrLen > 0 {
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
	//if arr.ptrLevel == 1 && arr.isPtr {
	//	//newArrPtr1 =
	//
	//	bh := (*reflect.SliceHeader)(unsafe.Pointer(&newArrPtr1))
	//	bh.Data, bh.Len, bh.Cap = uintptr((*emptyInterface)(unsafe.Pointer(&arr.dst)).ptr), size, size
	//
	//	for i := 0; i < size; i++ {
	//		newArrPtr1[i] = &newArr[i]
	//	}
	//
	//	arr.ptrLevel = 0
	//	return
	//}
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

	// 如果是数组 +++++++++++++++++++++++++
	if !arr.isPtr && arr.arrLen > 0 {
		dstArr := []string{}
		bh := (*reflect.SliceHeader)(unsafe.Pointer(&dstArr))
		bh.Data, bh.Len, bh.Cap = uintptr((*emptyInterface)(unsafe.Pointer(&arr.dst)).ptr), size, size
		copy(dstArr, val)
		return
	}

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
