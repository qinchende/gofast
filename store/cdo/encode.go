// Copyright 2024 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package cdo

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"unsafe"
)

type (
	//encKeyFunc func(bf *[]byte, ptr unsafe.Pointer)
	encValFunc  func(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type)
	encListFunc func(e *subEncode, listSize int)

	subEncode struct {
		srcPtr unsafe.Pointer // list or object 对象值地址（其指向的值不能为nil，也不能为指针）
		em     *encMeta       // Struct | Slice,Array
		bf     *[]byte        // 当解析数组时候用到的一系列临时队列
		//objIdx uint16         // TODO：涉及多少种不同的Struct
	}

	encMeta struct {
		// Struct
		ss        *dts.StructSchema
		fieldsEnc []encValFunc

		// array & slice & map
		itemType    reflect.Type
		itemKind    reflect.Kind
		itemEnc     encValFunc
		itemMemSize int         // item类型对应的内存字节大小
		arrLen      int         // 数组长度
		listEnc     encListFunc // List整体编码

		// map
		keyType reflect.Type
		keyKind reflect.Kind
		keyEnc  encValFunc
		keySize uint32

		// status
		isSuperKV bool // {} SuperKV
		isMap     bool // {} map
		isStruct  bool // {} struct
		isList    bool // [] array & slice
		isSlice   bool // [] slice
		isArray   bool // [] array

		isPtr    bool  // [] is list and item is pointer type
		ptrLevel uint8 // [] is list and item pointer level(max 256 deeps)
	}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func cdoEncode(v any) (bs []byte, err error) {
	if v == nil {
		return nil, nil
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
	se.bf = pool.GetBytesLarge()

	se.startEncode()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)

	pool.FreeBytes(se.bf)
	return
}

// Use SubEncode to encode Mix Item Value
func encMixedItem(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	se := subEncode{}
	se.getEncMeta(typ, ptr)
	se.bf = bf
	se.startEncode()
}

func (se *subEncode) startEncode() {
	// Note: add by sdx on 2024-06-06
	// 这里将数组和切片的情况合并考虑，简化了代码；
	// 但通常我们遇到的都是切片类型，如果分开处理，将能进一步提高约 10% 的性能。
	if se.em.isList {
		listSize := 0
		if !se.em.isArray {
			sh := (*rt.SliceHeader)(se.srcPtr)
			se.srcPtr = sh.DataPtr
			listSize = sh.Len
		} else {
			listSize = se.em.arrLen
		}
		se.em.listEnc(se, listSize)
		return
	}

	switch {
	default:
		se.encBasic()
	case se.em.isStruct:
		se.encStruct()
	case se.em.isMap:
		se.encMap()
	case se.em.isPtr:
		se.encPointer()
	}
}

func (se *subEncode) getEncMeta(rfType reflect.Type, ptr unsafe.Pointer) {
	// 最多只能剥掉一层指针
	if rfType.Kind() == reflect.Pointer {
		rfType = rfType.Elem()
		// Note：有些类型本质其实是指针，但是reflect.Kind() != reflect.Pointer
		// 比如：map | channel | func
		// 此时需要统一变量值 ptr 指向的内存
		if rfType.Kind() == reflect.Map {
			ptr = *(*unsafe.Pointer)(ptr)
		}
	}

	if meta := cacheGetEncMeta(rfType); meta != nil {
		se.em = meta
	} else {
		se.em = newEncMeta(rfType)
		cacheSetEncMeta(rfType, se.em)
	}
	se.srcPtr = ptr
}

// EncodeMeta
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func newEncMeta(rfType reflect.Type) *encMeta {
	em := encMeta{}

	switch kd := rfType.Kind(); kd {
	default:
		em.initBaseTypeMeta(rfType)
	case reflect.Struct:
		if rfType == cst.TypeTime {
			em.initBaseTypeMeta(rfType)
			break
		}
		em.initStructMeta(rfType)
	case reflect.Pointer:
		em.initPointerMeta(rfType)
	case reflect.Map:
		em.initMapMeta(rfType)
	case reflect.Array, reflect.Slice:
		em.initListMeta(rfType)
	}

	return &em
}

// 该类型中，项目值的类型解析
func (em *encMeta) peelPtr(rfType reflect.Type) {
	em.itemType = rfType.Elem()
	em.itemKind = em.itemType.Kind()
	em.itemMemSize = int(em.itemType.Size())

peelNext:
	if em.itemKind == reflect.Pointer {
		em.isPtr = true
		em.itemType = em.itemType.Elem()
		em.itemKind = em.itemType.Kind()
		em.ptrLevel++
		goto peelNext
	}
}

func (em *encMeta) initBaseTypeMeta(rfType reflect.Type) {
	em.itemKind = rfType.Kind()
	em.bindValueEnc()
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

	if rfType.Kind() == reflect.Array {
		em.isArray = true
		em.arrLen = rfType.Len() // 数组长度
	} else {
		em.isSlice = true
	}

	// List 项如果是 struct ，是本编解码方案重点处理的情况
	if em.itemKind == reflect.Struct && em.itemType != cst.TypeTime {
		em.ss = dts.SchemaAsReqByType(em.itemType)
		em.bindStructFieldsEnc()
	} else {
		em.bindValueEnc()
	}

	em.bindListEnc()
}

// ++++++++++++++++++++++++++++++ Struct
func (em *encMeta) initStructMeta(rfType reflect.Type) {
	em.isStruct = true
	em.ss = dts.SchemaAsReqByType(rfType)
	em.itemMemSize = int(rfType.Size())
	em.bindStructFieldsEnc()
}

// ++++++++++++++++++++++++++++++ Map
// Note: 当前只支持 map[string]any 形式
func (em *encMeta) initMapMeta(rfType reflect.Type) {
	em.isMap = true

	// 特殊的Map单独处理，提高性能, 当前只支持 map[string]any 形式
	if rfType == cst.TypeCstKV || rfType == cst.TypeStrAnyMap {
		em.isSuperKV = true
	}

	// Note: map 中的 key 只支持几种特定类型
	em.keyType = rfType.Key()
	em.keyKind = em.keyType.Kind()
	em.keySize = uint32(rfType.Key().Size())
	em.bindMapKeyEnc()

	// map 中的 value 可能是指针类型，需要拆包
	em.peelPtr(rfType)
	em.bindValueEnc()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (em *encMeta) bindValueEnc() {
	innerBindValueEnc(em.itemType, &em.itemEnc)
}

// 编码 map 的 Key 值，当前只支持如下的 Key 值类型。
// TODO: 这里只支持常见的 map 类型，暂时不支持复杂map
func (em *encMeta) bindMapKeyEnc() {
	switch {
	default:
		panic(errValueType)
	case em.keyKind <= reflect.Float64, em.keyKind == reflect.String:
		innerBindValueEnc(em.keyType, &em.keyEnc)
	}
	return
}

// Struct对象，各字段的编码函数
func (em *encMeta) bindStructFieldsEnc() {
	fLen := len(em.ss.FieldsAttr)
	em.fieldsEnc = make([]encValFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}
	innerBindValueEnc(em.ss.FieldsAttr[i].Type, &em.fieldsEnc[i])
	goto nextField
}

func innerBindValueEnc(typ reflect.Type, encFunc *encValFunc) {
	switch typ.Kind() {
	default:
		panic(errValueType)

	case reflect.Int:
		*encFunc = encInt[int]
	case reflect.Int8:
		*encFunc = encInt[int8]
	case reflect.Int16:
		*encFunc = encInt[int16]
	case reflect.Int32:
		*encFunc = encInt[int32]
	case reflect.Int64:
		*encFunc = encInt[int64]

	case reflect.Uint:
		*encFunc = encUint[uint]
	case reflect.Uint8:
		*encFunc = encUint[uint8]
	case reflect.Uint16:
		*encFunc = encUint[uint16]
	case reflect.Uint32:
		*encFunc = encUint[uint32]
	case reflect.Uint64:
		*encFunc = encUint[uint64]
	case reflect.Uintptr:
		*encFunc = encUint[uintptr]

	case reflect.Float32:
		*encFunc = encF32
	case reflect.Float64:
		*encFunc = encF64

	case reflect.String:
		*encFunc = encString
	case reflect.Bool:
		*encFunc = encBool
	case reflect.Interface:
		*encFunc = encAny

	case reflect.Pointer, reflect.Map, reflect.Array:
		*encFunc = encMixedItem
	case reflect.Slice:
		if typ == cst.TypeBytes {
			*encFunc = encBytes
		} else {
			*encFunc = encMixedItem
		}
	case reflect.Struct:
		if typ == cst.TypeTime {
			*encFunc = encTime
		} else {
			*encFunc = encMixedItem
		}
	}
}

func (em *encMeta) bindListEnc() {
	if !em.isPtr {
		switch em.itemKind {
		default:
			panic(errValueType)

		case reflect.Int:
			em.listEnc = encListInt[int]
		case reflect.Int8:
			em.listEnc = encListInt[int8]
		case reflect.Int16:
			em.listEnc = encListInt[int16]
		case reflect.Int32:
			em.listEnc = encListInt[int32]
		case reflect.Int64:
			em.listEnc = encListInt[int64]
		case reflect.Uint:
			em.listEnc = encListUint[uint]
		case reflect.Uint8:
			em.listEnc = encListUint[uint8]
		case reflect.Uint16:
			em.listEnc = encListUint[uint16]
		case reflect.Uint32:
			em.listEnc = encListUint[uint32]
		case reflect.Uint64:
			em.listEnc = encListUint[uint64]
		case reflect.Float32:
			em.listEnc = encListF32
		case reflect.Float64:
			em.listEnc = encListF64

		case reflect.String:
			em.listEnc = encListString
		case reflect.Bool:
			em.listEnc = encListBool
		case reflect.Interface:
			em.listEnc = encListAll

		//case reflect.Pointer:
		//	em.listEnc = encPointer // 这个分支不可能
		case reflect.Map, reflect.Array:
			em.listEnc = encListAll
		case reflect.Slice:
			em.listEnc = encListAll
		case reflect.Struct:
			if em.itemType == cst.TypeTime {
				em.listEnc = encListAll
			} else {
				em.listEnc = encListStruct
			}
		}
		return
	}

	// []*item 形式
	switch em.itemKind {
	default:
		panic(errValueType)

	case reflect.Int:
		em.listEnc = encListIntPtr[int]
	case reflect.Int8:
		em.listEnc = encListIntPtr[int8]
	case reflect.Int16:
		em.listEnc = encListIntPtr[int16]
	case reflect.Int32:
		em.listEnc = encListIntPtr[int32]
	case reflect.Int64:
		em.listEnc = encListIntPtr[int64]
	case reflect.Uint:
		em.listEnc = encListIntPtr[uint]
	case reflect.Uint8:
		em.listEnc = encListIntPtr[uint8]
	case reflect.Uint16:
		em.listEnc = encListIntPtr[uint16]
	case reflect.Uint32:
		em.listEnc = encListIntPtr[uint32]
	case reflect.Uint64:
		em.listEnc = encListIntPtr[uint64]
	case reflect.Float32:
		em.listEnc = encListAllPtr
	case reflect.Float64:
		em.listEnc = encListAllPtr

	case reflect.String:
		em.listEnc = encListAllPtr
	case reflect.Bool:
		em.listEnc = encListAllPtr
	case reflect.Interface:
		em.listEnc = encListAllPtr

	case reflect.Pointer:
		em.listEnc = encListAllPtr
	case reflect.Map, reflect.Array:
		em.listEnc = encListAllPtr
	case reflect.Slice:
		em.listEnc = encListAllPtr
	case reflect.Struct:
		if em.itemType == cst.TypeTime {
			em.listEnc = encListAllPtr
		} else {
			em.listEnc = encListStruct
		}
	}
	return
}
