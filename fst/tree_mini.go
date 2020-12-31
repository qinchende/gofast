// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

// 用新的数据结构重建整颗路由树，用数组实现的树结构
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 自定义数据结构存放 所有的 路由树相关信息，全部通过数组索引的方式来访问
type miniNode struct {
	// 前缀字符（3字节）
	matchLen   uint8
	matchStart uint16

	// 子节点（3字节）
	childLen   uint8
	childStart uint16

	// 分组事件索引+当前路由事件索引（4字节）
	hdsGroupIdx uint16
	hdsItemIdx  uint16

	// 节点类型 （1字节）
	nType uint8
	//wildChild bool // 下一个节点是否为通配符
}

// 重建生成 mini 版本的 路由树
func rebuildMethodTree(mTree *methodTree) {
	if mTree.root == nil {
		return
	}
	rootNodeIdx := fstMem.allMiniLen
	fstMem.allMiniLen += 1
	mTree.root.rebuildNode(rootNodeIdx)
	mTree.miniRoot = &fstMem.allMiniNodes[rootNodeIdx]
}

// 将原始路由树节点 转换成 mini版本的节点
func (n *radixNode) rebuildNode(idx uint16) {
	// 通过索引，找到新节点的存储位置
	newMini := &fstMem.allMiniNodes[idx]

	// 处理前缀字符串
	mLen := uint8(len(n.match))
	fstMem.treeCharT = append(fstMem.treeCharT, n.match...)
	newMini.matchLen = mLen
	newMini.matchStart = fstMem.treeCharsLen
	fstMem.treeCharsLen += uint16(mLen)

	// 为子节点分配内存位置
	newMini.childStart = fstMem.allMiniLen
	newMini.childLen = uint8(len(n.children))
	fstMem.allMiniLen += uint16(newMini.childLen)

	for i := uint16(0); i < uint16(newMini.childLen); i++ {
		n.children[i].rebuildNode(newMini.childStart + i)
	}

	// 第一种：如果为绑定事件的节点 (能匹配一个路由)
	if n.routerItem != nil {
		newMini.hdsGroupIdx = n.routerItem.parent.hdsGroupIdx
		newMini.hdsItemIdx = n.routerItem.rebuildHandlers()
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
// 将所有的 中间件函数 放入一个大的数组当中，以后都是通过数组索引来访问函数
func addCHandlers(hds CHandlers) (idxes []uint16) {
	// 所有处理函数的切片
	fstMem.allCHandlers = append(fstMem.allCHandlers, hds...)
	hLen := uint16(len(hds))
	idxes = make([]uint16, hLen, hLen)
	for i := uint16(0); i < hLen; i++ {
		idxes[i] = fstMem.allCHdsLen + i
	}
	fstMem.allCHdsLen += hLen
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 重新生成路由树相关数据结构
func (ft *Faster) rebuildRoutes() {
	// 分配内存空间
	allocateMemSpace(ft)
	// 合并分组事件到最后一个分组当中
	gpHandlersCombine(ft.home)

	// 重建路由树
	rebuildMethodTree(ft.treeGet)
	rebuildMethodTree(ft.treePost)
	for _, mTree := range ft.treeOthers {
		rebuildMethodTree(&mTree)
	}

	// 将临时字符串byte数组 一次性转换成 string
	fstMem.treeChars = string(fstMem.treeCharT)
	fstMem.treeCharT = nil

	// TODO: 释放掉原始树的资源
	ft.treeGet.root = nil
	ft.treePost.root = nil
	for _, mTree := range ft.treeOthers {
		mTree.root = nil
	}
	ft.home = nil
}

// 计算所有要预分配的内存空间
func allocateMemSpace(ft *Faster) {
	totalNodes := ft.treeGet.nodeCt
	nodeStrLen := ft.treeGet.nodeStrLen

	totalNodes += ft.treePost.nodeCt
	nodeStrLen += ft.treePost.nodeStrLen

	for _, mTree := range ft.treeOthers {
		totalNodes += mTree.nodeCt
		nodeStrLen += mTree.nodeStrLen
	}

	// 第一种：处理函数节点空间
	hdsNodesCt := fstMem.hdsGroupCt + fstMem.hdsItemCt
	fstMem.hdsNodes = make([]handlersNode, hdsNodesCt, hdsNodesCt)
	// 第二种：处理事件
	//fstMem.hdsNodesMini = make([]handlersNodeMini, fstMem.hdsItemCt, fstMem.hdsItemCt)

	// 新的函数指针数组
	fstMem.hdsList = make(CHandlers, fstMem.allCHdsLen, fstMem.allCHdsLen)
	// 路由树 字符串
	fstMem.treeCharT = make([]byte, 0, nodeStrLen)
	fstMem.allMiniNodes = make([]miniNode, totalNodes, totalNodes)
}

// 合并所有路由分组的事件到 最后的一个分组中
func gpHandlersCombine(gp *RouterGroup) {
	if gp.children == nil {
		gp.rebuildHandlers()
		fstMem.hdsGroupCt++
		return
	}
	for _, ch := range gp.children {
		// TODO: 这里要补充完整所有的事件类型
		ch.eValidHds = append(gp.eValidHds, ch.eValidHds...)
		ch.eBeforeHds = append(gp.eBeforeHds, ch.eBeforeHds...)
		// 分水岭
		ch.eAfterHds = append(ch.eAfterHds, gp.eAfterHds...)
		ch.eSendHds = append(ch.eSendHds, gp.eSendHds...)
		//ch.eResponseHds = append(ch.eResponseHds, gp.eResponseHds...)

		gpHandlersCombine(ch)
	}
}

// 第一套 分开方案
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 为每个最后一级的分组，生成mini版的事件节点
func (gp *RouterGroup) rebuildHandlers() {
	setNewNode(&gp.routeEvents)

	gp.hdsGroupIdx = fstMem.hdsNodesLen
	fstMem.hdsNodesLen++
}

func (ri *RouterItem) rebuildHandlers() (idx uint16) {
	setNewNode(&ri.routeEvents)
	idx = fstMem.hdsNodesLen
	fstMem.hdsNodesLen++
	return
}

func setNewNode(re *routeEvents) {
	node := &fstMem.hdsNodes[fstMem.hdsNodesLen]

	// 获取所有的 handlers
	node.hdsLen, node.hdsIdx = tidyEventHandlers(&re.eHds)
	node.beforeLen, node.beforeIdx = tidyEventHandlers(&re.eBeforeHds)
	node.afterLen, node.afterIdx = tidyEventHandlers(&re.eAfterHds)
	node.validLen, node.validIdx = tidyEventHandlers(&re.eValidHds)
	node.sendLen, node.sendIdx = tidyEventHandlers(&re.eSendHds)
	//node.responseLen, node.responseIdx = tidyEventHandlers(&re.eResponseHds)
}

func tidyEventHandlers(hds *[]uint16) (ct uint8, startIdx uint16) {
	ct = uint8(len(*hds))
	startIdx = fstMem.hdsListLen
	for i := uint8(0); i < ct; i++ {
		fstMem.hdsList[fstMem.hdsListLen] = fstMem.allCHandlers[(*hds)[i]]
		fstMem.hdsListLen++
	}
	return
}

// 第二套 合并 方案
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ri *RouterItem) combineHandlers() (idx uint16) {
	node := &fstMem.hdsNodesMini[fstMem.hdsNodesMiniLen]
	gp := ri.parent

	node.startIdx = fstMem.hdsListLen
	node.hdsLen += tidyEventHandlersMini(&gp.eValidHds)
	node.hdsLen += tidyEventHandlersMini(&ri.eValidHds)
	node.hdsLen += tidyEventHandlersMini(&gp.eBeforeHds)
	node.hdsLen += tidyEventHandlersMini(&ri.eBeforeHds)
	node.hdsLen += tidyEventHandlersMini(&ri.eHds)
	node.hdsLen += tidyEventHandlersMini(&ri.eAfterHds)
	node.hdsLen += tidyEventHandlersMini(&gp.eAfterHds)
	node.hdsLen += tidyEventHandlersMini(&ri.eSendHds)
	node.hdsLen += tidyEventHandlersMini(&gp.eSendHds)
	//node.hdsLen += tidyEventHandlersMini(&ri.eResponseHds)
	//node.hdsLen += tidyEventHandlersMini(&gp.eResponseHds)

	idx = fstMem.hdsNodesMiniLen
	fstMem.hdsNodesMiniLen++
	return
}

func tidyEventHandlersMini(hds *[]uint16) (ct uint8) {
	ct = uint8(len(*hds))
	for i := uint8(0); i < ct; i++ {
		fstMem.hdsList[fstMem.hdsListLen] = fstMem.allCHandlers[(*hds)[i]]
		fstMem.hdsListLen++
	}
	return
}
