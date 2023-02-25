// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 用新的数据结构重建整棵路由树，用数组实现的树结构
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 自定义数据结构存放 所有的 路由树相关信息，全部通过数组索引的方式来访问
// 目前是64位系统的两个字长，16字节
type radixMiniNode struct {
	// router index (2字节)
	routeIdx uint16

	// 前缀字符（3字节）
	matchStart uint16
	matchLen   uint8

	// 子节点（3字节）
	childStart uint16
	childLen   uint8

	// 分组事件索引+当前路由匹配节点事件索引（4字节）
	// 将来通过索引找到对应匹配的事件节点执行
	// 对应 fstMem.handlersNode 中的索引位置
	hdsGroupIdx int16
	hdsItemIdx  int16

	// 节点类型 （1字节）
	nType uint8
	// wildChild bool // 下一个节点是否为通配符

	// 因为下面这几种情况相对来说都少，预先判断，利于每次请求的执行效率
	hasAfterMatch bool // (1字节)
	hasBeforeSend bool // (1字节)
	hasAfterSend  bool // (1字节)
}

// 重建生成 mini 版本的 路由树
func rebuildMethodTree(fstMem *fstMemSpace, mTree *methodTree) {
	if mTree.root == nil {
		return
	}
	rootNodeIdx := fstMem.allRadixMiniLen
	fstMem.allRadixMiniLen += 1
	mTree.root.rebuildNode(fstMem, rootNodeIdx)
	mTree.miniRoot = &fstMem.allRadixMiniNodes[rootNodeIdx]
}

// 将原始路由树节点 转换成 mini版本的节点
func (n *radixNode) rebuildNode(fstMem *fstMemSpace, idx uint16) {
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
		n.children[i].rebuildNode(fstMem, newMini.childStart+i)
	}

	// 第一种：如果是一个路由叶子节点 (能匹配一个路由)
	if n.leafItem != nil {
		newMini.hdsGroupIdx = n.leafItem.group.hdsIdx     // 记录“分组”事件在 全局 事件队列中的 起始位置
		newMini.hdsItemIdx = n.leafItem.rebuildHandlers() // 记录“节点”事件在 全局 事件队列中的 起始位置
		newMini.routeIdx = n.leafItem.routeIdx
		combNodeHandlers(fstMem, newMini, true) // 构造执行链
	}
	// 释放掉资源
	n.leafItem = nil

	// 节点类型 和 是否通配符
	newMini.nType = n.nType
	//newMini.wildChild = n.wildChild

	if newMini.hdsGroupIdx > 0 && newMini.hdsItemIdx > 0 {
		hdsGroup := fstMem.hdsNodes[newMini.hdsGroupIdx]
		hdsItem := fstMem.hdsNodes[newMini.hdsItemIdx]
		if hdsGroup.afterMatchLen > 0 || hdsItem.afterMatchLen > 0 {
			newMini.hasAfterMatch = true
		}
		if hdsGroup.beforeSendLen > 0 || hdsItem.beforeSendLen > 0 {
			newMini.hasBeforeSend = true
		}
		if hdsGroup.afterSendLen > 0 || hdsItem.afterSendLen > 0 {
			newMini.hasAfterSend = true
		}
	}
}

// 默认特殊路径的路由执行链构造。
func combNodeHandlers(fstMem *fstMemSpace, miniNode *radixMiniNode, needGroup bool) {
	// 可以构造叶子节点的执行链切片，（注意：这里一定要用取地址符）
	eNode := &fstMem.hdsNodes[miniNode.hdsItemIdx]
	// 下面这两个不能取地址，而是值拷贝
	me := fstMem.hdsNodes[miniNode.hdsItemIdx]
	gp := fstMem.hdsNodes[miniNode.hdsGroupIdx]

	// 上级分组before + 自己before + 自己handler + 自己after + 上级分组after
	size := me.beforeLen + me.hdsLen + me.afterLen
	if needGroup {
		size += gp.beforeLen + gp.afterLen
	}
	eNode.hdsIdxChain = make([]uint16, size, size)

	// 2.before
	count := 0
	for needGroup && gp.beforeLen > 0 {
		eNode.hdsIdxChain[count] = gp.beforeIdx
		gp.beforeLen--
		count++
		gp.beforeIdx++
	}
	for me.beforeLen > 0 {
		eNode.hdsIdxChain[count] = me.beforeIdx
		me.beforeLen--
		count++
		me.beforeIdx++
	}
	// 3.handler
	for me.hdsLen > 0 {
		eNode.hdsIdxChain[count] = me.hdsIdx
		me.hdsLen--
		count++
		me.hdsIdx++
	}
	// 4.after
	for me.afterLen > 0 {
		eNode.hdsIdxChain[count] = me.afterIdx
		me.afterLen--
		count++
		me.afterIdx++
	}
	for needGroup && gp.afterLen > 0 {
		eNode.hdsIdxChain[count] = gp.afterIdx
		gp.afterLen--
		count++
		gp.afterIdx++
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 重建特殊节点
func rebuildSpecialHandlers(home *GoFast) {
	home.miniNodeAny = &radixMiniNode{routeIdx: home.allRoutes[0].routeIdx}
	home.miniNodeAny.hdsGroupIdx = home.allRoutes[0].group.hdsIdx
	home.miniNodeAny.hdsItemIdx = home.allRoutes[0].rebuildHandlers()
	combNodeHandlers(home.fstMem, home.miniNodeAny, true)

	home.miniNode404 = &radixMiniNode{routeIdx: home.allRoutes[1].routeIdx}
	home.miniNode404.hdsGroupIdx = home.allRoutes[1].group.hdsIdx
	home.miniNode404.hdsItemIdx = home.allRoutes[1].rebuildHandlers()
	combNodeHandlers(home.fstMem, home.miniNode404, true)

	home.miniNode405 = &radixMiniNode{routeIdx: home.allRoutes[2].routeIdx}
	home.miniNode405.hdsGroupIdx = home.allRoutes[2].group.hdsIdx
	home.miniNode405.hdsItemIdx = home.allRoutes[2].rebuildHandlers()
	combNodeHandlers(home.fstMem, home.miniNode405, true)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 将所有的 中间件函数 放入一个大的数组当中，以后都是通过数组索引来访问函数
func addCtxHandlers(fstMem *fstMemSpace, hds []CtxHandler) (idxes []uint16) {
	// 所有处理函数的切片
	fstMem.allCtxHandlers = append(fstMem.allCtxHandlers, hds...)

	hLen := uint16(len(hds))
	idxes = make([]uint16, hLen, hLen)
	for i := uint16(0); i < hLen; i++ {
		idxes[i] = fstMem.allCtxHdsLen + i
	}
	fstMem.allCtxHdsLen += hLen
	PanicIf(fstMem.allCtxHdsLen >= maxAllHandlers, "Too many handlers more than MaxUInt16.")
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// NOTE（ENTER）: 重新生成路由树相关数据结构
// 这里是生成数组版路由树的入口函数
func (gft *GoFast) buildMiniRoutes() {
	fstMem := gft.fstMem

	// TODO: 重置fstMem，方便重构路由树
	//fstMem.treeCharT = nil
	fstMem.hdsNodesLen = 0
	fstMem.treeCharsLen = 0
	fstMem.allRadixMiniLen = 0
	fstMem.routeGroupLen = 0
	fstMem.routeItemLen = uint16(len(gft.allRoutes))
	// end
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	// 合并分组事件到下一级分组当中，返回所有节点新增父级节点事件的和
	gft.combEvents = gft.RouteGroup.routeEvents
	parentHandlerSum := gpCombineHandlers(&gft.RouteGroup)
	gft.specialGroup.combEvents = gft.specialGroup.routeEvents
	parentHandlerSumSpecial := gpCombineHandlers(gft.specialGroup)
	// 合并之后，（多了多少个重复计算的事件 + 以前的所有事件个数）= 装事件数组的大小
	fstMem.tidyHdsLen = parentHandlerSum + parentHandlerSumSpecial + fstMem.allCtxHdsLen
	// 分配内存空间
	allocateMemSpace(gft)

	// 1. 将分组事件 转换到 新版全局数组中
	gpRebuildHandlers(&gft.RouteGroup)
	gpRebuildHandlers(gft.specialGroup)
	// 2. 重建路由树 （这里面将节点事件 转换到 新版全局数组中）
	for _, mTree := range gft.routerTrees {
		rebuildMethodTree(fstMem, mTree)
	}
	// 3. 重建特殊节点，比如 NoRoute | NoMethod
	rebuildSpecialHandlers(gft)

	// 将临时字符串byte数组 一次性转换成 string
	fstMem.treeChars = string(fstMem.treeCharT)

	// 所有路由节点URL
	gft.allPaths = gft.RoutePaths()

	// TODO: 释放掉原始树的资源，后面不可以根据这些树结构构造路由了。
	if !gft.IsDebugging() {
		fstMem.allCtxHandlers = nil
		fstMem.allCtxHdsLen = 0
		fstMem.treeCharT = nil

		for _, mTree := range gft.routerTrees {
			mTree.root = nil
		}

		// 将原始路由和分组删除
		gft.allRoutes = nil // allRoutes 不能清除，否则重要信息有丢失
		gft.RouteGroup.children = nil
	}
}

// 计算所有要预分配的内存空间
func allocateMemSpace(gft *GoFast) {
	fstMem := gft.fstMem

	var totalNodes, nodeStrLen uint16
	for _, mTree := range gft.routerTrees {
		totalNodes += mTree.nodeCt
		nodeStrLen += mTree.nodeStrLen
	}

	// 第一种：处理函数节点空间
	// 所有承载事件处理函数的Node个数（包括Group 和 Item）
	hdsNodesCt := fstMem.routeGroupLen + fstMem.routeItemLen
	fstMem.hdsNodes = make([]handlersNode, hdsNodesCt, hdsNodesCt)

	// 新的 handlers 指针数组
	PanicIf(fstMem.tidyHdsLen >= maxAllHandlers, "Chains tidy handlers more than MaxUInt16.")
	fstMem.tidyHandlers = make([]CtxHandler, fstMem.tidyHdsLen, fstMem.tidyHdsLen)
	fstMem.tidyHdsLen = 0 // 下标重置成0，后面从这里把事件加入 tidyHandlers

	// 路由树 字符串
	fstMem.treeCharT = make([]byte, 0, nodeStrLen)
	fstMem.allRadixMiniNodes = make([]radixMiniNode, totalNodes, totalNodes)
	for idx := 0; idx < len(fstMem.allRadixMiniNodes); idx++ {
		fstMem.allRadixMiniNodes[idx].routeIdx = 0
		fstMem.allRadixMiniNodes[idx].hdsGroupIdx = -1
		fstMem.allRadixMiniNodes[idx].hdsItemIdx = -1
	}
}

// 合并所有路由分组的事件到下一级的分组当中
// 返回所有节点新增加处理函数个数的和
func gpCombineHandlers(gp *RouteGroup) uint16 {
	// 所有分组个数
	gp.myApp.fstMem.routeGroupLen++
	if gp.children == nil {
		return gp.parentHdsLen
	}
	var allChildrenHdsCount uint16 = 0
	for _, chGroup := range gp.children {
		hdsCount := 0

		hdsCount += len(gp.combEvents.eAfterMatchHds)
		chGroup.combEvents.eAfterMatchHds = combineHandlers(gp.combEvents.eAfterMatchHds, chGroup.eAfterMatchHds)

		// TODO: 这里要补充完整所有的事件类型
		//hdsCount += len(gp.combEvents.ePreValidHds)
		//chGroup.combEvents.ePreValidHds = combineHandlers(gp.combEvents.ePreValidHds, chGroup.ePreValidHds)

		hdsCount += len(gp.combEvents.eBeforeHds)
		chGroup.combEvents.eBeforeHds = combineHandlers(gp.combEvents.eBeforeHds, chGroup.eBeforeHds)

		// 分水岭 -> item (not group) handlers

		hdsCount += len(gp.combEvents.eAfterHds)
		chGroup.combEvents.eAfterHds = append(chGroup.eAfterHds, gp.combEvents.eAfterHds...)

		hdsCount += len(gp.combEvents.eBeforeSendHds)
		chGroup.combEvents.eBeforeSendHds = append(chGroup.eBeforeSendHds, gp.combEvents.eBeforeSendHds...)

		hdsCount += len(gp.combEvents.eAfterSendHds)
		chGroup.combEvents.eAfterSendHds = append(chGroup.eAfterSendHds, gp.combEvents.eAfterSendHds...)

		chGroup.parentHdsLen = uint16(hdsCount)
		allChildrenHdsCount += gpCombineHandlers(chGroup)
	}
	return allChildrenHdsCount + gp.parentHdsLen
}

// 将分组事件 转换到 新版全局数组中
func gpRebuildHandlers(gp *RouteGroup) {
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
func (gp *RouteGroup) rebuildHandlers() {
	fstMem := gp.myApp.fstMem
	setNewNode(fstMem, &gp.combEvents)

	gp.hdsIdx = int16(fstMem.hdsNodesLen)
	fstMem.hdsNodesLen++
}

// 为每个路由节点，将 routeEvent 变成内存占用更小的 handlersNode
func (ri *RouteItem) rebuildHandlers() (idx int16) {
	fstMem := ri.group.myApp.fstMem
	setNewNode(fstMem, &ri.routeEvents)
	idx = int16(fstMem.hdsNodesLen)
	fstMem.hdsNodesLen++
	return
}

// 将 routeEvents 中不同的事件类型对应的事件处理函数，全部存入全局数组中。
func setNewNode(fstMem *fstMemSpace, re *routeEvents) {
	node := &fstMem.hdsNodes[fstMem.hdsNodesLen]

	node.afterMatchLen, node.afterMatchIdx = tidyEventHandlers(fstMem, &re.eAfterMatchHds)

	// 获取所有的 handlers (执行链)
	node.hdsLen, node.hdsIdx = tidyEventHandlers(fstMem, &re.eHds)
	node.beforeLen, node.beforeIdx = tidyEventHandlers(fstMem, &re.eBeforeHds)
	node.afterLen, node.afterIdx = tidyEventHandlers(fstMem, &re.eAfterHds)

	// 装饰器
	//node.validLen, node.validIdx = tidyEventHandlers(fstMem, &re.ePreValidHds)
	node.beforeSendLen, node.beforeSendIdx = tidyEventHandlers(fstMem, &re.eBeforeSendHds)
	node.afterSendLen, node.afterSendIdx = tidyEventHandlers(fstMem, &re.eAfterSendHds)
}

// allCtxHandlers 中无序存放的 handlers 转入 有序的 tidyHandlers 中
func tidyEventHandlers(fstMem *fstMemSpace, hds *[]uint16) (ct uint8, startIdx uint16) {
	// 多少个
	ct = uint8(len(*hds))
	// 新构造的所有函数指针切片，依次增长。
	startIdx = fstMem.tidyHdsLen
	for i := uint8(0); i < ct; i++ {
		fstMem.tidyHandlers[fstMem.tidyHdsLen] = fstMem.allCtxHandlers[(*hds)[i]]
		fstMem.tidyHdsLen++
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 合并两个 事件数组 到一个新的数组
func combineHandlers(a, b []uint16) []uint16 {
	size := len(a) + len(b)
	if size <= 0 {
		return nil
	}
	merge := make([]uint16, size)
	copy(merge, a)
	copy(merge[len(a):], b)
	return merge
}
