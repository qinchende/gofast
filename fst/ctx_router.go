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
