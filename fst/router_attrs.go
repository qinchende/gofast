// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

type UrlParam struct {
	Key   string
	Value string
}

type routeParams []UrlParam

func (ps *routeParams) Get(name string) (string, bool) {
	for _, item := range *ps {
		if item.Key == name {
			return item.Value, true
		}
	}
	return "", false
}

func (ps *routeParams) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

// 可以为路由自定义配置属性
type (
	RouteAttrs interface {
		SetRouteIndex(uint16)
	}
)
