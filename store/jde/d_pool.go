package jde

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"sync"
	"unsafe"
)

// cached dest value meta info
var cachedDestMeta sync.Map

func cacheSetMeta(typ *dataType, val *destMeta) {
	cachedDestMeta.Store(typ, val)
}

func cacheGetMeta(typ *dataType) *destMeta {
	if ret, ok := cachedDestMeta.Load(typ); ok {
		return ret.(*destMeta)
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var jdePool = sync.Pool{New: func() any { return &fastPool{} }}

// TODO: buffer pool 需要有个机制，释放那些某次偶发申请太大的buffer，而导致长时间不释放的问题
type fastPool struct {
	bufI64 []int64
	bufU64 []uint64
	bufF64 []float64
	bufStr []string
	bufBol []bool
	bufAny []any
	escPos []int // 存放转义字符'\'的索引位置
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) resetListPool() {
	sd.pl.bufI64 = sd.pl.bufI64[0:0]
	sd.pl.bufU64 = sd.pl.bufU64[0:0]
	sd.pl.bufF64 = sd.pl.bufF64[0:0]
	sd.pl.bufStr = sd.pl.bufStr[0:0]
	sd.pl.bufBol = sd.pl.bufBol[0:0]
	sd.pl.bufAny = sd.pl.bufAny[0:0]
}

func (sd *subDecode) flushListPool() {
	// 如果是定长数组，不会用到缓冲池，不需要转储
	if sd.isArray && !sd.isPtr {
		return
	}

	switch sd.dm.itemKind {
	case reflect.Int:
		sliceSetNum[int, int64](sd.pl.bufI64, sd)
	case reflect.Int8:
		sliceSetNum[int8, int64](sd.pl.bufI64, sd)
	case reflect.Int16:
		sliceSetNum[int16, int64](sd.pl.bufI64, sd)
	case reflect.Int32:
		sliceSetNum[int32, int64](sd.pl.bufI64, sd)
	case reflect.Int64:
		sliceSetNum[int64, int64](sd.pl.bufI64, sd)

	case reflect.Uint:
		sliceSetNum[int, int64](sd.pl.bufI64, sd)
	case reflect.Uint8:
		sliceSetNum[uint8, int64](sd.pl.bufI64, sd)
	case reflect.Uint16:
		sliceSetNum[uint16, int64](sd.pl.bufI64, sd)
	case reflect.Uint32:
		sliceSetNum[uint32, int64](sd.pl.bufI64, sd)
	case reflect.Uint64:
		sliceSetNum[uint64, int64](sd.pl.bufI64, sd)

	case reflect.Float32:
		sliceSetNum[float32, float64](sd.pl.bufF64, sd)
	case reflect.Float64:
		sliceSetNum[float64, float64](sd.pl.bufF64, sd)

	case reflect.Bool:
		sliceSetBool(sd.pl.bufBol, sd)

	case reflect.String:
		sliceSetString(sd.pl.bufStr, sd)
	case reflect.Interface:
		sliceSetAny(sd.pl.bufAny, sd)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 整形和浮点型
func sliceSetNum[T constraints.Integer | constraints.Float, T2 int64 | float64](val []T2, sd *subDecode) {
	size := len(val)

	ptrLevel := sd.dm.ptrLevel
	// 如果是 Slice ++++++++++++++++++++++
	newArr := make([]T, size)
	for i := 0; i < len(newArr); i++ {
		newArr[i] = T(val[i])
	}
	if ptrLevel <= 0 {
		*(sd.dst.(*[]T)) = newArr
		return
	}

	// 第一级指针
	var newArrPtr1 []*T
	newArrPtr1 = make([]*T, size)
	for i := 0; i < len(newArr); i++ {
		newArrPtr1[i] = &newArr[i]
	}
	ptrLevel--
	if ptrLevel <= 0 {
		*(sd.dst.(*[]*T)) = newArrPtr1
		return
	}

	// 第二级指针
	newArrPtr2 := make([]**T, size)
	for i := 0; i < len(newArrPtr1); i++ {
		newArrPtr2[i] = &newArrPtr1[i]
	}
	ptrLevel--
	if ptrLevel <= 0 {
		*(sd.dst.(*[]**T)) = newArrPtr2
		return
	}

	return
}

//
//func copyNumSlice[T string | *string | **string](sd *subDecode, ptrLevel uint8, arr []T) []*T {
//	size := len(arr)
//
//	newArr := make([]*T, size)
//	for i := 0; i < size; i++ {
//		newArr[i] = &arr[i]
//	}
//
//	if ptrLevel <= 0 {
//		if sd.isArray {
//			dstSnap := []*T{}
//			bh := (*reflect.SliceHeader)(unsafe.Pointer(&dstSnap))
//			bh.Data, bh.Len, bh.Cap = sd.dstPtr, size, size
//			copy(dstSnap, newArr)
//		} else {
//			*(sd.dst.(*[]*T)) = newArr
//		}
//		return nil
//	} else {
//		return newArr
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 字符串处理
func sliceSetString(val []string, sd *subDecode) {
	ptrLevel := sd.dm.ptrLevel

	// 如果绑定对象是字符串切片
	newArr := make([]string, len(val))
	copy(newArr, val)
	if ptrLevel <= 0 {
		*(sd.dst.(*[]string)) = newArr
		return
	}

	// 一级指针
	ptrLevel--
	ret1 := copySlice[string](sd, ptrLevel, newArr)
	if ret1 == nil {
		return
	}

	// 二级指针
	ptrLevel--
	ret2 := copySlice[*string](sd, ptrLevel, ret1)
	if ret2 == nil {
		return
	}

	// 三级指针
	ptrLevel--
	_ = copySlice[**string](sd, ptrLevel, ret2)
	return
}

// Bool处理
func sliceSetBool(val []bool, sd *subDecode) {
	ptrLevel := sd.dm.ptrLevel

	// 如果绑定对象是字符串切片
	newArr := make([]bool, len(val))
	copy(newArr, val)
	if ptrLevel <= 0 {
		*(sd.dst.(*[]bool)) = newArr
		return
	}

	// 一级指针
	ptrLevel--
	ret1 := copySlice[bool](sd, ptrLevel, newArr)
	if ret1 == nil {
		return
	}

	// 二级指针
	ptrLevel--
	ret2 := copySlice[*bool](sd, ptrLevel, ret1)
	if ret2 == nil {
		return
	}

	// 三级指针
	ptrLevel--
	_ = copySlice[**bool](sd, ptrLevel, ret2)
	return
}

func sliceSetAny(val []any, sd *subDecode) {
	ptrLevel := sd.dm.ptrLevel

	// 如果绑定对象是字符串切片
	newArr := make([]any, len(val))
	copy(newArr, val)
	if ptrLevel <= 0 {
		*(sd.dst.(*[]any)) = newArr
		return
	}

	// 一级指针
	ptrLevel--
	ret1 := copySlice[any](sd, ptrLevel, newArr)
	if ret1 == nil {
		return
	}

	// 二级指针
	ptrLevel--
	ret2 := copySlice[*any](sd, ptrLevel, ret1)
	if ret2 == nil {
		return
	}

	// 三级指针
	ptrLevel--
	_ = copySlice[**any](sd, ptrLevel, ret2)
	return
}

func copySlice[T string | *string | **string | bool | *bool | **bool | any | *any | **any](sd *subDecode, ptrLevel uint8, arr []T) []*T {
	size := len(arr)

	newArr := make([]*T, size)
	for i := 0; i < size; i++ {
		newArr[i] = &arr[i]
	}

	if ptrLevel <= 0 {
		if sd.isArray {
			dstSnap := []*T{}
			bh := (*reflect.SliceHeader)(unsafe.Pointer(&dstSnap))
			bh.Data, bh.Len, bh.Cap = sd.dstPtr, size, size
			copy(dstSnap, newArr)
		} else {
			*(sd.dst.(*[]*T)) = newArr
		}
		return nil
	} else {
		return newArr
	}
}
