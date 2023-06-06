package jde

import (
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"unsafe"
)

type (
	encValFunc    func(sb *subEncode, idx int)
	encKVPairFunc func(sb *subEncode, key string)

	subEncode struct {
		//share *subEncode // 共享的subEncode，用来解析子对象
		//mp     *cst.KV        // map
		//gr     *gson.GsonRow  // GsonRow

		dm     *encMeta       // Struct | Slice,Array
		dstPtr unsafe.Pointer // 目标值dst的地址

		// 当前解析JSON的状态信息 ++++++
		//str  string // 本段字符串
		//scan int    // 自己的扫描进度，当解析错误时，这个就是定位

		bs *[]byte // 当解析数组时候用到的一系列临时队列
		//escPos []int  // 存放转义字符'\'的索引位置
		//keyIdx int    // key index
		doIdx int // list解析的数量

		//skipValue bool // 跳过当前要解析的值
	}

	encMeta struct {
		// map & gson & struct
		kvPairEnc encKVPairFunc

		// Struct
		ss        *dts.StructSchema
		fieldsEnc []encValFunc

		// array & slice
		// itemType reflect.Type
		itemBaseType reflect.Type
		itemBaseKind reflect.Kind
		listItemEnc  encValFunc
		// only array
		itemBytes int // 数组属性，item类型对应的内存字节大小
		arrLen    int // 数组属性，数组长度

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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主解析入口

func startEncode(v any) (bs []byte, err error) {
	if v == nil {
		return nil, errValueIsNil
	}
	rfType := reflect.TypeOf(v)
	if rfType.Kind() != reflect.Pointer {
		return nil, errValueMustPtr
	}

	se := subEncode{}
	//se := jdeEncPool.Get().(*subEncode)
	se.initMeta(rfType.Elem(), (*emptyInterface)(unsafe.Pointer(&v)).dataPtr)

	se.bs = newBytes()
	err = se.warpErrorCode(se.encStart())
	bs = make([]byte, len(*se.bs))
	copy(bs, *se.bs)

	jdeBytesPool.Put(se.bs)

	//se.bs = nil
	//jdeEncPool.Put(se)

	return
}

//// 包含有子subDecode时，就递归调用
//func (se *subEncode) scanSubEncode(rfType reflect.Type, ptr unsafe.Pointer) {
//	if se.share == nil {
//		se.share = jdeDecPool.Get().(*subEncode)
//	} else {
//		se.share.reset()
//	}
//	se.share.str = se.str
//	se.share.scan = se.scan
//	se.share.initMeta(rfType, ptr)
//
//	if se.share.dm.isList {
//		se.share.scanList()
//	} else {
//		se.share.scanObject()
//	}
//
//	se.scan = se.share.scan
//	se.resetShareDecode()
//}

//func (se *subEncode) readyListMixItemDec(ptr unsafe.Pointer) {
//	if se.share == nil {
//		se.share = jdeDecPool.Get().(*subEncode)
//		se.share.str = se.str
//		se.share.scan = se.scan
//		se.share.initMeta(se.dm.itemBaseType, ptr)
//		return
//	}
//
//	se.share.scan = se.scan
//	if se.share.dm.isSuperKV {
//		if se.share.dm.isGson {
//			se.share.gr = (*gson.GsonRow)(ptr)
//		} else {
//			se.share.mp = (*cst.KV)(ptr)
//		}
//	} else {
//		se.share.dstPtr = ptr // 当前值的地址
//	}
//}
//
//func (se *subEncode) resetShareDecode() {
//	if se.share.share != nil {
//		se.share.share.reset()
//		jdeDecPool.Put(se.share.share)
//		se.share.share = nil
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// rfType 是 剥离 Pointer 之后的最终类型
func (se *subEncode) initMeta(rfType reflect.Type, ptr unsafe.Pointer) {
	typAddr := (*dataType)((*emptyInterface)(unsafe.Pointer(&rfType)).dataPtr)
	if meta := cacheGetEncMeta(typAddr); meta != nil {
		se.dm = meta
	} else {
		se.buildMeta(rfType)
		cacheSetEncMeta(typAddr, se.dm)
	}

	if se.dm.isSuperKV {
		//if se.dm.isGson {
		//	se.gr = (*gson.GsonRow)(ptr)
		//} else {
		//	se.mp = (*cst.KV)(ptr)
		//}
	} else {
		se.dstPtr = ptr // 当前值的地址
	}
	return
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (se *subEncode) buildMeta(rfType reflect.Type) {
	se.dm = &encMeta{}

	switch kd := rfType.Kind(); kd {
	case reflect.Array, reflect.Slice:
		se.initListMeta(rfType)
		//se.bindListEnc()
	case reflect.Struct:
		// 模拟泛型解析，提供性能
		if rfType.String() == "gson.GsonRow" {
			se.dm.isSuperKV = true
			se.dm.isGson = true
			//se.bindGsonDec()
			return
		}
		if rfType.String() == "time.Time" {
			panic(errValueType)
		}
		se.initStructMeta(rfType)
		//se.bindStructDec()
	case reflect.Map:
		// 常规泛型
		if rfType.String() == "cst.KV" || rfType.String() == "map[string]interface {}" {
			se.dm.isSuperKV = true
			se.dm.isMap = true
			//se.bindMapDec()
			return
		}
		panic(errValueType)
	default:
		panic(errValueType)
	}
}

func (se *subEncode) initStructMeta(rfType reflect.Type) {
	se.dm.isStruct = true
	se.dm.ss = dts.SchemaForInputByType(rfType)
}

func (se *subEncode) initListMeta(rfType reflect.Type) {
	se.dm.isList = true

	//se.dm.itemType = rfType.Elem()
	se.dm.itemBaseType = rfType.Elem()
	se.dm.itemBaseKind = se.dm.itemBaseType.Kind()
	se.dm.itemBytes = int(se.dm.itemBaseType.Size())

peelPtr:
	if se.dm.itemBaseKind == reflect.Pointer {
		se.dm.itemBaseType = se.dm.itemBaseType.Elem()
		se.dm.itemBaseKind = se.dm.itemBaseType.Kind()
		se.dm.isPtr = true
		se.dm.ptrLevel++
		// TODO：指针嵌套不能超过3层，这种很少见，也就是说此解码方案并不通用
		if se.dm.ptrLevel > 3 {
			panic(errPtrLevel)
		}
		goto peelPtr
	}

	// 是否是interface类型
	if se.dm.itemBaseKind == reflect.Interface {
		se.dm.isAny = true
	}

	// 进一步初始化数组
	if rfType.Kind() == reflect.Array {
		se.dm.isArray = true
		se.dm.arrLen = rfType.Len() // 数组长度
		if se.dm.isPtr {
			return
		}
		se.dm.itemBytes = int(se.dm.itemBaseType.Size())
		se.dm.isArrBind = true
	}
}

//func (sd *subEncode) bindStructEnc() {
//	sd.dm.kvPairEnc = scanStructValue
//
//	fLen := len(sd.dm.ss.FieldsAttr)
//	sd.dm.fieldsEnc = make([]encValFunc, fLen)
//
//	i := -1
//nextField:
//	i++
//	if i >= fLen {
//		return
//	}
//
//	// 字段不是指针类型
//	if sd.dm.ss.FieldsAttr[i].PtrLevel == 0 {
//		switch sd.dm.ss.FieldsAttr[i].Kind {
//		case reflect.Int:
//			sd.dm.fieldsEnc[i] = scanObjIntValue
//		case reflect.Int8:
//			sd.dm.fieldsEnc[i] = scanObjInt8Value
//		case reflect.Int16:
//			sd.dm.fieldsEnc[i] = scanObjInt16Value
//		case reflect.Int32:
//			sd.dm.fieldsEnc[i] = scanObjInt32Value
//		case reflect.Int64:
//			sd.dm.fieldsEnc[i] = scanObjInt64Value
//		case reflect.Uint:
//			sd.dm.fieldsEnc[i] = scanObjUintValue
//		case reflect.Uint8:
//			sd.dm.fieldsEnc[i] = scanObjUint8Value
//		case reflect.Uint16:
//			sd.dm.fieldsEnc[i] = scanObjUint16Value
//		case reflect.Uint32:
//			sd.dm.fieldsEnc[i] = scanObjUint32Value
//		case reflect.Uint64:
//			sd.dm.fieldsEnc[i] = scanObjUint64Value
//		case reflect.Float32:
//			sd.dm.fieldsEnc[i] = scanObjFloat32Value
//		case reflect.Float64:
//			sd.dm.fieldsEnc[i] = scanObjFloat64Value
//		case reflect.String:
//			sd.dm.fieldsEnc[i] = scanObjStrValue
//		case reflect.Bool:
//			sd.dm.fieldsEnc[i] = scanObjBoolValue
//		case reflect.Interface:
//			sd.dm.fieldsEnc[i] = scanObjAnyValue
//		case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
//			sd.dm.fieldsEnc[i] = scanObjMixValue
//		default:
//			panic(errValueType)
//		}
//		goto nextField
//	}
//
//	// 字段是指针类型，我们需要判断的是真实的数据类型
//	switch sd.dm.ss.FieldsAttr[i].Kind {
//	case reflect.Int:
//		sd.dm.fieldsEnc[i] = scanObjPtrIntValue
//	case reflect.Int8:
//		sd.dm.fieldsEnc[i] = scanObjPtrInt8Value
//	case reflect.Int16:
//		sd.dm.fieldsEnc[i] = scanObjPtrInt16Value
//	case reflect.Int32:
//		sd.dm.fieldsEnc[i] = scanObjPtrInt32Value
//	case reflect.Int64:
//		sd.dm.fieldsEnc[i] = scanObjPtrInt64Value
//	case reflect.Uint:
//		sd.dm.fieldsEnc[i] = scanObjPtrUintValue
//	case reflect.Uint8:
//		sd.dm.fieldsEnc[i] = scanObjPtrUint8Value
//	case reflect.Uint16:
//		sd.dm.fieldsEnc[i] = scanObjPtrUint16Value
//	case reflect.Uint32:
//		sd.dm.fieldsEnc[i] = scanObjPtrUint32Value
//	case reflect.Uint64:
//		sd.dm.fieldsEnc[i] = scanObjPtrUint64Value
//	case reflect.Float32:
//		sd.dm.fieldsEnc[i] = scanObjPtrFloat32Value
//	case reflect.Float64:
//		sd.dm.fieldsEnc[i] = scanObjPtrFloat64Value
//	case reflect.String:
//		sd.dm.fieldsEnc[i] = scanObjPtrStrValue
//	case reflect.Bool:
//		sd.dm.fieldsEnc[i] = scanObjPtrBoolValue
//	case reflect.Interface:
//		sd.dm.fieldsEnc[i] = scanObjPtrAnyValue
//	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
//		sd.dm.fieldsEnc[i] = scanObjPtrMixValue
//	default:
//		panic(errValueType)
//	}
//	goto nextField
//}

//func (se *subEncode) bindListEnc() {
//	// 如果是数组，而且数组项类型不是指针类型
//	if se.dm.isArrBind {
//		switch se.dm.itemBaseKind {
//		case reflect.Int:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Int8:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Int16:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Int32:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Int64:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Uint:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Uint8:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Uint16:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Uint32:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Uint64:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Float32:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Float64:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.String:
//			se.dm.listItemEnc = encArrStrValue
//		case reflect.Bool:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Interface:
//			se.dm.listItemEnc = encListIntValue
//		case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
//			se.dm.listItemEnc = encListIntValue
//		default:
//			panic(errValueType)
//		}
//		return
//	}
//
//	switch se.dm.itemBaseKind {
//	case reflect.Int:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Int8:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Int16:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Int32:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Int64:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Uint:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Uint8:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Uint16:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Uint32:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Uint64:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Float32:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Float64:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.String:
//		se.dm.listItemEnc = encListStrValue
//	case reflect.Bool:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Interface:
//		se.dm.listItemEnc = encListIntValue
//	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
//		if se.dm.isArray {
//			se.dm.listItemEnc = encListIntValue
//		} else {
//			se.dm.listItemEnc = encListIntValue // 这里只能是Slice
//		}
//		se.dm.isArrBind = true // Note：这些情况，无需用到缓冲池
//	default:
//		panic(errValueType)
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) warpErrorCode(errCode errType) error {
	if errCode >= 0 {
		return nil
	}

	//sta := sd.scan
	//end := sta + 20 // 输出标记后面 n 个字符
	//if end > len(sd.str) {
	//	end = len(sd.str)
	//}

	//errMsg := fmt.Sprintf("jde: %s, pos %d, character %q near ( %s )", errDescription[-errCode], sta, sd.str[sta], sd.str[sta:end])
	////errMsg := strings.Join([]string{"jsonx: error pos: ", strconv.Itoa(sta), ", near ", string(sd.str[sta]), " of (", sd.str[sta:end], ")"}, "")
	//return errors.New(errMsg)

	return nil
}
