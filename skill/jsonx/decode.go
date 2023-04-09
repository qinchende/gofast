package jsonx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"math"
)

const (
	maxJsonLength     = math.MaxInt32 - 1 // 最大解析2GB JSON字符串
	tempByteStackSize = 128               // 栈上分配一定空间，方便放临时字符串（不能太大，防止协程栈伸缩）| 或者单独申请内存并管理
)

const (
	noErr        int = 0
	errNotFound  int = -1
	errEOF       int = -2
	errJson      int = -3
	errChar      int = -4
	errEscape    int = -5
	errUnicode   int = -6
	errOverflow  int = -7
	errNumberFmt int = -8
	errExceedMax int = -9
	errInfinity  int = -10
	errMismatch  int = -11
	errUTF8      int = -12

	//errNotSupportType int = -13
)

//var errorStrings = []string{
//	0:                      "ok",
//	-(errEOF):              "eof",
//	ERR_INVALID_CHAR:       "invalid char",
//	ERR_INVALID_ESCAPE:     "invalid escape char",
//	ERR_INVALID_UNICODE:    "invalid unicode escape",
//	ERR_INTEGER_OVERFLOW:   "integer overflow",
//	ERR_INVALID_NUMBER_FMT: "invalid number format",
//	ERR_RECURSE_EXCEED_MAX: "recursion exceeded max depth",
//	ERR_FLOAT_INFINITY:     "float number is infinity",
//	ERR_MISMATCH:           "mismatched type with value",
//	ERR_INVALID_UTF8:       "invalid UTF8",
//}

var (
	//sErr            = errors.New("jsonx: json syntax error.")
	errJsonTooLarge = errors.New("jsonx: string too large.")
)

type fastDecode struct {
	//dst cst.SuperKV
	//gr  *gson.GsonRow // Gson 作为特殊解析对象，单独处理
	src string
	//head int // 头位置
	//tail int // 尾位置
	// 这里的内存分配不是在栈上，因为后面要用到，发生了逃逸。既然已经逃逸，可以考虑动态初始化
	// 即使逃逸也有一定意义，同一次解析中共享了内存
	//share []byte

	// 递归解析节点
	root subDecode
}

type subDecode struct {
	offSrc int // 相对于原始字符串的偏移
	share  []byte
	dst    cst.SuperKV
	//gr  *gson.GsonRow // Gson 作为特殊解析对象，单独处理

	// 解析对象过程中用到临时变量 +++++++++++++++++++++++++++++++++++++
	sub        string // 本段字符串
	scan       int    // 自己的扫描进度
	errPos     int    // 当前扫描位置做标记，方便错误定位
	isMixedVal bool   // 判断当前 {} 或者 [] 是一个字符串整体，不需要解析
}

func (dd *fastDecode) init(dst cst.SuperKV, src string) {
	dd.root.dst = dst
	dd.src = src
	// dd.head = 0
	// dd.tail = uint32(len(dd.src) - 1)
	//dd.gr, _ = dst.(*gson.GsonRow)

	dd.root.initSubDecode(dd.src, 0)
}

func (sd *subDecode) initSubDecode(subStr string, offSrc int) {
	//sd.flag = 0
	sd.isMixedVal = false

	th := trimHead(subStr)
	tt := trimTail(subStr)
	if th > tt {
		th = tt
	}
	sd.scan = th
	sd.offSrc = offSrc + th
	sd.sub = subStr[th : tt+1]
}

func (sd *subDecode) warpError(errCode int) error {
	if errCode >= 0 {
		return nil
	}

	//end := sd.flag + 20 // 输出标记后面 n 个字符
	//if end > sd.tail {
	//	end = sd.tail
	//}
	pos := sd.errPos + sd.offSrc
	errMsg := fmt.Sprintf("jsonx: error type %d. pos: %d, near: %q", errCode, pos, sd.sub[pos])
	return errors.New(errMsg)
}
