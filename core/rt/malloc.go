package rt

import (
	"unsafe"
)

func SliceAutoGrow(sh *SliceHeader, itemSize int) {
	oldLen := sh.Len
	// 如果已经没有待用内存空间，就执行扩容算法
	if oldLen == sh.Cap {
		var oldMem = []byte{}
		old := (*SliceHeader)(unsafe.Pointer(&oldMem))
		old.DataPtr, old.Len, old.Cap = sh.DataPtr, sh.Len*itemSize, sh.Cap*itemSize

		// TODO: 需要研究更高效的扩容算法
		newLen := int(float64(sh.Cap)*1.6) + 5 // 一种简易的动态扩容算法
		//fmt.Printf("growing len: %d, cap: %d \n\n", sh.Len*itemSize, newLen*itemSize)

		// 不管什么类型的数据，底层存放在内存字节序列当中，我们只需要申请足量的字节序列，
		// 之后想让这段序列代表啥数据类型都行。
		bsPtr := (*[]byte)(unsafe.Pointer(sh))
		*bsPtr = make([]byte, sh.Len*itemSize, newLen*itemSize)
		copy(*bsPtr, oldMem)

		sh.Len, sh.Cap = oldLen, newLen
	}
}

// 返回下一个值内存空间的地址
func SliceNextItem(sh *SliceHeader, itemSize int) (ptr unsafe.Pointer) {
	SliceAutoGrow(sh, itemSize)
	ptr = unsafe.Add(sh.DataPtr, sh.Len*itemSize) // unsafe.Pointer(uintptr(sh.DataPtr) + uintptr(sh.Len*itemSize))
	sh.Len++
	return
}

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
