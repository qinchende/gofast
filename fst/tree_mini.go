// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

// 用新的数据结构重建整棵路由树，用数组实现的树结构
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 自定义数据结构存放 所有的 路由树相关信息，全部通过数组索引的方式来访问
type radixMiniNode struct {
	// 前缀字符（3字节）
	matchLen   uint8
	matchStart uint16

	// 子节点（3字节）
	childLen   uint8
	childStart uint16

	// 分组事件索引+当前路由匹配节点事件索引（4字节）
	hdsGroupIdx int16
	hdsItemIdx  int16

	// 节点类型 （1字节）
	nType uint8
	//wildChild bool // 下一个节点是否为通配符
}

// 重建生成 mini 版本的 路由树
func rebuildMethodTree(mTree *methodTree) {
	if mTree.root == nil {
		return
	}
	rootNodeIdx := fstMem.allRadixMiniLen
	fstMem.allRadixMiniLen += 1
	mTree.root.rebuildNode(rootNodeIdx)
	mTree.miniRoot = &fstMem.allRadixMiniNodes[rootNodeIdx]
}

// 将原始路由树节点 转换成 mini版本的节点
func (n *radixNode) rebuildNode(idx uint16) {
	// 通过索引，找到新节点的存储位置
	newMini := &fstMem.allRadixMiniNodes[idx]

	// 处理前缀字符串
	mLen := uint8(len(n.match))
	fstMem.treeCharT = append(fstMem.treeCharT, n.match...)
	newMini.matchLen = mLen
	newMini.matchStart = fstMem.treeCharsLen
	fstMem.treeCharsLen += uint16(mLen)

	// 为子节点分配内存位置
	newMini.childStart = fstMem.allRadixMiniLen
	newMini.childLen = uint8(len(n.children))
	fstMem.allRadixMiniLen += uint16(newMini.childLen)

	for i := uint16(0); i < uint16(newMini.childLen); i++ {
		n.children[i].rebuildNode(newMini.childStart + i)
	}

	// 第一种：如果为绑定事件的节点 (能匹配一个路由)
	if n.routerItem != nil {
		newMini.hdsGroupIdx = n.routerItem.parent.hdsGroupIdx // 记录“分组”事件在 全局 事件队列中的 起始位置
		newMini.hdsItemIdx = n.routerItem.rebuildHandlers()   // 记录“节点”事件在 全局 事件队列中的 起始位置
	}
	// 第二种：按顺序合并所有事件
	//if n.routerItem != nil {
	//	newMini.hdsItemIdx = n.routerItem.combineHandlers()
	//}
	// 释放掉资源
	n.routerItem = nil

	// 节点类型 和 是否通配符
	newMini.nType = n.nType
	//newMini.wildChild = n.wildChild
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 重建特殊节点
func rebuildDefaultHandlers(home *GoFast) {
	// 第一种：如果为绑定事件的节点 (能匹配一个路由)
	home.miniNode404 = &radixMiniNode{}
	home.miniNode404.hdsGroupIdx = home.routerItem404.parent.hdsGroupIdx
	home.miniNode404.hdsItemIdx = home.routerItem404.rebuildHandlers()

	home.miniNode405 = &radixMiniNode{}
	home.miniNode405.hdsGroupIdx = home.routerItem405.parent.hdsGroupIdx
	home.miniNode405.hdsItemIdx = home.routerItem405.rebuildHandlers()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 将所有的 中间件函数 放入一个大的数组当中，以后都是通过数组索引来访问函数
func addCtxHandlers(hds CtxHandlers) (idxes []uint16) {
	// 所有处理函数的切片
	fstMem.allCtxHandlers = append(fstMem.allCtxHandlers, hds...)
	hLen := uint16(len(hds))
	idxes = make([]uint16, hLen, hLen)
	for i := uint16(0); i < hLen; i++ {
		idxes[i] = fstMem.allCtxHdsLen + i
	}
	fstMem.allCtxHdsLen += hLen
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 重新生成路由树相关数据结构
func (gft *GoFast) rebuildRoutes() {
	// 合并分组事件到下一级分组当中，返回所有节点新增父级节点事件的和
	parentHandlerSum := gpCombineHandlers(&gft.RouterGroup)
	// 合并之后，（多了多少个重复计算的事件 + 以前的所有事件个数）= 装事件数组的大小
	fstMem.hdsListLen = parentHandlerSum + fstMem.allCtxHdsLen
	// 分配内存空间
	allocateMemSpace(gft)

	// 1. 将分组事件 转换到 新版全局数组中
	gpRebuildHandlers(&gft.RouterGroup)
	// 2. 重建路由树 （这里面将节点事件 转换到 新版全局数组中）
	rebuildMethodTree(gft.treeGet)
	rebuildMethodTree(gft.treePost)
	for _, mTree := range gft.treeOthers {
		rebuildMethodTree(&mTree)
	}
	// 3. 重建特殊节点，比如 NoRoute | NoMethod
	rebuildDefaultHandlers(gft)

	// 将临时字符串byte数组 一次性转换成 string
	fstMem.treeChars = string(fstMem.treeCharT)
	fstMem.treeCharT = nil

	// TODO: 释放掉原始树的资源
	gft.treeGet.root = nil
	gft.treePost.root = nil
	for _, mTree := range gft.treeOthers {
		mTree.root = nil
	}
}

// 计算所有要预分配的内存空间
func allocateMemSpace(gft *GoFast) {
	totalNodes := gft.treeGet.nodeCt
	nodeStrLen := gft.treeGet.nodeStrLen

	totalNodes += gft.treePost.nodeCt
	nodeStrLen += gft.treePost.nodeStrLen

	for _, mTree := range gft.treeOthers {
		totalNodes += mTree.nodeCt
		nodeStrLen += mTree.nodeStrLen
	}

	// 第一种：处理函数节点空间
	// 所有承载事件处理函数的Node个数（包括Group 和 Item）
	hdsNodesCt := fstMem.hdsGroupCt + fstMem.hdsItemCt
	fstMem.hdsMiniNodes = make([]handlersNode, hdsNodesCt, hdsNodesCt)
	// 第二种：处理事件
	// fstMem.hdsNodesPlan2 = make([]handlersNodePlan2, fstMem.hdsItemCt, fstMem.hdsItemCt)

	// 新的 handlers 指针数组
	fstMem.hdsList = make(CtxHandlers, fstMem.hdsListLen, fstMem.hdsListLen)
	fstMem.hdsListLen = 0 // 下标重置成0，后面从这里把事件加入 hdsList

	// 路由树 字符串
	fstMem.treeCharT = make([]byte, 0, nodeStrLen)
	fstMem.allRadixMiniNodes = make([]radixMiniNode, totalNodes, totalNodes)
}

// 合并所有路由分组的事件到下一级的分组当中
// 返回所有节点新增加处理函数个数的和
func gpCombineHandlers(gp *RouterGroup) uint16 {
	// 所有分组个数
	fstMem.hdsGroupCt++
	if gp.children == nil {
		return gp.parentHdsLen
	}
	var allChildrenHdsCount uint16 = 0
	for _, chGroup := range gp.children {
		hdsCount := 0

		// TODO: 这里要补充完整所有的事件类型
		hdsCount += len(gp.ePreValidHds)
		chGroup.ePreValidHds = combineHandlers(gp.ePreValidHds, chGroup.ePreValidHds)

		hdsCount += len(gp.eBeforeHds)
		chGroup.eBeforeHds = combineHandlers(gp.eBeforeHds, chGroup.eBeforeHds)

		// 分水岭 -> item (not group) handlers

		hdsCount += len(gp.eAfterHds)
		chGroup.eAfterHds = append(chGroup.eAfterHds, gp.eAfterHds...)

		hdsCount += len(gp.ePreSendHds)
		chGroup.ePreSendHds = append(chGroup.ePreSendHds, gp.ePreSendHds...)

		hdsCount += len(gp.eAfterSendHds)
		chGroup.eAfterSendHds = append(chGroup.eAfterSendHds, gp.eAfterSendHds...)

		chGroup.parentHdsLen = uint16(hdsCount)
		allChildrenHdsCount += gpCombineHandlers(chGroup)
	}
	return allChildrenHdsCount + gp.parentHdsLen
}

// 将分组事件 转换到 新版全局数组中
func gpRebuildHandlers(gp *RouterGroup) {
	// 每一个分组的事件都去注册
	gp.rebuildHandlers()
	if gp.children == nil {
		return
	}
	for _, chGroup := range gp.children {
		gpRebuildHandlers(chGroup)
	}
}

// 第一套 分开方案
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 为每个最后一级的分组，将 routeEvent 变成内存占用更小的 handlersNode
func (gp *RouterGroup) rebuildHandlers() {
	setNewNode(&gp.routeEvents)

	gp.hdsGroupIdx = int16(fstMem.hdsMiniNodesLen)
	fstMem.hdsMiniNodesLen++
}

// 为每个路由节点，将 routeEvent 变成内存占用更小的 handlersNode
func (ri *RouterItem) rebuildHandlers() (idx int16) {
	setNewNode(&ri.routeEvents)
	idx = int16(fstMem.hdsMiniNodesLen)
	fstMem.hdsMiniNodesLen++
	return
}

// 将 routeEvents 中不同的事件类型对应的事件处理函数，全部存入全局数组中。
func setNewNode(re *routeEvents) {
	node := &fstMem.hdsMiniNodes[fstMem.hdsMiniNodesLen]

	// 获取所有的 handlers
	node.hdsLen, node.hdsIdx = tidyEventHandlers(&re.eHds)
	node.beforeLen, node.beforeIdx = tidyEventHandlers(&re.eBeforeHds)
	node.afterLen, node.afterIdx = tidyEventHandlers(&re.eAfterHds)
	node.validLen, node.validIdx = tidyEventHandlers(&re.ePreValidHds)
	node.preSendLen, node.preSendIdx = tidyEventHandlers(&re.ePreSendHds)
	node.afterSendLen, node.afterSendIdx = tidyEventHandlers(&re.eAfterSendHds)
}

// allCtxHandlers 中无序存放的 handlers 转入 有序的 hdsList 中
func tidyEventHandlers(hds *[]uint16) (ct uint8, startIdx uint16) {
	ct = uint8(len(*hds))
	startIdx = fstMem.hdsListLen
	for i := uint8(0); i < ct; i++ {
		fstMem.hdsList[fstMem.hdsListLen] = fstMem.allCtxHandlers[(*hds)[i]]
		fstMem.hdsListLen++
	}
	return
}

//// 第二套 合并 方案
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (ri *RouterItem) combineHandlers() (idx uint16) {
//	node := &fstMem.hdsNodesPlan2[fstMem.hdsNodesPlan2Len]
//	gp := ri.parent
//
//	node.startIdx = fstMem.hdsListLen
//	node.hdsLen += tidyEventHandlersMini(&gp.ePreValidHds)
//	node.hdsLen += tidyEventHandlersMini(&ri.ePreValidHds)
//	node.hdsLen += tidyEventHandlersMini(&gp.eBeforeHds)
//	node.hdsLen += tidyEventHandlersMini(&ri.eBeforeHds)
//	node.hdsLen += tidyEventHandlersMini(&ri.eHds)
//	node.hdsLen += tidyEventHandlersMini(&ri.eAfterHds)
//	node.hdsLen += tidyEventHandlersMini(&gp.eAfterHds)
//	node.hdsLen += tidyEventHandlersMini(&ri.ePreSendHds)
//	node.hdsLen += tidyEventHandlersMini(&gp.ePreSendHds)
//	//node.hdsLen += tidyEventHandlersMini(&ri.eResponseHds)
//	//node.hdsLen += tidyEventHandlersMini(&gp.eResponseHds)
//
//	idx = fstMem.hdsNodesPlan2Len
//	fstMem.hdsNodesPlan2Len++
//	return
//}
//
//func tidyEventHandlersMini(hds *[]uint16) (ct uint8) {
//	ct = uint8(len(*hds))
//	for i := uint8(0); i < ct; i++ {
//		fstMem.hdsList[fstMem.hdsListLen] = fstMem.allCtxHandlers[(*hds)[i]]
//		fstMem.hdsListLen++
//	}
//	return
//}
