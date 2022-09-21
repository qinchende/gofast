// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 强制路由匹配走404逻辑
func (c *Context) RouteAs404() {
	c.route.ptrNode = c.myApp.miniNode404
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 第一种方案：将可执行中间件分类，依次执行。
// 方案1. 依次执行分组和节点自己的事件中间件函数
func (c *Context) execHandlers() {
	if c.execIdx == maxRouteHandlers {
		return
	}

	c.RouteIdx = c.route.ptrNode.routeIdx
	c.handlers = c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	c.execIdx = -1
	c.Next()
}

// 执行下一个中间件函数
func (c *Context) Next() {
	c.execIdx++
	for c.execIdx < int8(len(c.handlers.hdsIdxChain)) {
		c.myApp.fstMem.tidyHandlers[c.handlers.hdsIdxChain[c.execIdx]](c)
		// 可能被设置成了 abort ，这样后面的 handlers 不用再调用了
		if c.execIdx == maxRouteHandlers {
			break
		}
		c.execIdx++
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c *Context) execAfterMatchHandlers() {
	if c.route.ptrNode == nil {
		return
	}
	it := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	for it.afterMatchLen > 0 {
		c.myApp.fstMem.tidyHandlers[it.afterMatchIdx](c)
		it.afterMatchLen--
		it.afterMatchIdx++
	}
}

// NOTE: 下面的钩子函数不需要中断执行链。
func (c *Context) execBeforeSendHandlers() {
	if c.route.ptrNode == nil {
		return
	}
	it := c.handlers // c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]

	// 5.preSend
	for it.beforeSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[it.beforeSendIdx](c)
		it.beforeSendLen--
		it.beforeSendIdx++
	}
	for gp.beforeSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[gp.beforeSendIdx](c)
		gp.beforeSendLen--
		gp.beforeSendIdx++
	}
	//over:
	//	return
}

func (c *Context) execAfterSendHandlers() {
	if c.route.ptrNode == nil {
		return
	}
	it := c.handlers // c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]

	// 6.afterSend
	for it.afterSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[it.afterSendIdx](c)
		it.afterSendLen--
		it.afterSendIdx++
	}
	for gp.afterSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[gp.afterSendIdx](c)
		gp.afterSendLen--
		gp.afterSendIdx++
	}
	//over:
	//	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (c CtxHandlers) Last() CtxHandler {
//	if length := len(c); length > 0 {
//		return c[length-1]
//	}
//	return nil
//}

//// TODO: 第二种方案：将所有中间件组成链式（暂时不用）（实际上也无法将所有类型的filter串联起来）
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 方案2. 基于已经将所有的事件函数组织成了一个有序的索引数组。只需要一次循环就执行所有的中间件函数
//// 这种实现其实不现实，不同类型的事件是在框架封装过程中分开执行的(比如 before-render | after-render 是无法和 handler 串联
//// 到一起的，它们只是 handler 中 render 的前后执行函数。)
//func (c *Context) execHandlersMini(ptrMini *radixMiniNode) {
//	it := c.myApp.fstMem.hdsNodesPlan2[ptrMini.hdsItemIdx]
//
//	for it.hdsLen > 0 {
//		if c.aborted {
//			goto over
//		}
//		c.myApp.fstMem.hdsList[it.startIdx](c)
//		it.hdsLen--
//		it.startIdx++
//	}
//over:
//	return
//}
