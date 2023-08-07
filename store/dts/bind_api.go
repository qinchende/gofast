package dts

import "github.com/qinchende/gofast/cst"

// NOTE(Important): 下面API中的第一个参数dst，最好都是指针类型，避免将来发生值拷贝
func BindKV(dst any, kvs cst.SuperKV, like int8) error {
	return bindKVToStruct(dst, kvs, AsOptions(like))
}

func BindKVX(dst any, kvs cst.KV, opts *BindOptions) error {
	return bindKVToStruct(dst, kvs, opts)
}

//// 根据结构体配置信息，优化字段值 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func Optimize(dst any, like int8) error {
//	return optimizeStruct(dst, AsOptions(like))
//
//}
//func OptimizeX(dst any, opts *BindOptions) error {
//	return optimizeStruct(dst, opts)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 提取对象的字段的column名
//func Columns(obj any, like int8) []string {
//	sm := Schema(obj, AsOptions(like))
//	return sm.columns
//}
