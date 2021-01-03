// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

// 绑定在 RouterGroup 和 RouterItem 上的 不同事件处理函数数组
// RouterGroup 上的事件处理函数 最后需要作用在 RouterItem 上才会有实际的意义
// 事件要尽量少一些，每个路由节点都要分配一个对象
type routeEvents struct {
	// 下面的事件类型，按照执行顺序排列
	eValidHds     []uint16
	eBeforeHds    []uint16
	eHds          []uint16
	eAfterHds     []uint16
	eSendHds      []uint16
	eAfterSendHds []uint16

	//eRequestHds          []uint16
	//ePreParsingHds       []uint16
	//ePreHandlerHds    []uint16
	//ePreSerializationHds []uint16
	//eAfterSendHds []uint16
	//eTimeoutHds          []uint16
	//eErrorHds            []uint16
}

type RouterGroup struct {
	routeEvents
	gftApp       *GoFast
	prefix       string
	children     []*RouterGroup
	hdsGroupIdx  int16 	// 记录当前分组 对应新事件数组中的起始位置索引
	selfHdsLen   uint16	//
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
	sendIdx      uint16
	afterSendIdx uint16

	validLen     uint8
	beforeLen    uint8
	afterLen     uint8
	hdsLen       uint8
	sendLen      uint8
	afterSendLen uint8
}

// ++++++++++++++++++++++++++++++++++++++++++++++
// 第二种方案（暂时不用）
type handlersNodePlan2 struct {
	startIdx uint16 // 2字节
	hdsLen   uint8  // 1字节
}
