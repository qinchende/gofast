// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"math"
	"net/http"
	"regexp"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 注册一个404处理函数, 只能设置一次
func (gft *GoFast) reg404Handler(hds []CtxHandler) {
	if hds == nil {
		return
	}
	GFPanicIf(len(gft.allRouters[1].eHds) > 0, "重复，你可能已经设置了NoRoute处理函数")
	if hds != nil {
		gft.allRouters[1].eHds = addCtxHandlers(gft.fstMem, hds)
	}
}

// 注册一个405处理函数, 只能设置一次
func (gft *GoFast) reg405Handler(hds []CtxHandler) {
	if hds == nil {
		return
	}
	GFPanicIf(len(gft.allRouters[2].eHds) > 0, "重复，你可能已经设置了NoMethod处理函数")
	if hds != nil {
		gft.allRouters[2].eHds = addCtxHandlers(gft.fstMem, hds)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 所有注册的 router handlers 都要通过此函数来注册
func (gp *RouterGroup) register(httpMethod, relPath string, hds []CtxHandler) *RouteItem {
	GFPanicIf(len(hds) <= 0, "there must be at least one handler")
	// 最终的路由绝对路径
	absPath := gp.fixAbsolutePath(relPath)
	GFPanicIf(absPath[0] != '/', "Path must begin with '/'")
	GFPanicIf(len(absPath) > math.MaxUint8, "The path is more than 255 chars")
	GFPanicIf(len(httpMethod) == 0, "HTTP method can not be empty")

	// 新添加一个 GroupItem，记录所有的处理函数
	ri := &RouteItem{
		method:    httpMethod,
		fullPath:  absPath,
		group:     gp,
		routerIdx: 0,
	}
	myApp := gp.myApp
	ri.eHds = addCtxHandlers(myApp.fstMem, hds)
	// 保存了所有的合法路由规则，暂不生成路由树，待所有环境初始化完成之后再构造路由前缀树
	ri.routerIdx = uint16(len(myApp.allRouters))
	myApp.allRouters = append(myApp.allRouters, ri)
	GFPanicIf(len(myApp.allRouters) > math.MaxInt16, "Too many routers more than MaxInt16.")
	return ri
}

// TODO: 有个问题，httpMethod参数没有做枚举校验，可以创建任意名称的method路由数，真要这么自由吗???
func (gp *RouterGroup) Handle(httpMethod, relPath string, handlers ...CtxHandler) *RouteItem {
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		panic("http method " + httpMethod + " is not valid")
	}
	return gp.register(httpMethod, relPath, handlers)
}

func (gp *RouterGroup) Get(relPath string, handlers ...CtxHandler) *RouteItem {
	return gp.register(http.MethodGet, relPath, handlers)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (gp *RouterGroup) Post(relPath string, handlers ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPost, relPath, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (gp *RouterGroup) Delete(relPath string, handlers ...CtxHandler) *RouteItem {
	return gp.register(http.MethodDelete, relPath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (gp *RouterGroup) Patch(relPath string, handlers ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPatch, relPath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (gp *RouterGroup) Put(relPath string, handlers ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPut, relPath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (gp *RouterGroup) Options(relPath string, handlers ...CtxHandler) *RouteItem {
	return gp.register(http.MethodOptions, relPath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (gp *RouterGroup) Head(relPath string, handlers ...CtxHandler) *RouteItem {
	return gp.register(http.MethodHead, relPath, handlers)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (gp *RouterGroup) All(relPath string, handlers ...CtxHandler) {
	gp.register(http.MethodGet, relPath, handlers)
	gp.register(http.MethodPost, relPath, handlers)
	gp.register(http.MethodPut, relPath, handlers)
	gp.register(http.MethodPatch, relPath, handlers)
	gp.register(http.MethodHead, relPath, handlers)
	gp.register(http.MethodOptions, relPath, handlers)
	gp.register(http.MethodDelete, relPath, handlers)
	gp.register(http.MethodConnect, relPath, handlers)
	gp.register(http.MethodTrace, relPath, handlers)
}

// 特殊类型
func (gp *RouterGroup) GetPost(relPath string, handlers ...CtxHandler) (get, post *RouteItem) {
	get = gp.register(http.MethodGet, relPath, handlers)
	post = gp.register(http.MethodPost, relPath, handlers)
	return
}
