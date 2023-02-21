// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"github.com/qinchende/gofast/cst"
)

// NOTE(Important): 下面API中的第一个参数dst，最好都是指针类型，避免将来发生值拷贝
func ApplyKV(dst any, kvs cst.KV, like int8) error {
	return applyKVToStruct(dst, kvs, matchOptions(like))
}

func ApplyKVX(dst any, kvs cst.KV, opts *ApplyOptions) error {
	return applyKVToStruct(dst, kvs, opts)
}

func ApplySlice(dst any, src any, like int8) error {
	fOpt := &fieldOptions{}
	return applyList(dst, src, fOpt, matchOptions(like))
}
func ApplySliceX(dst any, src any, opts *ApplyOptions) error {
	fOpt := &fieldOptions{}
	return applyList(dst, src, fOpt, opts)
}

// 根据结构体配置信息，优化字段值 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Optimize(dst any, like int8) error {
	return optimizeStruct(dst, matchOptions(like))

}
func OptimizeX(dst any, opts *ApplyOptions) error {
	return optimizeStruct(dst, opts)
}
