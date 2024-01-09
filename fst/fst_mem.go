// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 自定义内存数据库，存放路由树所有相关的数据
type fstMemSpace struct {
	app *GoFast

	// 我们需要自己定义一个切片，管理所有的 Context handlers.
	// 所有handler函数都需要加到这里来，形成一个全局的handler数组，以后任何路由都只记录这里的索引，执行时需要通过索引
	// 定位这里的handler函数，然后再执行。
	allCtxHandlers []CtxHandler // handler数组
	allCtxHdsLen   uint16       // 意味这所有 handler 不能超过 uint16 能标识的最大值

	// 新的handlers, 有序的, 按分组和事件类型排序
	// 将上面 allCtxHandlers (无序的)，转换成 tidyHandlers （有序的）, 按分组或路由节点整理处理函数
	tidyHandlers []CtxHandler
	tidyHdsLen   uint16

	routeGroupLen uint16 // 所有分组个数，网站根目录就是第一个分组
	routeItemLen  uint16 // 所有路由节点的个数，1个路由匹配就是一个ItemNode

	// 将路由树节点中的前缀字符 拼接 成一个大的字符串，以后所有路由查找都在这个字符串中
	treeCharT    []byte
	treeChars    string
	treeCharsLen uint16

	// 存放所有 radixMiniNode 数组，最终版的 Radix路由树数组实现方式（非链表）。
	allRadixMiniNodes []radixMiniNode
	allRadixMiniLen   uint16

	// handlers节点切片（专门用户记录 不同类型处理函数 的索引值）
	hdsNodes    []handlersNode // 专门记录事件处理函数的节点切片
	hdsNodesLen uint16         // 节点数量
}
