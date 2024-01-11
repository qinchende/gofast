// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/store/dts"
)

const (
	PmsToken = "tok"
)

type BaseFields struct {
	Tok string `v:"len=[64:128]"`
}

var _BasePms = []string{PmsToken}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 构造给定对象的字段名数组，加上公共的字段
func PmsKeys(obj any) []string {
	ss := dts.SchemaAsReq(obj)
	cls := ss.Columns
	newCls := make([]string, len(cls)+len(_BasePms))
	copy(newCls, cls)
	copy(newCls[len(cls):], _BasePms)
	return newCls // TODO: 可能需要考虑排重
}
