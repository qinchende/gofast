// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 绑定在 RouteGroup 和 RouteItem 上的 不同事件处理函数数组
// RouteGroup 上的事件处理函数 最后需要作用在 RouteItem 上才会有实际的意义
// 事件要尽量少一些，每个路由节点都要分配一个对象
// TODO: 此结构占用空间还是比较大的，可以考虑释放。
type routeEvents struct {
	// 下面的事件类型，按照执行顺序排列
	//ePreValidHds   []uint16
	eAfterMatchHds []uint16
	eBeforeHds     []uint16
	eHds           []uint16
	eAfterHds      []uint16
	eBeforeSendHds []uint16
	eAfterSendHds  []uint16
}

type RouteGroup struct {
	routeEvents              // 直接作用于本节点的事件可能为空
	combEvents   routeEvents // 合并父节点的分组事件，routeEvents可能为空，但是combEvents几乎不会为空
	myApp        *GoFast
	prefix       string
	children     []*RouteGroup
	hdsIdx       int16  // 记录当前分组 对应新事件数组中的起始位置索引
	selfHdsLen   uint16 // 记录当前分组中一共加入的处理函数个数（仅限于本分组加入的事件，不包含合并上级分组的）
	parentHdsLen uint16 // 记录所属上级分组的所有处理函数个数（仅包含上级分组，不含本分组的事件个数）
}

type RouteItem struct {
	group       *RouteGroup // router group
	method      string      // httpMethod
	fullPath    string      // 路由的完整路径
	routeEvents             // all handlers
	routeIdx    uint16      // 此路由在路由数组中的索引值
}

// 每一种事件类型需要占用3个字节(开始索引2字节 + 长度1字节(长度最大255))
// 这里抽象出N种事件类型，应该够用了，这样每个路由节点占用3*N字节空间，64位机器1字长是8字节
// RouteGroup 和 RouteItem 都用这一组数据结构记录事件处理函数
type handlersNode struct {
	hdsIdxChain []uint16 // 执行链的索引数组

	//validIdx      uint16
	afterMatchIdx uint16
	beforeIdx     uint16
	hdsIdx        uint16
	afterIdx      uint16
	beforeSendIdx uint16
	afterSendIdx  uint16

	//validLen      uint8
	afterMatchLen uint8
	beforeLen     uint8
	afterLen      uint8
	hdsLen        uint8
	beforeSendLen uint8
	afterSendLen  uint8
}
