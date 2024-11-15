// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/aid/logx"
	"net/http"
	"path"
)

// 如果没有配置，添加默认的处理函数
func (gft *GoFast) initDefaultHandlers() {
	if len(gft.allRoutes[1].eHds) == 0 {
		gft.Reg404(func(c *Context) { c.AbortDirect(http.StatusNotFound, "404 (Not Found)") })
	}
	if len(gft.allRoutes[2].eHds) == 0 {
		gft.Reg405(func(c *Context) { c.AbortDirect(http.StatusMethodNotAllowed, "405 (Method Not Allowed)") })
	}
}

// 每次设置都会替换掉以前设置好的方法
// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (gft *GoFast) Reg404(hds ...CtxHandler) {
	gft.regSpecialHandlers(hds, 1)
}

// 每次设置都会替换掉以前设置好的方法
// NoMethod sets the handlers called when...
func (gft *GoFast) Reg405(hds ...CtxHandler) {
	gft.regSpecialHandlers(hds, 2)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 请求结尾的 '/' 取消或者添加之后重定向，看是否能够匹配到相应路由
func redirectTrailingSlash(c *Context) {
	req := c.Req
	p := req.Raw.URL.Path
	if prefix := path.Clean(c.Req.Raw.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		p = prefix + "/" + req.Raw.URL.Path
	}
	req.Raw.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.Raw.URL.Path = p[:length-1]
	}

	// GET 和 非GET 请求重定向状态不一样
	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Raw.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}

	rURL := req.Raw.URL.String()
	logx.Info().MsgF("redirecting request %d: %s --> %s", code, req.Raw.URL.Path, rURL)
	c.AbortRedirect(code, rURL)
}
