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
		keyKind reflect.Kind
		keyEnc  encValFunc
		keySize uint32

		// status
		isSuperKV bool // {} SuperKV
		isMap     bool // {} map
		isStruct  bool // {} struct
		isList    bool // [] array & slice
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
	se.bf = pool.GetBytes()

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
	em.keyKind = rfType.Key().Kind()
	em.keySize = uint32(rfType.Key().Size())
	em.bindMapKeyEnc()

	// map 中的 value 可能是指针类型，需要拆包
	em.peelPtr(rfType)
	em.bindValueEnc()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (em *encMeta) bindValueEnc() {
	switch em.itemKind {
	default:
		panic(errValueType)

	case reflect.Int:
		em.itemEnc = encInt[int]
	case reflect.Int8:
		em.itemEnc = encInt[int8]
	case reflect.Int16:
		em.itemEnc = encInt[int16]
	case reflect.Int32:
		em.itemEnc = encInt[int32]
	case reflect.Int64:
		em.itemEnc = encInt[int64]
	case reflect.Uint:
		em.itemEnc = encUint[uint]
	case reflect.Uint8:
		em.itemEnc = encUint[uint8]
	case reflect.Uint16:
		em.itemEnc = encUint[uint16]
	case reflect.Uint32:
		em.itemEnc = encUint[uint32]
	case reflect.Uint64:
		em.itemEnc = encUint[uint64]
	case reflect.Float32:
		em.itemEnc = encFloat32
	case reflect.Float64:
		em.itemEnc = encFloat64

	case reflect.String:
		em.itemEnc = encString
	case reflect.Bool:
		em.itemEnc = encBool
	case reflect.Interface:
		em.itemEnc = encAny

	case reflect.Pointer:
		em.itemEnc = encMixedItem
	case reflect.Map, reflect.Array:
		em.itemEnc = encMixedItem
	case reflect.Slice:
		if em.itemType == cst.TypeBytes {
			em.itemEnc = encBytes
		} else {
			em.itemEnc = encMixedItem
		}
	case reflect.Struct:
		if em.itemType == cst.TypeTime {
			em.itemEnc = encTime
		} else {
			em.itemEnc = encMixedItem
		}
	}
	return
}

// 编码 map 的 Key 值，当前只支持如下的 Key 值类型。
// TODO: 这里只支持常见的 map 类型，暂时不支持复杂map
func (em *encMeta) bindMapKeyEnc() {
	switch em.keyKind {
	default:
		panic(errValueType)

	case reflect.Int:
		em.keyEnc = encInt[int]
	case reflect.Int8:
		em.keyEnc = encInt[int8]
	case reflect.Int16:
		em.keyEnc = encInt[int16]
	case reflect.Int32:
		em.keyEnc = encInt[int32]
	case reflect.Int64:
		em.keyEnc = encInt[int64]

	case reflect.Uint:
		em.keyEnc = encUint[uint]
	case reflect.Uint8:
		em.keyEnc = encUint[uint8]
	case reflect.Uint16:
		em.keyEnc = encUint[uint16]
	case reflect.Uint32:
		em.keyEnc = encUint[uint32]
	case reflect.Uint64:
		em.keyEnc = encUint[uint64]
	case reflect.Uintptr:
		em.keyEnc = encUint[uint64]

	case reflect.String:
		em.keyEnc = encString
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

	switch em.ss.FieldsAttr[i].Kind {
	default:
		panic(errValueType)

	case reflect.Int:
		em.fieldsEnc[i] = encInt[int]
	case reflect.Int8:
		em.fieldsEnc[i] = encInt[int8]
	case reflect.Int16:
		em.fieldsEnc[i] = encInt[int16]
	case reflect.Int32:
		em.fieldsEnc[i] = encInt[int32]
	case reflect.Int64:
		em.fieldsEnc[i] = encInt[int64]
	case reflect.Uint:
		em.fieldsEnc[i] = encUint[uint]
	case reflect.Uint8:
		em.fieldsEnc[i] = encUint[uint8]
	case reflect.Uint16:
		em.fieldsEnc[i] = encUint[uint16]
	case reflect.Uint32:
		em.fieldsEnc[i] = encUint[uint32]
	case reflect.Uint64:
		em.fieldsEnc[i] = encUint[uint64]
	case reflect.Float32:
		em.fieldsEnc[i] = encFloat32
	case reflect.Float64:
		em.fieldsEnc[i] = encFloat64

	case reflect.String:
		em.fieldsEnc[i] = encString
	case reflect.Bool:
		em.fieldsEnc[i] = encBool
	case reflect.Interface:
		em.fieldsEnc[i] = encAny

	case reflect.Pointer:
		em.fieldsEnc[i] = encMixedItem
	case reflect.Map, reflect.Array:
		em.fieldsEnc[i] = encMixedItem
	case reflect.Slice:
		if em.ss.FieldsAttr[i].Type == cst.TypeBytes {
			em.fieldsEnc[i] = encBytes
		} else {
			em.fieldsEnc[i] = encMixedItem
		}
	case reflect.Struct:
		if em.ss.FieldsAttr[i].Type == cst.TypeTime {
			em.fieldsEnc[i] = encTime
		} else {
			em.fieldsEnc[i] = encMixedItem
		}
	}

	goto nextField
}

func (em *encMeta) bindListEnc() {
	if !em.isPtr {
		switch em.itemKind {
		default:
			panic(errValueType)

		case reflect.Int:
			em.listEnc = encIntList[int]
		case reflect.Int8:
			em.listEnc = encIntList[int8]
		case reflect.Int16:
			em.listEnc = encIntList[int16]
		case reflect.Int32:
			em.listEnc = encIntList[int32]
		case reflect.Int64:
			em.listEnc = encIntList[int64]
		case reflect.Uint:
			em.listEnc = encUintList[uint]
		case reflect.Uint8:
			em.listEnc = encUintList[uint8]
		case reflect.Uint16:
			em.listEnc = encUintList[uint16]
		case reflect.Uint32:
			em.listEnc = encUintList[uint32]
		case reflect.Uint64:
			em.listEnc = encUintList[uint64]
		case reflect.Float32:
			em.listEnc = encAllList
		case reflect.Float64:
			em.listEnc = encAllList

		case reflect.String:
			em.listEnc = encStringList
		case reflect.Bool:
			em.listEnc = encAllList
		case reflect.Interface:
			em.listEnc = encAllList

		//case reflect.Pointer:
		//	em.listEnc = encPointer // 这个分支不可能
		case reflect.Map, reflect.Array:
			em.listEnc = encAllList
		case reflect.Slice:
			em.listEnc = encAllList
		case reflect.Struct:
			// 分情况，如果是时间类型，单独处理
			if em.itemType == cst.TypeTime {
				em.listEnc = encAllList
			} else {
				em.listEnc = encStructList
			}
		}
		return
	}

	// []*item 形式
	switch em.itemKind {
	default:
		panic(errValueType)

	case reflect.Int:
		em.listEnc = encIntListPtr[int]
	case reflect.Int8:
		em.listEnc = encIntListPtr[int8]
	case reflect.Int16:
		em.listEnc = encIntListPtr[int16]
	case reflect.Int32:
		em.listEnc = encIntListPtr[int32]
	case reflect.Int64:
		em.listEnc = encIntListPtr[int64]
	case reflect.Uint:
		em.listEnc = encIntListPtr[uint]
	case reflect.Uint8:
		em.listEnc = encIntListPtr[uint8]
	case reflect.Uint16:
		em.listEnc = encIntListPtr[uint16]
	case reflect.Uint32:
		em.listEnc = encIntListPtr[uint32]
	case reflect.Uint64:
		em.listEnc = encIntListPtr[uint64]
	case reflect.Float32:
		em.listEnc = encAllListPtr
	case reflect.Float64:
		em.listEnc = encAllListPtr

	case reflect.String:
		em.listEnc = encAllListPtr
	case reflect.Bool:
		em.listEnc = encAllListPtr
	case reflect.Interface:
		em.listEnc = encAllListPtr

	case reflect.Pointer:
		em.listEnc = encAllListPtr
	case reflect.Map, reflect.Array:
		em.listEnc = encAllListPtr
	case reflect.Slice:
		em.listEnc = encAllListPtr
	case reflect.Struct:
		// 分情况，如果是时间类型，单独处理
		if em.itemType == cst.TypeTime {
			em.listEnc = encAllListPtr
		} else {
			em.listEnc = encStructList
		}
	}
	return
}
