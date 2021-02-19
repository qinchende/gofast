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
	// 请求前置拦截器 要用到的上下文
	*GFResponse
	ReqW *http.Request

	//gftApp   *GoFast
	//resW     GFResponse
	//Reply    *ResWriteWrap
	Params   Params
	matchRst matchResult

	// This mutex protect Keys map
	mu sync.RWMutex

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]interface{}

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	// queryCache use url.ParseQuery cached the param query result from c.ReqW.URL.Query()
	queryCache url.Values

	// formCache use url.ParseQuery cached PostForm contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	formCache url.Values

	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite
}

/************************************/
/********** context creation ********/
/************************************/

func (c *Context) reset() {
	c.Params = c.Params[0:0]

	// add by sdx 2021.01.06
	c.matchRst.ptrNode = nil
	c.matchRst.params = c.Params
	c.matchRst.tsr = false

	c.Keys = nil
	c.Errors = c.Errors[0:0]
	c.Accepted = nil
	c.queryCache = nil
	c.formCache = nil
}
//
//// Copy returns a copy of the current context that can be safely used outside the request's scope.
//// This has to be used when the context has to be passed to a goroutine.
//func (c *Context) Copy() *Context {
//	cp := Context{
//		//gftApp:   c.gftApp,
//		GFResponse: c.GFResponse,
//		ReqW:       c.ReqW,
//		Params:     c.Params,
//		matchRst:   c.matchRst,
//	}
//	cp.ResW.ResponseWriter = nil
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
