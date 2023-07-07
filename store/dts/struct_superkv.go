package dts

import (
	"github.com/qinchende/gofast/core/rt"
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

func (skv *StructKV) Get(k string) (any, bool) {
	return nil, false
}

func (skv *StructKV) Set(k string, v any) {
}

func (skv *StructKV) Del(k string) {
}

func (skv *StructKV) Len() int {
	return len(skv.ss.columns)
}
