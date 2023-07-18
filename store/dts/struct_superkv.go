package dts

import (
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"unsafe"
)

//type SuperKV interface {
//	Get(k string) (any, bool)
//	Set(k string, v any)
//	Del(k string)
//	Len() int
//	//GetString(k string) (string, bool)
//	//SetString(k string, v string)
//}

type StructKV struct {
	ss  *StructSchema
	Ptr unsafe.Pointer
}

func AsSuperKV(v any) (ret *StructKV) {
	ret = &StructKV{}
	ret.ss = SchemaForInput(v)
	ret.Ptr = (*rt.AFace)(unsafe.Pointer(&v)).DataPtr
	return
}

func (skv *StructKV) Get(k string) (v any, tf bool) {
	idx := skv.ss.ColumnIndex(k)
	p := unsafe.Pointer(uintptr(skv.Ptr) + skv.ss.FieldsAttr[idx].Offset)

	tf = true
	switch skv.ss.FieldsAttr[idx].Kind {
	case reflect.Int:
		v = *(*int)(p)
	case reflect.Int8:
		v = *(*int8)(p)
	case reflect.Int16:
		v = *(*int16)(p)
	case reflect.Int32:
		v = *(*int32)(p)
	case reflect.Int64:
		v = *(*int64)(p)

	case reflect.Uint:
		v = *(*uint)(p)
	case reflect.Uint8:
		v = *(*uint8)(p)
	case reflect.Uint16:
		v = *(*uint16)(p)
	case reflect.Uint32:
		v = *(*uint32)(p)
	case reflect.Uint64:
		v = *(*uint64)(p)

	case reflect.String:
		v = *(*string)(p)
	case reflect.Bool:
		v = *(*bool)(p)

	case reflect.Interface:
		*(*any)(p) = v

	default:
		tf = false
	}
	return
}

func (skv *StructKV) Set(k string, v any) {
}

func (skv *StructKV) Del(k string) {
}

func (skv *StructKV) Len() int {
	return len(skv.ss.columns)
}
