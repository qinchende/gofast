package test

//
//import (
//	"github.com/qinchende/gofast/core/rt"
//	"reflect"
//	"unsafe"
//)
//
//type GoStructType struct {
//	rt.GoType
//	Pkg    *byte
//	Fields []GoStructField
//}
//
//type GoStructField struct {
//	Name     *byte
//	Type     *rt.GoType
//	OffEmbed uintptr
//}
//
//type GoInterfaceType struct {
//	rt.GoType
//	PkgPath *byte
//	Methods []GoInterfaceMethod
//}
//
//type GoInterfaceMethod struct {
//	Name int32
//	Type int32
//}
//
//type GoSlice struct {
//	Ptr unsafe.Pointer
//	Len int
//	Cap int
//}
//
//type GoString struct {
//	Ptr unsafe.Pointer
//	Len int
//}
//
//func PtrElem(t *rt.GoType) *rt.GoType {
//	return (*rt.GoPtrType)(unsafe.Pointer(t)).Elem
//}
//
//func MapType(t *rt.GoType) *rt.GoMapType {
//	return (*rt.GoMapType)(unsafe.Pointer(t))
//}
//
//func IfaceType(t *rt.GoType) *GoInterfaceType {
//	return (*GoInterfaceType)(unsafe.Pointer(t))
//}
//
//func UnpackType(t reflect.Type) *rt.GoType {
//	return (*rt.GoType)((*rt.IFace)(unsafe.Pointer(&t)).DataPtr)
//}
//
//func UnpackEface(v interface{}) rt.EFace {
//	return *(*rt.EFace)(unsafe.Pointer(&v))
//}
//
//func UnpackIface(v interface{}) rt.IFace {
//	return *(*rt.IFace)(unsafe.Pointer(&v))
//}
