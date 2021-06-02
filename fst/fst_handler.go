// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/logx"
	"net/http"
	"path"
)

// 系统默认错误处理函数，可以设置 code 和 message.
func defErrorHandler(code int, defaultMessage []byte) CtxHandler {
	return func(c *Context) {
		c.ResWrap.WriteHeader(code)
		if c.ResWrap.Written() {
			return
		}
		if c.ResWrap.Status() == code {
			c.ResWrap.Header()["Content-Type"] = mimePlain
			_, err := c.ResWrap.Write(defaultMessage)
			if err != nil {
				logx.DebugPrint("Cannot write message to writer during serve error: %v", err)
			}
			return
		}
		c.ResWrap.WriteHeaderNow()
		return
	}
}

// 特殊函数处理
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 如果没有配置，添加默认的处理函数
func (gft *GoFast) checkDefaultHandler() {
	if gft.routerItem404 == nil {
		if gft.DisableDefNoRoute {
			gft.NoRoute()
		} else {
			gft.NoRoute(defErrorHandler(http.StatusNotFound, default404Body))
		}
	}
	if gft.routerItem405 == nil {
		if gft.DisableDefNotAllowed {
			gft.NoMethod()
		} else {
			gft.NoMethod(defErrorHandler(http.StatusMethodNotAllowed, default405Body))
		}
	}
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (gft *GoFast) NoRoute(handlers ...CtxHandler) {
	gft.reg404Handler(handlers)
}

// NoMethod sets the handlers called when...
func (gft *GoFast) NoMethod(handlers ...CtxHandler) {
	gft.reg405Handler(handlers)
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
	redirectRequest(c)
}

func redirectRequest(c *Context) {
	req := c.ReqRaw
	rPath := req.URL.Path
	rURL := req.URL.String()

	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}
	logx.DebugPrint("redirecting request %d: %s --> %s", code, rPath, rURL)
	http.Redirect(c.ResWrap, req, rURL, code)
	c.ResWrap.WriteHeaderNow()
}
