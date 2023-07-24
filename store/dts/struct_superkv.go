package dts

import (
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"unsafe"
)

// +++ gofast/cst/SuperKV
//// 可能用map，也可能自定义数组等合适的数据结构存取。
//// 比如上下文中用来保存解析到的请求数据，主要是KV形式
//type SuperKV interface {
//	Set(k string, v any)
//	Get(k string) (any, bool)
//	Del(k string)
//	Len() int
//	SetString(k string, v string)
//	GetString(k string) (string, bool)
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type StructKV struct {
	SS  *StructSchema
	Ptr unsafe.Pointer
}

func AsSuperKV(v any) (ret *StructKV) {
	ret = &StructKV{}
	ret.SS = SchemaForInput(v)
	ret.Ptr = (*rt.AFace)(unsafe.Pointer(&v)).DataPtr
	return
}

func (skv *StructKV) Get(k string) (v any, tf bool) {
	idx := skv.SS.ColumnIndex(k)
	p := unsafe.Pointer(uintptr(skv.Ptr) + skv.SS.FieldsAttr[idx].Offset)

	tf = true
	switch skv.SS.FieldsAttr[idx].Kind {
	//case reflect.Int:
	//	v = *(*int)(p)
	//case reflect.Int8:
	//	v = *(*int8)(p)
	//case reflect.Int16:
	//	v = *(*int16)(p)
	//case reflect.Int32:
	//	v = *(*int32)(p)
	//case reflect.Int64:
	//	v = *(*int64)(p)
	//
	//case reflect.Uint:
	//	v = *(*uint)(p)
	//case reflect.Uint8:
	//	v = *(*uint8)(p)
	//case reflect.Uint16:
	//	v = *(*uint16)(p)
	//case reflect.Uint32:
	//	v = *(*uint32)(p)
	//case reflect.Uint64:
	//	v = *(*uint64)(p)
	//
	//case reflect.Float32:
	//	v = *(*float32)(p)
	//case reflect.Float64:
	//	v = *(*float64)(p)

	case reflect.String:
		v = *(*string)(p)
	//case reflect.Bool:
	//	v = *(*bool)(p)
	//
	//case reflect.Interface:
	//	*(*any)(p) = v

	default:
		tf = false
	}
	return
}

func (skv *StructKV) GetString(k string) (v string, tf bool) {
	tmp, tf := skv.Get(k)
	v = tmp.(string)
	return
}

func (skv *StructKV) Set(k string, v any) {
	idx := skv.SS.ColumnIndex(k)
	p := unsafe.Pointer(uintptr(skv.Ptr) + skv.SS.FieldsAttr[idx].Offset)
	str := v.(string)

	switch skv.SS.FieldsAttr[idx].Kind {
	//case reflect.Int:
	//	BindInt(p, lang.ParseInt(str))
	//case reflect.Int8:
	//	BindInt8(p, lang.ParseInt(str))
	//case reflect.Int16:
	//	BindInt16(p, lang.ParseInt(str))
	//case reflect.Int32:
	//	BindInt32(p, lang.ParseInt(str))
	//case reflect.Int64:
	//	BindInt64(p, lang.ParseInt(str))
	//
	//case reflect.Uint:
	//	BindUint(p, lang.ParseUint(str))
	//case reflect.Uint8:
	//	BindUint8(p, lang.ParseUint(str))
	//case reflect.Uint16:
	//	BindUint16(p, lang.ParseUint(str))
	//case reflect.Uint32:
	//	BindUint32(p, lang.ParseUint(str))
	//case reflect.Uint64:
	//	BindUint64(p, lang.ParseUint(str))
	//
	//case reflect.Float32:
	//	BindFloat32(p, lang.ParseFloat(str))
	//case reflect.Float64:
	//	BindFloat64(p, lang.ParseFloat(str))

	case reflect.String:
		BindString(p, str)
	//case reflect.Bool:
	//	BindBool(p, lang.ParseBool(str))
	//
	//case reflect.Interface:
	//	BindAny(p, str)

	default:
		panic(errNotSupportType)
	}
}

func (skv *StructKV) SetString(k string, v string) {
	skv.Set(k, v)
}

func (skv *StructKV) Del(k string) {
}

func (skv *StructKV) Len() int {
	return len(skv.SS.columns)
}
