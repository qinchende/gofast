package dts

import "unsafe"

func (ss *StructSchema) BindString(ptr uintptr, idx int, val string) {
	offset := ss.fieldsOpts[idx].sField.Offset
	*(*string)(unsafe.Pointer(ptr + offset)) = val
}

func (ss *StructSchema) BindBool(ptr uintptr, idx int, val bool) {
	offset := ss.fieldsOpts[idx].sField.Offset
	*(*bool)(unsafe.Pointer(ptr + offset)) = val
}

func (ss *StructSchema) BindInt(ptr uintptr, idx int, val int64) {
	offset := ss.fieldsOpts[idx].sField.Offset
	*(*int)(unsafe.Pointer(ptr + offset)) = int(val)
}

func (ss *StructSchema) BindFloat(ptr uintptr, idx int, val float64) {
	offset := ss.fieldsOpts[idx].sField.Offset
	*(*float64)(unsafe.Pointer(ptr + offset)) = val
}

//func (sd *structPost) setStringByIndex(idx int, val string) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetString(val)
//}
//
//func (sd *structPost) setBoolByIndex(idx int, val bool) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetBool(val)
//}
//
//func (sd *structPost) setIntByIndex(idx int, val int64) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetInt(val)
//}
//
//func (sd *structPost) setFloatByIndex(idx int, val float64) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetFloat(val)
//}
