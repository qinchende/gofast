// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 获取所在节点的Path |
func (gft *GoFast) FullPath(idx uint16) string {
	if gft.allRouters[idx] != nil {
		return gft.allRouters[idx].fullPath
	} else {
		return ""
	}
}

func (gft *GoFast) RouteLength() uint16 {
	return uint16(len(gft.allRouters))
}

func (c *Context) FullPath() string {
	if c.match.ptrNode != nil {
		return c.gftApp.allRouters[c.match.ptrNode.routerIdx].fullPath
	} else {
		return ""
	}
}

// 获取当前路由节点
func (c *Context) CurrRoute() *RouteItem {
	if c.RouteID <= 0 || c.RouteID >= uint16(len(c.gftApp.allRouters)) {
		return nil
	}
	return c.gftApp.allRouters[c.RouteID]
}

func (ri *RouteItem) FullPath() string {
	return ri.fullPath
}

//func (c *Context) RouteIndex() int16 {
//	var nodeIdx int16 = -1
//	if c != nil && c.match.ptrNode != nil {
//		nodeIdx = c.match.ptrNode.routerIdx
//	}
//	return nodeIdx
//}
