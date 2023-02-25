// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"math"
	"net/http"
	"regexp"
)

// idx 0. not match any 1. 404 handlers 2. 405 handlers
func (gft *GoFast) regSpecialHandlers(hds []CtxHandler, idx int) {
	PanicIf(len(hds) <= 0, "there must be at least one handler")
	PanicIf(len(gft.allRoutes[idx].eHds) > 0, "handlers already exists.")
	gft.allRoutes[idx].eHds = addCtxHandlers(gft.fstMem, hds)
}

// 所有注册的 router handlers 都要通过此函数来注册
func (gp *RouteGroup) register(httpMethod, relPath string, hds []CtxHandler) *RouteItem {
	PanicIf(len(hds) <= 0, "there must be at least one handler")
	// 最终的路由绝对路径
	absPath := gp.fixAbsolutePath(relPath)
	PanicIf(absPath[0] != '/', "Path must begin with '/'")
	PanicIf(len(absPath) > math.MaxUint8, "The path is more than 255 chars")
	PanicIf(len(httpMethod) == 0, "HTTP method can not be empty")

	// 新添加一个 GroupItem，记录所有的处理函数
	ri := &RouteItem{
		method:   httpMethod,
		fullPath: absPath,
		group:    gp,
		routeIdx: 0,
	}
	myApp := gp.myApp
	ri.eHds = addCtxHandlers(myApp.fstMem, hds)
	// 保存了所有的合法路由规则，暂不生成路由树，待所有环境初始化完成之后再构造路由前缀树
	ri.routeIdx = uint16(len(myApp.allRoutes))
	myApp.allRoutes = append(myApp.allRoutes, ri)
	PanicIf(len(myApp.allRoutes) > math.MaxInt16, "Too many routers more than MaxInt16.")
	return ri
}

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

// POST is a shortcut for router.Handle("POST", path, handle).
func (gp *RouteGroup) Post(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPost, relPath, hds)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (gp *RouteGroup) Delete(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodDelete, relPath, hds)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (gp *RouteGroup) Patch(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPatch, relPath, hds)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (gp *RouteGroup) Put(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPut, relPath, hds)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (gp *RouteGroup) Options(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodOptions, relPath, hds)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (gp *RouteGroup) Head(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodHead, relPath, hds)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (gp *RouteGroup) All(relPath string, hds ...CtxHandler) {
	gp.register(http.MethodGet, relPath, hds)
	gp.register(http.MethodPost, relPath, hds)
	gp.register(http.MethodPut, relPath, hds)
	gp.register(http.MethodPatch, relPath, hds)
	gp.register(http.MethodHead, relPath, hds)
	gp.register(http.MethodOptions, relPath, hds)
	gp.register(http.MethodDelete, relPath, hds)
	gp.register(http.MethodConnect, relPath, hds)
	gp.register(http.MethodTrace, relPath, hds)
}

// 特殊类型
func (gp *RouteGroup) GetPost(relPath string, hds ...CtxHandler) (get, post *RouteItem) {
	get = gp.register(http.MethodGet, relPath, hds)
	post = gp.register(http.MethodPost, relPath, hds)
	return
}
