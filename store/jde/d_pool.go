package jde

import "reflect"

var pl fastPool

var sap arrPet
var ssp structPet

func init() {
	pl.arrStr = make([]string, 0, 500)
	pl.arrStrPtr = make([]*string, 0, 500)
}

type fastPool struct {
	arrStr    []string
	arrStrPtr []*string
	arrBool   []bool
	arrAny    []any
}

type arrayFunc func(arr *arrPet)
type sliceFunc func(arr *arrPet)

var arraySetFunc = [32][2]arrayFunc{
	reflect.String: {setArrayString, setArrayStringPtr},
}

var sliceSetFunc = [32][2]sliceFunc{
	reflect.String: {setSliceString, setSliceStringPtr},
}

//var structSetFunc = [32]arrayFunc{
//	reflect.String: setArrayString,
//}

func (sd *subDecode) startListPool() {
	pl.arrStr = pl.arrStr[0:0]
	pl.arrStrPtr = pl.arrStrPtr[0:0]
}

func (sd *subDecode) endListPool() {
	if sd.arr.arrType.Kind() == reflect.Slice {
		if sd.arr.isPtr {
			fun := sliceSetFunc[sd.arr.recKind][1]
			fun(sd.arr)
		} else {
			fun := sliceSetFunc[sd.arr.recKind][0]
			fun(sd.arr)
		}
	} else {
		if sd.arr.isPtr {
			fun := arraySetFunc[sd.arr.recKind][1]
			fun(sd.arr)
		} else {
			fun := arraySetFunc[sd.arr.recKind][0]
			fun(sd.arr)
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func setArrayString(arr *arrPet) {
	ct := arr.val.Len()
	if ct > len(pl.arrStr) {
		ct = len(pl.arrStr)
	}

	newArr := make([]string, len(pl.arrStr))
	copy(newArr, pl.arrStr)
	*(arr.dst.(*[]string)) = newArr
}

func setArrayStringPtr(arr *arrPet) {
	newArr := make([]string, len(pl.arrStr))
	copy(newArr, pl.arrStr)
	*(arr.dst.(*[]string)) = newArr
}

// ++++++++++++++++++++++++++++++++++++++++
func setSliceString(arr *arrPet) {
	newArr := make([]string, len(pl.arrStr))
	copy(newArr, pl.arrStr)
	*(arr.dst.(*[]string)) = newArr
}

func setSliceStringPtr(arr *arrPet) {
	newArr := make([]string, len(pl.arrStr))
	copy(newArr, pl.arrStr)

	newArrPtr := make([]*string, len(pl.arrStr))
	for i := 0; i < len(newArr); i++ {
		newArrPtr[i] = &newArr[i]
	}
	*(arr.dst.(*[]*string)) = newArrPtr
}

//
//func (sd *subDecode) startListPool() {
//	pl.arrStr = pl.arrStr[0:0]
//	pl.arrStrPtr = pl.arrStrPtr[0:0]
//	//if sd.arr.recKind == reflect.String {
//	//	pl.arrStr = make([]string, 0, 16)
//	//}
//}
//
//func (sd *subDecode) endListPool() {
//	//dstNew := reflect.MakeSlice(sliceTyp, srcVal.Len(), srcVal.Len())
//	//dstVal.Set(dstNew)
//
//	newArr := make([]string, len(pl.arrStr))
//	copy(newArr, pl.arrStr)
//	//*(sd.arr.dst.(*[]string)) = newArr
//
//	newArrPtr := make([]*string, len(pl.arrStr))
//	for i := 0; i < len(newArr); i++ {
//		newArrPtr[i] = &newArr[i]
//	}
//	*(sd.arr.dst.(*[]*string)) = newArrPtr
//
//	//sd.arr.val.Set(reflect.ValueOf(newArr))
//	//*(sd.arr.dst.(*[]string)) = pl.arrStr
//}
