package jde

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
	bytesNull  = "null"
	bytesTrue  = "true"
	bytesFalse = "false"
)

const (
	noErr        int = 0  // 没有错误
	scanEOF      int = -1 // 扫描结束
	errNormal    int = -2 // 没找到期望的字符
	errJson      int = -3 // 非法JSON格式
	errChar      int = -4 // 非预期的字符
	errEscape    int = -5
	errUnicode   int = -6
	errOverflow  int = -7
	errNumberFmt int = -8
	errExceedMax int = -9
	errInfinity  int = -10
	errMismatch  int = -11
	errUTF8      int = -12
	errKey       int = -13
	errValue     int = -14
	errKV        int = -15
	errNull      int = -16
	errObject    int = -17
	errArray     int = -18
	errTrue      int = -19
	errFalse     int = -20

	//errNotSupportType int = -13
)

//var errorStrings = []string{
//	0:                      "ok",
//	-(scanEOF):              "eof",
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
	src  string    // 原始字符串
	root subDecode // 当前解析片段，用于递归

	//head int // 头位置
	//tail int // 尾位置
	// 这里的内存分配不是在栈上，因为后面要用到，发生了逃逸。既然已经逃逸，可以考虑动态初始化
	// 即使逃逸也有一定意义，同一次解析中共享了内存
	//share []byte
}

type subDecode struct {
	dst  cst.SuperKV
	gr   *gson.GsonRow
	list []any

	//directString bool   // 判断当前 {} 或者 [] 是一个字符串整体，不需要解析
	//offSrc       int    // 相对于原始字符串的偏移
	// 解析对象过程中用到临时变量 +++++++++++++++++++++++++++++++++++++
	str       string // 本段字符串
	scan      int    // 自己的扫描进度，当解析错误时，这个就是定位
	skipValue bool   // 跳过当前要解析的值
	skipTotal bool   // 跳过所有项目
}

func (dd *fastDecode) init(dst cst.SuperKV, src string) error {
	if dst == nil {
		return errors.New("target value can't nil.")
	}
	if len(src) == 0 {
		return errors.New("json content empty.")
	}

	dd.root.dst = dst
	dd.src = src
	dd.root.initSubDecode(dd.src)
	return nil
}

func (sd *subDecode) initSubDecode(subStr string) {
	sd.str = subStr
	sd.scan = 0
	sd.gr, _ = sd.dst.(*gson.GsonRow)
}

func (sd *subDecode) warpError(errCode int) error {
	if errCode >= 0 {
		return nil
	}

	sta := sd.scan
	end := sta + 20 // 输出标记后面 n 个字符
	if end > len(sd.str) {
		end = len(sd.str)
	}

	errMsg := fmt.Sprintf("jsonx: error pos: %d, near %q of ( %s )", sta, sd.str[sta], sd.str[sta:end])
	return errors.New(errMsg)
}
