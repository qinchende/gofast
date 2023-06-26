package test

//
//import (
//	"reflect"
//	"unsafe"
//)
//
//// +// go:linkname TypeAgent reflect.rtype
//// +// go:noescape
//type TypeAgent struct{}
//
//// Type representing reflect.rtype for noescape trick
////type TypeAgent struct{}
//
////go:linkname rtype_Align reflect.(*rtype).Align
////go:noescape
//func rtype_Align(*TypeAgent) int
//
//func (t *TypeAgent) Align() int {
//	return rtype_Align(t)
//}
//
////go:linkname rtype_FieldAlign reflect.(*rtype).FieldAlign
////go:noescape
//func rtype_FieldAlign(*TypeAgent) int
//
//func (t *TypeAgent) FieldAlign() int {
//	return rtype_FieldAlign(t)
//}
//
////go:linkname rtype_Method reflect.(*rtype).Method
////go:noescape
//func rtype_Method(*TypeAgent, int) reflect.Method
//
//func (t *TypeAgent) Method(a0 int) reflect.Method {
//	return rtype_Method(t, a0)
//}
//
////go:linkname rtype_MethodByName reflect.(*rtype).MethodByName
////go:noescape
//func rtype_MethodByName(*TypeAgent, string) (reflect.Method, bool)
//
//func (t *TypeAgent) MethodByName(a0 string) (reflect.Method, bool) {
//	return rtype_MethodByName(t, a0)
//}
//
////go:linkname rtype_NumMethod reflect.(*rtype).NumMethod
////go:noescape
//func rtype_NumMethod(*TypeAgent) int
//
//func (t *TypeAgent) NumMethod() int {
//	return rtype_NumMethod(t)
//}
//
////go:linkname rtype_Name reflect.(*rtype).Name
////go:noescape
//func rtype_Name(*TypeAgent) string
//
//func (t *TypeAgent) Name() string {
//	return rtype_Name(t)
//}
//
////go:linkname rtype_PkgPath reflect.(*rtype).PkgPath
////go:noescape
//func rtype_PkgPath(*TypeAgent) string
//
//func (t *TypeAgent) PkgPath() string {
//	return rtype_PkgPath(t)
//}
//
////go:linkname rtype_Size reflect.(*rtype).Size
//
////go:noescape
//func rtype_Size(*TypeAgent) uintptr
//
//func (t *TypeAgent) Size() uintptr {
//	return rtype_Size(t)
//}
//
////go:linkname rtype_String reflect.(*rtype).String
////go:noescape
//func rtype_String(*TypeAgent) string
//
//func (t *TypeAgent) String() string {
//	return rtype_String(t)
//}
//
////go:linkname rtype_Kind reflect.(*rtype).Kind
////go:noescape
//func rtype_Kind(*TypeAgent) reflect.Kind
//
//func (t *TypeAgent) Kind() reflect.Kind {
//	return rtype_Kind(t)
//}
//
////go:linkname rtype_Implements reflect.(*rtype).Implements
////go:noescape
//func rtype_Implements(*TypeAgent, reflect.Type) bool
//
//func (t *TypeAgent) Implements(u reflect.Type) bool {
//	return rtype_Implements(t, u)
//}
//
////go:linkname rtype_AssignableTo reflect.(*rtype).AssignableTo
////go:noescape
//func rtype_AssignableTo(*TypeAgent, reflect.Type) bool
//
//func (t *TypeAgent) AssignableTo(u reflect.Type) bool {
//	return rtype_AssignableTo(t, u)
//}
//
////go:linkname rtype_ConvertibleTo reflect.(*rtype).ConvertibleTo
////go:noescape
//func rtype_ConvertibleTo(*TypeAgent, reflect.Type) bool
//
//func (t *TypeAgent) ConvertibleTo(u reflect.Type) bool {
//	return rtype_ConvertibleTo(t, u)
//}
//
////go:linkname rtype_Comparable reflect.(*rtype).Comparable
////go:noescape
//func rtype_Comparable(*TypeAgent) bool
//
//func (t *TypeAgent) Comparable() bool {
//	return rtype_Comparable(t)
//}
//
////go:linkname rtype_Bits reflect.(*rtype).Bits
////go:noescape
//func rtype_Bits(*TypeAgent) int
//
//func (t *TypeAgent) Bits() int {
//	return rtype_Bits(t)
//}
//
////go:linkname rtype_ChanDir reflect.(*rtype).ChanDir
////go:noescape
//func rtype_ChanDir(*TypeAgent) reflect.ChanDir
//
//func (t *TypeAgent) ChanDir() reflect.ChanDir {
//	return rtype_ChanDir(t)
//}
//
////go:linkname rtype_IsVariadic reflect.(*rtype).IsVariadic
////go:noescape
//func rtype_IsVariadic(*TypeAgent) bool
//
//func (t *TypeAgent) IsVariadic() bool {
//	return rtype_IsVariadic(t)
//}
//
////go:linkname rtype_Elem reflect.(*rtype).Elem
////go:noescape
//func rtype_Elem(*TypeAgent) reflect.Type
//
//func (t *TypeAgent) Elem() *TypeAgent {
//	return Type2RType(rtype_Elem(t))
//}
//
////go:linkname rtype_Field reflect.(*rtype).Field
////go:noescape
//func rtype_Field(*TypeAgent, int) reflect.StructField
//
//func (t *TypeAgent) Field(i int) reflect.StructField {
//	return rtype_Field(t, i)
//}
//
////go:linkname rtype_FieldByIndex reflect.(*rtype).FieldByIndex
////go:noescape
//func rtype_FieldByIndex(*TypeAgent, []int) reflect.StructField
//
//func (t *TypeAgent) FieldByIndex(index []int) reflect.StructField {
//	return rtype_FieldByIndex(t, index)
//}
//
////go:linkname rtype_FieldByName reflect.(*rtype).FieldByName
////go:noescape
//func rtype_FieldByName(*TypeAgent, string) (reflect.StructField, bool)
//
//func (t *TypeAgent) FieldByName(name string) (reflect.StructField, bool) {
//	return rtype_FieldByName(t, name)
//}
//
////go:linkname rtype_FieldByNameFunc reflect.(*rtype).FieldByNameFunc
////go:noescape
//func rtype_FieldByNameFunc(*TypeAgent, func(string) bool) (reflect.StructField, bool)
//
//func (t *TypeAgent) FieldByNameFunc(match func(string) bool) (reflect.StructField, bool) {
//	return rtype_FieldByNameFunc(t, match)
//}
//
////go:linkname rtype_In reflect.(*rtype).In
////go:noescape
//func rtype_In(*TypeAgent, int) reflect.Type
//
//func (t *TypeAgent) In(i int) reflect.Type {
//	return rtype_In(t, i)
//}
//
////go:linkname rtype_Key reflect.(*rtype).Key
////go:noescape
//func rtype_Key(*TypeAgent) reflect.Type
//
//func (t *TypeAgent) Key() *TypeAgent {
//	return Type2RType(rtype_Key(t))
//}
//
////go:linkname rtype_Len reflect.(*rtype).Len
////go:noescape
//func rtype_Len(*TypeAgent) int
//
//func (t *TypeAgent) Len() int {
//	return rtype_Len(t)
//}
//
////go:linkname rtype_NumField reflect.(*rtype).NumField
////go:noescape
//func rtype_NumField(*TypeAgent) int
//
//func (t *TypeAgent) NumField() int {
//	return rtype_NumField(t)
//}
//
////go:linkname rtype_NumIn reflect.(*rtype).NumIn
////go:noescape
//func rtype_NumIn(*TypeAgent) int
//
//func (t *TypeAgent) NumIn() int {
//	return rtype_NumIn(t)
//}
//
////go:linkname rtype_NumOut reflect.(*rtype).NumOut
////go:noescape
//func rtype_NumOut(*TypeAgent) int
//
//func (t *TypeAgent) NumOut() int {
//	return rtype_NumOut(t)
//}
//
////go:linkname rtype_Out reflect.(*rtype).Out
////go:noescape
//func rtype_Out(*TypeAgent, int) reflect.Type
//
////go:linkname PtrTo reflect.(*rtype).ptrTo
////go:noescape
//func PtrTo(*TypeAgent) *TypeAgent
//
//func (t *TypeAgent) Out(i int) reflect.Type {
//	return rtype_Out(t, i)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
////go:linkname Rt_rtype reflect.rtype
//type Rt_rtype struct{}
//
////go:linkname rtype_common reflect.(*rtype).common
////go:noescape
//func rtype_common(*TypeAgent) *Rt_rtype
//
//func (t *TypeAgent) common() *Rt_rtype {
//	return rtype_common(t)
//}
//
////go:linkname Rt_uncommonType reflect.uncommonType
//type Rt_uncommonType struct{}
//
////go:linkname rtype_uncommon reflect.(*rtype).uncommon
////go:noescape
//func rtype_uncommon(*TypeAgent) *Rt_uncommonType
//
//func (t *TypeAgent) uncommon() *Rt_uncommonType {
//	return rtype_uncommon(t)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
////
////go:linkname IfaceIndir reflect.ifaceIndir
////go:noescape
//func IfaceIndir(*TypeAgent) bool
//
////go:linkname RType2Type reflect.toType
////go:noescape
//func RType2Type(t *TypeAgent) reflect.Type
//
//func Type2RType(t reflect.Type) *TypeAgent {
//	return (*TypeAgent)(((*AFace)(unsafe.Pointer(&t))).DataPtr)
//}
