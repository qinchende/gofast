package dts

//
//import (
//	"math"
//	"reflect"
//	"unsafe"
//)
//
//func (ss *StructSchema) BindString(ptr uintptr, idx int, v string) {
//	p := unsafe.Pointer(ptr + ss.FieldsOffset[idx])
//	switch ss.FieldsKind[idx] {
//	case reflect.Interface:
//		*(*any)(p) = v
//	case reflect.String:
//		*(*string)(p) = v
//	}
//}
//
//func (ss *StructSchema) BindBool(ptr uintptr, idx int, v bool) {
//	if ss.FieldsKind[idx] == reflect.Interface {
//		*(*any)(unsafe.Pointer(ptr + ss.FieldsOffset[idx])) = v
//	} else {
//		*(*bool)(unsafe.Pointer(ptr + ss.FieldsOffset[idx])) = v
//	}
//}
//
//func (ss *StructSchema) BindInt(ptr uintptr, idx int, v int64) {
//	p := unsafe.Pointer(ptr + ss.FieldsOffset[idx])
//	switch ss.FieldsKind[idx] {
//	case reflect.Int:
//		*(*int)(p) = int(v)
//	case reflect.Int8:
//		if v < math.MinInt8 || v > math.MaxInt8 {
//			goto errPanic
//		}
//		*(*int8)(p) = int8(v)
//	case reflect.Int16:
//		if v < math.MinInt16 || v > math.MaxInt16 {
//			goto errPanic
//		}
//		*(*int16)(p) = int16(v)
//	case reflect.Int32:
//		if v < math.MinInt32 || v > math.MaxInt32 {
//			goto errPanic
//		}
//		*(*int32)(p) = int32(v)
//	case reflect.Int64:
//		*(*int64)(p) = v
//	case reflect.Interface:
//		*(*any)(p) = v
//	}
//	return
//
//errPanic:
//	panic(errNumOutOfRange)
//}
//
//func (ss *StructSchema) BindUint(ptr uintptr, idx int, v uint64) {
//	p := unsafe.Pointer(ptr + ss.FieldsOffset[idx])
//	switch ss.FieldsKind[idx] {
//	case reflect.Uint:
//		*(*uint)(p) = uint(v)
//	case reflect.Uint8:
//		if v > math.MaxUint8 {
//			goto errPanic
//		}
//		*(*uint8)(p) = uint8(v)
//	case reflect.Uint16:
//		if v > math.MaxUint16 {
//			goto errPanic
//		}
//		*(*uint16)(p) = uint16(v)
//	case reflect.Uint32:
//		if v > math.MaxUint32 {
//			goto errPanic
//		}
//		*(*uint32)(p) = uint32(v)
//	case reflect.Uint64:
//		*(*uint64)(p) = v
//	case reflect.Interface:
//		*(*any)(p) = v
//	}
//	return
//
//errPanic:
//	panic(errNumOutOfRange)
//}
//
//func (ss *StructSchema) BindFloat(ptr uintptr, idx int, v float64) {
//	p := unsafe.Pointer(ptr + ss.FieldsOffset[idx])
//	switch ss.FieldsKind[idx] {
//	case reflect.Float32:
//		if v < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
//			goto errPanic
//		}
//		*(*float32)(p) = float32(v)
//	case reflect.Float64:
//		*(*float64)(p) = v
//	case reflect.Interface:
//		*(*any)(p) = v
//	}
//	return
//
//errPanic:
//	panic(errNumOutOfRange)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (ss *StructSchema) BindColumn(ptr uintptr, k string, v any) {
//	ss.BindValue(ptr, ss.ColumnIndex(k), v)
//}
//
//func (ss *StructSchema) BindField(ptr uintptr, k string, v any) {
//	ss.BindValue(ptr, ss.FieldIndex(k), v)
//}
//
//func (ss *StructSchema) BindValue(ptr uintptr, idx int, v any) {
//	p := unsafe.Pointer(ptr + ss.FieldsOffset[idx])
//	switch ss.FieldsKind[idx] {
//	case reflect.Int:
//		*(*int)(p) = v.(int)
//	case reflect.Int8:
//		*(*int8)(p) = v.(int8)
//	case reflect.Int16:
//		*(*int16)(p) = v.(int16)
//	case reflect.Int32:
//		*(*int32)(p) = v.(int32)
//	case reflect.Int64:
//		*(*int64)(p) = v.(int64)
//
//	case reflect.Uint:
//		*(*uint)(p) = v.(uint)
//	case reflect.Uint8:
//		*(*uint8)(p) = v.(uint8)
//	case reflect.Uint16:
//		*(*uint16)(p) = v.(uint16)
//	case reflect.Uint32:
//		*(*uint32)(p) = v.(uint32)
//	case reflect.Uint64:
//		*(*uint64)(p) = v.(uint64)
//
//	case reflect.String:
//		*(*string)(p) = v.(string)
//	case reflect.Bool:
//		*(*bool)(p) = v.(bool)
//
//	case reflect.Interface:
//		*(*any)(p) = v
//	}
//}
