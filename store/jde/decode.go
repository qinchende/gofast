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
	emptyInterface struct { // emptyInterface is the header for an interface{} value. (ignore method)
		typAddr *dataType
		dataPtr unsafe.Pointer
	}

	decValFunc    func(sb *subDecode)
	decKVPairFunc func(sb *subDecode, key string)

	subDecode struct {
		mp     *cst.KV       // map
		gr     *gson.GsonRow // GsonRow
		dm     *destMeta     // Struct | Slice,Array
		dstPtr uintptr       // 数组首值地址

		// 当前解析JSON的状态信息 ++++++
		str  string // 本段字符串
		scan int    // 自己的扫描进度，当解析错误时，这个就是定位

		pl     *listPool // 当解析数组时候用到的一系列临时队列
		escPos []int     // 存放转义字符'\'的索引位置
		keyIdx int       // key index
		arrIdx int       // 数组索引

		skipValue bool // 跳过当前要解析的值
		//skipTotal bool // 跳过所有项目
	}

	destMeta struct {
		// map & gson & struct
		kvPairDec decKVPairFunc

		// Struct
		ss        *dts.StructSchema
		fieldsDec []decValFunc

		// array & slice
		//listType    reflect.Type
		//listKind    reflect.Kind
		itemType    reflect.Type
		itemKind    reflect.Kind
		listItemDec decValFunc
		// only array
		arrItemBytes int // 数组属性，item类型对应的内存字节大小
		arrLen       int // 数组属性，数组长度

		// status
		isSuperKV bool // {} SuperKV
		isGson    bool // {} gson
		isMap     bool // {} map
		isStruct  bool // {} struct

		isList    bool  // [] array & slice
		isAny     bool  // [] is list and item is interface type in the final
		isPtr     bool  // [] is list and item is pointer type
		ptrLevel  uint8 // [] is list and item pointer level
		isArray   bool  // [] array
		isArrBind bool  // [] is array and item not pointer type
	}
)

// 默认值，用于缓存对象的重置
var _subDecodeDefValues subDecode

func (sd *subDecode) reset() {
	tp := sd.escPos
	*sd = _subDecodeDefValues
	sd.escPos = tp[0:0]
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主解析入口
func startDecode(dst any, source string) (err error) {
	if dst == nil {
		return errValueIsNil
	}
	rfTyp := reflect.TypeOf(dst)
	if rfTyp.Kind() != reflect.Pointer {
		return errValueMustPtr
	}

	sd := jdeDecPool.Get().(*subDecode)
	sd.str = source
	sd.scan = 0
	sd.initMeta(rfTyp.Elem(), (*emptyInterface)(unsafe.Pointer(&dst)).dataPtr)

	err = sd.warpErrorCode(sd.scanStart())

	// TODO：此时 sd 中指针指向的对象没有被释放，存在一定风险，所以要先释放再回收
	sd.reset()
	jdeDecPool.Put(sd)
	return
}

// 包含有子subDecode时，就递归调用
func (sd *subDecode) scanSubDecode(rfTyp reflect.Type, ptr unsafe.Pointer) {
	nsd := jdeDecPool.Get().(*subDecode)

	nsd.str = sd.str
	nsd.scan = sd.scan
	nsd.initMeta(rfTyp, ptr)
	if nsd.dm.isList {
		nsd.scanList()
	} else {
		nsd.scanObject()
	}
	sd.scan = nsd.scan

	nsd.reset()
	jdeDecPool.Put(nsd)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// rfTyp 是 剥离 Pointer 之后的最终类型
func (sd *subDecode) initMeta(rfTyp reflect.Type, ptr unsafe.Pointer) {
	typAddr := (*dataType)((*emptyInterface)(unsafe.Pointer(&rfTyp)).dataPtr)
	if meta := cacheGetMeta(typAddr); meta != nil {
		sd.dm = meta
	} else {
		sd.buildMeta(rfTyp)
		cacheSetMeta(typAddr, sd.dm)
	}

	if sd.dm.isSuperKV {
		if sd.dm.isGson {
			sd.gr = (*gson.GsonRow)(ptr)
		} else {
			sd.mp = (*cst.KV)(ptr)
		}
	} else {
		sd.dstPtr = uintptr(ptr) // 当前值的地址
	}
	return
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) buildMeta(rfTyp reflect.Type) {
	sd.dm = &destMeta{}

	switch kd := rfTyp.Kind(); kd {
	case reflect.Array, reflect.Slice:
		sd.initListMeta(rfTyp)
		sd.bindListDec()
	case reflect.Struct:
		// 模拟泛型解析，提供性能
		if rfTyp.String() == "gson.GsonRow" {
			sd.dm.isSuperKV = true
			sd.dm.isGson = true
			sd.bindGsonDec()
			return
		}
		if rfTyp.String() == "time.Time" {
			panic(errValueType)
		}
		sd.initStructMeta(rfTyp)
		sd.bindStructDec()
	case reflect.Map:
		// 常规泛型
		if rfTyp.String() == "cst.KV" || rfTyp.String() == "map[string]interface {}" {
			sd.dm.isSuperKV = true
			sd.dm.isMap = true
			sd.bindMapDec()
			return
		}
		panic(errValueType)
	default:
		panic(errValueType)
	}
}

func (sd *subDecode) initStructMeta(rfType reflect.Type) {
	sd.dm.isStruct = true
	sd.dm.ss = dts.SchemaForInputByType(rfType)
}

func (sd *subDecode) initListMeta(rfType reflect.Type) {
	sd.dm.isList = true

	//sd.dm.listType = rfType
	//sd.dm.listKind = rfType.Kind()
	sd.dm.itemType = rfType.Elem()
	sd.dm.itemKind = sd.dm.itemType.Kind()

peelPtr:
	if sd.dm.itemKind == reflect.Pointer {
		sd.dm.itemType = sd.dm.itemType.Elem()
		sd.dm.itemKind = sd.dm.itemType.Kind()
		sd.dm.isPtr = true
		sd.dm.ptrLevel++
		// TODO：指针嵌套不能超过3层，这种很少见，也就是说此解码方案并不通用
		if sd.dm.ptrLevel > 3 {
			panic(errPtrLevel)
		}
		goto peelPtr
	}

	// 是否是interface类型
	if sd.dm.itemKind == reflect.Interface {
		sd.dm.isAny = true
	}

	// 进一步初始化数组
	if rfType.Kind() == reflect.Array {
		sd.dm.isArray = true
		sd.dm.arrLen = rfType.Len() // 数组长度
		if sd.dm.isPtr {
			return
		}
		sd.dm.isArrBind = true
		sd.dm.arrItemBytes = int(sd.dm.itemType.Size())
	}
}

// JSON数据主要分 object {} + list [] 两种类型
// map & gson & struct 都需要解析 {} 他们都是 kvPair 形式的数据
// array & slice 需要解析 [] ，他们都是 List 形式的数据
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) bindMapDec() {
	sd.dm.kvPairDec = scanMapAnyValue
}

func (sd *subDecode) bindGsonDec() {
	sd.dm.kvPairDec = scanGsonValue
}

func (sd *subDecode) bindStructDec() {
	sd.dm.kvPairDec = scanStructValue

	fLen := len(sd.dm.ss.FieldsAttr)
	sd.dm.fieldsDec = make([]decValFunc, fLen)

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
		case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
			sd.dm.fieldsDec[i] = scanObjMixValue
		default:
			panic(errValueType)
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
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		sd.dm.fieldsDec[i] = scanObjPtrMixValue
	default:
		panic(errValueType)
	}
	goto nextField
}

func (sd *subDecode) bindListDec() {
	// 如果是数组，而且数组项类型不是指针类型
	if sd.dm.isArrBind {
		switch sd.dm.itemKind {
		case reflect.Int:
			sd.dm.listItemDec = scanArrIntValue
		case reflect.Int8:
			sd.dm.listItemDec = scanArrInt8Value
		case reflect.Int16:
			sd.dm.listItemDec = scanArrInt16Value
		case reflect.Int32:
			sd.dm.listItemDec = scanArrInt32Value
		case reflect.Int64:
			sd.dm.listItemDec = scanArrInt64Value
		case reflect.Uint:
			sd.dm.listItemDec = scanArrUintValue
		case reflect.Uint8:
			sd.dm.listItemDec = scanArrUint8Value
		case reflect.Uint16:
			sd.dm.listItemDec = scanArrUint16Value
		case reflect.Uint32:
			sd.dm.listItemDec = scanArrUint32Value
		case reflect.Uint64:
			sd.dm.listItemDec = scanArrUint64Value
		case reflect.Float32:
			sd.dm.listItemDec = scanArrFloat32Value
		case reflect.Float64:
			sd.dm.listItemDec = scanArrFloat64Value
		case reflect.String:
			sd.dm.listItemDec = scanArrStrValue
		case reflect.Bool:
			sd.dm.listItemDec = scanArrBoolValue
		case reflect.Interface:
			sd.dm.listItemDec = scanArrAnyValue
		case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
			sd.dm.listItemDec = scanListMixValue
		default:
			panic(errValueType)
		}
		return
	}

	switch sd.dm.itemKind {
	case reflect.Int:
		sd.dm.listItemDec = scanListIntValue
	case reflect.Int8:
		sd.dm.listItemDec = scanListInt8Value
	case reflect.Int16:
		sd.dm.listItemDec = scanListInt16Value
	case reflect.Int32:
		sd.dm.listItemDec = scanListInt32Value
	case reflect.Int64:
		sd.dm.listItemDec = scanListInt64Value
	case reflect.Uint:
		sd.dm.listItemDec = scanListUintValue
	case reflect.Uint8:
		sd.dm.listItemDec = scanListUint8Value
	case reflect.Uint16:
		sd.dm.listItemDec = scanListUint16Value
	case reflect.Uint32:
		sd.dm.listItemDec = scanListUint32Value
	case reflect.Uint64:
		sd.dm.listItemDec = scanListUint64Value
	case reflect.Float32:
		sd.dm.listItemDec = scanListFloat32Value
	case reflect.Float64:
		sd.dm.listItemDec = scanListFloat64Value
	case reflect.String:
		sd.dm.listItemDec = scanListStrValue
	case reflect.Bool:
		sd.dm.listItemDec = scanListBoolValue
	case reflect.Interface:
		sd.dm.listItemDec = scanListAnyValue
	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		sd.dm.listItemDec = scanListMixValue
	default:
		panic(errValueType)
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
