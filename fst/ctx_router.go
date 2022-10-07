// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 获取所在节点的Path |
func (gft *GoFast) FullPath(idx uint16) string {
	if gft.allRoutes[idx] != nil {
		return gft.allRoutes[idx].fullPath
	} else {
		return ""
	}
}

func (gft *GoFast) RouteLength() uint16 {
	return uint16(len(gft.allRoutes))
}

func (c *Context) FullPath() string {
	if c.route.ptrNode != nil {
		return c.myApp.allRoutes[c.route.ptrNode.routeIdx].fullPath
	} else {
		return ""
	}
}

// 获取当前路由节点
func (c *Context) CurrRoute() *RouteItem {
	if c.RouteIdx <= 0 || c.RouteIdx >= uint16(len(c.myApp.allRoutes)) {
		return nil
	}
	return c.myApp.allRoutes[c.RouteIdx]
}

func (ri *RouteItem) FullPath() string {
	return ri.fullPath
}

//func (c *Context) RouteIndex() int16 {
//	var nodeIdx int16 = -1
//	if c != nil && c.route.ptrNode != nil {
//		nodeIdx = c.route.ptrNode.routeIdx
//	}
//	return nodeIdx
//}
