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
	//encValFunc  func(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type)
	encValFunc  func(bs []byte, ptr unsafe.Pointer, typ reflect.Type) []byte
	encListFunc func(e *encoder)

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

	encoder struct {
		srcPtr unsafe.Pointer // list or object 对象值地址（其指向的值不能为nil，也不能为指针）
		slice  rt.SliceHeader // 用于将数组模拟成切片
		em     *encMeta       // Struct | Slice,Array
		bf     *[]byte        // 当解析数组时候用到的一系列临时队列
		bs     []byte         // 用来辅助上面的bf指针，防止24个字节的切片对象堆分配
		//objIdx uint16         // TODO：涉及多少种不同的Struct
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

	//se := encoder{}
	se := cdoEncPool.Get().(*encoder)
	se.fetchEncMeta(v)

	se.bf = pool.GetBytes()
	se.run()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)
	pool.FreeBytes(se.bf)

	se.reset()
	cdoEncPool.Put(se)
	return
}

func encMixedItemRet(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	se := cdoEncPool.Get().(*encoder)
	se.applyEncMeta(typ, ptr)

	se.bs = bf
	se.bf = &se.bs
	se.run()

	se.reset()
	cdoEncPool.Put(se)
	return se.bs
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *encoder) fetchEncMeta(v any) {
	vType := reflect.TypeOf(v)
	ptr := (*rt.AFace)(unsafe.Pointer(&v)).DataPtr

	// v可能是一个值的地址，最多只能剥掉一层指针
	if vType.Kind() == reflect.Pointer {
		vType = vType.Elem()
		// Note：有些类型本质其实是指针，但是reflect.Kind() != reflect.Pointer
		// 比如：map | channel | func
		// 此时需要统一变量值 ptr 指向的内存
		if vType.Kind() == reflect.Map {
			ptr = *(*unsafe.Pointer)(ptr)
		}
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

func (se *encoder) run() {
	switch {
	default:
		se.encBasic()
	case se.em.isList:
		se.encList()
	case se.em.isStruct:
		se.encStruct()
	case se.em.isMap:
		se.encMap()
	case se.em.isPtr:
		se.encPointer()
	}
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
		em.bindFieldsEnc()
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
	em.bindFieldsEnc()
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
func (em *encMeta) bindFieldsEnc() {
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
	case reflect.Interface:
		*encFunc = encAnyRet

	//case reflect.Pointer:
	//	*encFunc = encPointer
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
		switch em.itemKind {
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
		case reflect.Interface:
			em.listEnc = encListAll

		//case reflect.Pointer: // 此时分支不可能
		//	em.listEnc = encPointer
		case reflect.Map, reflect.Array:
			em.listEnc = encListAll
		case reflect.Slice:
			em.listEnc = encListAll
		case reflect.Struct:
			if em.itemType == cst.TypeTime {
				em.listEnc = encListAll
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
	switch em.itemKind {
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
			em.listEnc = encListStructPtr
		}
	}
	return
}
