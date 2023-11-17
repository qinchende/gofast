package rt

import (
	"reflect"
	"unsafe"
)

func SliceAutoGrow(sh *reflect.SliceHeader, itemSize int) {
	oldLen := sh.Len
	// 如果已经没有待用内存空间，就执行扩容算法
	if oldLen == sh.Cap {
		var oldMem = []byte{}
		old := (*reflect.SliceHeader)(unsafe.Pointer(&oldMem))
		old.Data, old.Len, old.Cap = sh.Data, sh.Len*itemSize, sh.Cap*itemSize

		// TODO: 需要研究更高效的扩容算法
		newLen := int(float64(sh.Cap)*1.6) + 5 // 一种简易的动态扩容算法
		//fmt.Printf("growing len: %d, cap: %d \n\n", sh.Len*itemSize, newLen*itemSize)

		bsPtr := (*[]byte)(unsafe.Pointer(sh))
		*bsPtr = make([]byte, sh.Len*itemSize, newLen*itemSize)
		copy(*bsPtr, oldMem)

		sh.Len, sh.Cap = oldLen, newLen
	}
}

// 返回下一个值内存空间的地址
func SliceNextItem(sh *reflect.SliceHeader, itemSize int) (ptr unsafe.Pointer) {
	SliceAutoGrow(sh, itemSize)
	ptr = unsafe.Pointer(sh.Data + uintptr(sh.Len*itemSize))
	sh.Len++
	return
}
