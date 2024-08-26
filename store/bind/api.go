package bind

import (
	"errors"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/dts"
	"reflect"
	"unsafe"
)

const (
	AsDef    = dts.AsDef
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

// 根据结构体配置信息，优化字段值 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Note：比如用在配置文件读取后默认值的设置
// 1. 没有赋值的字段，如果指定了默认值，就自动设置成默认值
// 2. 字段值根据valid规则，做合法性验证
func Optimize(dst any, model int8) error {
	return optimizeStruct(dst, dts.AsOptions(model))
}
func OptimizeX(dst any, opts *dts.BindOptions) error {
	return optimizeStruct(dst, opts)
}
