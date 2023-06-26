package rt

import (
	"reflect"
	"unsafe"
)

var (
	reflectRtypeItab = findReflectRtypeItab()
)

const (
	F_direct    = 1 << 5
	F_kind_mask = (1 << 5) - 1
)

type (
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

	GoPtrType struct {
		GoType
		Elem *GoType
	}

	GoItab struct {
		it unsafe.Pointer
		Vt *GoType
		hv uint32
		_  [4]byte
		fn [1]uintptr
	}
)

func (self *GoType) Kind() reflect.Kind {
	return reflect.Kind(self.KindFlags & F_kind_mask)
}

func (self *GoType) Pack() (t reflect.Type) {
	(*IFace)(unsafe.Pointer(&t)).ItabPtr = reflectRtypeItab
	(*IFace)(unsafe.Pointer(&t)).DataPtr = unsafe.Pointer(self)
	return
}

func (self *GoType) String() string {
	return self.Pack().String()
}

func (self *GoType) Indirect() bool {
	return self.KindFlags&F_direct == 0
}

func findReflectRtypeItab() *GoItab {
	v := reflect.TypeOf(struct{}{})
	return (*IFace)(unsafe.Pointer(&v)).ItabPtr
}
