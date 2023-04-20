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

type fastDecode struct {
	dst       any    // 指向原始目标值
	source    string // 原始字符串
	subDecode        // 当前解析片段，用于递归
}

type subDecode struct {
	pl  *fastPool
	mp  *cst.KV       // 解析到map
	gr  *gson.GsonRow // 解析到GsonRow
	arr *listMeta     // array pet (Slice|Array)
	obj *structMeta   // struct pet

	str       string // 本段字符串
	scan      int    // 自己的扫描进度，当解析错误时，这个就是定位
	key       string // 当前KV对的Key值
	keyIdx    int    // key index
	skipValue bool   // 跳过当前要解析的值
	skipTotal bool   // 跳过所有项目
	isList    bool   // 区分 [] 或者 {}
	isArray   bool   // 不是slice
	isStruct  bool   // {} 可能目标是 一个 struct 对象
	isSuperKV bool   // {} 可能目标是 cst.SuperKV 类型
}

func (dd *fastDecode) init(dst any, src string) error {
	if dst == nil {
		return errValueIsNil
	}
	if len(src) == 0 {
		return errJsonEmpty
	}
	rfVal := reflect.ValueOf(dst)
	if rfVal.Kind() != reflect.Pointer {
		return errValueMustPtr
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
		if dd.mp, ok = dst.(*cst.KV); !ok {
			if mpt, ok := dst.(*map[string]any); ok {
				*dd.mp = *mpt
			}
		}
	}
	if dd.gr != nil || dd.mp != nil {
		dd.isSuperKV = true
		return nil
	} else {
		dd.subDecode.pl = jdePool.Get().(*fastPool)
		dd.subDecode.pl.initMem()
		dd.subDecode.pl.arr.dst = dst

		return dd.subDecode.initListStruct(rfVal.Elem())
	}
}

func (dd *fastDecode) finish() {
	jdePool.Put(dd.subDecode.pl)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

func decodeFromString(dst any, source string) (err error) {
	if len(source) > maxJsonLength {
		return errJsonTooLarge
	}

	dd := fastDecode{}
	if err = dd.init(dst, source); err != nil {
		return
	}
	errInt := dd.parseJson()
	dd.finish()
	return dd.warpError(errInt)
}
