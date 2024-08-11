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
