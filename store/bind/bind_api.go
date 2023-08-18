package bind

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/dts"
)

// NOTE(Important): 下面API中的第一个参数dst，必须是对象指针
func BindKV(dst any, kvs cst.SuperKV, model int8) error {
	return bindKVToStruct(dst, kvs, dts.AsOptions(model))
}

func BindKVX(dst any, kvs cst.SuperKV, opts *dts.BindOptions) error {
	return bindKVToStruct(dst, kvs, opts)
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
