// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Context is the most important part of GoFast. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	//*GFResponse // response (请求前置拦截器 要用到的上下文)
	//ResWrap   *ResWriterWrap
	//Ctx       *Context
	EnterTime time.Duration // 请求起始时间
	gftApp    *GoFast       // 用于上下文
	//fitIdx    int
	Errors errMessages // []*Error

	ResWrap *ResWriterWrap
	ReqRaw  *http.Request // request

	Params   *Params                // : 或 * 对应的参数
	Pms      map[string]interface{} // 所有Request参数的map（queryCache + formCache）一般用于构造model对象
	match    matchResult            // 路由匹配结果，[Params] ? 一般用于确定相应资源
	handlers handlersNode           // 匹配到的执行链标记
	execIdx  int8                   // 执行链的索引 不能大于 127 个
	aborted  bool                   // 设置成 true ，将中断后面的所有handlers
	rendered bool                   // 是否已经执行了Render

	queryCache url.Values // param query result from c.ReqRaw.URL.Query()
	formCache  url.Values // the parsed form data from POST, PATCH, or PUT body parameters.

	// Session数据，数据存储部分可以自定义
	Sess *CtxSession

	// This mutex protect Keys map
	mu sync.RWMutex

	// -----------------------------Context对象占用内存越少越好，以下待定
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
	c.execIdx = math.MaxInt8

	c.Pms = nil
	c.queryCache = nil
	c.formCache = nil
	c.aborted = false
	c.rendered = false
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
