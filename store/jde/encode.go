package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
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

		// array & slice & map
		itemBaseType reflect.Type
		itemBaseKind reflect.Kind
		itemPick     encValFunc
		itemMemSize  int // item类型对应的内存字节大小
		arrLen       int // 数组长度

		// map
		keyBaseType reflect.Type
		keyPick     encValFunc

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
	se.initMeta(reflect.TypeOf(v), (*rt.AFace)(unsafe.Pointer(&v)).DataPtr)
	se.newBytesBuf()

	se.encStart()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)

	se.freeBytesBuf()
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) initMeta(rfType reflect.Type, ptr unsafe.Pointer) {
	typAddr := (*rt.TypeAgent)((*rt.AFace)(unsafe.Pointer(&rfType)).DataPtr)
	if meta := cacheGetEncMeta(typAddr); meta != nil {
		se.em = meta
	} else {
		se.buildEncMeta(rfType)
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

func (se *subEncode) buildEncMeta(rfType reflect.Type) {
	if rfType.Kind() == reflect.Pointer {
		rfType = rfType.Elem()
	}
	em := encMeta{}
	se.em = &em

	switch kd := rfType.Kind(); kd {
	case reflect.Array, reflect.Slice:
		se.initListMeta(rfType)
		se.bindPick(em.itemBaseKind, &em.itemPick)
	case reflect.Struct:
		// GoFast Special type GsonRow
		if rfType.String() == gson.StrTypeOfGsonRow {
			em.isSuperKV = true
			em.isGson = true
			//se.bindGsonDec()
			return
		}
		if rfType.String() == cst.StrTypeOfTime {
			panic(errValueType)
		}
		se.initStructMeta(rfType)
	case reflect.Map:
		// Map type
		se.initMapMeta(rfType)
	case reflect.Pointer:
		// Pointer type
		se.initPointerMeta(rfType)
	default:
		// Others normal types
		se.bindPick(kd, &em.itemPick)
	}
}

func (se *subEncode) initPointerMeta(rfType reflect.Type) {
	se.em.isPtr = true
	se.em.itemBaseType = rfType.Elem()
	se.em.itemBaseKind = se.em.itemBaseType.Kind()
	se.em.ptrLevel++

peelPtr:
	if se.em.itemBaseKind == reflect.Pointer {
		se.em.itemBaseType = se.em.itemBaseType.Elem()
		se.em.itemBaseKind = se.em.itemBaseType.Kind()
		se.em.ptrLevel++
		goto peelPtr
	}
}

func (se *subEncode) initListMeta(rfType reflect.Type) {
	se.em.isList = true
	se.em.itemBaseType = rfType.Elem()
	se.em.itemBaseKind = se.em.itemBaseType.Kind()
	se.em.itemMemSize = int(se.em.itemBaseType.Size())

peelPtr:
	if se.em.itemBaseKind == reflect.Pointer {
		se.em.isPtr = true
		se.em.itemBaseType = se.em.itemBaseType.Elem()
		se.em.itemBaseKind = se.em.itemBaseType.Kind()
		se.em.ptrLevel++
		goto peelPtr
	}

	// 如果原始类型是Array，进一步提取数组信息
	if rfType.Kind() == reflect.Array {
		se.em.isArray = true
		se.em.arrLen = rfType.Len() // 数组长度

		if se.em.isPtr {
			return
		}
		se.em.itemMemSize = int(se.em.itemBaseType.Size())
	}
}

func (se *subEncode) initStructMeta(rfType reflect.Type) {
	se.em.isStruct = true
	se.em.ss = dts.SchemaForInputByType(rfType)

	se.bindStructPick()
}

// Note: 当前只支持 map[string]any 形式
func (se *subEncode) initMapMeta(rfType reflect.Type) {
	se.em.isMap = true

	// Note: map 中的 key 不能是 指针类耐
	se.em.keyBaseType = rfType.Key()
	if se.em.keyBaseType.Kind() != reflect.String {
		panic(errValueType)
	}
	se.bindPick(se.em.keyBaseType.Kind(), &se.em.keyPick)

	// map 中的 value 可以是指针类型，需要拆包
	se.em.itemBaseType = rfType.Elem()
	se.em.itemBaseKind = se.em.itemBaseType.Kind()

peelPtr:
	if se.em.itemBaseKind == reflect.Pointer {
		se.em.isPtr = true
		se.em.itemBaseType = se.em.itemBaseType.Elem()
		se.em.itemBaseKind = se.em.itemBaseType.Kind()
		se.em.ptrLevel++
		goto peelPtr
	}

	se.bindPick(se.em.itemBaseKind, &se.em.itemPick)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) bindPick(kind reflect.Kind, pick *encValFunc) {
	switch kind {
	case reflect.Int:
		*pick = encInt[int]
	case reflect.Int8:
		*pick = encInt[int8]
	case reflect.Int16:
		*pick = encInt[int16]
	case reflect.Int32:
		*pick = encInt[int32]
	case reflect.Int64:
		*pick = encInt[int64]
	case reflect.Uint:
		*pick = encUint[uint]
	case reflect.Uint8:
		*pick = encUint[uint8]
	case reflect.Uint16:
		*pick = encUint[uint16]
	case reflect.Uint32:
		*pick = encUint[uint32]
	case reflect.Uint64:
		*pick = encUint[uint64]
	case reflect.Float32:
		*pick = encFloat[float32]
	case reflect.Float64:
		*pick = encFloat[float64]

	case reflect.String:
		*pick = encString
	case reflect.Bool:
		*pick = encBool

	case reflect.Interface:
		*pick = encAny
	//case reflect.Pointer:
	//	*pick = encPointer

	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		*pick = encMixItem
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
