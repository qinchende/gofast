// Copyright 2024 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package cdo

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/aid/iox"
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/store/dts"
	"io"
	"reflect"
	"runtime/debug"
	"unsafe"
)

type (
	decValFunc    func(d *decoder)
	decKVPairFunc func(d *decoder, key string)
	decListFunc   func(d *decoder, tLen int)

	decMeta struct {
		// map & struct
		kvPairDec decKVPairFunc

		// struct
		ss        *dts.StructSchema
		fieldsDec []decValFunc

		// array & slice
		listType    reflect.Type
		itemType    reflect.Type
		itemKind    reflect.Kind
		itemDec     decValFunc  // include baseType
		itemMemSize int         // item类型对应的内存字节大小
		arrLen      int         // 如果是数组，记录其长度
		listDec     decListFunc // List整体解码

		// status
		isSuperKV bool // {} SuperKV
		isMap     bool // {} map
		isWebKV   bool // {} WebKV
		isStruct  bool // {} struct
		isList    bool // [] array & slice
		isArray   bool // [] array
		isSlice   bool // [] slice
		isAny     bool // [] is list and item is interface type in the final

		isUnsafe bool  // 分配内存容易出现安全漏洞
		isPtr    bool  // (curr-val | list-item-val | map-value) is ptr
		ptrLevel uint8 // ptr deep (max 256)
	}

	decoder struct {
		sub *decoder // 共享的subDecode，用来解析子对象

		mp     *cst.KV        // KV
		wk     *cst.WebKV     // WebKV
		dm     *decMeta       // Struct | Slice,Array
		dstPtr unsafe.Pointer // 目标对象的内存地址
		slice  rt.SliceHeader // 指向一个切片类型
		//bOpts  *dts.BindOptions // 绑定控制

		str  string // 数据源
		scan int    // 当前扫描定位

		fIdx   int        // field index
		fIdxes [256]int16 // 不多于 256 个字段，暂不支持更多字段

		skipValue   bool // 跳过当前要解析的值
		isNeedValid bool // 在绑定到对象时，是否需要验证字段
	}
)

// 默认值，用于缓存对象的重置
var _subDecodeDefValues decoder

func (d *decoder) reset() {
	*d = _subDecodeDefValues
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decodeFromReader(dst any, reader io.Reader, ctSize int64) error {
	// 一次性读取完成，或者遇到EOF标记或者其它错误
	if ctSize > maxCdoStrLen {
		ctSize = maxCdoStrLen
	}
	if bytes, err := iox.ReadAll(reader, ctSize); err == nil {
		return decodeFromString(dst, lang.BTS(bytes))
	} else {
		return err
	}
}

func decodeFromString(dst any, source string) error {
	if len(source) > maxCdoStrLen {
		return errCdoTooLarge
	}
	if sk, ok := dst.(*dts.StructKV); ok {
		return startDecEx(sk.SS.Type, sk.Ptr, source)
	} else {
		return startDec(dst, source)
	}
}

func startDec(v any, source string) error {
	if v == nil {
		return errValueIsNil
	}
	if len(source) == 0 {
		return errEmptyCdoStr
	}
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Pointer {
		return errValueMustPtr
	}
	return startDecEx(typ.Elem(), (*rt.AFace)(unsafe.Pointer(&v)).DataPtr, source)
}

func startDecEx(typ reflect.Type, ptr unsafe.Pointer, source string) error {
	d := cdoDecPool.Get().(*decoder)
	d.str = source
	d.scan = 0
	d.applyDecMeta(typ, ptr)

	innErr := d.run()
	err := d.warpErrorCode(innErr)

	if d.sub != nil {
		d.sub.reset()
		cdoDecPool.Put(d.sub)
		d.sub = nil
	}
	d.reset() // 此时 sd 中指针指向的对象没有被释放，要先释放再回收
	cdoDecPool.Put(d)
	return err
}

func (d *decoder) warpErrorCode(errCode errType) error {
	if errCode >= 0 {
		return nil
	}

	sta := d.scan
	end := sta + 20 // 输出标记后面 n 个字符
	if end > len(d.str) {
		end = len(d.str)
	}

	errMsg := fmt.Sprintf("Cdo: %s, pos %d, character %q near ( %s )", errDescription[-errCode], sta, d.str[sta], d.str[sta:end])
	return errors.New(errMsg)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (d *decoder) run() (err errType) {
	// 解析过程中异常，这里统一截获处理，返回解析错误编号
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				err = code
			} else if stdErr, yes := pic.(error); yes {
				fmt.Println(stdErr)
				err = errCdoChar
			} else {
				fmt.Printf("%s\n%s", pic, debug.Stack()) // 调试的时候打印错误信息
				err = errCdoChar
			}
		}
	}()

	switch {
	default:
		d.dm.itemDec(d)
	case d.dm.isList:
		d.decList()
	case d.dm.isStruct, d.dm.isSuperKV:
		d.scanKVS()
	}

	if d.scan != len(d.str) {
		err = errEOF // 数据源和目标对象需要完全匹配
	}
	return
}

func (d *decoder) runSub(typ reflect.Type, ptr unsafe.Pointer) {
	dSub := d.sub

	if dSub == nil {
		dSub = cdoDecPool.Get().(*decoder)
	} else {
		dSub.reset()
	}
	dSub.str = d.str
	dSub.scan = d.scan
	dSub.applyDecMeta(typ, ptr)

	switch {
	default:
		dSub.dm.itemDec(dSub)
	case dSub.dm.isList:
		dSub.decList()
	case dSub.dm.isStruct, dSub.dm.isSuperKV:
		dSub.scanKVS()
	}

	d.scan = dSub.scan
	d.resetSub()
}

func (d *decoder) initSub(ptr unsafe.Pointer) {
	if d.sub == nil {
		d.sub = cdoDecPool.Get().(*decoder)
		d.sub.str = d.str
		d.sub.scan = d.scan
		d.sub.applyDecMeta(d.dm.itemType, ptr)
		return
	}

	d.sub.scan = d.scan
	if d.sub.dm.isSuperKV {
		if d.dm.isWebKV {
			d.wk = (*cst.WebKV)(ptr)
		} else {
			d.sub.mp = (*cst.KV)(ptr)
		}
	} else {
		d.sub.dstPtr = ptr
	}
}

func (d *decoder) resetSub() {
	if d.sub.sub != nil {
		d.sub.sub.reset()
		cdoDecPool.Put(d.sub.sub)
		d.sub.sub = nil
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// typ 是 剥离 Pointer 之后的最终类型
func (d *decoder) applyDecMeta(typ reflect.Type, ptr unsafe.Pointer) {
	if meta := cacheGetDecMeta(typ); meta != nil {
		d.dm = meta
	} else {
		d.dm = newDecMeta(typ)
		cacheSetDecMeta(typ, d.dm)
	}

	if d.dm.isSuperKV {
		if d.dm.isWebKV {
			d.wk = (*cst.WebKV)(ptr)
		} else {
			d.mp = (*cst.KV)(ptr)
		}
	} else {
		d.dstPtr = ptr
	}
	return
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func newDecMeta(typ reflect.Type) (dm *decMeta) {
	dm = &decMeta{}

	switch kd := typ.Kind(); kd {
	default:
		dm.initBaseTypeMeta(typ)
	case reflect.Pointer:
		dm.initPointerMeta(typ)
	case reflect.Struct:
		if typ == cst.TypeTime {
			dm.initBaseTypeMeta(typ)
		}
		dm.initStructMeta(typ)
	case reflect.Map:
		if typ == cst.TypeCstKV || typ == cst.TypeStrAnyMap {
			dm.isSuperKV = true
			dm.isMap = true
			dm.bindCstKVDec()
			return
		}
		if typ == cst.TypeWebKV {
			dm.isSuperKV = true
			dm.isWebKV = true
			dm.bindWebKVDec()
			return
		}
		panic(errValueType)
	case reflect.Array, reflect.Slice:
		dm.initListMeta(typ)
	}
	return
}

func (dm *decMeta) peelPtr(typ reflect.Type) {
	dm.itemType = typ.Elem()
	dm.itemKind = dm.itemType.Kind()
	dm.itemMemSize = int(dm.itemType.Size())

ptrLevel:
	if dm.itemKind == reflect.Pointer {
		dm.itemType = dm.itemType.Elem()
		dm.itemKind = dm.itemType.Kind()

		dm.isPtr = true
		dm.ptrLevel++
		// NOTE：指针嵌套不能超过3层，这种很少见，也就是说此解码方案并不通用
		if dm.ptrLevel > 3 {
			panic(errPtrLevel)
		}
		goto ptrLevel
	}
}

func (dm *decMeta) initPointerMeta(typ reflect.Type) {
	dm.isPtr = true
	dm.ptrLevel++
	dm.peelPtr(typ)
	dm.itemDec = scanPointerValue
}

func (dm *decMeta) initStructMeta(typ reflect.Type) {
	dm.isStruct = true
	dm.ss = dts.SchemaAsReqByType(typ)
	if dm.isPtr && dm.ss.HasPtrField {
		dm.isUnsafe = true
	}
	dm.itemMemSize = int(typ.Size())
	dm.bindFieldsDec()
}

func (dm *decMeta) initBaseTypeMeta(typ reflect.Type) {
	dm.itemType = typ
	dm.itemKind = dm.itemType.Kind()
	dm.itemDec = bindBaseValDec
}

func (dm *decMeta) initListMeta(typ reflect.Type) {
	dm.isList = true
	dm.listType = typ
	dm.peelPtr(typ)

	if dm.itemKind == reflect.Interface {
		dm.isAny = true
	}

	// 不是数组就是切片
	if typ.Kind() == reflect.Array {
		dm.isArray = true
		dm.arrLen = typ.Len() // 数组长度
		if dm.isPtr {
			return
		}
		dm.itemMemSize = int(dm.itemType.Size())
	} else {
		dm.isSlice = true
	}

	// List 项如果是 struct ，是本编解码方案重点处理的情况
	if dm.itemKind == reflect.Struct && dm.itemType != cst.TypeTime {
		dm.ss = dts.SchemaAsReqByType(dm.itemType)
		if dm.isPtr && dm.ss.HasPtrField {
			dm.isUnsafe = true
		}
		dm.bindFieldsDec()
	}
	dm.bindListDec()
}

// JSON数据主要分 object {} + list [] 两种类型
// map & gson & struct 都需要解析 {} 他们都是 kvPair 形式的数据
// array & slice 需要解析 [] ，他们都是 List 形式的数据
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (dm *decMeta) bindCstKVDec() {
	dm.kvPairDec = scanCstKVValue
}

func (dm *decMeta) bindWebKVDec() {
	dm.kvPairDec = scanWebKVValue
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func bindBaseValDec(d *decoder) {
	// NOTE：只能是数值类型
	switch d.dm.itemKind {
	case reflect.Int:

	case reflect.Int8:

	case reflect.Int16:

	case reflect.Int32:

	case reflect.Int64:

	case reflect.Uint:

	case reflect.Uint8:

	case reflect.Uint16:

	case reflect.Uint32:

	case reflect.Uint64:

	case reflect.Float32:
		v := scanF32Val(d.str[d.scan:])
		d.scan += 4
		bindF32(d.dstPtr, v)
	case reflect.Float64:
		v := scanF64Val(d.str[d.scan:])
		d.scan += 8
		bindF64(d.dstPtr, v)
	default:
		panic(errValueType)
	}
}

func (dm *decMeta) bindFieldsDec() {
	dm.kvPairDec = decField

	fLen := len(dm.ss.FieldsAttr)
	dm.fieldsDec = make([]decValFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}

	// 非指针字段
	if dm.ss.FieldsAttr[i].PtrLevel == 0 {
		switch fa := dm.ss.FieldsAttr[i]; fa.Kind {
		default:
			panic(errValueType)

		case reflect.Int:
			dm.fieldsDec[i] = func(d *decoder) { bindInt(decVarIntField(d)) }
		case reflect.Int8:
			dm.fieldsDec[i] = func(d *decoder) { bindInt8(decVarIntField(d)) }
		case reflect.Int16:
			dm.fieldsDec[i] = func(d *decoder) { bindInt16(decVarIntField(d)) }
		case reflect.Int32:
			dm.fieldsDec[i] = func(d *decoder) { bindInt32(decVarIntField(d)) }
		case reflect.Int64:
			dm.fieldsDec[i] = func(d *decoder) { bindInt64(decVarIntField(d)) }
		case reflect.Uint:
			dm.fieldsDec[i] = func(d *decoder) { bindUint(decVarIntField(d)) }
		case reflect.Uint8:
			dm.fieldsDec[i] = func(d *decoder) { bindUint8(decVarIntField(d)) }
		case reflect.Uint16:
			dm.fieldsDec[i] = func(d *decoder) { bindUint16(decVarIntField(d)) }
		case reflect.Uint32:
			dm.fieldsDec[i] = func(d *decoder) { bindUint32(decVarIntField(d)) }
		case reflect.Uint64:
			dm.fieldsDec[i] = func(d *decoder) { bindUint64(decVarIntField(d)) }

		case reflect.Float32:
			dm.fieldsDec[i] = decF32Field
		case reflect.Float64:
			dm.fieldsDec[i] = decF64Field
		case reflect.String:
			dm.fieldsDec[i] = decStrField
		case reflect.Bool:
			dm.fieldsDec[i] = decBoolField
		case reflect.Interface:
			dm.fieldsDec[i] = decAnyField

		case reflect.Struct:
			if dm.ss.FieldsAttr[i].Type == cst.TypeTime {
				dm.fieldsDec[i] = decTimeField
			} else {
				dm.fieldsDec[i] = decMixField
			}
		case reflect.Slice:
			if fa.Type == cst.TypeBytes {
				dm.fieldsDec[i] = decBytesField
			} else {
				dm.fieldsDec[i] = decMixField
			}
		case reflect.Map, reflect.Array:
			dm.fieldsDec[i] = decMixField
		}
		goto nextField
	}

	// 指针字段
	switch fa := dm.ss.FieldsAttr[i]; fa.Kind {
	default:
		panic(errValueType)

	case reflect.Int:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindInt)
	case reflect.Int8:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindInt8)
	case reflect.Int16:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindInt16)
	case reflect.Int32:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindInt32)
	case reflect.Int64:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindInt64)
	case reflect.Uint:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindUint)
	case reflect.Uint8:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindUint8)
	case reflect.Uint16:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindUint16)
	case reflect.Uint32:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindUint32)
	case reflect.Uint64:
		dm.fieldsDec[i] = decVarIntFieldPtr(bindUint64)

	case reflect.Float32:
		dm.fieldsDec[i] = decF32FieldPtr
	case reflect.Float64:
		dm.fieldsDec[i] = decF64FieldPtr
	case reflect.String:
		dm.fieldsDec[i] = decStrFieldPtr
	case reflect.Bool:
		dm.fieldsDec[i] = decBoolFieldPtr
	case reflect.Interface:
		dm.fieldsDec[i] = decAnyFieldPtr

	case reflect.Struct:
		if dm.ss.FieldsAttr[i].Type == cst.TypeTime {
			dm.fieldsDec[i] = decTimeFieldPtr
		} else {
			dm.fieldsDec[i] = decMixFieldPtr
		}
	case reflect.Slice:
		if fa.Type == cst.TypeBytes {
			dm.fieldsDec[i] = decBytesFieldPtr
		} else {
			dm.fieldsDec[i] = decMixFieldPtr
		}
	case reflect.Map, reflect.Array:
		dm.fieldsDec[i] = decMixFieldPtr
	}
	goto nextField
}

// 特定类型 List +++++++++++++++++++++++++++++
func (dm *decMeta) bindListDec() {
	if !dm.isPtr {
		switch dm.itemKind {
		default:
			panic(errValueType)

		case reflect.Int:
			dm.listDec = decListInt
		case reflect.Int8:
			dm.listDec = decListInt8
		case reflect.Int16:
			dm.listDec = decListInt16
		case reflect.Int32:
			dm.listDec = decListInt32
		case reflect.Int64:
			dm.listDec = decListInt64
		case reflect.Uint:
			dm.listDec = decListUint
		case reflect.Uint8:
			dm.listDec = decListUint8
		case reflect.Uint16:
			dm.listDec = decListUint16
		case reflect.Uint32:
			dm.listDec = decListUint32
		case reflect.Uint64:
			dm.listDec = decListUint64

		case reflect.Float32:
			dm.listDec = decListF32
		case reflect.Float64:
			dm.listDec = decListF64

		case reflect.String:
			dm.listDec = decListStr
		case reflect.Bool:
			dm.listDec = decListBool
		case reflect.Struct:
			if dm.itemType == cst.TypeTime {
				dm.listDec = decListTime
			} else {
				dm.listDec = decListStruct
			}

		case reflect.Pointer:

		case reflect.Interface:

		case reflect.Map, reflect.Array, reflect.Slice:
			dm.listDec = decListAll
			dm.itemDec = decListMixItem
		}
		return
	}

	// 数据项 为 指针类型
	switch dm.itemKind {
	default:
		panic(errValueType)

	case reflect.Int:
		dm.listDec = decListIntPtr
	case reflect.Int8:
		dm.listDec = decListInt8Ptr
	case reflect.Int16:
		dm.listDec = decListInt16Ptr
	case reflect.Int32:
		dm.listDec = decListInt32Ptr
	case reflect.Int64:
		dm.listDec = decListInt64Ptr
	case reflect.Uint:
		dm.listDec = decListUintPtr
	case reflect.Uint8:
		dm.listDec = decListUint8Ptr
	case reflect.Uint16:
		dm.listDec = decListUint16Ptr
	case reflect.Uint32:
		dm.listDec = decListUint32Ptr
	case reflect.Uint64:
		dm.listDec = decListUint64Ptr

	case reflect.Float32:
		dm.listDec = decListF32Ptr
	case reflect.Float64:
		dm.listDec = decListF64Ptr

	case reflect.String:
		dm.listDec = decListStrPtr
	case reflect.Bool:
		dm.listDec = decListBoolPtr
	case reflect.Struct:
		if dm.itemType == cst.TypeTime {
			dm.listDec = decListTimePtr
		} else {
			dm.listDec = decListStruct
		}

	case reflect.Pointer:

	case reflect.Interface:

	case reflect.Map, reflect.Array, reflect.Slice:
		dm.listDec = decListAll
		dm.itemDec = decListMixItem
	}
}
