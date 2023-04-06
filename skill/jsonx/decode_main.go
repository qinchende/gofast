package jsonx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/gson"
)

const (
	tempByteStackSize = 128 // 栈上分配一定空间，方便放临时字符串（不能太大，防止协程栈伸缩）| 或者单独申请内存并管理
)

type fastDecode struct {
	dst  cst.SuperKV
	gr   *gson.GsonRow // Gson 作为特殊解析对象，单独处理
	src  string
	head uint32
	tail uint32

	// 这里的内存分配不是在栈上，因为后面要用到，发生了逃逸。既然已经逃逸，可以考虑动态初始化
	// 即使逃逸也有一定意义，同一次解析中共享了内存
	share []byte
	//braces  bracesMark  // 大括号
	//squares squaresMark // 中括号

	seg segment // 当前解析中的片段
}

type segment struct {
	data   string // 本段字符串
	offset uint32 // 相对原始的偏移量
	step   uint32 // 自己的扫描进度
}

func (dd *fastDecode) init(dst cst.SuperKV, src string) error {
	//if err := dd.searchBrackets(); err != nil {
	//	return err
	//}
	dd.dst = dst
	dd.src = src
	dd.head = 0
	dd.tail = uint32(len(dd.src) - 1)
	// 直接明确具体的解析对象
	if gr, ok := dst.(*gson.GsonRow); ok {
		dd.gr = gr
	}
	return nil
}

func (dd *fastDecode) setSeg(str string, offset uint32) {
	dd.seg.data = str
	dd.seg.offset = offset
	dd.seg.step = 0
}

func (dd *fastDecode) warpError(err error) error {
	if err != nil {
		end := dd.head + 20 // 输出标记后面 n 个字符
		if end > dd.tail {
			end = dd.tail
		}
		err = errors.New(err.Error() + " near: " + dd.src[dd.head:end])
	}
	return err
}
