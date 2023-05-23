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

	reflectType struct {
		data *dataType
	}

	decodeFunc func(sb *subDecode)
	//decodePtr  func(sd *subDecode) uintptr

	subDecode struct {
		pl     *listPool // 当解析数组时候用到的一系列临时队列
		escPos []int     // 存放转义字符'\'的索引位置

		// 直接两种 SupperKV +++++++++++
		mp *cst.KV       // 解析到map
		gr *gson.GsonRow // 解析到GsonRow

		// Struct | Slice,Array ++++++++
		//dst    any     // 原始值
		dm     *destMeta
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

	destMeta struct {
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

// 默认值，用于缓存对象的重置
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

	rfTyp := reflect.TypeOf(dst)
	if rfTyp.Kind() != reflect.Pointer {
		return errValueMustPtr
	}

	ei := (*emptyInterface)(unsafe.Pointer(&dst))
	return sd.initDecodeInner(rfTyp, ei.ptr)
}

func (sd *subDecode) initDecodeInner(rfTyp reflect.Type, ptr unsafe.Pointer) (err error) {

	//ei := (*emptyInterface)(unsafe.Pointer(&dst))
	//tmpType := (*dataType)(*(*unsafe.Pointer)(unsafe.Pointer(&rfTyp)))

	var rt reflectType

	*(*unsafe.Pointer)(unsafe.Pointer(&rt)) = unsafe.Pointer(&rfTyp)

	//rt := (reflectType)((*unsafe.Pointer)(unsafe.Pointer(&rfTyp)))

	meta := cacheGetMeta(rt.data)
	if meta != nil {
		sd.dm = meta
	} else {
		if err = sd.buildMeta(rfTyp); err != nil {
			return
		}
		cacheSetMeta(rt.data, sd.dm)
	}

	// 当前值的地址等信息
	sd.dstPtr = uintptr(ptr)
	//if sd.dm.isList {
	//	sd.dst = dst
	//}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) buildMeta(rfTyp reflect.Type) (err error) {
	sd.dm = &destMeta{}

	switch kd := rfTyp.Kind(); kd {
	case reflect.Struct:
		if rfTyp.String() == "time.Time" {
			return errValueType
		}
		if err = sd.initStructMeta(rfTyp); err != nil {
			return
		}
		sd.bindStructDec()
	case reflect.Array, reflect.Slice:
		if err = sd.initListMeta(rfTyp); err != nil {
			return
		}
		// 进一步初始化数组
		if kd == reflect.Array {
			sd.initArrayMeta()
			//sd.dm.arrLen = rfVal.Len()
		}
		sd.bindListDec()
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) bindStructDec() {
	fLen := len(sd.dm.ss.FieldsAttr)
	sd.dm.fieldsDec = make([]decodeFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}

	// 字段不是指针类型
	if sd.dm.ss.FieldsAttr[i].PtrLevel == 0 {
		switch sd.dm.ss.FieldsAttr[i].Kind {
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

		case reflect.Pointer: // 这种不可能
		}
		goto nextField
	}

	// 字段是指针类型，我们需要判断的是真实的数据类型
	switch sd.dm.ss.FieldsAttr[i].Kind {
	case reflect.Int:
		sd.dm.fieldsDec[i] = scanObjPtrIntValue
	case reflect.Int8:
		sd.dm.fieldsDec[i] = scanObjPtrInt8Value
	case reflect.Int16:
		sd.dm.fieldsDec[i] = scanObjPtrInt16Value
	case reflect.Int32:
		sd.dm.fieldsDec[i] = scanObjPtrInt32Value
	case reflect.Int64:
		sd.dm.fieldsDec[i] = scanObjPtrInt64Value
	case reflect.Uint:
		sd.dm.fieldsDec[i] = scanObjPtrUintValue
	case reflect.Uint8:
		sd.dm.fieldsDec[i] = scanObjPtrUint8Value
	case reflect.Uint16:
		sd.dm.fieldsDec[i] = scanObjPtrUint16Value
	case reflect.Uint32:
		sd.dm.fieldsDec[i] = scanObjPtrUint32Value
	case reflect.Uint64:
		sd.dm.fieldsDec[i] = scanObjPtrUint64Value
	case reflect.Float32:
		sd.dm.fieldsDec[i] = scanObjPtrFloat32Value
	case reflect.Float64:
		sd.dm.fieldsDec[i] = scanObjPtrFloat64Value
	case reflect.String:
		sd.dm.fieldsDec[i] = scanObjPtrStrValue
	case reflect.Bool:
		sd.dm.fieldsDec[i] = scanObjPtrBoolValue
	case reflect.Interface:
		sd.dm.fieldsDec[i] = scanObjPtrAnyValue
	case reflect.Map:

	case reflect.Slice:

	case reflect.Array:

	case reflect.Struct:

	case reflect.Pointer: // 这种不可能
	}
	goto nextField
}

func (sd *subDecode) bindListDec() {
	// 如果是数组，而且数组项类型不是指针类型
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

		case reflect.Pointer: // 这种不可能
		}
		return
	}

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

	case reflect.Pointer: // 这种不可能
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
