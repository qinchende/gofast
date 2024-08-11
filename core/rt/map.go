package rt

import (
	"unsafe"
)

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
