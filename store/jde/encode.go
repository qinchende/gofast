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

type (
	encValFunc func(bf []byte, ptr unsafe.Pointer, typ reflect.Type) (nbf []byte)

	subEncode struct {
		srcPtr unsafe.Pointer // list or object 对象值地址（其指向的值不能为nil，也不能为指针）
		em     *encMeta       // Struct | Slice,Array
		bf     *[]byte        // 当解析数组时候用到的一系列临时队列
		mp     *cst.KV        // map
		gr     *gson.GsonRow  // GsonRow
	}

	encMeta struct {
		// Struct
		ss         *dts.StructSchema
		fieldsPick []encValFunc

		// array & slice
		itemBaseType reflect.Type
		itemBaseKind reflect.Kind
		itemPick     encValFunc
		itemMemSize  int // item类型对应的内存字节大小
		arrLen       int // 数组长度

		// status
		isSuperKV bool // {} SuperKV
		isGson    bool // {} gson
		isMap     bool // {} map
		isStruct  bool // {} struct

		//isAny    bool  // [] is list and item is interface type in the final
		isList   bool  // [] array & slice
		isArray  bool  // [] array
		isPtr    bool  // [] is list and item is pointer type
		ptrLevel uint8 // [] is list and item pointer level
	}
)

// 主解析入口
func startEncode(v any) (bs []byte, err error) {
	if v == nil {
		return nullBytes, nil
	}

	defer func() {
		if pic := recover(); pic != nil {
			err = errors.New(fmt.Sprint(pic))
		}
	}()

	se := subEncode{}
	se.initMeta(reflect.TypeOf(v), (*emptyInterface)(unsafe.Pointer(&v)).dataPtr)
	se.newBytesBuf()

	se.encStart()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)

	se.freeBytesBuf()
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// rfType 是 剥离 Pointer 之后的最终类型
func (se *subEncode) initMeta(rfType reflect.Type, ptr unsafe.Pointer) {
	typAddr := (*dataType)((*emptyInterface)(unsafe.Pointer(&rfType)).dataPtr)
	if meta := cacheGetEncMeta(typAddr); meta != nil {
		se.em = meta
	} else {
		se.buildEncodeMeta(rfType)
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
func (se *subEncode) buildEncodeMeta(rfType reflect.Type) {
	if rfType.Kind() == reflect.Pointer {
		rfType = rfType.Elem()
	}
	em := encMeta{}
	se.em = &em

	switch kd := rfType.Kind(); kd {
	case reflect.Array, reflect.Slice:
		se.initListMeta(rfType)
		se.bindItemPick(em.itemBaseKind)
	case reflect.Struct:
		// 模拟泛型解析，提供性能
		if rfType.String() == "gson.GsonRow" {
			em.isSuperKV = true
			em.isGson = true
			//se.bindGsonDec()
			return
		}
		if rfType.String() == "time.Time" {
			panic(errValueType)
		}
		se.initStructMeta(rfType)
		se.bindStructPick()
	case reflect.Map:
		// 常规泛型
		if rfType.String() == "cst.KV" || rfType.String() == "map[string]interface {}" {
			em.isSuperKV = true
			em.isMap = true
			//se.bindMapDec()
			return
		}
		panic(errValueType)
	case reflect.Pointer:
		se.initPointerMeta(rfType)
	default:
		se.bindItemPick(kd)
	}
}

func (se *subEncode) initPointerMeta(rfType reflect.Type) {
	se.em.itemBaseType = rfType.Elem()
	se.em.itemBaseKind = se.em.itemBaseType.Kind()
	se.em.ptrLevel++
	se.em.isPtr = true
peelPtr:
	if se.em.itemBaseKind == reflect.Pointer {
		se.em.itemBaseType = se.em.itemBaseType.Elem()
		se.em.itemBaseKind = se.em.itemBaseType.Kind()
		se.em.ptrLevel++
		goto peelPtr
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
	se.em.itemMemSize = int(se.em.itemBaseType.Size())

peelPtr:
	if se.em.itemBaseKind == reflect.Pointer {
		se.em.itemBaseType = se.em.itemBaseType.Elem()
		se.em.itemBaseKind = se.em.itemBaseType.Kind()
		se.em.isPtr = true
		se.em.ptrLevel++
		goto peelPtr
	}

	//// 是否是interface类型
	//if se.em.itemBaseKind == reflect.Interface {
	//	se.em.isAny = true
	//}

	// 进一步初始化数组
	if rfType.Kind() == reflect.Array {
		se.em.isArray = true
		se.em.arrLen = rfType.Len() // 数组长度

		if se.em.isPtr {
			return
		}
		se.em.itemMemSize = int(se.em.itemBaseType.Size())
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) bindItemPick(kind reflect.Kind) {
	switch kind {
	case reflect.Int:
		se.em.itemPick = encInt[int]
	case reflect.Int8:
		se.em.itemPick = encInt[int8]
	case reflect.Int16:
		se.em.itemPick = encInt[int16]
	case reflect.Int32:
		se.em.itemPick = encInt[int32]
	case reflect.Int64:
		se.em.itemPick = encInt[int64]
	case reflect.Uint:
		se.em.itemPick = encUint[uint]
	case reflect.Uint8:
		se.em.itemPick = encUint[uint8]
	case reflect.Uint16:
		se.em.itemPick = encUint[uint16]
	case reflect.Uint32:
		se.em.itemPick = encUint[uint32]
	case reflect.Uint64:
		se.em.itemPick = encUint[uint64]
	case reflect.Float32:
		se.em.itemPick = encFloat[float32]
	case reflect.Float64:
		se.em.itemPick = encFloat[float64]

	case reflect.String:
		se.em.itemPick = encString
	case reflect.Bool:
		se.em.itemPick = encBool

	case reflect.Interface:
		se.em.itemPick = encAny
	//case reflect.Pointer:
	//	se.em.itemPick = encPointer

	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		se.em.itemPick = encMixItem
	default:
		panic(errValueType)
	}
	return
}

func (se *subEncode) bindStructPick() {
	fLen := len(se.em.ss.FieldsAttr)
	se.em.fieldsPick = make([]encValFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}

	switch se.em.ss.FieldsAttr[i].Kind {
	case reflect.Int:
		se.em.fieldsPick[i] = encInt[int]
	case reflect.Int8:
		se.em.fieldsPick[i] = encInt[int8]
	case reflect.Int16:
		se.em.fieldsPick[i] = encInt[int16]
	case reflect.Int32:
		se.em.fieldsPick[i] = encInt[int32]
	case reflect.Int64:
		se.em.fieldsPick[i] = encInt[int64]
	case reflect.Uint:
		se.em.fieldsPick[i] = encUint[uint]
	case reflect.Uint8:
		se.em.fieldsPick[i] = encUint[uint8]
	case reflect.Uint16:
		se.em.fieldsPick[i] = encUint[uint16]
	case reflect.Uint32:
		se.em.fieldsPick[i] = encUint[uint32]
	case reflect.Uint64:
		se.em.fieldsPick[i] = encUint[uint64]
	case reflect.Float32:
		se.em.fieldsPick[i] = encFloat[float32]
	case reflect.Float64:
		se.em.fieldsPick[i] = encFloat[float64]

	case reflect.String:
		se.em.fieldsPick[i] = encString
	case reflect.Bool:
		se.em.fieldsPick[i] = encBool

	case reflect.Interface:
		se.em.fieldsPick[i] = encAny
	//case reflect.Pointer:
	//	se.em.fieldsPick[i] = encPointer

	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		se.em.fieldsPick[i] = encMixItem
	default:
		panic(errValueType)
	}
	goto nextField
}
