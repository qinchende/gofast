// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"net/http"
	"net/url"
	"sync"
)

// Context is the most important part of GoFast. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	*GFResponse               // response (请求前置拦截器 要用到的上下文)
	ReqW        *http.Request // request
	matchRst    matchResult   // 路由匹配结果

	Pms        url.Values // 所有Request参数的map（Params + queryCache + formCache）
	Params     Params     // : 或 * 对应的参数
	queryCache url.Values // param query result from c.ReqW.URL.Query()
	formCache  url.Values // the parsed form data from POST, PATCH, or PUT body parameters.

	// This mutex protect Keys map
	mu sync.RWMutex
	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string
	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite
	// Keys is a key/value pair exclusively for the context of each request.
	// 上下文传值
	Keys map[string]interface{}
}

/************************************/
/********** context creation ********/
/************************************/

func (c *Context) reset() {
	// add by sdx 2021.01.06
	c.matchRst.ptrNode = nil
	c.matchRst.params = c.Params
	c.matchRst.tsr = false

	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.Accepted = nil

	c.Pms = nil
	c.Params = c.Params[0:0]
	c.queryCache = nil
	c.formCache = nil
}

// 如果在当前请求上下文中需要新建goroutine，那么新的 goroutine 中必须要用 copy 后的 Context
// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This has to be used when the context has to be passed to a goroutine.
func (c *Context) Copy() *Context {
	cp := Context{
		GFResponse: c.GFResponse,
		ReqW:       c.ReqW,
		Params:     c.Params,
		matchRst:   c.matchRst,
		Pms:        c.Pms,
	}
	cp.ResW.ResponseWriter = nil

	cp.Keys = map[string]interface{}{}
	for k, v := range c.Keys {
		cp.Keys[k] = v
	}
	paramCopy := make([]Param, len(cp.Params))
	copy(paramCopy, cp.Params)
	cp.Params = paramCopy
	return &cp
}
