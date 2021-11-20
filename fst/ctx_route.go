// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 获取所在节点的Path |
func (gft *GoFast) FullPath(idx int16) string {
	if idx >= 0 && gft.allRouters[idx] != nil {
		return gft.allRouters[idx].fullPath
	} else {
		return ""
	}
}

func (gft *GoFast) RoutesLen() uint16 {
	return uint16(len(gft.allRouters))
}

func (c *Context) FullPath() string {
	if c.match.ptrNode != nil && c.match.ptrNode.routerIdx >= 0 {
		return c.gftApp.allRouters[c.match.ptrNode.routerIdx].fullPath
	} else {
		return ""
	}
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
