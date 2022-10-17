// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

type (
	// 可以为路由自定义配置属性
	RouteAttrs interface {
		SetRouteIndex(uint16)
	}

	UrlParam struct {
		Key   string
		Value string
	}

	routeParams []UrlParam
)

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
