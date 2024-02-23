package dts

import (
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"unsafe"
)

// @@@ +++++++++++++++++++++++++++++++++
// Important Note:
// StructKV 实现 SuperKV 接口是带有很大局限性的。特别是Set和Get函数，只支持类似Web数据提交这种KV都是string的特殊场景。
// 切记不可随意使用
// @@@ +++++++++++++++++++++++++++++++++
type StructKV struct {
	SS  *StructSchema  // struct 预解析的 schema
	Ptr unsafe.Pointer // struct 对象对应的地址
}

// Note: 这里返回StructKV的指针类型，而不是值类型。主要是因为要讲这个变量赋值给接口 fst.Context.Pms ，不希望发生值拷贝
func AsSuperKV(v any) (ret *StructKV) {
	ret = &StructKV{
		SS:  SchemaAsReq(v),
		Ptr: (*rt.AFace)(unsafe.Pointer(&v)).DataPtr,
	}
	return
}

// 为 StructKV 实现 gofast/core/cst/SuperKV 接口
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Note: 这里的v只支持string类型。因为此结构只是为了实现接口SuperKV，使其能加载Web请求中URL、Header等包含的参数
func (skv *StructKV) Get(k string) (v any, tf bool) {
	idx := skv.SS.ColumnIndex(k)

	switch skv.SS.FieldsAttr[idx].Kind {
	case reflect.String:
		p := unsafe.Pointer(uintptr(skv.Ptr) + skv.SS.FieldsAttr[idx].Offset)
		v = *(*string)(p)
		tf = true
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

// Note: 这里的v只支持string类型。因为此结构只是为了实现接口SuperKV，使其能加载Web请求中URL、Header等包含的参数
func (skv *StructKV) Set(k string, v any) {
	idx := skv.SS.ColumnIndex(k)

	// NOTE：目前只支持API请求提交的字节数据，KV都是string类型
	switch skv.SS.FieldsAttr[idx].Kind {
	case reflect.String:
		p := unsafe.Pointer(uintptr(skv.Ptr) + skv.SS.FieldsAttr[idx].Offset)
		// BindString(p, v.(string))
		*(*string)(p) = v.(string)
	default:
		panic(errNotSupportType)
	}
}

func (skv *StructKV) SetString(k string, v string) {
	skv.Set(k, v)
}

// 不需要删除任何项目
func (skv *StructKV) Del(k string) {
}

func (skv *StructKV) Len() int {
	return len(skv.SS.Columns)
}
