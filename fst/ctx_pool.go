// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/httpx"
	"github.com/qinchende/gofast/store/gson"
	"sync"
)

// Request Context 中可能用到的 对象资源池，复用内存对象，避免更多的GC，提高性能。
type webPools struct {
	myApp    *GoFast
	ctxPool  sync.Pool    // 第二级所有的handler有context对象传递(第一级是标准http handler，不需要缓冲池)
	pmsPools []*sync.Pool // 单项值可能为nil, 不是所有路由都需要Pms缓冲，访问量不大，用Pool缓冲cst.KV无意义
}

func (wp *webPools) initWebPools(gft *GoFast) {
	wp.myApp = gft

	wp.ctxPool.New = func() any {
		return &Context{
			myApp: wp.myApp,
			Res:   &httpx.ResponseWrap{},
			Req:   &httpx.RequestWrap{},
		}
	}

	wp.pmsPools = make([]*sync.Pool, wp.myApp.RoutesLen())
	pms := wp.pmsPools
	for i := range pms {
		if routesAttrs[i] == nil {
			continue
		}

		pms[i] = &sync.Pool{}
		pms[i].New = func() any {
			return &gson.GsonRow{} // 存放请求数据的对象
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Context Pool
func (wp *webPools) getContext() *Context {
	return wp.ctxPool.Get().(*Context)
}

func (wp *webPools) putContext(c *Context) {
	wp.putPms(c)
	wp.ctxPool.Put(c)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Pms Pool get
func (c *Context) newPms() cst.SuperKV {
	pmsPool := c.myApp.pools.pmsPools[c.RouteIdx]
	if pmsPool != nil {
		gr := pmsPool.Get().(*gson.GsonRow)
		if gr.Cls == nil {
			gr.Cls = routesAttrs[c.RouteIdx].PmsFields
			gr.Row = make([]any, gr.Len())
		} else {
			for i := range gr.Row {
				gr.Row[i] = nil // gr.Row reset value
			}
		}
		return gr // 如果Pms是GsonRow类型，从缓冲池中取出对象复用
	}
	return make(cst.KV) // 默认使用map类型保存KV值
}

func (wp *webPools) putPms(c *Context) {
	pmsPool := c.myApp.pools.pmsPools[c.RouteIdx]
	if pmsPool != nil {
		pmsPool.Put(c.Pms)
	}
}
