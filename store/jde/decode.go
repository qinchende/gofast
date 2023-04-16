package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/dts"
	"github.com/qinchende/gofast/store/gson"
	"math"
	"reflect"
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
	errJsonTooLarge = errors.New("jde: string too large")
	errValueType    = errors.New("jde: target value type error")
	errValueIsNil   = errors.New("jde: target value is nil")
	errJsonEmpty    = errors.New("jde: json content empty")
)

type fastDecode struct {
	dst       any    // 指向原始目标值
	source    string // 原始字符串
	subDecode        // 当前解析片段，用于递归

	// 这里的内存分配不是在栈上，因为后面要用到，发生了逃逸。既然已经逃逸，可以考虑动态初始化
	// 即使逃逸也有一定意义，同一次解析中共享了内存
	//share []byte
}

var shareAP arrPet
var shareSP structPet

type arrPet struct {
	arrType reflect.Type
	recType reflect.Type
	isPtr   bool
	val     reflect.Value // 反射值
}

type structPet struct {
	sm  *dts.StructSchema // 目标值是一个Struct时候
	val reflect.Value     // 反射值
}

type subDecode struct {
	//kind reflect.Kind  // 这里只能是：Struct|Slice|Array (Kind uint)
	mp cst.KV        // 解析到map
	gr *gson.GsonRow // 解析到GsonRow
	ap *arrPet       // array pet
	sp *structPet    // struct pet

	str       string // 本段字符串
	scan      int    // 自己的扫描进度，当解析错误时，这个就是定位
	key       string // 当前KV对的Key值
	keyIdx    int    // key index
	skipValue bool   // 跳过当前要解析的值
	skipTotal bool   // 跳过所有项目

	isList    bool // 区分 [] 或者 {}
	isStruct  bool // {} 可能目标是 一个 struct 对象
	isSuperKV bool // {} 可能目标是 cst.SuperKV 类型
}

func (dd *fastDecode) init(dst any, src string) error {
	if dst == nil {
		return errValueIsNil
	}
	if len(src) == 0 {
		return errJsonEmpty
	}

	// origin
	dd.dst = dst
	dd.source = src

	// subDecode
	dd.str = src
	dd.scan = 0

	// 先确定是否是 cst.SuperKV 类型
	var ok bool
	if dd.gr, ok = dst.(*gson.GsonRow); !ok {
		if dd.mp, ok = dst.(cst.KV); !ok {
			dd.mp, _ = dst.(map[string]any)
		}
	}
	if dd.gr != nil || dd.mp != nil {
		dd.isSuperKV = true
	}

	// 初始化其它类型
	return dd.subDecode.init(dst)
}

func (sd *subDecode) init(dst any) error {
	if sd.isSuperKV {
		return nil
	}

	// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
	val := reflect.ValueOf(dst)
	typ := val.Type()
	kind := typ.Kind()
	if kind != reflect.Pointer {
		return errValueType
	}
	kind = typ.Elem().Kind()

	// 只支持特定类型解析
	if kind == reflect.Struct {
		if typ.String() == "time.Time" {
			return errValueType
		}
		sd.isStruct = true
	} else if kind == reflect.Slice || kind == reflect.Array {
		sd.isList = true

		shareAP.arrType = typ
		shareAP.recType = typ.Elem()
		if shareAP.recType.Kind() == reflect.Pointer {
			shareAP.isPtr = true
		}
		sd.ap = &shareAP

	} else {
		return errValueType
	}
	return nil
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
