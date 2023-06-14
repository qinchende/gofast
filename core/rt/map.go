package rt

import "unsafe"

//go:linkname MapIterInit runtime.mapiterinit
//go:noescape
func MapIterInit(mapType *GoType, m unsafe.Pointer, it *GoMapIter)

//go:linkname MapIterKey reflect.mapiterkey
//go:noescape
func MapIterKey(it *GoMapIter) unsafe.Pointer

//go:linkname MapIterNext reflect.mapiternext
//go:noescape
func MapIterNext(it *GoMapIter)

//go:linkname MapLen reflect.maplen
//go:noescape
func MapLen(m unsafe.Pointer) int
