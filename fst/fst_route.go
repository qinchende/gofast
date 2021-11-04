// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "net/http"

// 一次性注册所有路由项
func (gft *GoFast) regAllRouters() {
	gft.treeGet = &methodTree{method: http.MethodGet}
	gft.treePost = &methodTree{method: http.MethodPost}
	gft.treeOthers = make(methodTrees, 0, 9)

	// TODO：启动server之前，注册的路由只是做了记录在allRouters变量中，这里开始一次性构造路由前缀树
	for _, it := range gft.allRouters {
		gft.regRouterItem(it)
	}

	// 设置 treeAll
	lenTreeOthers := len(gft.treeOthers)
	ifPanic(lenTreeOthers > 7, "More than seven methods.")
	gft.treeAll = gft.treeOthers[:lenTreeOthers:9]
	gft.treeAll = append(gft.treeAll, gft.treeGet)
	gft.treeAll = append(gft.treeAll, gft.treePost)

	// 打印底层构造的路由树
	if gft.PrintRouteTrees {
		gft.printRouteTrees()
	}
	// 这里开始完整的重建整个路由树的数据结构
	gft.buildMiniRoutes()
}

// 注册每一条的路由，生成 原始的 Radix 树
func (gft *GoFast) regRouterItem(ri *RouteItem) {
	// Debug模式下打印新添加的路由
	if gft.PrintRouteTrees {
		debugPrintRoute(ri)
	}

	mTree := gft.getMethodTree(ri.method)
	if mTree == nil {
		mTree = &methodTree{method: ri.method, root: nil}
		gft.treeOthers = append(gft.treeOthers, mTree)
	}
	mTree.regRoute(ri.fullPath, ri)
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
