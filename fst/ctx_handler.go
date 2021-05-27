// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// TODO： 底层大致有两种设计思路，目前采用第一种方案， 没有优化
// 方案1. 依次执行分组和节点自己的事件中间件函数
func (c *Context) execHandlers(ptrMini *radixMiniNode) {
	it := c.gftApp.fstMem.hdsMiniNodes[ptrMini.hdsItemIdx]
	gp := c.gftApp.fstMem.hdsMiniNodes[ptrMini.hdsGroupIdx]

	// 2.before
	for gp.beforeLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[gp.beforeIdx](c)
		gp.beforeLen--
		gp.beforeIdx++
	}
	for it.beforeLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[it.beforeIdx](c)
		it.beforeLen--
		it.beforeIdx++
	}

	// 3.handler
	for it.hdsLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[it.hdsIdx](c)
		it.hdsLen--
		it.hdsIdx++
	}

	// 4.after
	for it.afterLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[it.afterIdx](c)
		it.afterLen--
		it.afterIdx++
	}
	for gp.afterLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[gp.afterIdx](c)
		gp.afterLen--
		gp.afterIdx++
	}
over:
	return
}

// 可以指定任何路由节点中的 handlers 来执行
func (c *Context) execJustHandlers(ptrMini *radixMiniNode) {
	it := c.gftApp.fstMem.hdsMiniNodes[ptrMini.hdsItemIdx]

	// 3.handler
	for it.hdsLen > 0 {
		c.gftApp.fstMem.hdsList[it.hdsIdx](c)
		it.hdsLen--
		it.hdsIdx++
	}
}

//
//func (c *Context) execPreBindHandlers() {
//	if c.matchRst.ptrNode == nil {
//		return
//	}
//	it := c.gftApp.fstMem.hdsMiniNodes[c.matchRst.ptrNode.hdsItemIdx]
//	gp := c.gftApp.fstMem.hdsMiniNodes[c.matchRst.ptrNode.hdsIdx]
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

//func (c *Context) execPreSendHandlers(code int, r render.Render) {
func (c *Context) execPreSendHandlers() {
	if c.matchRst.ptrNode == nil {
		return
	}
	it := c.gftApp.fstMem.hdsMiniNodes[c.matchRst.ptrNode.hdsItemIdx]
	gp := c.gftApp.fstMem.hdsMiniNodes[c.matchRst.ptrNode.hdsGroupIdx]

	// 5.preSend
	for it.preSendLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[it.preSendIdx](c)
		it.preSendLen--
		it.preSendIdx++
	}
	for gp.preSendLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[gp.preSendIdx](c)
		gp.preSendLen--
		gp.preSendIdx++
	}
over:
	return
}

func (c *Context) execAfterSendHandlers() {
	if c.matchRst.ptrNode == nil {
		return
	}
	it := c.gftApp.fstMem.hdsMiniNodes[c.matchRst.ptrNode.hdsItemIdx]
	gp := c.gftApp.fstMem.hdsMiniNodes[c.matchRst.ptrNode.hdsGroupIdx]

	// 6.afterSend
	for it.afterSendLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[it.afterSendIdx](c)
		it.afterSendLen--
		it.afterSendIdx++
	}
	for gp.afterSendLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[gp.afterSendIdx](c)
		gp.afterSendLen--
		gp.afterSendIdx++
	}
over:
	return
}

// TODO: 暂时不用
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 方案2. 基于已经将所有的事件函数组织成了一个有序的索引数组。只需要一次循环就执行所有的中间件函数
// 这种实现其实不现实，不同类型的事件是在框架封装过程中分开执行的
func (c *Context) execHandlersMini(ptrMini *radixMiniNode) {
	it := c.gftApp.fstMem.hdsNodesPlan2[ptrMini.hdsItemIdx]

	for it.hdsLen > 0 {
		if c.aborted {
			goto over
		}
		c.gftApp.fstMem.hdsList[it.startIdx](c)
		it.hdsLen--
		it.startIdx++
	}
over:
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c CtxHandlers) Last() CtxHandler {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
