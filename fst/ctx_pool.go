// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/gson"
	"sync"
)

// Request Context 中可能用到的 对象资源池

type webPools struct {
	ctxPool   sync.Pool    // 第二级：Handler context pools (第一级是标准形式，不需要缓冲池)
	pmsPools  []*sync.Pool // 单项值可能为nil, 不是所有路由都需要Pms缓冲，访问量不大，用Pool缓冲cst.KV无意义
	pmsFields [][]string
}

func (hr *HomeRouter) initPools() {
	hr.pools.ctxPool.New = func() any {
		return &Context{
			myApp: hr.myApp,
			Res:   &ResponseWrap{},
		}
	}

	hr.pools.pmsPools = make([]*sync.Pool, hr.myApp.RoutesLen())
	pms := hr.pools.pmsPools
	pmsFls := hr.pools.pmsFields
	for i := range pms {
		pms[i] = &sync.Pool{}
		// 这里放的是
		pms[i].New = func() any {
			return &gson.GsonRow{
				Cls: pmsFls[i],
				Row: make([]any, len(pmsFls[i])),
			}
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Context Pool
func (wp *webPools) getContext() *Context {
	return wp.ctxPool.Get().(*Context)
}

func (wp *webPools) putContext(c *Context) {
	// request.context pms

	// request.context
	wp.ctxPool.Put(c)

}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Pms Pool
func (c *Context) getPms() cst.SuperKV {
	pmsPool := c.myApp.pools.pmsPools[c.RouteIdx]
	if pmsPool != nil {
		gr := pmsPool.Get().(*gson.GsonRow)
		// TODO: gr.Row reset value
		for i := range gr.Row {
			gr.Row[i] = nil
		}
		return gr
	}
	return make(cst.KV)
}

//func (wp *webPools) putPms(c *Context) {
//
//}
