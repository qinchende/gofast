package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/iox"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/gson"
	"io"
	"reflect"
)

func decodeFromReader(dst any, reader io.Reader, ctSize int64) error {
	// 一次性读取完成，或者遇到EOF标记或者其它错误
	if ctSize > maxJsonLength {
		ctSize = maxJsonLength
	}
	bytes, err1 := iox.ReadAll(reader, ctSize)
	if err1 != nil {
		return err1
	}
	return decodeFromString(dst, lang.BTS(bytes))
}

var shareDecode = fastDecode{}

func decodeFromString(dst any, source string) (err error) {
	if len(source) > maxJsonLength {
		return errJsonTooLarge
	}

	//fd := fastDecode{}
	fd := &shareDecode
	if err = fd.init(dst, source); err != nil {
		return
	}
	errCode := fd.scanJson()
	//fd.subDecode.putPool()
	//return fd.warpErrorCode(errCode)
	return fd.warpErrorCode(errCode)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type fastDecode struct {
	//source    string // 原始字符串
	//dst       any    // 指向原始目标值
	subDecode // 当前解析片段，用于递归
}

type subDecode struct {
	pl  *fastPool
	mp  *cst.KV       // 解析到map
	gr  *gson.GsonRow // 解析到GsonRow
	arr *listPost     // array pet (Slice|Array)
	obj *structPost   // struct pet

	str    string // 本段字符串
	scan   int    // 自己的扫描进度，当解析错误时，这个就是定位
	key    string // 当前KV对的Key值
	keyIdx int    // key index

	skipValue bool // 跳过当前要解析的值
	skipTotal bool // 跳过所有项目
	isList    bool // 区分 [] 或者 {}
	isArray   bool // 不是slice
	isStruct  bool // {} 可能目标是 一个 struct 对象
	isSuperKV bool // {} 可能目标是 cst.SuperKV 类型
}

func (fd *fastDecode) init(dst any, src string) error {
	//if dst == nil {
	//	return errValueIsNil
	//}
	//if len(src) == 0 {
	//	return errJsonEmpty
	//}

	// origin
	//fd.dst = dst
	//fd.source = src

	// subDecode
	fd.str = src
	fd.scan = 0

	//// 先确定是否是 cst.SuperKV 类型
	//var ok bool
	//if fd.gr, ok = dst.(*gson.GsonRow); !ok {
	//	if fd.mp, ok = dst.(*cst.KV); !ok {
	//		if mpt, ok := dst.(*map[string]any); ok {
	//			*fd.mp = *mpt
	//		}
	//	}
	//}
	//if fd.gr != nil || fd.mp != nil {
	//	fd.isSuperKV = true
	//	return nil
	//}
	// 目标对象不是 KV 型，那么后面只能是 List or Struct
	//return fd.subDecode.initListStruct(dst)

	if fd.subDecode.obj == nil {
		_ = fd.subDecode.initListStruct(dst)
	} else {
		fd.subDecode.obj.objPtr = reflect.ValueOf(dst).Elem().Addr().Pointer()
	}
	return nil
}

func (sd *subDecode) getPool() {
	if sd.pl == nil {
		sd.pl = jdePool.Get().(*fastPool)
	}
}

func (sd *subDecode) putPool() {
	jdePool.Put(sd.pl)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) warpErrorCode(errCode int) error {
	if errCode >= 0 {
		return nil
	}

	sta := sd.scan
	end := sta + 20 // 输出标记后面 n 个字符
	if end > len(sd.str) {
		end = len(sd.str)
	}

	errMsg := fmt.Sprintf("jsonx: error pos: %d, near %q of ( %s )", sta, sd.str[sta], sd.str[sta:end])
	//errMsg := strings.Join([]string{"jsonx: error pos: ", strconv.Itoa(sta), ", near ", string(sd.str[sta]), " of (", sd.str[sta:end], ")"}, "")
	return errors.New(errMsg)
}
