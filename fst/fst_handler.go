// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/logx"
	"net/http"
	"path"
)

// 特殊函数处理
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 系统默认错误处理函数，可以设置 code 和 message.
func specialHandler(resStatus int, defaultMessage []byte) CtxHandler {
	return func(c *Context) {
		c.AbortDirect(resStatus, defaultMessage)
	}
}

// 如果没有配置，添加默认的处理函数
func (gft *GoFast) initDefaultHandlers() {
	if gft.DefNoRouteHandler && len(gft.allRoutes[1].eHds) == 0 {
		gft.NoRoute(specialHandler(http.StatusNotFound, default404Body))
	}
	if gft.DefNotAllowedHandler && len(gft.allRoutes[2].eHds) == 0 {
		gft.NoMethod(specialHandler(http.StatusMethodNotAllowed, default405Body))
	}
}

// 每次设置都会替换掉以前设置好的方法
// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (gft *GoFast) NoRoute(hds ...CtxHandler) {
	gft.reg404Handler(hds)
}

// 每次设置都会替换掉以前设置好的方法
// NoMethod sets the handlers called when...
func (gft *GoFast) NoMethod(hds ...CtxHandler) {
	gft.reg405Handler(hds)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 请求结尾的 '/' 取消或者添加之后重定向，看是否能够匹配到相应路由
func redirectTrailingSlash(c *Context) {
	req := c.ReqRaw
	p := req.URL.Path
	if prefix := path.Clean(c.ReqRaw.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}

	// GET 和 非GET 请求重定向状态不一样
	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}

	rURL := req.URL.String()
	logx.InfoF("redirecting request %d: %s --> %s", code, req.URL.Path, rURL)
	c.AbortRedirect(code, rURL)
}
