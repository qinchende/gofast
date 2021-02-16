// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 绑定在 RouterGroup 和 RouterItem 上的 不同事件处理函数数组
// RouterGroup 上的事件处理函数 最后需要作用在 RouterItem 上才会有实际的意义
// 事件要尽量少一些，每个路由节点都要分配一个对象
type routeEvents struct {
	// 下面的事件类型，按照执行顺序排列
	ePreValidHds  []uint16
	eBeforeHds    []uint16
	eHds          []uint16
	eAfterHds     []uint16
	ePreSendHds   []uint16
	eAfterSendHds []uint16
}

type RouterGroup struct {
	routeEvents
	gftApp       *GoFast
	prefix       string
	children     []*RouterGroup
	hdsGroupIdx  int16  // 记录当前分组 对应新事件数组中的起始位置索引
	selfHdsLen   uint16 // 记录分组中一共加入的 处理 函数个数
	parentHdsLen uint16 //
}

type RouterItem struct {
	routeEvents
	parent *RouterGroup
}

// 每一种事件类型需要占用3个字节(开始索引2字节 + 长度1字节(长度最大255))
// 这里抽象出N种事件类型，应该够用了，这样每个路由节点占用3*N字节空间，64位机器1字长是8字节
// RouterGroup 和 RouterItem 都用这一组数据结构记录事件处理函数
type handlersNode struct {
	validIdx     uint16
	beforeIdx    uint16
	hdsIdx       uint16
	afterIdx     uint16
	preSendIdx   uint16
	afterSendIdx uint16

	validLen     uint8
	beforeLen    uint8
	afterLen     uint8
	hdsLen       uint8
	preSendLen   uint8
	afterSendLen uint8
}

// ++++++++++++++++++++++++++++++++++++++++++++++
// 第二种方案（暂时不用）
// 将某个路由节点的所有处理函数按顺序全部排序成数组，请求匹配到路由节点之后直接执行这里的队列即可
// 当节点多的时候这种方式相对第一种占用更多内存。
type handlersNodePlan2 struct {
	startIdx uint16 // 2字节
	hdsLen   uint8  // 1字节
}
