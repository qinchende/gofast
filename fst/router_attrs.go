// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "github.com/qinchende/gofast/cst"

type (
	// 可以为路由自定义配置属性
	RouteAttrs interface {
		SetIndex(uint16)
	}

	UrlParam struct {
		Key   string
		Value string
	}

	urlParams []UrlParam
)

func (ps *urlParams) Get(name string) (string, bool) {
	kvs := *ps
	for i := range kvs {
		if kvs[i].Key == name {
			return kvs[i].Value, true
		}
	}
	return "", false
}

func (ps *urlParams) Value(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

func (ps *urlParams) ValueMust(name string) (va string) {
	v, ok := ps.Get(name)
	if !ok {
		cst.Panic("没有找到参数：" + name)
	}
	return v
}
