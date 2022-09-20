// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 可以为路由自定义配置属性
type (
	RouteAttrs interface {
		SetRouteIndex(uint16)
	}
)
