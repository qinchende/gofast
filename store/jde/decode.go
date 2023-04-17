package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/iox"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/dts"
	"github.com/qinchende/gofast/store/gson"
	"io"
	"reflect"
)

type fastDecode struct {
	dst       any    // 指向原始目标值
	source    string // 原始字符串
	subDecode        // 当前解析片段，用于递归
}

type arrPet struct {
	dst     any
	arrType reflect.Type
	recType reflect.Type
	recKind reflect.Kind
	isPtr   bool
	val     reflect.Value // 反射值
}

type structPet struct {
	sm *dts.StructSchema // 目标值是一个Struct时候
	//val reflect.Value     // 反射值
}

type subDecode struct {
	//kind reflect.Kind  // 这里只能是：Struct|Slice|Array (Kind uint)
	mp  cst.KV        // 解析到map
	gr  *gson.GsonRow // 解析到GsonRow
	arr *arrPet       // array pet (Slice|Array)
	obj *structPet    // struct pet

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
	sap.val = reflect.Indirect(val)

	// 只支持特定类型解析
	if kind == reflect.Struct {
		if typ.String() == "time.Time" {
			return errValueType
		}
		sd.isStruct = true
	} else if kind == reflect.Slice || kind == reflect.Array {
		sd.isList = true

		sap.arrType = typ
		sap.recType = typ.Elem()
		sap.recKind = sap.recType.Elem().Kind()
		if sap.recKind == reflect.Pointer {
			sap.recKind = sap.recType.Elem().Elem().Kind()
			sap.isPtr = true
		}

		sap.dst = dst
		sd.arr = &sap

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

func decodeFromString(dst any, source string) error {
	if len(source) > maxJsonLength {
		return errJsonTooLarge
	}

	dd := fastDecode{}
	if err := dd.init(dst, source); err != nil {
		return err
	}
	return dd.warpError(dd.parseJson())
}
