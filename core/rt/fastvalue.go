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

type GoMap struct {
	Count      int
	Flags      uint8
	B          uint8
	Overflow   uint16
	Hash0      uint32
	Buckets    unsafe.Pointer
	OldBuckets unsafe.Pointer
	Evacuate   uintptr
	Extra      unsafe.Pointer
}

type GoMapIter struct {
	K           unsafe.Pointer
	V           unsafe.Pointer
	T           *GoMapType
	H           *GoMap
	Buckets     unsafe.Pointer
	Bptr        *unsafe.Pointer
	Overflow    *[]unsafe.Pointer
	OldOverflow *[]unsafe.Pointer
	StartBucket uintptr
	Offset      uint8
	Wrapped     bool
	B           uint8
	I           uint8
	Bucket      uintptr
	CheckBucket uintptr
}

func (self EFace) Pack() (v interface{}) {
	*(*EFace)(unsafe.Pointer(&v)) = self
	return
}

type GoPtrType struct {
	GoType
	Elem *GoType
}

type GoMapType struct {
	GoType
	Key        *GoType
	Elem       *GoType
	Bucket     *GoType
	Hasher     func(unsafe.Pointer, uintptr) uintptr
	KeySize    uint8
	ElemSize   uint8
	BucketSize uint16
	Flags      uint32
}

func (self *GoMapType) IndirectElem() bool {
	return self.Flags&2 != 0
}

type GoStructType struct {
	GoType
	Pkg    *byte
	Fields []GoStructField
}

type GoStructField struct {
	Name     *byte
	Type     *GoType
	OffEmbed uintptr
}

type GoInterfaceType struct {
	GoType
	PkgPath *byte
	Methods []GoInterfaceMethod
}

type GoInterfaceMethod struct {
	Name int32
	Type int32
}

type GoSlice struct {
	Ptr unsafe.Pointer
	Len int
	Cap int
}

type GoString struct {
	Ptr unsafe.Pointer
	Len int
}

func PtrElem(t *GoType) *GoType {
	return (*GoPtrType)(unsafe.Pointer(t)).Elem
}

func MapType(t *GoType) *GoMapType {
	return (*GoMapType)(unsafe.Pointer(t))
}

func IfaceType(t *GoType) *GoInterfaceType {
	return (*GoInterfaceType)(unsafe.Pointer(t))
}

func UnpackType(t reflect.Type) *GoType {
	return (*GoType)((*IFace)(unsafe.Pointer(&t)).DataPtr)
}

func UnpackEface(v interface{}) EFace {
	return *(*EFace)(unsafe.Pointer(&v))
}

func UnpackIface(v interface{}) IFace {
	return *(*IFace)(unsafe.Pointer(&v))
}

func findReflectRtypeItab() *GoItab {
	v := reflect.TypeOf(struct{}{})
	return (*IFace)(unsafe.Pointer(&v)).ItabPtr
}
