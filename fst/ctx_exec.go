// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

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
// 此时还没有执行execHandlers, c.handlers还没有赋值
func (c *Context) execAfterMatchHandlers() {
	it := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]

	for gp.afterMatchLen > 0 {
		c.myApp.fstMem.tidyHandlers[gp.afterMatchIdx](c)
		gp.afterMatchLen--
		gp.afterMatchIdx++
	}
	for it.afterMatchLen > 0 {
		c.myApp.fstMem.tidyHandlers[it.afterMatchIdx](c)
		it.afterMatchLen--
		it.afterMatchIdx++
	}
}

// NOTE: 下面的钩子函数不需要中断执行链。
func (c *Context) execBeforeSendHandlers() {
	it := c.handlers // c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]

	for gp.beforeSendLen > 0 {
		c.myApp.fstMem.tidyHandlers[gp.beforeSendIdx](c)
		gp.beforeSendLen--
		gp.beforeSendIdx++
	}
	for it.beforeSendLen > 0 {
		c.myApp.fstMem.tidyHandlers[it.beforeSendIdx](c)
		it.beforeSendLen--
		it.beforeSendIdx++
	}
}

func (c *Context) execAfterSendHandlers() {
	it := c.handlers // c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]

	for it.afterSendLen > 0 {
		c.myApp.fstMem.tidyHandlers[it.afterSendIdx](c)
		it.afterSendLen--
		it.afterSendIdx++
	}
	for gp.afterSendLen > 0 {
		c.myApp.fstMem.tidyHandlers[gp.afterSendIdx](c)
		gp.afterSendLen--
		gp.afterSendIdx++
	}
}

//
//// 执行异常处理之前的需要调用的函数
//func (c *Context) ExecAfterPanicHandlers() {
//	it := c.handlers // c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
//	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]
//
//	for it.afterPanicLen > 0 {
//		c.myApp.fstMem.tidyHandlers[it.afterPanicIdx](c)
//		it.afterPanicLen--
//		it.afterPanicIdx++
//	}
//	for gp.afterPanicLen > 0 {
//		c.myApp.fstMem.tidyHandlers[gp.afterPanicIdx](c)
//		gp.afterPanicLen--
//		gp.afterPanicIdx++
//	}
//}
