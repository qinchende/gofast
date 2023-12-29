// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/httpx"
	"github.com/qinchende/gofast/store/dts"
	"github.com/qinchende/gofast/store/gson"
	"sync"
	"unsafe"
)

// Request Context 中可能用到的 对象资源池，复用内存对象，避免更多的GC，提高性能。
type reqPools struct {
	ctxPool   sync.Pool    // 第二级所有的handler有context对象传递(第一级是标准http handler，不需要缓冲池)
	pmsPools  []*sync.Pool // 单项值可能为nil, 不是所有路由都需要Pms缓冲，访问量不大，用Pool缓冲cst.KV无意义
	nilValues []any        // 类型零值数据，将来给脏内存初始化
}

func (gft *GoFast) initRoutePools() {
	rp := &gft.pools

	// 所有有效请求都需要用到 fst.Context 上下文信息
	rp.ctxPool.New = func() any {
		return &Context{
			myApp: gft,
			Res:   &httpx.ResponseWrap{},
			Req:   &httpx.RequestWrap{},
		}
	}

	// 可能有的路由需要用到 Pms收集器 的内存共用
	pls := make([]*sync.Pool, gft.RoutesLen())
	maxSize := 0
	for i := range pls {
		rh := rHandlers[i]

		// 不需要用到对象池
		// 因为 1. pms为一般的cst.KV  2. pms对象已有NewPms自定义创建函数
		if rh == nil || rh.pmsFunc != nil {
			continue
		}
		if len(rh.pmsFields) > 0 {
			if maxSize < len(rh.pmsFields) {
				maxSize = len(rh.pmsFields)
			}
			pls[i] = &sync.Pool{}
			pls[i].New = func() any {
				return &gson.GsonRow{} // 存放请求数据的对象
			}
		}
	}
	rp.nilValues = make([]any, maxSize)
	rp.pmsPools = pls
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Context Pool
func (rp *reqPools) getContext() *Context {
	return rp.ctxPool.Get().(*Context)
}

func (rp *reqPools) putContext(c *Context) {
	rp.putPms(c)
	rp.ctxPool.Put(c)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Pms Pool get
func (c *Context) createPms() {
	// 有自定义Pms解析对象
	ra := rHandlers[c.RouteIdx]
	if ra != nil && ra.pmsFunc != nil {
		c.Pms = ra.pmsFunc()
		return
	}

	// 找 GsonRow 的缓存
	pp := c.myApp.pools.pmsPools[c.RouteIdx]
	if pp == nil {
		newMP := make(cst.KV)
		c.Pms = &newMP // 默认使用map类型保存KV值，c.Pms必须是指针类型
		return
	}

	// 如果Pms是GsonRow类型，从缓冲池中取出对象复用
	gr := pp.Get().(*gson.GsonRow)
	if gr.Cls == nil {
		gr.Init(ra.pmsFields)
	} else {
		copy(gr.Row, c.myApp.pools.nilValues[:len(gr.Cls)]) // reset member
	}
	c.Pms = gr
}

func (rp *reqPools) putPms(c *Context) {
	pp := c.myApp.pools.pmsPools[c.RouteIdx]
	if pp != nil {
		pp.Put(c.Pms)
	}
}

// 直接绑定struct对应的内存地址
func (c *Context) PmsPointer() unsafe.Pointer {
	return (c.Pms).(*dts.StructKV).Ptr
}
