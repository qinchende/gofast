package jde

import (
	"github.com/qinchende/gofast/core/pool"
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

// 解析Slice数据的时候，先用缓冲池装载，等全部解析完成之后再从缓冲池中提取结果。
type listPool struct {
	nulPos []int // 用来记录JSON中某一项为 null 的索引位置

	bufInt []int
	bufI8  []int8
	bufI16 []int16
	bufI32 []int32
	bufI64 []int64

	bufUint []uint
	bufU8   []uint8
	bufU16  []uint16
	bufU32  []uint32
	bufU64  []uint64

	bufF32 []float32
	bufF64 []float64

	bufStr []string
	bufBol []bool
	bufAny []any

	_memBytes *[]byte
}

// 默认值
var _listPoolDefValue listPool

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) resetListPool() {
	// 如果是数组，本来就具备内存空间，不需要借助缓存
	if sd.dm.isArrBind {
		return
	}

	pl := jdeBufPool.Get().(*listPool)
	// 获取缓存内存空间
	pl._memBytes = pool.GetBytes()
	shMem := (*reflect.SliceHeader)(unsafe.Pointer(pl._memBytes))

	var sh *reflect.SliceHeader
	switch sd.dm.itemKind {
	case reflect.Int:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufInt))
		sh.Cap = shMem.Cap / 8
	case reflect.Int8:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufI8))
		sh.Cap = shMem.Cap
	case reflect.Int16:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufI16))
		sh.Cap = shMem.Cap / 2
	case reflect.Int32:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufI32))
		sh.Cap = shMem.Cap / 4
	case reflect.Int64:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufI64))
		sh.Cap = shMem.Cap / 8

	case reflect.Uint:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufUint))
		sh.Cap = shMem.Cap / 8
	case reflect.Uint8:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufU8))
		sh.Cap = shMem.Cap
	case reflect.Uint16:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufU16))
		sh.Cap = shMem.Cap / 2
	case reflect.Uint32:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufU32))
		sh.Cap = shMem.Cap / 4
	case reflect.Uint64:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufU64))
		sh.Cap = shMem.Cap / 8

	case reflect.Float32:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufF32))
		sh.Cap = shMem.Cap / 4
	case reflect.Float64:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufF64))
		sh.Cap = shMem.Cap / 8

	case reflect.String:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufStr))
		sh.Cap = shMem.Cap / 16
	case reflect.Bool:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufBol))
		sh.Cap = shMem.Cap
	case reflect.Interface:
		sh = (*reflect.SliceHeader)(unsafe.Pointer(&pl.bufAny))
		sh.Cap = shMem.Cap / 16
	}
	sh.Data = shMem.Data
	sh.Len = 0

	sd.pl = pl
}

func (sd *subDecode) flushListPool() {
	// 如果是定长数组，不会用到缓冲池，不需要转储
	if sd.dm.isArrBind {
		return
	}

	switch sd.dm.itemKind {
	case reflect.Int:
		flushNoCast[int](sd, sd.pl.bufInt)
	case reflect.Int8:
		flushNoCast[int8](sd, sd.pl.bufI8)
	case reflect.Int16:
		flushNoCast[int16](sd, sd.pl.bufI16)
	case reflect.Int32:
		flushNoCast[int32](sd, sd.pl.bufI32)
	case reflect.Int64:
		flushNoCast[int64](sd, sd.pl.bufI64)

	case reflect.Uint:
		flushNoCast[uint](sd, sd.pl.bufUint)
	case reflect.Uint8:
		flushNoCast[uint8](sd, sd.pl.bufU8)
	case reflect.Uint16:
		flushNoCast[uint16](sd, sd.pl.bufU16)
	case reflect.Uint32:
		flushNoCast[uint32](sd, sd.pl.bufU32)
	case reflect.Uint64:
		flushNoCast[uint64](sd, sd.pl.bufU64)

	case reflect.Float32:
		flushNoCast[float32](sd, sd.pl.bufF32)
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

	// 回收数组的内存空间
	pool.FreeBytes(sd.pl._memBytes)

	// Reset pl 对象
	if cap(sd.pl.nulPos) > 0 {
		// 保留已分配的内存
		tp := sd.pl.nulPos
		*sd.pl = _listPoolDefValue
		sd.pl.nulPos = tp[0:0]
	} else {
		*sd.pl = _listPoolDefValue
	}

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
