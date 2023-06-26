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
	AFace     struct {
		TypePtr *TypeAgent
		DataPtr unsafe.Pointer
	}

	EFace struct {
		TypePtr *GoType
		DataPtr unsafe.Pointer
	}

	IFace struct {
		ItabPtr *GoItab
		DataPtr unsafe.Pointer
	}
)

func (self EFace) Pack() (v interface{}) {
	*(*EFace)(unsafe.Pointer(&v)) = self
	return
}
