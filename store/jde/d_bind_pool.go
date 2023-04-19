package jde

import (
	"sync"
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
	arr listDest
	obj structDest
}

func (pl *fastPool) initMem() {
	//pl.arrStr = make([]string, 0, 512)
	//pl.arrStrPtr = make([]*string, 0, 512)
	//
	//pl.arrI64 = make([]int64, 0, 512)
	//pl.arrI64Ptr = make([]*int64, 0, 512)
}

//type arrayFunc func(pl *fastPool, arr *listDest)
//type sliceFunc func(pl *fastPool, arr *listDest)

//var arraySetFunc = [kindsCount][2]arrayFunc{
//	reflect.String: {setArrayString, setArrayStringPtr},
//}

//var sliceSetFunc = [kindsCount][2]sliceFunc{
//	reflect.String: {setSliceString, setSliceStringPtr},
//}

//var structSetFunc = [32]arrayFunc{
//	reflect.String: setArrayString,
//}

func (sd *subDecode) startListPool() {
	if isNumKind(sd.pl.arr.itemKind) {
		sd.pl.arrI64 = sd.pl.arrI64[0:0]
		//sd.pl.arrI64Ptr = sd.pl.arrI64Ptr[0:0]
	} else {
		sd.pl.arrStr = sd.pl.arrStr[0:0]
		//sd.pl.arrStrPtr = sd.pl.arrStrPtr[0:0]
	}
}

func (sd *subDecode) endListPool() {
	if !sd.isArray {
		if isNumKind(sd.pl.arr.itemKind) {
			setSliceInt[int](sd.pl, sd.arr)
		} else {
			setSliceString(sd.pl, sd.arr)
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func setSliceInt[T int8 | int16 | int32 | int | int64 | uint8 | uint16 | uint32 | uint | uint64](pl *fastPool, arr *listDest) {
	newArr := make([]T, len(pl.arrI64))
	for i := 0; i < len(newArr); i++ {
		newArr[i] = T(pl.arrI64[i])
	}
	if arr.ptrLevel <= 0 {
		//arr.refVal.Set(reflect.ValueOf(newArr))
		*(arr.dst.(*[]T)) = newArr
		return
	}

	// 第一级指针
	newArrPtr1 := make([]*T, len(pl.arrI64))
	for i := 0; i < len(newArr); i++ {
		newArrPtr1[i] = &newArr[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		//arr.refVal.Set(reflect.ValueOf(newArrPtr1))
		*(arr.dst.(*[]*T)) = newArrPtr1
		return
	}

	// 第二级指针
	newArrPtr2 := make([]**T, len(pl.arrI64))
	for i := 0; i < len(newArrPtr1); i++ {
		newArrPtr2[i] = &newArrPtr1[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		//arr.refVal.Set(reflect.ValueOf(newArrPtr2))
		*(arr.dst.(*[]**T)) = newArrPtr2
		return
	}

	// 第三级指针
	newArrPtr3 := make([]***T, len(pl.arrI64))
	for i := 0; i < len(newArrPtr2); i++ {
		newArrPtr3[i] = &newArrPtr2[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		//arr.refVal.Set(reflect.ValueOf(newArrPtr3))
		*(arr.dst.(*[]***T)) = newArrPtr3
	}
	return
}

func setSliceString(pl *fastPool, arr *listDest) {
	newArr := make([]string, len(pl.arrStr))
	copy(newArr, pl.arrStr)
	if arr.ptrLevel <= 0 {
		// arr.refVal.Set(reflect.ValueOf(newArr))
		*(arr.dst.(*[]string)) = newArr
		return
	}

	// 第一级指针
	newArrPtr1 := make([]*string, len(pl.arrStr))
	for i := 0; i < len(newArr); i++ {
		newArrPtr1[i] = &newArr[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		// arr.refVal.Set(reflect.ValueOf(newArrPtr1))
		*(arr.dst.(*[]*string)) = newArrPtr1
		return
	}

	// 第二级指针
	newArrPtr2 := make([]**string, len(pl.arrStr))
	for i := 0; i < len(newArrPtr1); i++ {
		newArrPtr2[i] = &newArrPtr1[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		// arr.refVal.Set(reflect.ValueOf(newArrPtr2))
		*(arr.dst.(*[]**string)) = newArrPtr2
		return
	}

	// 第三级指针
	newArrPtr3 := make([]***string, len(pl.arrStr))
	for i := 0; i < len(newArrPtr2); i++ {
		newArrPtr3[i] = &newArrPtr2[i]
	}
	arr.ptrLevel--
	if arr.ptrLevel <= 0 {
		// arr.refVal.Set(reflect.ValueOf(newArrPtr3))
		*(arr.dst.(*[]***string)) = newArrPtr3
	}
	return
}

//func setSliceString(pl *fastPool, arr *listDest) {
//	newArr := make([]T, len(pl.arrStr))
//	copy(newArr, pl.arrStr)
//
//	newArrPtr := make([]*string, len(pl.arrStr))
//	for i := 0; i < len(newArr); i++ {
//		newArrPtr[i] = &newArr[i]
//	}
//	//arr.refVal.Set(reflect.ValueOf(newArrPtr))
//	*(arr.dst.(*[]*string)) = newArrPtr
//
//
//
//
//	newArr := make([]T, len(pl.arrI64))
//	for i := 0; i < len(newArr); i++ {
//		newArr[i] = T(pl.arrI64[i])
//	}
//	if arr.ptrLevel <= 0 {
//		arr.refVal.Set(reflect.ValueOf(newArr))
//		//*(arr.dst.(*[]int)) = newArr
//		return
//	}
//
//	// 第一级指针
//	newArrPtr1 := make([]*T, len(pl.arrI64))
//	for i := 0; i < len(newArr); i++ {
//		newArrPtr1[i] = &newArr[i]
//	}
//	arr.ptrLevel--
//	if arr.ptrLevel <= 0 {
//		arr.refVal.Set(reflect.ValueOf(newArrPtr1))
//		return
//	}
//
//	// 第二级指针
//	newArrPtr2 := make([]**T, len(pl.arrI64))
//	for i := 0; i < len(newArrPtr1); i++ {
//		newArrPtr2[i] = &newArrPtr1[i]
//	}
//	arr.ptrLevel--
//	if arr.ptrLevel <= 0 {
//		arr.refVal.Set(reflect.ValueOf(newArrPtr2))
//		return
//	}
//
//	// 第三级指针
//	newArrPtr3 := make([]***T, len(pl.arrI64))
//	for i := 0; i < len(newArrPtr2); i++ {
//		newArrPtr3[i] = &newArrPtr2[i]
//	}
//	arr.ptrLevel--
//	if arr.ptrLevel <= 0 {
//		arr.refVal.Set(reflect.ValueOf(newArrPtr3))
//	}
//	return
//}
