// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/tips"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// 异常处理逻辑的接口定义
type (
	PanicHandler interface {
		Callback()
	}
	PanicFunc struct {
		Func func()
	}
)

func (pw PanicFunc) Callback() { pw.Func() }

func NewPanicPet(fn func()) *PanicFunc {
	return &PanicFunc{Func: fn}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// Context is the most important part of GoFast. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	myApp *GoFast // 用于上下文

	EnterTime  time.Duration  // 请求起始时间
	ResWrap    *ResponseWrap  // 被封装后的Response
	ReqRaw     *http.Request  // 原始 request
	Sess       SessionKeeper  // Session数据，数据存储部分可以自定义
	UrlParams  *routeParams   // : 或 * 对应的参数
	Pms        cst.KV         // 所有Request参数的map（queryCache + formCache）一般用于构造model对象
	PanicPet   PanicHandler   // 业务逻辑异常之后的处理
	CarryItems tips.CarryList // []*CarryItem，可以携带扩展的自定义数据

	queryCache url.Values   // param query result from c.ReqRaw.URL.Query()
	formCache  url.Values   // the parsed form data from POST, PATCH, or PUT body parameters.
	rwLock     sync.RWMutex // This mutex protect context

	route    matchRoute   // 路由匹配结果，[UrlParams] ? 一般用于确定相应资源
	handlers handlersNode // 匹配到的执行链标记
	execIdx  int8         // 执行链的索引 不能大于 127 个
	rendered bool         // 是否已经执行了Render

	RouteIdx uint16 // route的唯一标识ID，方便区分不同的route
}

/************************************/
/********** context creation ********/
/************************************/

func (c *Context) reset() {
	//c.EnterTime = timex.Now()
	//c.ResWrap = nil
	//c.ReqRaw = nil
	c.Sess = nil
	c.UrlParams = nil
	c.Pms = nil
	c.CarryItems = c.CarryItems[0:0]
	c.PanicPet = nil
	c.RouteIdx = 0

	// add by sdx 2021.01.06
	c.route.ptrNode = nil
	if c.route.params == nil {
		c.route.params = new(routeParams)
	}
	*c.route.params = (*c.route.params)[0:0]
	c.route.rts = false
	//c.handlers = nil
	c.execIdx = -1 // 当前不处于任何执行函数
	c.rendered = false

	c.queryCache = nil
	c.formCache = nil
	//c.mu = nil
}
