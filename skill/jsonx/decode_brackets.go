package jsonx

//type fastDecode struct {
//	dst  cst.SuperKV
//	gr   *gson.GsonRow // Gson 作为特殊解析对象，单独处理
//	src  string
//	head uint32
//	tail uint32
//
//	// 这里的内存分配不是在栈上，因为后面要用到，发生了逃逸。既然已经逃逸，可以考虑动态初始化
//	// 即使逃逸也有一定意义，同一次解析中共享了内存
//	share []byte
//	//braces  bracesMark  // 大括号
//	//squares squaresMark // 中括号
//
//	// 解析过程中用到临时变量
//	seg segment // 当前解析中的片段
//}

//// 大括号
//type bracesMark struct {
//	left  []uint32
//	right []uint32
//	lIdx  uint32
//	rIdx  uint32
//}
//
//// 中括号，方括号
//type squaresMark struct {
//	left  []uint32
//	right []uint32
//	lIdx  uint32
//	rIdx  uint32
//}

//// 先扫描字符串，获取所有左右括号的位置信息
//func (dd *fastDecode) searchBrackets() error {
//	dd.braces.left = make([]uint32, 0, 8)
//	dd.braces.right = make([]uint32, 0, 8)
//
//	dd.squares.left = make([]uint32, 0, 8)
//	dd.squares.right = make([]uint32, 0, 8)
//
//	strLen := uint32(len(dd.src))
//	for i := uint32(0); i < strLen; i++ {
//		switch dd.src[i] {
//		case '{':
//			dd.braces.left = append(dd.braces.left, i)
//		case '}':
//			dd.braces.right = append(dd.braces.right, i)
//		case '[':
//			dd.squares.left = append(dd.squares.left, i)
//		case ']':
//			dd.squares.right = append(dd.squares.right, i)
//		}
//	}
//
//	// 左右大括号数量不一致，格式错误(object)
//	if len(dd.braces.left) != len(dd.braces.right) {
//		return sErr
//	}
//	// 左右中括号数量不一致，格式错误(array)
//	if len(dd.squares.left) != len(dd.squares.right) {
//		return sErr
//	}
//	return nil
//}
