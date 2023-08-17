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

	// placeholder interface (to get memory value)
	AFace struct {
		TypePtr *TypeAgent
		DataPtr unsafe.Pointer
	}

	// empty interface (has no function)
	EFace struct {
		TypePtr *GoType
		DataPtr unsafe.Pointer
	}

	// typed interface (define some function)
	IFace struct {
		ItabPtr *GoItab
		DataPtr unsafe.Pointer
	}
)

func (self EFace) Pack() (v interface{}) {
	*(*EFace)(unsafe.Pointer(&v)) = self
	return
}
