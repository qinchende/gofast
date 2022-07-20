// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Context is the most important part of GoFast. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	myApp *GoFast // 用于上下文

	Errors    errMessages   // []*Error
	EnterTime time.Duration // 请求起始时间
	ResWrap   *ResponseWrap
	ReqRaw    *http.Request // request
	Sess      *CtxSession   // Session数据，数据存储部分可以自定义
	Params    *Params       // : 或 * 对应的参数
	Pms       cst.KV        // 所有Request参数的map（queryCache + formCache）一般用于构造model对象
	//PmsCarry  cst.KV        // 后台上下文传递的数据

	queryCache url.Values   // param query result from c.ReqRaw.URL.Query()
	formCache  url.Values   // the parsed form data from POST, PATCH, or PUT body parameters.
	mu         sync.RWMutex // This mutex protect Keys map

	handlers  handlersNode // 匹配到的执行链标记
	match     matchResult  // 路由匹配结果，[Params] ? 一般用于确定相应资源
	execIdx   int8         // 执行链的索引 不能大于 127 个
	rendered  bool         // 是否已经执行了Render
	IsTimeout bool         // 请求是否超时了
	RouteIdx  uint16       // router的唯一标识ID，方便区分不同的router
}

/************************************/
/********** context creation ********/
/************************************/

func (c *Context) reset() {
	c.Errors = nil
	//c.EnterTime = timex.Now()
	//c.ResWrap = nil
	//c.ReqRaw = nil
	c.Sess = nil
	c.Params = c.match.params
	c.Pms = nil
	//c.PmsCarry = nil
	c.RouteIdx = 0
	c.IsTimeout = false

	// add by sdx 2021.01.06
	c.match.ptrNode = nil
	if c.match.params == nil {
		c.match.params = new(Params)
	}
	*c.match.params = (*c.match.params)[0:0]
	c.match.rts = false
	c.match.allowRTS = c.myApp.RedirectTrailingSlash
	//c.handlers = nil
	c.execIdx = math.MaxInt8
	c.rendered = false
	c.queryCache = nil
	c.formCache = nil
	//c.mu = nil
}
