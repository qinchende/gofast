package jsonx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/gson"
)

const (
	tempByteStackSize = 128 // 栈上分配一定空间，方便放临时字符串（不能太大，防止协程栈伸缩）| 或者单独申请内存并管理
)

var (
	sErr = errors.New("jsonx: json syntax error.")
)

type fastDecode struct {
	dst  cst.SuperKV
	gr   *gson.GsonRow
	src  string
	head int
	tail int

	// 这里的内存分配不是在栈上，因为后面要用到，发生了逃逸。既然已经逃逸，可以考虑动态初始化
	// 即使逃逸也有一定意义，同一次解析中共享了内存
	share []byte
	//braces  bracesMark  // 大括号
	//squares squaresMark // 中括号
}

func (dd *fastDecode) init(dst cst.SuperKV, src string) error {
	//if err := dd.searchBrackets(); err != nil {
	//	return err
	//}
	//dd.dst = dst
	dd.src = src
	dd.head = 0
	dd.tail = len(dd.src) - 1
	// 直接明确具体的解析对象
	if gr, ok := dst.(*gson.GsonRow); ok {
		dd.gr = gr
	}
	return nil
}

func (dd *fastDecode) warpError(err error) error {
	if err != nil {
		end := dd.head + 10
		if end > dd.tail {
			end = dd.tail
		}
		err = errors.New(err.Error() + " near: " + dd.src[dd.head:end])
	}
	return err
}
