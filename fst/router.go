package fst

type RouterGroup struct {
	routeEvents
	faster      *Faster
	prefix      string
	children    []*RouterGroup
	hdsGroupIdx uint16 // 记录当前组 对应 新事件结构的索引
}

type RouterItem struct {
	routeEvents
	parent *RouterGroup
}

// 每一种事件类型需要占用3个字节(开始索引+长度)
// 这里抽象出N种事件类型，应该够用了，这样每个路由节点占用3*N字节空间，64位机器1字长是8字节
type handlersNode struct {
	validIdx    uint16
	beforeIdx   uint16
	hdsIdx      uint16
	afterIdx    uint16
	sendIdx     uint16
	responseIdx uint16

	validLen    uint8
	beforeLen   uint8
	afterLen    uint8
	hdsLen      uint8
	sendLen     uint8
	responseLen uint8
}

type handlersNodeMini struct {
	startIdx uint16 // 2字节
	hdsLen   uint8  // 1字节
}

// 事件要尽量少一些，每个路由节点都要分配一个对象
type routeEvents struct {
	eHds          []uint16
	eBeforeHds    []uint16
	eAfterHds     []uint16
	eValidHds     []uint16
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
