// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "net/http"

// 一次性构建Mini内存版的所有路由项
func (gft *GoFast) buildAllRoutes() {
	gft.routerTrees = make(methodTrees, 0, 9)
	gft.routerTrees = append(gft.routerTrees, &methodTree{method: http.MethodGet})
	gft.routerTrees = append(gft.routerTrees, &methodTree{method: http.MethodPost})

	// TODO：启动server之前，注册的路由只是做了记录在allRouters变量中，这里开始一次性构造路由前缀树
	// Note: 前面三个路由项是系统默认的特殊路由，不参与具体的路由树构造
	for i := int(gft.speRoutesLen); i < len(gft.allRoutes); i++ {
		gft.regRouteItem(gft.allRoutes[i])
	}

	// 打印底层构造的路由树
	if gft.WebConfig.PrintRouteTrees {
		gft.printRouteTrees()
	}
	// 这里开始完整的重建整个路由树的数据结构
	gft.buildMiniRoutes()
}

// 注册每一条的路由，生成 原始的 Radix 树
func (gft *GoFast) regRouteItem(ri *RouteItem) {
	// Debug模式下打印新添加的路由
	if gft.WebConfig.PrintRouteTrees {
		debugPrintRoute(gft, ri)
	}

	mTree := gft.getMethodTree(ri.method)
	if mTree == nil {
		mTree = &methodTree{method: ri.method, root: nil}
		gft.routerTrees = append(gft.routerTrees, mTree)
	}
	mTree.regRoute(ri.fullPath, ri)
}

// 获取method树的根节点
func (gft *GoFast) getMethodMiniRoot(method string) (tRoot *radixMiniNode) {
	switch method[0] {
	case 'G':
		tRoot = gft.routerTrees[0].miniRoot
	case 'P':
		if method[1] == 'O' {
			tRoot = gft.routerTrees[1].miniRoot
		} else {
			tRoot = gft.routerTrees.getTreeMiniRoot(method)
		}
	default:
		tRoot = gft.routerTrees.getTreeMiniRoot(method)
	}
	return
}

func (gft *GoFast) getMethodTree(method string) (tree *methodTree) {
	switch method[0] {
	case 'G':
		tree = gft.routerTrees[0]
	case 'P':
		if method[1] == 'O' {
			tree = gft.routerTrees[1]
		} else {
			tree = gft.routerTrees.getTree(method)
		}
	default:
		tree = gft.routerTrees.getTree(method)
	}
	return
}
