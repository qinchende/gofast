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
	encValFunc  func(bs []byte, ptr unsafe.Pointer, typ reflect.Type) []byte
	encListFunc func(e *encoder)

	encMeta struct {
		// struct
		ss        *dts.StructSchema
		fieldsEnc []encValFunc

		// array & slice & map & baseType
		itemType reflect.Type
		itemKind reflect.Kind
		itemEnc  encValFunc

		// array & slice & map
		itemMemSize int         // item类型对应的内存字节大小
		arrLen      int         // 数组长度
		listEnc     encListFunc // List整体编码

		// map
		keyType reflect.Type
		keyKind reflect.Kind
		keyEnc  encValFunc
		keySize uint32

		// type status
		isMap       bool // {} map
		isMapStrStr bool // {} map[str]str
		isMapStrAll bool // {} map[str]all-type
		isStruct    bool // {} struct
		isList      bool // [] array & slice
		isSlice     bool // [] slice
		isArray     bool // [] array

		// ext status
		isPtr    bool  // (curr-val | list-item-val | map-value) is ptr
		ptrLevel uint8 // ptr deep (max 256)
	}

	encoder struct {
		depth int      // 当前层级，有限制不能无限递归
		em    *encMeta // Struct | Slice,Array

		srcPtr unsafe.Pointer // list or object 对象值地址（其指向的值不能为nil，也不能为指针）
		slice  rt.SliceHeader // 用于将数组模拟成切片
		bf     *[]byte        // 当解析数组时候用到的一系列临时队列
		bs     []byte         // 用来辅助上面的bf指针，防止24个字节的切片对象堆分配
	}
)

// 默认值，用于缓存对象的重置
var _subEncodeDefValues encoder

func (se *encoder) reset() {
	*se = _subEncodeDefValues
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func cdoEncode(v any) (bs []byte, err error) {
	if v == nil {
		return nil, errValueIsNil
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

	se := cdoEncPool.Get().(*encoder)
	se.fetchEncMeta(v)

	se.bf = pool.GetBytes()
	se.execEnc()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)
	pool.FreeBytes(se.bf)

	se.reset()
	cdoEncPool.Put(se)
	return
}

func encMixedItemRet(bf []byte, ptr unsafe.Pointer, typ reflect.Type) (bs []byte) {
	se := cdoEncPool.Get().(*encoder)
	se.applyEncMeta(typ, ptr)

	se.bs = bf
	se.bf = &se.bs
	se.execEnc()
	bs = se.bs // must return bytes before se.reset()

	se.reset()
	cdoEncPool.Put(se)
	return
}

func (se *encoder) execEnc() {
	switch {
	default:
		if se.em.isPtr {
			se.encPointer()
			return
		}
		se.encBasic()
	case se.em.isList:
		se.encList()
	case se.em.isStruct:
		se.encStruct()
	case se.em.isMap:
		se.encMap()
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *encoder) fetchEncMeta(v any) {
	vType := reflect.TypeOf(v)
	ptr := (*rt.AFace)(unsafe.Pointer(&v)).DataPtr

	// v is a ptr value, need peel ptr
	if vType.Kind() == reflect.Pointer {
		vType = vType.Elem()
		ptr = realPtr(vType.Kind(), ptr)
	}

	se.applyEncMeta(vType, ptr)
}

func (se *encoder) applyEncMeta(vType reflect.Type, ptr unsafe.Pointer) {
	if meta := cacheGetEncMeta(vType); meta != nil {
		se.em = meta
	} else {
		se.em = newEncMeta(vType)
		cacheSetEncMeta(vType, se.em)
	}
	se.srcPtr = ptr
}

func newEncMeta(typ reflect.Type) *encMeta {
	em := encMeta{}

	switch kd := typ.Kind(); kd {
	default:
		em.initBaseTypeMeta(typ)
	case reflect.Pointer:
		em.initPointerMeta(typ)
	case reflect.Struct:
		if typ == cst.TypeTime {
			em.initBaseTypeMeta(typ)
			break
		}
		em.initStructMeta(typ)
	case reflect.Map:
		em.initMapMeta(typ)
	case reflect.Array, reflect.Slice:
		em.initListMeta(typ)
	}

	return &em
}

// 该类型中，项目值的类型解析
func (em *encMeta) peelPtr(typ reflect.Type) {
	em.itemType = typ.Elem()
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

func (em *encMeta) initBaseTypeMeta(typ reflect.Type) {
	em.itemType = typ
	em.itemKind = typ.Kind()
	bindEnc(em.itemType, &em.itemEnc)
}

// ++++++++++++++++++++++++++++++ Pointer
func (em *encMeta) initPointerMeta(typ reflect.Type) {
	em.isPtr = true
	em.ptrLevel++
	em.peelPtr(typ)
	bindEnc(em.itemType, &em.itemEnc)
}

// ++++++++++++++++++++++++++++++ Array & Slice
func (em *encMeta) initListMeta(typ reflect.Type) {
	em.isList = true
	em.peelPtr(typ)

	if typ.Kind() == reflect.Array {
		em.isArray = true
		em.arrLen = typ.Len() // 数组长度
	} else {
		em.isSlice = true
	}

	// List 项如果是 struct ，是本编解码方案重点处理的情况
	if em.itemKind == reflect.Struct && em.itemType != cst.TypeTime {
		em.ss = dts.SchemaAsReqByType(em.itemType)
		em.bindFieldsEnc()
	} else {
		bindEnc(em.itemType, &em.itemEnc)
	}

	em.bindListEnc()
}

// ++++++++++++++++++++++++++++++ Struct
func (em *encMeta) initStructMeta(typ reflect.Type) {
	em.isStruct = true
	em.ss = dts.SchemaAsReqByType(typ)
	em.itemMemSize = int(typ.Size())
	em.bindFieldsEnc()
}

// ++++++++++++++++++++++++++++++ Map
func (em *encMeta) initMapMeta(typ reflect.Type) {
	em.isMap = true

	// Note: map 中的 key 只支持几种特定类型
	em.keyType = typ.Key()
	em.keyKind = em.keyType.Kind()
	em.keySize = uint32(typ.Key().Size())

	if em.keyKind == reflect.String {
		em.isMapStrAll = true
		if typ.Elem().Kind() == reflect.String {
			em.isMapStrStr = true
		}
	}

	// key 的 encode 函数
	// Note: 目前 key 只支持一些基本类型
	switch {
	default:
		panic(errMap)
	case em.keyKind <= reflect.Float64, em.keyKind == reflect.String:
		bindEnc(em.keyType, &em.keyEnc)
	}

	// value 的 encode 函数
	em.peelPtr(typ) // value 可能是指针类型，需要拆包
	bindEnc(em.itemType, &em.itemEnc)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Struct对象，各字段的编码函数
func (em *encMeta) bindFieldsEnc() {
	fLen := len(em.ss.FieldsAttr)
	em.fieldsEnc = make([]encValFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}
	bindEnc(em.ss.FieldsAttr[i].Type, &em.fieldsEnc[i])
	goto nextField
}

func bindEnc(typ reflect.Type, encFunc *encValFunc) {
	switch typ.Kind() {
	default:
		panic(errValueType)

	case reflect.Int:
		*encFunc = encIntRet[int]
	case reflect.Int8:
		*encFunc = encIntRet[int8]
	case reflect.Int16:
		*encFunc = encIntRet[int16]
	case reflect.Int32:
		*encFunc = encIntRet[int32]
	case reflect.Int64:
		*encFunc = encIntRet[int64]

	case reflect.Uint:
		*encFunc = encUintRet[uint]
	case reflect.Uint8:
		*encFunc = encUintRet[uint8]
	case reflect.Uint16:
		*encFunc = encUintRet[uint16]
	case reflect.Uint32:
		*encFunc = encUintRet[uint32]
	case reflect.Uint64:
		*encFunc = encUintRet[uint64]
	case reflect.Uintptr:
		*encFunc = encUintRet[uintptr]

	case reflect.Float32:
		*encFunc = encF32Ret
	case reflect.Float64:
		*encFunc = encF64Ret
	case reflect.String:
		*encFunc = encStringRet
	case reflect.Bool:
		*encFunc = encBoolRet

	case reflect.Pointer:
		*encFunc = encMixedItemRet
	case reflect.Interface:
		*encFunc = encAnyRet
	case reflect.Map, reflect.Array:
		*encFunc = encMixedItemRet
	case reflect.Slice:
		if typ == cst.TypeBytes {
			*encFunc = encBytesRet
		} else {
			*encFunc = encMixedItemRet
		}
	case reflect.Struct:
		if typ == cst.TypeTime {
			*encFunc = encTimeRet
		} else {
			*encFunc = encMixedItemRet
		}
	}
}

func (em *encMeta) bindListEnc() {
	// 数据项是非指针类型
	if !em.isPtr {
		switch em.itemType.Kind() {
		default:
			panic(errValueType)

		case reflect.Int:
			em.listEnc = encListVarInt[int]
		case reflect.Int8:
			em.listEnc = encListVarInt[int8]
		case reflect.Int16:
			em.listEnc = encListVarInt[int16]
		case reflect.Int32:
			em.listEnc = encListVarInt[int32]
		case reflect.Int64:
			em.listEnc = encListVarInt[int64]
		case reflect.Uint:
			em.listEnc = encListVarUint[uint]
		case reflect.Uint8:
			em.listEnc = encListVarUint[uint8]
		case reflect.Uint16:
			em.listEnc = encListVarUint[uint16]
		case reflect.Uint32:
			em.listEnc = encListVarUint[uint32]
		case reflect.Uint64:
			em.listEnc = encListVarUint[uint64]

		case reflect.Float32:
			em.listEnc = encListF32
		case reflect.Float64:
			em.listEnc = encListF64
		case reflect.String:
			em.listEnc = encListString
		case reflect.Bool:
			em.listEnc = encListBool

		case reflect.Pointer: // 不可能
		case reflect.Interface:
			em.listEnc = encListAll
		case reflect.Map, reflect.Array:
			em.listEnc = encListAll
		case reflect.Slice:
			em.listEnc = encListAll
		case reflect.Struct:
			if em.itemType == cst.TypeTime {
				em.listEnc = encListTime
			} else {
				if em.ss.HasPtrField {
					em.listEnc = encListStructPtr
				} else {
					em.listEnc = encListStruct
				}
			}
		}
		return
	}

	// 数据项是指针类型
	// []*item 形式
	switch em.itemType.Kind() {
	default:
		panic(errValueType)

	case reflect.Int:
		em.listEnc = encListVarIntPtr[int]
	case reflect.Int8:
		em.listEnc = encListVarIntPtr[int8]
	case reflect.Int16:
		em.listEnc = encListVarIntPtr[int16]
	case reflect.Int32:
		em.listEnc = encListVarIntPtr[int32]
	case reflect.Int64:
		em.listEnc = encListVarIntPtr[int64]
	case reflect.Uint:
		em.listEnc = encListVarIntPtr[uint]
	case reflect.Uint8:
		em.listEnc = encListVarIntPtr[uint8]
	case reflect.Uint16:
		em.listEnc = encListVarIntPtr[uint16]
	case reflect.Uint32:
		em.listEnc = encListVarIntPtr[uint32]
	case reflect.Uint64:
		em.listEnc = encListVarIntPtr[uint64]

	case reflect.Float32:
		em.listEnc = encListAllPtr
	case reflect.Float64:
		em.listEnc = encListAllPtr
	case reflect.String:
		em.listEnc = encListAllPtr
	case reflect.Bool:
		em.listEnc = encListAllPtr

	case reflect.Pointer:
		em.listEnc = encListAllPtr
	case reflect.Interface:
		em.listEnc = encListAllPtr
	case reflect.Map, reflect.Array:
		em.listEnc = encListAllPtr
	case reflect.Slice:
		em.listEnc = encListAllPtr
	case reflect.Struct:
		if em.itemType == cst.TypeTime {
			em.listEnc = encListAllPtr
		} else {
			em.listEnc = encListStructPtr
		}
	}
	return
}
