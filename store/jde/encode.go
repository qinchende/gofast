package jde

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/dts"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"unsafe"
)

type (
	encValFunc    func(bf []byte, ptr unsafe.Pointer) []byte
	encMixValFunc func(bf []byte, ptr unsafe.Pointer, rfType reflect.Type) []byte

	subEncode struct {
		srcPtr unsafe.Pointer // 对象值地址
		mp     *cst.KV        // map
		gr     *gson.GsonRow  // GsonRow
		em     *encMeta       // Struct | Slice,Array
		bs     *[]byte        // 当解析数组时候用到的一系列临时队列
	}

	encMeta struct {
		// Struct
		ss           *dts.StructSchema
		fieldsEnc    []encValFunc
		fieldsEncMix []encMixValFunc

		// array & slice
		itemBaseType   reflect.Type
		itemBaseKind   reflect.Kind
		listItemEnc    encValFunc
		listItemEncMix encMixValFunc
		itemBytes      int // 数组属性，item类型对应的内存字节大小
		arrLen         int // 数组属性，数组长度

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
	se.encStart()
	bs = make([]byte, len(*se.bs))
	copy(bs, *se.bs)

	jdeBytesPool.Put(se.bs)

	//se.bs = nil
	//jdeEncPool.Put(se)

	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// rfType 是 剥离 Pointer 之后的最终类型
func (se *subEncode) initMeta(rfType reflect.Type, ptr unsafe.Pointer) {
	typAddr := (*dataType)((*emptyInterface)(unsafe.Pointer(&rfType)).dataPtr)
	if meta := cacheGetEncMeta(typAddr); meta != nil {
		se.em = meta
	} else {
		se.buildMeta(rfType)
		cacheSetEncMeta(typAddr, se.em)
	}
	se.srcPtr = ptr

	//if se.em.isSuperKV {
	//	//if se.em.isGson {
	//	//	se.gr = (*gson.GsonRow)(ptr)
	//	//} else {
	//	//	se.mp = (*cst.KV)(ptr)
	//	//}
	//} else {
	//	se.srcPtr = ptr // 当前值的地址
	//}
	return
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (se *subEncode) buildMeta(rfType reflect.Type) {
	se.em = &encMeta{}

	switch kd := rfType.Kind(); kd {
	case reflect.Array, reflect.Slice:
		se.initListMeta(rfType)
		se.bindListEnc()
	case reflect.Struct:
		// 模拟泛型解析，提供性能
		if rfType.String() == "gson.GsonRow" {
			se.em.isSuperKV = true
			se.em.isGson = true
			//se.bindGsonDec()
			return
		}
		if rfType.String() == "time.Time" {
			panic(errValueType)
		}
		se.initStructMeta(rfType)
		se.bindStructEnc()
	case reflect.Map:
		// 常规泛型
		if rfType.String() == "cst.KV" || rfType.String() == "map[string]interface {}" {
			se.em.isSuperKV = true
			se.em.isMap = true
			//se.bindMapDec()
			return
		}
		panic(errValueType)
	default:
		panic(errValueType)
	}
}

func (se *subEncode) initStructMeta(rfType reflect.Type) {
	se.em.isStruct = true
	se.em.ss = dts.SchemaForInputByType(rfType)
}

func (se *subEncode) initListMeta(rfType reflect.Type) {
	se.em.isList = true

	se.em.itemBaseType = rfType.Elem()
	se.em.itemBaseKind = se.em.itemBaseType.Kind()
	se.em.itemBytes = int(se.em.itemBaseType.Size())

peelPtr:
	if se.em.itemBaseKind == reflect.Pointer {
		se.em.itemBaseType = se.em.itemBaseType.Elem()
		se.em.itemBaseKind = se.em.itemBaseType.Kind()
		se.em.isPtr = true
		se.em.ptrLevel++
		// TODO：指针嵌套不能超过3层，这种很少见，也就是说此解码方案并不通用
		if se.em.ptrLevel > 3 {
			panic(errPtrLevel)
		}
		goto peelPtr
	}

	// 是否是interface类型
	if se.em.itemBaseKind == reflect.Interface {
		se.em.isAny = true
	}

	// 进一步初始化数组
	if rfType.Kind() == reflect.Array {
		se.em.isArray = true
		se.em.arrLen = rfType.Len() // 数组长度

		if se.em.isPtr {
			return
		}
		se.em.itemBytes = int(se.em.itemBaseType.Size())
		se.em.isArrBind = true
	}
}

func (se *subEncode) bindStructEnc() {
	fLen := len(se.em.ss.FieldsAttr)
	se.em.fieldsEnc = make([]encValFunc, fLen)
	se.em.fieldsEncMix = make([]encMixValFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}

	switch se.em.ss.FieldsAttr[i].Kind {
	case reflect.Int:
		se.em.fieldsEnc[i] = encInt[int]
	case reflect.Int8:
		se.em.fieldsEnc[i] = encInt[int8]
	case reflect.Int16:
		se.em.fieldsEnc[i] = encInt[int16]
	case reflect.Int32:
		se.em.fieldsEnc[i] = encInt[int32]
	case reflect.Int64:
		se.em.fieldsEnc[i] = encInt[int64]
	case reflect.Uint:
		se.em.fieldsEnc[i] = encUint[uint]
	case reflect.Uint8:
		se.em.fieldsEnc[i] = encUint[uint8]
	case reflect.Uint16:
		se.em.fieldsEnc[i] = encUint[uint16]
	case reflect.Uint32:
		se.em.fieldsEnc[i] = encUint[uint32]
	case reflect.Uint64:
		se.em.fieldsEnc[i] = encUint[uint64]
	case reflect.Float32:
		se.em.fieldsEnc[i] = encFloat[float32]
	case reflect.Float64:
		se.em.fieldsEnc[i] = encFloat[float64]
	case reflect.String:
		se.em.fieldsEnc[i] = encString
	case reflect.Bool:
		se.em.fieldsEnc[i] = encBool
	case reflect.Interface:
		se.em.fieldsEnc[i] = encAny

	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		se.em.fieldsEncMix[i] = encMixItem
	default:
		panic(errValueType)
	}
	goto nextField
}

func (se *subEncode) bindListEnc() {
	// 如果是数组，而且数组项类型不是指针类型
	//if se.em.isArrBind {
	switch se.em.itemBaseKind {
	case reflect.Int:
		se.em.listItemEnc = encInt[int]
	case reflect.Int8:
		se.em.listItemEnc = encInt[int8]
	case reflect.Int16:
		se.em.listItemEnc = encInt[int16]
	case reflect.Int32:
		se.em.listItemEnc = encInt[int32]
	case reflect.Int64:
		se.em.listItemEnc = encInt[int64]
	case reflect.Uint:
		se.em.listItemEnc = encUint[uint]
	case reflect.Uint8:
		se.em.listItemEnc = encUint[uint8]
	case reflect.Uint16:
		se.em.listItemEnc = encUint[uint16]
	case reflect.Uint32:
		se.em.listItemEnc = encUint[uint32]
	case reflect.Uint64:
		se.em.listItemEnc = encUint[uint64]
	case reflect.Float32:
		se.em.listItemEnc = encFloat[float32]
	case reflect.Float64:
		se.em.listItemEnc = encFloat[float64]
	case reflect.String:
		se.em.listItemEnc = encString
	case reflect.Bool:
		se.em.listItemEnc = encBool
	case reflect.Interface:
		se.em.listItemEnc = encAny

	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		se.em.listItemEncMix = encMixItem
	default:
		panic(errValueType)
	}
	return
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (se *subEncode) warpErrorCode(errCode errType) error {
//	if errCode >= 0 {
//		return nil
//	}
//
//	//sta := sd.scan
//	//end := sta + 20 // 输出标记后面 n 个字符
//	//if end > len(sd.str) {
//	//	end = len(sd.str)
//	//}
//
//	//errMsg := fmt.Sprintf("jde: %s, pos %d, character %q near ( %s )", errDescription[-errCode], sta, sd.str[sta], sd.str[sta:end])
//	////errMsg := strings.Join([]string{"jsonx: error pos: ", strconv.Itoa(sta), ", near ", string(sd.str[sta]), " of (", sd.str[sta:end], ")"}, "")
//	//return errors.New(errMsg)
//
//	return nil
//}
