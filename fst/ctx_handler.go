// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// TODO: 第一种方案：将可执行中间件分类，依次执行。
// 方案1. 依次执行分组和节点自己的事件中间件函数
func (c *Context) execHandlers() {
	c.handlers = c.gftApp.fstMem.hdsNodes[c.match.ptrNode.hdsItemIdx]
	c.execIdx = 0
	c.Next()
}

// 执行下一个拦截器
func (c *Context) Next() {
	for c.execIdx < uint8(len(c.handlers.hdsIdxChain)) {
		c.gftApp.fstMem.tidyHandlers[c.handlers.hdsIdxChain[c.execIdx]](c)
		c.execIdx++
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// NOTE: 下面的钩子函数不需要中断执行链。

func (c *Context) execPreSendHandlers() {
	if c.match.ptrNode == nil {
		return
	}
	it := c.handlers // c.gftApp.fstMem.hdsNodes[c.match.ptrNode.hdsItemIdx]
	gp := c.gftApp.fstMem.hdsNodes[c.match.ptrNode.hdsGroupIdx]

	// 5.preSend
	for it.preSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.gftApp.fstMem.tidyHandlers[it.preSendIdx](c)
		it.preSendLen--
		it.preSendIdx++
	}
	for gp.preSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.gftApp.fstMem.tidyHandlers[gp.preSendIdx](c)
		gp.preSendLen--
		gp.preSendIdx++
	}
	//over:
	//	return
}

func (c *Context) execAfterSendHandlers() {
	if c.match.ptrNode == nil {
		return
	}
	it := c.handlers // c.gftApp.fstMem.hdsNodes[c.match.ptrNode.hdsItemIdx]
	gp := c.gftApp.fstMem.hdsNodes[c.match.ptrNode.hdsGroupIdx]

	// 6.afterSend
	for it.afterSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.gftApp.fstMem.tidyHandlers[it.afterSendIdx](c)
		it.afterSendLen--
		it.afterSendIdx++
	}
	for gp.afterSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.gftApp.fstMem.tidyHandlers[gp.afterSendIdx](c)
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
//	it := c.gftApp.fstMem.hdsNodesPlan2[ptrMini.hdsItemIdx]
//
//	for it.hdsLen > 0 {
//		if c.aborted {
//			goto over
//		}
//		c.gftApp.fstMem.hdsList[it.startIdx](c)
//		it.hdsLen--
//		it.startIdx++
//	}
//over:
//	return
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (c *Context) execHandlers() {
//	c.handlers = c.gftApp.fstMem.hdsNodes[c.match.ptrNode.hdsItemIdx]
//	c.execIdx = 0
//	c.Next()
//
//	//it := c.gftApp.fstMem.hdsNodes[c.match.ptrNode.hdsItemIdx]
//	//gp := c.gftApp.fstMem.hdsMiniNodes[c.match.ptrNode.hdsGroupIdx]
//
//	//// 2.before
//	//for gp.beforeLen > 0 {
//	//	if c.aborted {
//	//		goto over
//	//	}
//	//	c.gftApp.fstMem.hdsList[gp.beforeIdx](c)
//	//	gp.beforeLen--
//	//	gp.beforeIdx++
//	//}
//	//for it.beforeLen > 0 {
//	//	if c.aborted {
//	//		goto over
//	//	}
//	//	c.gftApp.fstMem.hdsList[it.beforeIdx](c)
//	//	it.beforeLen--
//	//	it.beforeIdx++
//	//}
//
//	//// 3.handler
//	//for it.hdsLen > 0 {
//	//	if c.aborted {
//	//		goto over
//	//	}
//	//	c.gftApp.fstMem.tidyHandlers[it.hdsIdx](c)
//	//	it.hdsLen--
//	//	it.hdsIdx++
//	//}
//
//	//// 4.after
//	//for it.afterLen > 0 {
//	//	if c.aborted {
//	//		goto over
//	//	}
//	//	c.gftApp.fstMem.hdsList[it.afterIdx](c)
//	//	it.afterLen--
//	//	it.afterIdx++
//	//}
//	//for gp.afterLen > 0 {
//	//	if c.aborted {
//	//		goto over
//	//	}
//	//	c.gftApp.fstMem.hdsList[gp.afterIdx](c)
//	//	gp.afterLen--
//	//	gp.afterIdx++
//	//}
//	//over:
//	//	return
//}

//// 可以指定任何路由节点中的 handlers 来执行
//func (c *Context) execJustHandlers(ptrMini *radixMiniNode) {
//	it := c.gftApp.fstMem.hdsNodes[ptrMini.hdsItemIdx]
//
//	// 3.handler
//	for it.hdsLen > 0 {
//		if c.aborted {
//			return
//		}
//		c.gftApp.fstMem.tidyHandlers[it.hdsIdx](c)
//		it.hdsLen--
//		it.hdsIdx++
//	}
//}

//func (c *Context) execPreBindHandlers() {
//	if c.match.ptrNode == nil {
//		return
//	}
//	it := c.gftApp.fstMem.hdsMiniNodes[c.match.ptrNode.hdsItemIdx]
//	gp := c.gftApp.fstMem.hdsMiniNodes[c.match.ptrNode.hdsIdx]
//
//	// 1.valid
//	for gp.validLen > 0 {
//		if c.aborted {
//			goto over
//		}
//		c.gftApp.fstMem.hdsList[gp.validIdx](c)
//		gp.validLen--
//		gp.validIdx++
//	}
//	for it.validLen > 0 {
//		if c.aborted {
//			goto over
//		}
//		c.gftApp.fstMem.hdsList[it.validIdx](c)
//		it.validLen--
//		it.validIdx++
//	}
//over:
//	return
//}
