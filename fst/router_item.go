// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/core/cst"
	"math"
)

// 当前节点设置必要路由属性
func (ri *RouteItem) Handle(hds ...CtxHandler) *RouteItem {
	// 本方法，每个节点只能执行一次
	cst.PanicIf(ri.routeIdx > 0, "this route already registered.")
	cst.PanicIf(len(hds) <= 0, "there must be at least one handler")

	server := ri.group.app
	ri.eHds = addCtxHandlers(server.fstMem, hds)
	// 保存了所有的合法路由规则，暂不生成路由树，待所有环境初始化完成之后再构造路由前缀树
	ri.routeIdx = uint16(len(server.allRoutes))
	server.allRoutes = append(server.allRoutes, ri)
	cst.PanicIf(len(server.allRoutes) > math.MaxInt16, "Too many routers more than MaxInt16.")
	return ri
}

// 所有路由节点都设置同样的路由属性
func (ris RouteItems) Handle(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].Handle(hds...)
	}
	return ris
}

// 所有路由节点都设置同样的路由属性
func (ri *RouteItem) Index() uint16 {
	return ri.routeIdx
}
