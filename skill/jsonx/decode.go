package jsonx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/gson"
	"math"
)

const (
	maxJsonLength     = math.MaxInt32 - 1 // 最大解析2GB JSON字符串
	tempByteStackSize = 128               // 栈上分配一定空间，方便放临时字符串（不能太大，防止协程栈伸缩）| 或者单独申请内存并管理
)

const (
	noErr        int = 0
	errEOF       int = -1
	errJson      int = -11
	errChar      int = -2
	errEscape    int = -3
	errUnicode   int = -4
	errOverflow  int = -5
	errNumberFmt int = -6
	errExceedMax int = -7
	errInfinity  int = -8
	errMismatch  int = -9
	errUTF8      int = -10

	//errNotFound       int = -33
	//errNotSupportType int = -34
)

var (
	//sErr            = errors.New("jsonx: json syntax error.")
	errJsonTooLarge = errors.New("jsonx: string too large.")
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

	// 解析过程中用到临时变量 +++
	sub  string // 本段字符串
	scan uint32 // 自己的扫描进度
	flag uint32 // 当前扫描位置做标记，方便错误定位

	isMixedVal bool
}

func (dd *fastDecode) init(dst cst.SuperKV, src string) error {
	dd.dst = dst
	dd.src = src
	dd.head = 0
	dd.tail = uint32(len(dd.src) - 1)
	dd.gr, _ = dst.(*gson.GsonRow)
	return nil
}

func (dd *fastDecode) changeSub(str string) {
	dd.sub = str
	dd.scan = 0
}

func (dd *fastDecode) warpError(code int) error {
	if code >= 0 {
		return nil
	}

	//end := dd.flag + 20 // 输出标记后面 n 个字符
	//if end > dd.tail {
	//	end = dd.tail
	//}
	errMsg := fmt.Sprintf("jsonx: error type %d. pos: %d, near: %q", code, dd.flag, dd.src[dd.flag])
	return errors.New(errMsg)
}
