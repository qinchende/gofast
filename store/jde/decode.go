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

//type fastDecode struct {
//	subDecode // 当前解析片段，用于递归
//}

type (
	dataType       struct{}
	emptyInterface struct {
		typ *dataType
		ptr unsafe.Pointer
	}

	subDecode struct {
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
		keyIdx int    // key index
		//key    string // 当前KV对的Key值

		isSuperKV bool // {} 可能目标是 cst.SuperKV 类型
		skipValue bool // 跳过当前要解析的值
		skipTotal bool // 跳过所有项目
	}

	decodeFunc func(sb *subDecode)
	decodePtr  func(sd *subDecode) uintptr
	destMeta   struct {
		// struct
		ss        *dts.StructSchema // 目标值是一个Struct时候
		fieldsDec []decodeFunc

		// array & slice
		listType reflect.Type
		listKind reflect.Kind
		itemType reflect.Type
		itemKind reflect.Kind
		itemDec  decodeFunc

		// array
		arrItemBytes int // 数组属性，item类型对应的内存字节大小
		arrLen       int // 数组属性，数组长度

		// status
		isList    bool // 区分 [] 或者 {}
		isArray   bool
		isStruct  bool // {} 可能目标是 一个 struct 对象
		isAny     bool
		isArrBind bool //isArray  bool // 不是slice
		isPtr     bool
		ptrLevel  uint8
	}
)

//type destStatus struct {
//	isList    bool // 区分 [] 或者 {}
//	isArray   bool
//	isStruct  bool // {} 可能目标是 一个 struct 对象
//	isAny     bool
//	isPtr     bool
//	isArrBind bool //isArray  bool // 不是slice
//}

var _subDecodeDefValues subDecode

func (sd *subDecode) reset() {
	//sd.pl = nil
	//sd.mp = nil
	//sd.gr = nil
	//sd.dm = nil
	//sd.arrIdx = 0
	//sd.skipTotal = false
	//sd.skipValue = false
	*sd = _subDecodeDefValues
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

	fd.getPool()
	errCode := fd.scanJson()
	fd.putPool()
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
	//sd.destStatus = sd.dm.destStatus

	// 当前值的地址等信息
	sd.dstPtr = uintptr(ei.ptr)
	if sd.dm.isList {
		sd.dst = dst
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) getPool() {
	if sd.dm.isList && sd.pl == nil {
		sd.pl = jdeBufPool.Get().(*fastPool)
	}
}

func (sd *subDecode) putPool() {
	if sd.dm.isList {
		jdeBufPool.Put(sd.pl)
		//sd.pl = nil
	}
}

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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
		// 初始化解析函数 ++++++++++++++++++
		sd.dm.fieldsDec = make([]decodeFunc, len(sd.dm.ss.FieldsKind))
		for i := 0; i < len(sd.dm.ss.FieldsKind); i++ {

			switch sd.dm.ss.FieldsKind[i] {
			case reflect.Int:
				sd.dm.fieldsDec[i] = scanObjIntValue
			case reflect.Int8:
				sd.dm.fieldsDec[i] = scanObjInt8Value
			case reflect.Int16:
				sd.dm.fieldsDec[i] = scanObjInt16Value
			case reflect.Int32:
				sd.dm.fieldsDec[i] = scanObjInt32Value
			case reflect.Int64:
				sd.dm.fieldsDec[i] = scanObjInt64Value
			case reflect.Uint:
				sd.dm.fieldsDec[i] = scanObjUintValue
			case reflect.Uint8:
				sd.dm.fieldsDec[i] = scanObjUint8Value
			case reflect.Uint16:
				sd.dm.fieldsDec[i] = scanObjUint16Value
			case reflect.Uint32:
				sd.dm.fieldsDec[i] = scanObjUint32Value
			case reflect.Uint64:
				sd.dm.fieldsDec[i] = scanObjUint64Value
			case reflect.Float32:
				sd.dm.fieldsDec[i] = scanObjFloat32Value
			case reflect.Float64:
				sd.dm.fieldsDec[i] = scanObjFloat64Value
			case reflect.String:
				sd.dm.fieldsDec[i] = scanObjStrValue
			case reflect.Bool:
				sd.dm.fieldsDec[i] = scanObjBoolValue
			case reflect.Interface:
				sd.dm.fieldsDec[i] = scanObjAnyValue
			case reflect.Map:

			case reflect.Slice:

			case reflect.Array:

			case reflect.Struct:

			case reflect.Pointer:
				sd.dm.fieldsDec[i] = sd.scanObjPtrValue(scanObjAnyValue)
			}
		}
	case reflect.Array, reflect.Slice:
		if err = sd.initListMeta(rfTyp); err != nil {
			return
		}

		// 进一步初始化数组
		if kd == reflect.Array {
			sd.initArrayMeta()
			sd.dm.arrLen = rfVal.Len()
		}

		if sd.dm.isArrBind {

			switch sd.dm.itemKind {
			case reflect.Int:
				sd.dm.itemDec = scanArrIntValue
			case reflect.Int8:
				sd.dm.itemDec = scanArrInt8Value
			case reflect.Int16:
				sd.dm.itemDec = scanArrInt16Value
			case reflect.Int32:
				sd.dm.itemDec = scanArrInt32Value
			case reflect.Int64:
				sd.dm.itemDec = scanArrInt64Value
			case reflect.Uint:
				sd.dm.itemDec = scanArrUintValue
			case reflect.Uint8:
				sd.dm.itemDec = scanArrUint8Value
			case reflect.Uint16:
				sd.dm.itemDec = scanArrUint16Value
			case reflect.Uint32:
				sd.dm.itemDec = scanArrUint32Value
			case reflect.Uint64:
				sd.dm.itemDec = scanArrUint64Value
			case reflect.Float32:
				sd.dm.itemDec = scanArrFloat32Value
			case reflect.Float64:
				sd.dm.itemDec = scanArrFloat64Value
			case reflect.String:
				sd.dm.itemDec = scanArrStrValue
			case reflect.Bool:
				sd.dm.itemDec = scanArrBoolValue
			case reflect.Interface:
				sd.dm.itemDec = scanArrAnyValue
			case reflect.Map:

			case reflect.Slice:

			case reflect.Array:

			case reflect.Struct:

			case reflect.Pointer:
				// 不可能
			}

		} else {

			switch sd.dm.itemKind {
			case reflect.Int:
				sd.dm.itemDec = scanListIntValue
			case reflect.Int8:
				sd.dm.itemDec = scanListInt8Value
			case reflect.Int16:
				sd.dm.itemDec = scanListInt16Value
			case reflect.Int32:
				sd.dm.itemDec = scanListInt32Value
			case reflect.Int64:
				sd.dm.itemDec = scanListInt64Value
			case reflect.Uint:
				sd.dm.itemDec = scanListUintValue
			case reflect.Uint8:
				sd.dm.itemDec = scanListUint8Value
			case reflect.Uint16:
				sd.dm.itemDec = scanListUint16Value
			case reflect.Uint32:
				sd.dm.itemDec = scanListUint32Value
			case reflect.Uint64:
				sd.dm.itemDec = scanListUint64Value
			case reflect.Float32:
				sd.dm.itemDec = scanListFloat32Value
			case reflect.Float64:
				sd.dm.itemDec = scanListFloat64Value
			case reflect.String:
				sd.dm.itemDec = scanListStrValue
			case reflect.Bool:
				sd.dm.itemDec = scanListBoolValue
			case reflect.Interface:
				sd.dm.itemDec = scanListAnyValue
			case reflect.Map:

			case reflect.Slice:

			case reflect.Array:

			case reflect.Struct:
			case reflect.Pointer:
				// 已经统一处理
			}

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

	sd.dm.listType = rfType
	sd.dm.listKind = rfType.Kind()
	sd.dm.itemType = sd.dm.listType.Elem()
	sd.dm.itemKind = sd.dm.itemType.Kind()

peelPtr:
	if sd.dm.itemKind == reflect.Pointer {
		sd.dm.itemType = sd.dm.itemType.Elem()
		sd.dm.itemKind = sd.dm.itemType.Kind()
		sd.dm.isPtr = true
		sd.dm.ptrLevel++
		// TODO：指针嵌套不能超过3层，这种很少见，也就是说此解码方案并不通用
		if sd.dm.ptrLevel > 3 {
			return errPtrLevel
		}
		goto peelPtr
	}

	// 是否是interface类型
	if sd.dm.itemKind == reflect.Interface {
		sd.dm.isAny = true
	}
	return nil
}

func (sd *subDecode) initArrayMeta() {
	sd.dm.isArray = true
	if sd.dm.isPtr {
		return
	}
	sd.dm.isArrBind = true
	sd.dm.arrItemBytes = int(sd.dm.itemType.Size())
}
