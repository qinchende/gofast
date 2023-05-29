package jde

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"sync"
	"unsafe"
)

var (
	jdeDecPool     = sync.Pool{New: func() any { return &subDecode{} }}
	jdeBufPool     = sync.Pool{New: func() any { return &listPool{} }}
	cachedDestMeta sync.Map // cached dest value meta info
)

func cacheSetMeta(typ *dataType, val *destMeta) {
	cachedDestMeta.Store(typ, val)
}

func cacheGetMeta(typ *dataType) *destMeta {
	if ret, ok := cachedDestMeta.Load(typ); ok {
		return ret.(*destMeta)
	}
	return nil
}

// TODO: buffer pool 需要有个机制，释放那些某次偶发申请太大的buffer，而导致长时间不释放的问题
type listPool struct {
	bufPtr []unsafe.Pointer
	bufI64 []int64
	bufU64 []uint64
	bufF64 []float64
	bufStr []string
	bufBol []bool
	bufAny []any
	nulPos []int // 指针类型值，可能是nil
	//escPos []int // 存放转义字符'\'的索引位置
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) resetListPool() {
	if sd.dm.isArrBind {
		return
	}

	// 获取缓存空间
	sd.pl = jdeBufPool.Get().(*listPool)

	sd.pl.bufPtr = sd.pl.bufPtr[0:0]
	sd.pl.bufI64 = sd.pl.bufI64[0:0]
	sd.pl.bufU64 = sd.pl.bufU64[0:0]
	sd.pl.bufF64 = sd.pl.bufF64[0:0]
	sd.pl.bufStr = sd.pl.bufStr[0:0]
	sd.pl.bufBol = sd.pl.bufBol[0:0]
	sd.pl.bufAny = sd.pl.bufAny[0:0]

	sd.pl.nulPos = sd.pl.nulPos[0:0] // 必须要初始化
}

func (sd *subDecode) flushListPool() {
	// 如果是定长数组，不会用到缓冲池，不需要转储
	if sd.dm.isArrBind {
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
		sliceSetNum[int, uint64](sd.pl.bufU64, sd)
	case reflect.Uint8:
		sliceSetNum[uint8, uint64](sd.pl.bufU64, sd)
	case reflect.Uint16:
		sliceSetNum[uint16, uint64](sd.pl.bufU64, sd)
	case reflect.Uint32:
		sliceSetNum[uint32, uint64](sd.pl.bufU64, sd)
	case reflect.Uint64:
		sliceSetNum[uint64, uint64](sd.pl.bufU64, sd)

	case reflect.Float32:
		sliceSetNum[float32, float64](sd.pl.bufF64, sd)
	case reflect.Float64:
		sliceSetNum[float64, float64](sd.pl.bufF64, sd)

	case reflect.Bool:
		sliceSetNotNum[bool](sd.pl.bufBol, sd)
	case reflect.String:
		sliceSetNotNum[string](sd.pl.bufStr, sd)
	case reflect.Interface:
		sliceSetNotNum[any](sd.pl.bufAny, sd)
	}

	// 用完了就归还
	jdeBufPool.Put(sd.pl)
	sd.pl = nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 整形和浮点型
func sliceSetNum[T constraints.Integer | constraints.Float, T2 int64 | uint64 | float64](val []T2, sd *subDecode) {
	size := len(val)

	ptrLevel := sd.dm.ptrLevel
	newArr := make([]T, size)
	for i := 0; i < len(newArr); i++ {
		newArr[i] = T(val[i])
	}

	// 这里只可能是slice
	if ptrLevel <= 0 {
		if len(sd.pl.nulPos) > 0 {
			dstSnap := []T{}
			*(*reflect.SliceHeader)(unsafe.Pointer(&dstSnap)) = *(*reflect.SliceHeader)(unsafe.Pointer(sd.dstPtr))
			for i := 0; i < len(sd.pl.nulPos); i++ {
				idx := sd.pl.nulPos[i]
				if idx >= len(dstSnap) {
					break
				}
				newArr[idx] = dstSnap[idx]
			}
		}
		//*(sd.dst.(*[]T)) = newArr
		*(*[]T)(unsafe.Pointer(sd.dstPtr)) = newArr
		return
	}

	// 第一级指针 ( 可能是 slice 或者 pointer array )
	ptrLevel--
	if sd.dm.isArray && ptrLevel <= 0 {
		for i := 0; i < len(newArr); i++ {
			*((**T)(unsafe.Pointer(sd.dstPtr + uintptr(i*ptrByteSize)))) = &newArr[i]
		}
		for i := 0; i < len(sd.pl.nulPos); i++ {
			*((**T)(unsafe.Pointer(sd.dstPtr + uintptr(sd.pl.nulPos[i]*ptrByteSize)))) = nil
		}
		for i := size; i < sd.dm.arrLen; i++ {
			*((**T)(unsafe.Pointer(sd.dstPtr + uintptr(i*ptrByteSize)))) = nil
		}
		return
	}
	var newArrPtr1 []*T
	newArrPtr1 = make([]*T, size)
	for i := 0; i < len(newArr); i++ {
		newArrPtr1[i] = &newArr[i]
	}
	if ptrLevel <= 0 {
		for i := 0; i < len(sd.pl.nulPos); i++ {
			newArrPtr1[sd.pl.nulPos[i]] = nil
		}

		//*(sd.dst.(*[]*T)) = newArrPtr1
		*(*[]*T)(unsafe.Pointer(sd.dstPtr)) = newArrPtr1
		return
	}

	// 第二级指针
	ptrLevel--
	if sd.dm.isArray && ptrLevel <= 0 {
		for i := 0; i < len(newArrPtr1); i++ {
			*((***T)(unsafe.Pointer(sd.dstPtr + uintptr(i*ptrByteSize)))) = &newArrPtr1[i]
		}
		for i := 0; i < len(sd.pl.nulPos); i++ {
			*((***T)(unsafe.Pointer(sd.dstPtr + uintptr(sd.pl.nulPos[i]*ptrByteSize)))) = nil
		}
		for i := size; i < sd.dm.arrLen; i++ {
			*((***T)(unsafe.Pointer(sd.dstPtr + uintptr(i*ptrByteSize)))) = nil
		}
		return
	}
	newArrPtr2 := make([]**T, size)
	for i := 0; i < len(newArrPtr1); i++ {
		newArrPtr2[i] = &newArrPtr1[i]
	}
	if ptrLevel <= 0 {
		for i := 0; i < len(sd.pl.nulPos); i++ {
			newArrPtr2[sd.pl.nulPos[i]] = nil
		}

		//*(sd.dst.(*[]**T)) = newArrPtr2
		*(*[]**T)(unsafe.Pointer(sd.dstPtr)) = newArrPtr2
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
//		if sd.dm.isArray {
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
// 三种特殊类型单独处理，因为和number不同，这里的几种不存在类型转换，固单独处理
func sliceSetNotNum[T bool | string | any](val []T, sd *subDecode) {
	ptrLevel := sd.dm.ptrLevel

	list := make([]T, len(val))
	copy(list, val)

	// 这里只可能是slice
	if ptrLevel <= 0 {
		if len(sd.pl.nulPos) > 0 {
			dstSnap := []T{}
			*(*reflect.SliceHeader)(unsafe.Pointer(&dstSnap)) = *(*reflect.SliceHeader)(unsafe.Pointer(sd.dstPtr))
			for i := 0; i < len(sd.pl.nulPos); i++ {
				idx := sd.pl.nulPos[i]
				if idx >= len(dstSnap) {
					break
				}
				list[idx] = dstSnap[idx]
			}
		}
		//*(sd.dst.(*[]T)) = list
		*(*[]T)(unsafe.Pointer(sd.dstPtr)) = list
		return
	}

	// 一级指针
	ptrLevel--
	ret1 := copySlice[T](sd, ptrLevel, list)
	if ret1 == nil {
		return
	}

	// 二级指针
	ptrLevel--
	ret2 := copySlice[*T](sd, ptrLevel, ret1)
	if ret2 == nil {
		return
	}

	// 三级指针
	ptrLevel--
	_ = copySlice[**T](sd, ptrLevel, ret2)
	return
}

func copySlice[T bool | *bool | **bool | string | *string | **string | any | *any | **any](sd *subDecode, ptrLevel uint8, list []T) []*T {
	size := len(list)

	newList := make([]*T, size)
	for i := 0; i < size; i++ {
		newList[i] = &list[i]
	}

	if ptrLevel <= 0 {
		for i := 0; i < len(sd.pl.nulPos); i++ {
			newList[sd.pl.nulPos[i]] = nil
		}
		// 只能是指针数组才可能到这里的逻辑
		if sd.dm.isArray {
			dstSnap := []*T{}
			bh := (*reflect.SliceHeader)(unsafe.Pointer(&dstSnap))
			bh.Data, bh.Len, bh.Cap = sd.dstPtr, size, size
			copy(dstSnap, newList)
			// array，没有匹配到值的项，给初始化为nil
			for i := size; i < sd.dm.arrLen; i++ {
				*((**T)(unsafe.Pointer(sd.dstPtr + uintptr(i*ptrByteSize)))) = nil
			}
		} else {
			//*(sd.dst.(*[]*T)) = newList
			*(*[]*T)(unsafe.Pointer(sd.dstPtr)) = newList
		}
		return nil
	} else {
		return newList
	}
}
