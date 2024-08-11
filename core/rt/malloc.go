package rt

import (
	"reflect"
	"unsafe"
)

// TODO: 需要研究更高效的扩容算法
// 一种简易的动态扩容算法
func growSize(cap int) int {
	return int(float64(cap)*1.6) + 5
}

func SliceAutoGrow(sh *SliceHeader, itemSize int) {
	oldLen := sh.Len
	// 如果已经没有待用内存空间，就执行扩容算法
	if oldLen == sh.Cap {
		var oldMem []byte
		old := (*SliceHeader)(unsafe.Pointer(&oldMem))
		old.DataPtr, old.Len, old.Cap = sh.DataPtr, sh.Len*itemSize, sh.Cap*itemSize

		newLen := growSize(sh.Cap)
		// 不管什么类型的数据，底层存放在内存字节序列当中，我们只需要申请足量的字节序列，
		// 之后想让这段序列代表啥数据类型都行。
		bsPtr := (*[]byte)(unsafe.Pointer(sh))
		*bsPtr = make([]byte, sh.Len*itemSize, newLen*itemSize)
		copy(*bsPtr, oldMem)

		sh.Len, sh.Cap = oldLen, newLen
	}
}

// 返回下一个值内存空间的地址
func SliceNextItem(sh *SliceHeader, itemSize int) unsafe.Pointer {
	SliceAutoGrow(sh, itemSize)
	ptr := unsafe.Add(sh.DataPtr, sh.Len*itemSize)
	sh.Len++
	return ptr
}

func SliceNextItemSafe(sh *SliceHeader, itemSize int, itemType reflect.Type) unsafe.Pointer {
	if sh.Cap == sh.Len {
		var oBytes []byte
		osh := (*SliceHeader)(unsafe.Pointer(&oBytes))
		osh.DataPtr, osh.Len = sh.DataPtr, sh.Len*itemSize
		osh.Cap = osh.Len

		newLen := growSize(sh.Cap)
		sh.DataPtr = reflect.MakeSlice(itemType, newLen, newLen).UnsafePointer()
		sh.Cap = newLen

		var nBytes []byte
		nsh := (*SliceHeader)(unsafe.Pointer(&nBytes))
		nsh.DataPtr, nsh.Len, nsh.Cap = sh.DataPtr, osh.Len, osh.Cap

		copy(nBytes, oBytes)
	}
	sh.Len++
	return unsafe.Add(sh.DataPtr, sh.Len*itemSize)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 为Slice对象分配足够的内存空间，并像Array一样，返回第一个值的地址
func SliceToArray(slicePtr unsafe.Pointer, itemSize int, sliceLen int) unsafe.Pointer {
	sh := (*SliceHeader)(slicePtr)
	if sh.Cap < sliceLen {
		newMem := make([]byte, itemSize*sliceLen)
		sh.DataPtr = (*SliceHeader)(unsafe.Pointer(&newMem)).DataPtr
		sh.Len, sh.Cap = sliceLen, sliceLen
	} else {
		sh.Len = sliceLen
	}
	return sh.DataPtr
}

// 为Slice对象分配足够的内存空间，并像Array一样，返回第一个值的地址
// 此版本要保证内存安全性
func SliceToArraySafe(slicePtr unsafe.Pointer, sliceLen int, itemType reflect.Type) unsafe.Pointer {
	sh := (*SliceHeader)(slicePtr)
	if sh.Cap < sliceLen {
		sh.DataPtr = reflect.MakeSlice(itemType, sliceLen, sliceLen).UnsafePointer()
		sh.Len, sh.Cap = sliceLen, sliceLen
	} else {
		sh.Len = sliceLen
	}
	return sh.DataPtr
}
