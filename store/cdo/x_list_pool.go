package cdo

//
//import (
//	"github.com/qinchende/gofast/core/pool"
//	"github.com/qinchende/gofast/core/rt"
//	"golang.org/x/exp/constraints"
//	"reflect"
//	"unsafe"
//)
//
//// 解析Slice数据的时候，先用缓冲池装载，等全部解析完成之后再从缓冲池中提取结果。
//type listPool struct {
//	nulPos []int // 用来记录JSON中某一项为 null 的索引位置
//
//	bufInt []int
//	bufI8  []int8
//	bufI16 []int16
//	bufI32 []int32
//	bufI64 []int64
//
//	bufUint []uint
//	bufU8   []uint8
//	bufU16  []uint16
//	bufU32  []uint32
//	bufU64  []uint64
//
//	bufF32 []float32
//	bufF64 []float64
//
//	bufStr []string
//	bufBol []bool
//	bufAny []any
//
//	_memBytes *[]byte
//}
//
//// 默认值
//var _listPoolInitializer listPool
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (d *subDecode) resetListPool() {
//	// 如果是数组，本来就具备内存空间，不需要借助缓存
//	if d.dm.isArrBind {
//		return
//	}
//
//	pl := jdeBufPool.Get().(*listPool)
//	// 获取缓存内存空间
//	pl._memBytes = pool.GetBytes()
//	shMem := (*rt.SliceHeader)(unsafe.Pointer(pl._memBytes))
//
//	// 当 slice 项是下面几种基本数据类型的时候，自动准备一定大小的空间
//	var sh *rt.SliceHeader
//	switch d.dm.itemKind {
//	default:
//	case reflect.Int:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufInt))
//		sh.Cap = shMem.Cap / 8
//	case reflect.Int8:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufI8))
//		sh.Cap = shMem.Cap
//	case reflect.Int16:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufI16))
//		sh.Cap = shMem.Cap / 2
//	case reflect.Int32:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufI32))
//		sh.Cap = shMem.Cap / 4
//	case reflect.Int64:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufI64))
//		sh.Cap = shMem.Cap / 8
//
//	case reflect.Uint:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufUint))
//		sh.Cap = shMem.Cap / 8
//	case reflect.Uint8:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufU8))
//		sh.Cap = shMem.Cap
//	case reflect.Uint16:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufU16))
//		sh.Cap = shMem.Cap / 2
//	case reflect.Uint32:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufU32))
//		sh.Cap = shMem.Cap / 4
//	case reflect.Uint64:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufU64))
//		sh.Cap = shMem.Cap / 8
//
//	case reflect.Float32:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufF32))
//		sh.Cap = shMem.Cap / 4
//	case reflect.Float64:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufF64))
//		sh.Cap = shMem.Cap / 8
//
//	case reflect.String:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufStr))
//		sh.Cap = shMem.Cap / 16
//	case reflect.Bool:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufBol))
//		sh.Cap = shMem.Cap
//	case reflect.Interface:
//		sh = (*rt.SliceHeader)(unsafe.Pointer(&pl.bufAny))
//		sh.Cap = shMem.Cap / 16
//	}
//	sh.DataPtr = shMem.DataPtr
//	sh.Len = 0
//
//	d.pl = pl
//}
//
//func (d *subDecode) flushListPool() {
//	// 如果是定长数组，不会用到缓冲池，不需要转储
//	if d.dm.isArrBind {
//		return
//	}
//
//	switch d.dm.itemKind {
//	default:
//	case reflect.Int:
//		flushNoCast[int](d, d.pl.bufInt)
//	case reflect.Int8:
//		flushNoCast[int8](d, d.pl.bufI8)
//	case reflect.Int16:
//		flushNoCast[int16](d, d.pl.bufI16)
//	case reflect.Int32:
//		flushNoCast[int32](d, d.pl.bufI32)
//	case reflect.Int64:
//		flushNoCast[int64](d, d.pl.bufI64)
//
//	case reflect.Uint:
//		flushNoCast[uint](d, d.pl.bufUint)
//	case reflect.Uint8:
//		flushNoCast[uint8](d, d.pl.bufU8)
//	case reflect.Uint16:
//		flushNoCast[uint16](d, d.pl.bufU16)
//	case reflect.Uint32:
//		flushNoCast[uint32](d, d.pl.bufU32)
//	case reflect.Uint64:
//		flushNoCast[uint64](d, d.pl.bufU64)
//
//	case reflect.Float32:
//		flushNoCast[float32](d, d.pl.bufF32)
//	case reflect.Float64:
//		flushNoCast[float64](d, d.pl.bufF64)
//
//	case reflect.Bool:
//		flushNoCast[bool](d, d.pl.bufBol)
//	case reflect.String:
//		flushNoCast[string](d, d.pl.bufStr)
//	case reflect.Interface:
//		flushNoCast[any](d, d.pl.bufAny)
//	}
//	//case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
//	// 上面这几种情况，通过特殊方法处理
//
//	// 回收数组的内存空间
//	// Note: 解析的工程中slice可能要扩容，这里回收的内存是原始内存还是扩容后的内存呢？
//	pool.FreeBytes(d.pl._memBytes)
//
//	// Reset pl 对象
//	if cap(d.pl.nulPos) > 0 {
//		// 保留已分配的内存
//		tp := d.pl.nulPos
//		*d.pl = _listPoolInitializer
//		d.pl.nulPos = tp[0:0]
//	} else {
//		*d.pl = _listPoolInitializer
//	}
//
//	// 用完了就归还
//	jdeBufPool.Put(d.pl)
//	d.pl = nil
//}
//
//func flushNoCast[T any](d *subDecode, val []T) {
//	// 必须先Copy数据，才能使用
//	values := make([]T, len(val))
//	copy(values, val)
//	listSetValues[T](d, values)
//}
//
//func flushCast[T constraints.Integer | constraints.Float, T2 int64 | uint64 | float64](d *subDecode, val []T2) {
//	values := make([]T, len(val))
//	for i := 0; i < len(values); i++ {
//		values[i] = T(val[i])
//	}
//	listSetValues[T](d, values)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// NOTE: 目前本解码方案只支持指针层级在三级以内的基础数据类型（实际应用基本上也不会有层级太多的情况）
//func listSetValues[T any](d *subDecode, values []T) {
//	ptrLevel := d.dm.ptrLevel
//
//	// 这里只可能是slice，因为array的ptrLevel不可能是0（这种情况直接绑定结果了，不用缓冲池）
//	if ptrLevel == 0 {
//		if len(d.pl.nulPos) > 0 {
//			oriSlice := *(*[]T)(d.dstPtr)
//			for i := 0; i < len(d.pl.nulPos); i++ {
//				idx := d.pl.nulPos[i]
//				if idx >= len(oriSlice) {
//					break
//				}
//				values[idx] = oriSlice[idx]
//			}
//		}
//		*(*[]T)(d.dstPtr) = values
//		return
//	}
//
//	// 一级指针
//	ptrLevel--
//	ret1 := copySlice[T](d, ptrLevel, values)
//	if ret1 == nil {
//		return
//	}
//
//	// 二级指针
//	ptrLevel--
//	ret2 := copySlice[*T](d, ptrLevel, ret1)
//	if ret2 == nil {
//		return
//	}
//
//	// 三级指针
//	ptrLevel--
//	_ = copySlice[**T](d, ptrLevel, ret2)
//	return
//}
//
//func copySlice[T any | *any | **any](d *subDecode, ptrLevel uint8, sList []T) []*T {
//	size := len(sList)
//
//	// 如果是ptr类型的array，而且已到最后一级指针
//	if ptrLevel == 0 && d.dm.isArray {
//		var oriArr []*T
//		bh := (*rt.SliceHeader)(unsafe.Pointer(&oriArr))
//		bh.DataPtr, bh.Len, bh.Cap = d.dstPtr, d.dm.arrLen, d.dm.arrLen
//
//		for i := 0; i < size; i++ {
//			oriArr[i] = &sList[i]
//		}
//		for i := size; i < d.dm.arrLen; i++ {
//			oriArr[i] = nil // 此时array item是指针，给剩余的item重置为nil
//		}
//		for i := 0; i < len(d.pl.nulPos); i++ {
//			oriArr[d.pl.nulPos[i]] = nil
//		}
//		return nil
//	}
//
//	newList := make([]*T, size)
//	for i := 0; i < size; i++ {
//		newList[i] = &sList[i]
//	}
//
//	if ptrLevel == 0 {
//		for i := 0; i < len(d.pl.nulPos); i++ {
//			newList[d.pl.nulPos[i]] = nil
//		}
//		*(*[]*T)(d.dstPtr) = newList
//		return nil
//	}
//	return newList
//}
