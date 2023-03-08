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

// Context is the most important part of GoFast. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	myApp *GoFast // 用于上下文

	EnterTime time.Duration // 请求起始时间
	Res       *ResponseWrap // 被封装后的Response
	Req       *http.Request // 原始 request
	Sess      SessionKeeper // Session数据，数据存储部分可以自定义

	Pms        cst.KV     // 所有Request参数的map（queryCache + formCache）一般用于构造model对象
	formCache  url.Values // the parsed form data from POST, PATCH, or PUT body parameters.
	queryCache url.Values // param query result from c.Req.URL.Query()

	route    matchRoute   // 路由匹配结果，一般用于确定相应资源
	handlers handlersNode // 匹配到的执行链标记
	RouteIdx uint16       // route的唯一标识ID，方便区分不同的route
	execIdx  int8         // 执行链的索引 不能大于 127 个
	rendered bool         // 当前请求是否已经Render

	CarryItems tips.CarryList // []*CarryItem，可以携带扩展的自定义数据
	PanicPet   panicHandler   // 业务逻辑异常之后的处理
	rwLock     sync.RWMutex   // This mutex protect context
}

/************************************/
/********** context creation ********/
/************************************/

func (c *Context) reset() {
	// c.EnterTime = timex.Now()
	// c.Res = nil
	// c.Req = nil
	c.Sess = nil

	c.formCache = nil
	c.queryCache = nil
	c.Pms = nil

	// add by sdx 2021.01.06
	c.route.ptrNode = nil
	if c.route.params == nil {
		c.route.params = new(urlParams)
	}
	*c.route.params = (*c.route.params)[0:0]
	c.route.rts = false
	// c.handlers
	c.RouteIdx = 0
	c.execIdx = -1 // 当前不处于任何执行函数
	c.rendered = false

	c.CarryItems = c.CarryItems[0:0]
	c.PanicPet = nil
	// c.rwLock
}
