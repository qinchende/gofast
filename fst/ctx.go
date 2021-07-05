// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/render"
	"net/http"
	"net/url"
	"sync"
)

// Context is the most important part of GoFast. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	*GFResponse               // response (请求前置拦截器 要用到的上下文)
	ReqRaw      *http.Request // request

	Params     *Params                // : 或 * 对应的参数
	match      matchResult            // 路由匹配结果，[Params] ? 一般用于确定相应资源
	Pms        map[string]interface{} // 所有Request参数的map（queryCache + formCache）一般用于构造model对象
	queryCache url.Values             // param query result from c.ReqRaw.URL.Query()
	formCache  url.Values             // the parsed form data from POST, PATCH, or PUT body parameters.

	// Session数据，这里不规定Session的载体，可以自定义
	Sess *CtxSession
	// 设置成 true ，将中断后面的所有handlers
	aborted bool
	// render.Render 对象
	PRender *render.Render // render 对象
	PCode   *int           // status code

	// -----------------------------
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
	c.Keys = nil
	c.Sess = nil
	c.Accepted = nil

	// add by sdx 2021.01.06
	c.match.ptrNode = nil
	if c.match.params == nil {
		c.match.params = new(Params)
	}
	*c.match.params = (*c.match.params)[0:0]
	c.match.rts = false
	c.match.allowRTS = c.gftApp.RedirectTrailingSlash
	c.Params = c.match.params

	c.Pms = nil
	c.queryCache = nil
	c.formCache = nil
	c.aborted = false
}

//// 如果在当前请求上下文中需要新建goroutine，那么新的 goroutine 中必须要用 copy 后的 Context
//// Copy returns a copy of the current context that can be safely used outside the request's scope.
//// This has to be used when the context has to be passed to a goroutine.
//func (c *Context) Copy() *Context {
//	cp := Context{
//		GFResponse: c.GFResponse,
//		ReqRaw:     c.ReqRaw,
//		match:   c.match,
//		Pms:        c.Pms,
//		Sess:       c.Sess,
//		aborted:    c.aborted,
//	}
//	cp.ResWrap.ResponseWriter = nil
//
//	cp.Keys = map[string]interface{}{}
//	for k, v := range c.Keys {
//		cp.Keys[k] = v
//	}
//	paramCopy := make([]Param, len(cp.Params))
//	copy(paramCopy, cp.Params)
//	cp.Params = paramCopy
//	return &cp
//}

// 直接向上抛出异常，交给全局Recover函数处理
func (c *Context) panic() {
	panic("Handler exception!")
}
