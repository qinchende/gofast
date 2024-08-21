package rt

import "unsafe"

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(unsafe.Pointer) unsafe.Pointer
func UnsafeNew(p unsafe.Pointer) unsafe.Pointer {
	return unsafe_New(p)
}

//go:linkname makemap reflect.makemap
func makemap(unsafe.Pointer, int) unsafe.Pointer
func MakeMap(p unsafe.Pointer, n int) unsafe.Pointer {
	return makemap(p, n)
}

//go:linkname mapassign_faststr runtime.mapassign_faststr
//go:noescape
func mapassign_faststr(t unsafe.Pointer, m unsafe.Pointer, s string) unsafe.Pointer
func MapAssignFastStr(t unsafe.Pointer, m unsafe.Pointer, s string) unsafe.Pointer {
	return mapassign_faststr(t, m, s)
}

//go:linkname mapassign reflect.mapassign
//go:noescape
func mapassign(t unsafe.Pointer, m unsafe.Pointer, k, v unsafe.Pointer)
func MapAssign(t unsafe.Pointer, m unsafe.Pointer, k, v unsafe.Pointer) {
	mapassign(t, m, k, v)
}

// mapIter
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//go:linkname mapiterinit runtime.mapiterinit
func mapiterinit(t, h, it unsafe.Pointer)
func MapIterInit(t, h unsafe.Pointer, it *HIter) {
	mapiterinit(t, h, unsafe.Pointer(it))
}

//go:linkname mapiterkey reflect.mapiterkey
func mapiterkey(it unsafe.Pointer) unsafe.Pointer
func MapIterKey(it *HIter) unsafe.Pointer {
	return mapiterkey(unsafe.Pointer(it))
}

//go:linkname mapiterelem reflect.mapiterelem
func mapiterelem(it unsafe.Pointer) unsafe.Pointer
func MapIterValue(it *HIter) unsafe.Pointer {
	return mapiterelem(unsafe.Pointer(it))
}

//go:linkname mapiternext reflect.mapiternext
func mapiternext(it unsafe.Pointer)
func MapIterNext(it *HIter) {
	mapiternext(unsafe.Pointer(it))
}

//go:linkname maplen reflect.maplen
func maplen(m unsafe.Pointer) int
func MapLen(m unsafe.Pointer) int {
	return maplen(m)
}
