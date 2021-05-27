// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 一次性注册所有路由项
func (gft *GoFast) regAllRouters() {
	for _, it := range gft.allRouters {
		gft.regRouterItem(it)
	}
}

// 注册每一条的路由，生成 原始的 Radix 树
func (gft *GoFast) regRouterItem(ri *RouterItem) {
	// Debug模式下打印新添加的路由
	DebugPrintRoute(ri)

	mTree := gft.getMethodTree(ri.method)
	if mTree == nil {
		mTree = &methodTree{method: ri.method, root: nil}
		gft.treeOthers = append(gft.treeOthers, mTree)
	}
	mTree.regRoute(ri.fullPath, ri)
	gft.fstMem.hdsItemCt++
}

// 获取method树的根节点
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
