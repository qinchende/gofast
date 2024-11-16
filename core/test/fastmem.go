package test

//
//import (
//	"unsafe"
//)
//
////go:nosplit
//func Get16(v []byte) int16 {
//	return *(*int16)((*GoSlice)(unsafe.Pointer(&v)).Ptr)
//}
//
////go:nosplit
//func Get32(v []byte) int32 {
//	return *(*int32)((*GoSlice)(unsafe.Pointer(&v)).Ptr)
//}
//
////go:nosplit
//func Get64(v []byte) int64 {
//	return *(*int64)((*GoSlice)(unsafe.Pointer(&v)).Ptr)
//}
//
////go:nosplit
//func Mem2Str(v []byte) (s string) {
//	(*GoString)(unsafe.Pointer(&s)).Len = (*GoSlice)(unsafe.Pointer(&v)).Len
//	(*GoString)(unsafe.Pointer(&s)).Ptr = (*GoSlice)(unsafe.Pointer(&v)).Ptr
//	return
//}
//
////go:nosplit
//func Str2Mem(s string) (v []byte) {
//	(*GoSlice)(unsafe.Pointer(&v)).Cap = (*GoString)(unsafe.Pointer(&s)).Len
//	(*GoSlice)(unsafe.Pointer(&v)).Len = (*GoString)(unsafe.Pointer(&s)).Len
//	(*GoSlice)(unsafe.Pointer(&v)).Ptr = (*GoString)(unsafe.Pointer(&s)).Ptr
//	return
//}
//
//func BytesFrom(p unsafe.Pointer, n int, c int) (r []byte) {
//	(*GoSlice)(unsafe.Pointer(&r)).Ptr = p
//	(*GoSlice)(unsafe.Pointer(&r)).Len = n
//	(*GoSlice)(unsafe.Pointer(&r)).Cap = c
//	return
//}
//
////func FuncAddr(f interface{}) unsafe.Pointer {
////	if vv := UnpackEface(f); vv.TypePtr.Kind() != reflect.Func {
////		panic("f is not a function")
////	} else {
////		return *(*unsafe.Pointer)(vv.DataPtr)
////	}
////}
//
//func IndexChar(src string, index int) unsafe.Pointer {
//	return unsafe.Pointer(uintptr((*GoString)(unsafe.Pointer(&src)).Ptr) + uintptr(index))
//}
//
//func IndexByte(ptr []byte, index int) unsafe.Pointer {
//	return unsafe.Pointer(uintptr((*GoSlice)(unsafe.Pointer(&ptr)).Ptr) + uintptr(index))
//}
//
////go:nosplit
//func GuardSlice(buf *[]byte, n int) {
//	c := cap(*buf)
//	l := len(*buf)
//	if c-l < n {
//		c = c>>1 + n + l
//		if c < 32 {
//			c = 32
//		}
//		tmp := make([]byte, l, c)
//		copy(tmp, *buf)
//		*buf = tmp
//	}
//}
//
////go:nosplit
//func Ptr2SlicePtr(s unsafe.Pointer, l int, c int) unsafe.Pointer {
//	slice := &GoSlice{
//		Ptr: s,
//		Len: l,
//		Cap: c,
//	}
//	return unsafe.Pointer(slice)
//}
//
////go:nosplit
//func StrPtr(s string) unsafe.Pointer {
//	return (*GoString)(unsafe.Pointer(&s)).Ptr
//}
//
////go:nosplit
//func StrFrom(p unsafe.Pointer, n int64) (s string) {
//	(*GoString)(unsafe.Pointer(&s)).Ptr = p
//	(*GoString)(unsafe.Pointer(&s)).Len = int(n)
//	return
//}
