package jde

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

// TODO: buffer pool 需要有个机制，释放那些某次偶发申请太大的buffer，而导致长时间不释放的问题
type listPool struct {
	bufI64 []int64
	bufU64 []uint64
	bufF64 []float64
	bufStr []string
	bufBol []bool
	bufAny []any
	nulPos []int // 指针类型值，可能是nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) resetListPool() {
	if sd.dm.isArrBind {
		return
	}

	// 获取缓存空间
	sd.pl = jdeBufPool.Get().(*listPool)

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
		flushCast[int, int64](sd, sd.pl.bufI64)
	case reflect.Int8:
		flushCast[int8, int64](sd, sd.pl.bufI64)
	case reflect.Int16:
		flushCast[int16, int64](sd, sd.pl.bufI64)
	case reflect.Int32:
		flushCast[int32, int64](sd, sd.pl.bufI64)
	case reflect.Int64:
		flushNoCast[int64](sd, sd.pl.bufI64)

	case reflect.Uint:
		flushCast[int, uint64](sd, sd.pl.bufU64)
	case reflect.Uint8:
		flushCast[uint8, uint64](sd, sd.pl.bufU64)
	case reflect.Uint16:
		flushCast[uint16, uint64](sd, sd.pl.bufU64)
	case reflect.Uint32:
		flushCast[uint32, uint64](sd, sd.pl.bufU64)
	case reflect.Uint64:
		flushNoCast[uint64](sd, sd.pl.bufU64)

	case reflect.Float32:
		flushCast[float32, float64](sd, sd.pl.bufF64)
	case reflect.Float64:
		flushNoCast[float64](sd, sd.pl.bufF64)

	case reflect.Bool:
		flushNoCast[bool](sd, sd.pl.bufBol)
	case reflect.String:
		flushNoCast[string](sd, sd.pl.bufStr)
	case reflect.Interface:
		flushNoCast[any](sd, sd.pl.bufAny)
	}
	//case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
	// 上面这几种情况，通过特殊方法处理

	// 用完了就归还
	jdeBufPool.Put(sd.pl)
	sd.pl = nil
}

func flushNoCast[T any](sd *subDecode, val []T) {
	// 必须先Copy数据，才能使用
	values := make([]T, len(val))
	copy(values, val)
	listSetValues[T](sd, values)
}

func flushCast[T constraints.Integer | constraints.Float, T2 int64 | uint64 | float64](sd *subDecode, val []T2) {
	values := make([]T, len(val))
	for i := 0; i < len(values); i++ {
		values[i] = T(val[i])
	}
	listSetValues[T](sd, values)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// NOTE: 目前本解码方案只支持指针层级在三级以内的基础数据类型（实际应用基本上也不会有层级太多的情况）
func listSetValues[T any](sd *subDecode, values []T) {
	ptrLevel := sd.dm.ptrLevel

	// 这里只可能是slice，因为array的ptrLevel不可能是0（这种情况直接绑定结果了，不用缓冲池）
	if ptrLevel == 0 {
		if len(sd.pl.nulPos) > 0 {
			oriSlice := *(*[]T)(sd.dstPtr)
			for i := 0; i < len(sd.pl.nulPos); i++ {
				idx := sd.pl.nulPos[i]
				if idx >= len(oriSlice) {
					break
				}
				values[idx] = oriSlice[idx]
			}
		}
		*(*[]T)(sd.dstPtr) = values
		return
	}

	// 一级指针
	ptrLevel--
	ret1 := copySlice[T](sd, ptrLevel, values)
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

func copySlice[T any | *any | **any](sd *subDecode, ptrLevel uint8, sList []T) []*T {
	size := len(sList)

	// 如果是ptr类型的array，而且已到最后一级指针
	if ptrLevel == 0 && sd.dm.isArray {
		oriArr := []*T{}
		bh := (*reflect.SliceHeader)(unsafe.Pointer(&oriArr))
		bh.Data, bh.Len, bh.Cap = uintptr(sd.dstPtr), sd.dm.arrLen, sd.dm.arrLen

		for i := 0; i < size; i++ {
			oriArr[i] = &sList[i]
		}
		for i := size; i < sd.dm.arrLen; i++ {
			oriArr[i] = nil // 此时array item是指针，给剩余的item重置为nil
		}
		for i := 0; i < len(sd.pl.nulPos); i++ {
			oriArr[sd.pl.nulPos[i]] = nil
		}
		return nil
	}

	newList := make([]*T, size)
	for i := 0; i < size; i++ {
		newList[i] = &sList[i]
	}

	if ptrLevel == 0 {
		for i := 0; i < len(sd.pl.nulPos); i++ {
			newList[sd.pl.nulPos[i]] = nil
		}
		*(*[]*T)(sd.dstPtr) = newList
		return nil
	}
	return newList
}
