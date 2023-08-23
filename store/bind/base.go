package bind

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"unsafe"
)

const (
	AsConfig = dts.AsConfig
	AsReq    = dts.AsReq
	AsLoadDB = dts.AsDB
)

// NOTE(Important): 下面API中的第一个参数dst，必须是对象指针
func BindKV(dst any, kvs cst.SuperKV, model int8) error {
	return bindKVToStruct(dst, kvs, dts.AsOptions(model))
}

func BindKVX(dst any, kvs cst.SuperKV, opts *dts.BindOptions) error {
	return bindKVToStruct(dst, kvs, opts)
}

func BindList(dst any, src any, model int8) error {
	return BindListX(dst, src, dts.AsOptions(model))
}

func BindListX(dst any, src any, opts *dts.BindOptions) error {
	dstT := reflect.TypeOf(dst)

	dstKind := dstT.Kind()
	if dstKind != reflect.Array && dstKind != reflect.Slice {
		return errors.New("Dest value must be array or slice type.")
	}

	return bindList((unsafe.Pointer)(&dst), dstT, src, opts)
}

//// 根据结构体配置信息，优化字段值 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func Optimize(dst any, model int8) error {
//	return optimizeStruct(dst, AsOptions(model))
//
//}
//func OptimizeX(dst any, opts *BindOptions) error {
//	return optimizeStruct(dst, opts)
//}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 提取对象的字段的column名
//func Columns(obj any, like int8) []string {
//	sm := Schema(obj, AsOptions(like))
//	return sm.columns
//}
