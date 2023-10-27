package dts

import "unsafe"

func (fa *fieldAttr) intScanner(sPtr unsafe.Pointer) any {
	return (*int)(unsafe.Pointer(uintptr(sPtr) + fa.Offset))
}
