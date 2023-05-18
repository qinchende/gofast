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

type structDecode func(sb *subDecode)

type dataType struct{}
type emptyInterface struct {
	typ *dataType
	ptr unsafe.Pointer
}

//type fastDecode struct {
//	subDecode // 当前解析片段，用于递归
//}

type subDecode struct {
	pl *fastPool

	// 直接两种 SupperKV +++++++++++
	mp *cst.KV       // 解析到map
	gr *gson.GsonRow // 解析到GsonRow

	// Struct | Slice,Array ++++++++
	dm     *destMeta
	dst    any     // 原始值
	dstPtr uintptr // 数组首值地址
	arrIdx int     // 数组索引

	// 当前解析JSON的状态信息 ++++++
	str    string // 本段字符串
	scan   int    // 自己的扫描进度，当解析错误时，这个就是定位
	key    string // 当前KV对的Key值
	keyIdx int    // key index

	skipValue bool // 跳过当前要解析的值
	skipTotal bool // 跳过所有项目
	isSuperKV bool // {} 可能目标是 cst.SuperKV 类型
	destStatus
}

func (sd *subDecode) reset() {
	sd.pl = nil
	sd.mp = nil
	sd.gr = nil
	sd.dm = nil
	sd.arrIdx = 0
	sd.skipTotal = false
	sd.skipValue = false
}

type destMeta struct {
	ss     *dts.StructSchema // 目标值是一个Struct时候
	ssFunc []structDecode

	listType reflect.Type
	listKind reflect.Kind
	itemType reflect.Type
	itemKind reflect.Kind
	itemFunc structDecode
	itemSize int // 数组属性，item类型对应的内存字节大小
	arrLen   int // 数组属性，数组长度

	nextPtr func(sb *subDecode) uintptr

	destStatus
	ptrLevel uint8
}

type destStatus struct {
	isList    bool // 区分 [] 或者 {}
	isArray   bool
	isStruct  bool // {} 可能目标是 一个 struct 对象
	isAny     bool
	isPtr     bool
	isArrBind bool //isArray  bool // 不是slice
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//go:inline
func startDecode(dst any, source string) (err error) {
	if dst == nil {
		return errValueIsNil
	}

	fd := jdeDecPool.Get().(*subDecode)
	fd.reset()
	fd.str = source
	fd.scan = 0
	if err = fd.initDecode(dst); err != nil {
		return
	}

	//fd.getPool()
	errCode := fd.scanJson()
	//fd.putPool()
	err = fd.warpErrorCode(errCode)
	jdeDecPool.Put(fd)
	return
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
		return
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
	sd.destStatus = sd.dm.destStatus

	// 当前值的地址等信息
	sd.dstPtr = uintptr(ei.ptr)
	if sd.isList {
		sd.dst = dst
	}
	return
}

func (sd *subDecode) getPool() {
	if sd.isList && sd.pl == nil {
		sd.pl = jdeBufPool.Get().(*fastPool)
	}
}

func (sd *subDecode) putPool() {
	if sd.isList {
		jdeBufPool.Put(sd.pl)
		//sd.pl = nil
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) warpErrorCode(errCode errType) error {
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
