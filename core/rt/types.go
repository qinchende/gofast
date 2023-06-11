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

	//nolint:structcheck,unused
	MapIter struct {
		key         unsafe.Pointer
		elem        unsafe.Pointer
		t           unsafe.Pointer
		h           unsafe.Pointer
		buckets     unsafe.Pointer
		bptr        unsafe.Pointer
		overflow    unsafe.Pointer
		oldoverflow unsafe.Pointer
		startBucket uintptr
		offset      uint8
		wrapped     bool
		B           uint8
		i           uint8
		bucket      uintptr
		checkBucket uintptr
	}
)
