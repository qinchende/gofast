package rt

import (
	"unsafe"
)

type (
	StringHeader struct {
		DataPtr unsafe.Pointer
		Len     int
	}

	SliceHeader struct {
		DataPtr unsafe.Pointer
		Len     int
		Cap     int
	}

	TypeAgent struct{}

	// 0. placeholder interface (to get memory value)
	// 为了得到 interface 变量内存对应的两个字段值
	AFace struct {
		TypePtr *TypeAgent
		DataPtr unsafe.Pointer
	}

	// 1. empty interface (has no function)
	EFace struct {
		TypePtr *GoType
		DataPtr unsafe.Pointer
	}

	// 2. typed interface (define some function)
	IFace struct {
		ItabPtr *GoItab
		DataPtr unsafe.Pointer
	}
)

func (self EFace) Pack() (v interface{}) {
	*(*EFace)(unsafe.Pointer(&v)) = self
	return
}
