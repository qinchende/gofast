// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

func (c *Context) FullPath() string {
	if c.route.ptrNode != nil {
		return c.myApp.allPaths[c.route.ptrNode.routeIdx]
	} else {
		return ""
	}
}

func (c *Context) SetRouteToAny() {
	c.route.ptrNode = c.myApp.miniNodeAny
}

func (c *Context) SetRouteTo404() {
	c.route.ptrNode = c.myApp.miniNode404
}

func (c *Context) SetRouteTo405() {
	c.route.ptrNode = c.myApp.miniNode405
}
