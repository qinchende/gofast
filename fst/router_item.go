// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"math"
	"net/http"
	"regexp"
)

// idx 0. not match any 1. 404 handlers 2. 405 handlers
func (gft *GoFast) regSpecialHandlers(hds []CtxHandler, idx int) {
	cst.PanicIf(len(hds) <= 0, "there must be at least one handler")
	cst.PanicIf(len(gft.allRoutes[idx].eHds) > 0, "handlers already exists.")
	gft.allRoutes[idx].eHds = addCtxHandlers(gft.fstMem, hds)
}

// 在当前分组注册一个新节点
func (gp *RouteGroup) register(httpMethod, relPath string, hds []CtxHandler) *RouteItem {
	// 最终的路由绝对路径
	absPath := gp.fixAbsolutePath(relPath)
	cst.PanicIf(absPath[0] != '/', "Path must begin with '/'")
	cst.PanicIf(len(absPath) > math.MaxUint8, "The path is more than 255 chars")
	cst.PanicIf(len(httpMethod) == 0, "HTTP method can not be empty")

	// 新添加一个 GroupItem，记录所有的处理函数
	ri := &RouteItem{
		method:   httpMethod,
		fullPath: absPath,
		group:    gp,
		routeIdx: 0,
	}

	// 可以 RouteItem 只创建对象，不指定处理函数。等后面再设置
	if len(hds) == 0 {
		return ri
	} else {
		return ri.Handle(hds)
	}
}

// 当前节点设置必要路由属性
func (ri *RouteItem) Handle(hds []CtxHandler) *RouteItem {
	cst.PanicIf(ri.routeIdx > 0, "this route already registered.")
	cst.PanicIf(len(hds) <= 0, "there must be at least one handler")

	myApp := ri.group.myApp
	ri.eHds = addCtxHandlers(myApp.fstMem, hds)
	// 保存了所有的合法路由规则，暂不生成路由树，待所有环境初始化完成之后再构造路由前缀树
	ri.routeIdx = uint16(len(myApp.allRoutes))
	myApp.allRoutes = append(myApp.allRoutes, ri)
	cst.PanicIf(len(myApp.allRoutes) > math.MaxInt16, "Too many routers more than MaxInt16.")
	return ri
}

// 所有路由节点都设置同样的路由属性
func (ris RouteItems) Handle(hds []CtxHandler) RouteItems {
	for i := range ris {
		ris[i].Handle(hds)
	}
	return ris
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 有个问题，httpMethod参数没有做枚举校验，可以创建任意名称的method路由数，真要这么自由吗???
func (gp *RouteGroup) Handle(httpMethod, relPath string, hds ...CtxHandler) *RouteItem {
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		panic("http method " + httpMethod + " is not valid")
	}
	return gp.register(httpMethod, relPath, hds)
}

func (gp *RouteGroup) Get(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodGet, relPath, hds)
}

func (gp *RouteGroup) Post(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPost, relPath, hds)
}

func (gp *RouteGroup) Delete(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodDelete, relPath, hds)
}

func (gp *RouteGroup) Patch(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPatch, relPath, hds)
}

func (gp *RouteGroup) Put(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPut, relPath, hds)
}

func (gp *RouteGroup) Options(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodOptions, relPath, hds)
}

func (gp *RouteGroup) Head(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodHead, relPath, hds)
}

// 特殊类型
func (gp *RouteGroup) GetPost(relPath string, hds ...CtxHandler) RouteItems {
	get := gp.register(http.MethodGet, relPath, hds)
	post := gp.register(http.MethodPost, relPath, hds)
	return RouteItems{get, post}
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (gp *RouteGroup) All(relPath string, hds ...CtxHandler) RouteItems {
	get := gp.register(http.MethodGet, relPath, hds)
	post := gp.register(http.MethodPost, relPath, hds)
	put := gp.register(http.MethodPut, relPath, hds)
	patch := gp.register(http.MethodPatch, relPath, hds)
	head := gp.register(http.MethodHead, relPath, hds)
	opt := gp.register(http.MethodOptions, relPath, hds)
	del := gp.register(http.MethodDelete, relPath, hds)
	conn := gp.register(http.MethodConnect, relPath, hds)
	trace := gp.register(http.MethodTrace, relPath, hds)
	return RouteItems{get, post, put, patch, head, opt, del, conn, trace}
}
