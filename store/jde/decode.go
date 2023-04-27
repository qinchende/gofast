package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/dts"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"unsafe"
)

type dataType struct{}
type emptyInterface struct {
	typ *dataType
	ptr unsafe.Pointer
}

type fastDecode struct {
	//source    string // 原始字符串
	//dst       any    // 指向原始目标值
	subDecode // 当前解析片段，用于递归
}

type subDecode struct {
	pl *fastPool

	// 直接两种 SupperKV
	mp *cst.KV       // 解析到map
	gr *gson.GsonRow // 解析到GsonRow

	// 或者 Struct | Slice,Array
	dm     *destMeta
	dst    any     // 原始值
	dstPtr uintptr // 数组首值地址
	arrLen int     // 数组长度
	arrIdx int     // 数组索引

	str    string // 本段字符串
	scan   int    // 自己的扫描进度，当解析错误时，这个就是定位
	key    string // 当前KV对的Key值
	keyIdx int    // key index

	skipValue bool // 跳过当前要解析的值
	skipTotal bool // 跳过所有项目
	isSuperKV bool // {} 可能目标是 cst.SuperKV 类型
}

type destMeta struct {
	ss *dts.StructSchema // 目标值是一个Struct时候

	listType reflect.Type
	listKind reflect.Kind
	itemType reflect.Type
	itemKind reflect.Kind
	itemSize int // item类型对应的内存字节大小

	ptrLevel uint8
	isPtr    bool
	isList   bool // 区分 [] 或者 {}
	isArray  bool // 不是slice
	isStruct bool // {} 可能目标是 一个 struct 对象
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//go:inline
func startDecode(dst any, source string) (err error) {
	if dst == nil {
		return errValueIsNil
	}

	fd := &fastDecode{}
	fd.str = source
	fd.scan = 0
	if err = fd.initDecode(dst); err != nil {
		return
	}

	errCode := fd.scanJson()
	//fd.putPool()
	return fd.warpErrorCode(errCode)
}

func (sd *subDecode) initDecode(dst any) (err error) {
	// 先确定是否是 cst.SuperKV 类型
	var ok bool
	if sd.gr, ok = dst.(*gson.GsonRow); !ok {
		if sd.mp, ok = dst.(*cst.KV); !ok {
			if mpt, ok := dst.(*map[string]any); ok {
				*sd.mp = *mpt
			}
		}
	}

	if sd.gr != nil || sd.mp != nil {
		sd.isSuperKV = true
		return nil
	}

	// 目标对象不是 KV 型，那么后面只能是 List or Struct
	ei := (*emptyInterface)(unsafe.Pointer(&dst))
	meta := cacheGetMeta(ei.typ)
	if meta != nil {
		sd.dm = meta
	} else {
		if err = sd.buildMeta(dst); err != nil {
			return
		}
		cacheSetMeta(ei.typ, sd.dm)
	}

	// 当前值的地址等信息
	if sd.dm.isList {
		sd.dst = dst
		if sd.dm.isArray {
			sd.arrLen = reflect.Indirect(reflect.ValueOf(dst)).Len()
		}
	}
	sd.dstPtr = uintptr(ei.ptr)
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 绑定 array

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) buildMeta(dst any) (err error) {
	sd.dm = &destMeta{}

	rfVal := reflect.ValueOf(dst)
	if rfVal.Kind() != reflect.Pointer {
		return errValueMustPtr
	}
	rfVal = reflect.Indirect(rfVal)

	rfTyp := rfVal.Type()
	switch kd := rfTyp.Kind(); kd {
	case reflect.Struct:
		if rfTyp.String() == "time.Time" {
			return errValueType
		}
		if err = sd.initStructMeta(rfTyp); err != nil {
			return
		}
	case reflect.Array, reflect.Slice:
		if err = sd.initListMeta(rfTyp); err != nil {
			return
		}
		// 进一步初始化数组
		if kd == reflect.Array {
			sd.initArrayMeta()
		}
	default:
		return errValueType
	}
	return nil
}

func (sd *subDecode) initStructMeta(rfType reflect.Type) error {
	sd.dm.isStruct = true
	sd.dm.ss = dts.SchemaForInputByType(rfType)
	return nil
}

func (sd *subDecode) initListMeta(rfType reflect.Type) error {
	sd.dm.isList = true

	a := sd.dm
	a.listType = rfType
	a.itemType = a.listType.Elem()
	a.itemKind = a.itemType.Kind()

peelPtr:
	if a.itemKind == reflect.Pointer {
		a.itemType = a.itemType.Elem()
		a.itemKind = a.itemType.Kind()
		a.isPtr = true
		a.ptrLevel++
		// TODO：指针嵌套不能超过3层
		if a.ptrLevel > 3 {
			return errPtrLevel
		}
		goto peelPtr
	}
	return nil
}

func (sd *subDecode) initArrayMeta() {
	sd.dm.isArray = true

	sd.arrIdx = 0
	if sd.dm.isPtr {
		return
	}
	sd.dm.itemSize = int(sd.dm.itemType.Size())
}

//func (sd *subDecode) getPool() {
//	if sd.pl == nil {
//		sd.pl = jdePool.Get().(*fastPool)
//	}
//}
//
//func (sd *subDecode) putPool() {
//	jdePool.Put(sd.pl)
//}

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
