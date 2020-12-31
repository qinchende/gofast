package fst

// 自定义内存数据库，存放路由树所有相关的数据
type fstMemSpace struct {
	// 我们需要自己定义一个切片，管理所有的 Context handlers.
	allCHandlers CHandlers
	allCHdsLen   uint16
	// 新的handlers,有序的，按分组和事件类型排序
	hdsList    CHandlers
	hdsListLen uint16
	// 路由节点对应的处理方法索引结构
	hdsGroupCt      uint16 // 所有最后一级分组的个数
	hdsItemCt       uint16 // 所有路由节点的个数
	hdsNodes        []handlersNode
	hdsNodesLen     uint16
	hdsNodesMini    []handlersNodeMini
	hdsNodesMiniLen uint16

	treeCharT    []byte
	treeChars    string
	treeCharsLen uint16

	allMiniNodes []miniNode
	allMiniLen   uint16
}

var fstMem fstMemSpace

func init() {
	fstMem = fstMemSpace{}
}
