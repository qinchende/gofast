// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"github.com/qinchende/gofast/cst"
)

// NOTE(Important): 下面API中的第一个参数dst，最好都是指针类型，避免将来发生值拷贝
func BindKV(dst any, kvs cst.SuperKV, like int8) error {
	return bindKVToStruct(dst, kvs, matchOptions(like))
}

func BindKVX(dst any, kvs cst.KV, opts *BindOptions) error {
	return bindKVToStruct(dst, kvs, opts)
}

func BindSlice(dst any, src any, like int8) error {
	fOpt := &fieldOptions{}
	return bindList(dst, src, fOpt, matchOptions(like))
}
func BindSliceX(dst any, src any, opts *BindOptions) error {
	fOpt := &fieldOptions{}
	return bindList(dst, src, fOpt, opts)
}

// 根据结构体配置信息，优化字段值 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Optimize(dst any, like int8) error {
	return optimizeStruct(dst, matchOptions(like))

}
func OptimizeX(dst any, opts *BindOptions) error {
	return optimizeStruct(dst, opts)
}
