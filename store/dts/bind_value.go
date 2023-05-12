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

func (ss *StructSchema) BindUint(ptr uintptr, idx int, val uint64) {
	offset := ss.fieldsOpts[idx].sField.Offset
	*(*uint)(unsafe.Pointer(ptr + offset)) = uint(val)
}

func (ss *StructSchema) BindFloat(ptr uintptr, idx int, val float64) {
	offset := ss.fieldsOpts[idx].sField.Offset
	*(*float64)(unsafe.Pointer(ptr + offset)) = val
}
