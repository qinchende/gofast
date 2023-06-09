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

	GoType struct {
		Size       uintptr
		DataPtr    uintptr
		Hash       uint32
		Flags      uint8
		Align      uint8
		FieldAlign uint8
		KindFlags  uint8
		Traits     unsafe.Pointer
		GCData     *byte
		Str        int32
		PtrToSelf  int32
	}
	EFace struct {
		TypePtr *GoType
		DataPtr unsafe.Pointer
	}

	GoItab struct {
		it unsafe.Pointer
		Vt *GoType
		hv uint32
		_  [4]byte
		fn [1]uintptr
	}
	IFace struct {
		ItabPtr *GoItab
		DataPtr unsafe.Pointer
	}
)
