// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"math"
)

// 所有注册路由方法都走这个函数
func (gft *GoFast) addRoute(ri *RouterItem) {
	ifPanic(ri.fullPath[0] != '/', "Path must begin with '/'")
	ifPanic(len(ri.fullPath) > math.MaxUint8, "The path is more than 255 chars")
	ifPanic(len(ri.method) == 0, "HTTP method can not be empty")

	// 保存了所有的合法路由规则
	gft.allRouters = append(gft.allRouters, ri)

	mTree := gft.getMethodTree(ri.method)
	if mTree == nil {
		mTree = &methodTree{method: ri.method, root: nil}
		gft.treeOthers = append(gft.treeOthers, mTree)
	}
	mTree.regRoute(ri.fullPath, ri)
}

func (gft *GoFast) getMethodMiniRoot(method string) (tRoot *radixMiniNode) {
	switch method[0] {
	case 'G':
		tRoot = gft.treeGet.miniRoot
	case 'P':
		if method[1] == 'O' {
			tRoot = gft.treePost.miniRoot
		} else {
			tRoot = gft.treeOthers.getTreeMiniRoot(method)
		}
	default:
		tRoot = gft.treeOthers.getTreeMiniRoot(method)
	}
	return
}

func (gft *GoFast) getMethodTree(method string) (tree *methodTree) {
	switch method[0] {
	case 'P':
		if method[1] == 'O' {
			tree = gft.treePost
		} else {
			tree = gft.treeOthers.getTree(method)
		}
	case 'G':
		tree = gft.treeGet
	default:
		tree = gft.treeOthers.getTree(method)
	}
	return
}

//func redirectTrailingSlash(c *Context) {
//	req := c.ReqRaw
//	p := req.URL.Path
//	if prefix := path.Clean(c.ReqRaw.Header.Get("X-Forwarded-Prefix")); prefix != "." {
//		p = prefix + "/" + req.URL.Path
//	}
//	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
//	if req.Method != "GET" {
//		code = http.StatusTemporaryRedirect
//	}
//
//	req.URL.Path = p + "/"
//	if length := len(p); length > 1 && p[length-1] == '/' {
//		req.URL.Path = p[:length-1]
//	}
//	skill.DebugPrint("redirecting request %d: %s --> %s", code, p, req.URL.String())
//	http.Redirect(c.ResWrap, req, req.URL.String(), code)
//	c.ResWrap.WriteHeaderNow()
//}
