package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"unsafe"
)

type (
	encValFunc func(bf []byte, ptr unsafe.Pointer, typ reflect.Type) (nbf []byte)

	subEncode struct {
		srcPtr unsafe.Pointer // list or object 对象值地址（其指向的值不能为nil，也不能为指针）
		em     *encMeta       // Struct | Slice,Array
		bf     *[]byte        // 当解析数组时候用到的一系列临时队列
		//mp     *cst.KV        // map
		//gr     *gson.GsonRow  // GsonRow
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
			if err1, ok := pic.(error); ok {
				err = err1
			} else {
				err = errors.New(fmt.Sprint(pic))
			}
		}
	}()

	se := subEncode{}
	se.getEncMeta(reflect.TypeOf(v), (*rt.AFace)(unsafe.Pointer(&v)).DataPtr)
	se.newBytesBuf()

	se.encStart()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)

	se.freeBytesBuf()
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) getEncMeta(rfType reflect.Type, ptr unsafe.Pointer) {
	typAddr := (*rt.TypeAgent)((*rt.AFace)(unsafe.Pointer(&rfType)).DataPtr)
	if meta := cacheGetEncMeta(typAddr); meta != nil {
		se.em = meta
	} else {
		se.em = newEncodeMeta(rfType)
		cacheSetEncMeta(typAddr, se.em)
	}
	se.srcPtr = ptr
}

func newEncodeMeta(rfType reflect.Type) *encMeta {
	if rfType.Kind() == reflect.Pointer {
		rfType = rfType.Elem()
	}
	em := encMeta{}

	switch kd := rfType.Kind(); kd {
	case reflect.Array, reflect.Slice:
		em.initListMeta(rfType)
	case reflect.Struct:
		// GoFast Special type GsonRow
		//if rfType.String() == gson.StrTypeOfGsonRow {
		//	em.isSuperKV = true
		//	em.isGson = true
		//	return
		//}
		if rfType.String() == cst.StrTypeOfTime {
			panic(errValueType)
		}
		em.initStructMeta(rfType)
	case reflect.Map:
		// Map type
		em.initMapMeta(rfType)
	case reflect.Pointer:
		// Pointer type
		em.initPointerMeta(rfType)
	default:
		// Others normal types
		bindPick(kd, &em.itemPick)
	}

	return &em
}

func (em *encMeta) peelPtr(rfType reflect.Type) {
	em.itemBaseType = rfType.Elem()
	em.itemBaseKind = em.itemBaseType.Kind()
	em.itemMemSize = int(em.itemBaseType.Size())

peelLoop:
	if em.itemBaseKind == reflect.Pointer {
		em.isPtr = true
		em.itemBaseType = em.itemBaseType.Elem()
		em.itemBaseKind = em.itemBaseType.Kind()
		em.ptrLevel++
		goto peelLoop
	}
}

// ++++++++++++++++++++++++++++++ Pointer
func (em *encMeta) initPointerMeta(rfType reflect.Type) {
	em.isPtr = true
	em.ptrLevel++
	em.peelPtr(rfType)
}

// ++++++++++++++++++++++++++++++ Array & Slice
func (em *encMeta) initListMeta(rfType reflect.Type) {
	em.isList = true
	em.peelPtr(rfType)

	// 如果原始类型是Array，进一步提取数组信息
	if rfType.Kind() == reflect.Array {
		em.isArray = true
		em.arrLen = rfType.Len() // 数组长度

		if em.isPtr {
			return
		}
		em.itemMemSize = int(em.itemBaseType.Size())
	}

	bindPick(em.itemBaseKind, &em.itemPick)
}

// ++++++++++++++++++++++++++++++ Struct
func (em *encMeta) initStructMeta(rfType reflect.Type) {
	em.isStruct = true
	em.ss = dts.SchemaForInputByType(rfType)

	em.bindStructPick()
}

// ++++++++++++++++++++++++++++++ Map
// Note: 当前只支持 map[string]any 形式
func (em *encMeta) initMapMeta(rfType reflect.Type) {
	em.isMap = true

	// 特殊的Map单独处理，提高性能, 当前只支持 map[string]any 形式
	typStr := rfType.String()
	if typStr == cst.StrTypeOfKV || typStr == cst.StrTypeOfStrAnyMap {
		em.isSuperKV = true
	} else {
		panic(errMapType)
	}

	// Note: map 中的 key 只支持几种特定类型
	em.keyBaseType = rfType.Key()
	switch em.keyBaseType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		panic(errMapType)
	}
	bindPick(em.keyBaseType.Kind(), &em.keyPick)

	// map 中的 value 可能是指针类型，需要拆包
	em.peelPtr(rfType)
	bindPick(em.itemBaseKind, &em.itemPick)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindPick(kind reflect.Kind, pick *encValFunc) {
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
		*pick = encFloat32
	case reflect.Float64:
		*pick = encFloat64

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

func (em *encMeta) bindStructPick() {
	fLen := len(em.ss.FieldsAttr)
	em.fieldsPick = make([]encValFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}

	switch em.ss.FieldsAttr[i].Kind {
	case reflect.Int:
		em.fieldsPick[i] = encInt[int]
	case reflect.Int8:
		em.fieldsPick[i] = encInt[int8]
	case reflect.Int16:
		em.fieldsPick[i] = encInt[int16]
	case reflect.Int32:
		em.fieldsPick[i] = encInt[int32]
	case reflect.Int64:
		em.fieldsPick[i] = encInt[int64]
	case reflect.Uint:
		em.fieldsPick[i] = encUint[uint]
	case reflect.Uint8:
		em.fieldsPick[i] = encUint[uint8]
	case reflect.Uint16:
		em.fieldsPick[i] = encUint[uint16]
	case reflect.Uint32:
		em.fieldsPick[i] = encUint[uint32]
	case reflect.Uint64:
		em.fieldsPick[i] = encUint[uint64]
	case reflect.Float32:
		em.fieldsPick[i] = encFloat32
	case reflect.Float64:
		em.fieldsPick[i] = encFloat64

	case reflect.String:
		em.fieldsPick[i] = encString
	case reflect.Bool:
		em.fieldsPick[i] = encBool

	case reflect.Interface:
		em.fieldsPick[i] = encAny
	//case reflect.Pointer:
	//	em.fieldsPick[i] = encPointer

	case reflect.Map, reflect.Struct, reflect.Array, reflect.Slice:
		em.fieldsPick[i] = encMixItem
	default:
		panic(errValueType)
	}
	goto nextField
}
