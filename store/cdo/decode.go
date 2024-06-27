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
	decValFunc    func(d *subDecode)
	decKVPairFunc func(d *subDecode, key string)
	decListFunc   func(d *subDecode, tLen int)

	decMeta struct {
		// map & struct
		kvPairDec decKVPairFunc

		// Struct
		ss        *dts.StructSchema
		fieldsDec []decValFunc

		// array & slice
		itemType    reflect.Type
		itemKind    reflect.Kind
		itemDec     decValFunc  // include baseType
		itemMemSize int         // 数组属性，item类型对应的内存字节大小
		arrLen      int         // 数组属性，数组长度
		listDec     decListFunc // List整体解码

		// status
		isSuperKV bool // {} SuperKV
		isMap     bool // {} map
		isWebKV   bool // {} WebKV
		isStruct  bool // {} struct
		isList    bool // [] array & slice
		isArray   bool // [] array
		isSlice   bool // [] slice
		isArrBind bool // [] is array and item not pointer type

		isAny    bool  // [] is list and item is interface type in the final
		isPtr    bool  // [] is list and item is pointer type
		ptrLevel uint8 // [] is list and item pointer level
	}

	subDecode struct {
		share *subDecode // 共享的subDecode，用来解析子对象

		mp     *cst.KV        // KV
		wk     *cst.WebKV     // WebKV
		dm     *decMeta       // Struct | Slice,Array
		dstPtr unsafe.Pointer // 目标对象的内存地址
		slice  rt.SliceHeader // 指向一个切片类型
		//bOpts  *dts.BindOptions // 绑定控制

		str  string // 数据源
		scan int    // 当前扫描定位

		// 辅助变量
		arrIdx int        // list解析的数量
		fIdx   int        // key index
		fIdxes [256]int16 // 不多于 256 个字段，暂不支持更多字段

		skipValue   bool // 跳过当前要解析的值
		isNeedValid bool // 在绑定到对象时，是否需要验证字段
	}
)

// 默认值，用于缓存对象的重置
var _subDecodeDefValues subDecode

func (d *subDecode) reset() {
	*d = _subDecodeDefValues
}

// private enter
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decodeFromReader(dst any, reader io.Reader, ctSize int64) error {
	// 一次性读取完成，或者遇到EOF标记或者其它错误
	if ctSize > maxCdoStrLen {
		ctSize = maxCdoStrLen
	}
	bytes, err1 := iox.ReadAll(reader, ctSize)
	if err1 != nil {
		return err1
	}
	return decodeFromString(dst, lang.BTS(bytes))
}

func decodeFromString(dst any, source string) (err error) {
	if len(source) > maxCdoStrLen {
		return errCdoTooLarge
	}
	if sk, ok := dst.(*dts.StructKV); ok {
		return startDecodeInner(sk.SS.Attrs.Type, sk.Ptr, source)
	} else {
		return startDecode(dst, source)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主解析入口
func startDecode(dst any, source string) (err error) {
	if dst == nil {
		return errValueIsNil
	}
	if len(source) == 0 {
		return errEmptyCdoStr
	}
	rfType := reflect.TypeOf(dst)
	if rfType.Kind() != reflect.Pointer {
		return errValueMustPtr
	}
	return startDecodeInner(rfType.Elem(), (*rt.AFace)(unsafe.Pointer(&dst)).DataPtr, source)
}

func startDecodeInner(typ reflect.Type, ptr unsafe.Pointer, source string) (err error) {
	d := cdoDecPool.Get().(*subDecode)
	d.str = source
	d.scan = 0
	d.getDecMeta(typ, ptr)

	innErr := d.startDecode()
	err = d.warpErrorCode(innErr)

	if d.share != nil {
		d.share.reset()
		cdoDecPool.Put(d.share)
		d.share = nil
	}
	// 此时 sd 中指针指向的对象没有被释放，存在一定风险，所以要先释放再回收
	d.reset()
	cdoDecPool.Put(d)
	return
}

func (d *subDecode) warpErrorCode(errCode errType) error {
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

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (d *subDecode) startDecode() (err errType) {
	// 解析过程中异常，这里统一截获处理，返回解析错误编号
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				err = code
			} else if stdErr, yes := pic.(error); yes {
				fmt.Println(stdErr)
				err = errCdo
			} else {
				fmt.Printf("%s\n%s", pic, debug.Stack()) // 调试的时候打印错误信息
				err = errCdo
			}
		}
	}()

	switch {
	default:
		d.dm.itemDec(d)
	case d.dm.isList:
		d.scanList()
	case d.dm.isStruct, d.dm.isSuperKV:
		d.scanKVS()
	}

	if d.scan != len(d.str) {
		err = errCdo // 不能再有未解析内容了
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (d *subDecode) initShareDecode(ptr unsafe.Pointer) {
	if d.share == nil {
		d.share = cdoDecPool.Get().(*subDecode)
		d.share.str = d.str
		d.share.scan = d.scan
		d.share.getDecMeta(d.dm.itemType, ptr)
		return
	}

	d.share.scan = d.scan
	if d.share.dm.isSuperKV {
		if d.dm.isWebKV {
			d.wk = (*cst.WebKV)(ptr)
		} else {
			d.share.mp = (*cst.KV)(ptr)
		}
	} else {
		d.share.dstPtr = ptr // 当前值的地址
	}
}

// 包含有子subDecode时，就递归调用
func (d *subDecode) startSubDecode(rfType reflect.Type, ptr unsafe.Pointer) {
	dSub := d.share

	if dSub == nil {
		dSub = cdoDecPool.Get().(*subDecode)
	} else {
		dSub.reset()
	}
	dSub.str = d.str
	dSub.scan = d.scan
	dSub.getDecMeta(rfType, ptr)

	switch {
	default:
		dSub.dm.itemDec(dSub)
	case dSub.dm.isList:
		dSub.scanList()
	case dSub.dm.isStruct, dSub.dm.isSuperKV:
		dSub.scanKVS()
	}

	d.scan = dSub.scan
	d.resetShareDecode()
}

func (d *subDecode) resetShareDecode() {
	if d.share.share != nil {
		d.share.share.reset()
		cdoDecPool.Put(d.share.share)
		d.share.share = nil
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// rfType 是 剥离 Pointer 之后的最终类型
func (d *subDecode) getDecMeta(rfType reflect.Type, ptr unsafe.Pointer) {
	if meta := cacheGetDecMeta(rfType); meta != nil {
		d.dm = meta
	} else {
		d.dm = newDecodeMeta(rfType)
		cacheSetDecMeta(rfType, d.dm)
	}

	if d.dm.isSuperKV {
		if d.dm.isWebKV {
			d.wk = (*cst.WebKV)(ptr)
		} else {
			d.mp = (*cst.KV)(ptr)
		}
	} else {
		d.dstPtr = ptr // 当前值的地址
	}
	return
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func newDecodeMeta(rfType reflect.Type) (dm *decMeta) {
	dm = &decMeta{}

	switch kd := rfType.Kind(); kd {
	default:
		dm.initBaseValueMeta(rfType)
	case reflect.Pointer:
		dm.initPointerMeta(rfType)
	case reflect.Struct:
		if rfType == cst.TypeTime {
			panic(errValueType)
		}
		dm.initStructMeta(rfType)
	case reflect.Map:
		if rfType == cst.TypeCstKV || rfType == cst.TypeStrAnyMap {
			dm.isSuperKV = true
			dm.isMap = true
			dm.bindCstKVDec()
			return
		}
		if rfType == cst.TypeWebKV {
			dm.isSuperKV = true
			dm.isWebKV = true
			dm.bindWebKVDec()
			return
		}
		panic(errValueType)
	case reflect.Array, reflect.Slice:
		dm.initListMeta(rfType)

	}
	return
}

func (dm *decMeta) peelPtr(rfType reflect.Type) {
	dm.itemType = rfType.Elem()
	dm.itemKind = dm.itemType.Kind()
	dm.itemMemSize = int(dm.itemType.Size())

peelLoop:
	if dm.itemKind == reflect.Pointer {
		dm.itemType = dm.itemType.Elem()
		dm.itemKind = dm.itemType.Kind()

		dm.isPtr = true
		dm.ptrLevel++
		// NOTE：指针嵌套不能超过3层，这种很少见，也就是说此解码方案并不通用
		if dm.ptrLevel > 3 {
			panic(errPtrLevel)
		}
		goto peelLoop
	}
}

func (dm *decMeta) initPointerMeta(rfType reflect.Type) {
	dm.isPtr = true
	dm.ptrLevel++
	dm.peelPtr(rfType)

	dm.itemDec = scanPointerValue
}

func (dm *decMeta) initStructMeta(rfType reflect.Type) {
	dm.isStruct = true
	dm.ss = dts.SchemaAsReqByType(rfType)
	dm.itemMemSize = int(rfType.Size())

	dm.bindStructDec()
}

func (dm *decMeta) initBaseValueMeta(rfType reflect.Type) {
	dm.itemType = rfType
	dm.itemKind = dm.itemType.Kind()
	dm.itemDec = scanJustBaseValue
}

func (dm *decMeta) initListMeta(rfType reflect.Type) {
	dm.isList = true
	dm.peelPtr(rfType)

	// 是否是interface类型
	if dm.itemKind == reflect.Interface {
		dm.isAny = true
	}

	// 进一步初始化数组
	if rfType.Kind() == reflect.Array {
		dm.isArray = true
		dm.arrLen = rfType.Len() // 数组长度
		if dm.isPtr {
			return
		}
		dm.itemMemSize = int(dm.itemType.Size())
		dm.isArrBind = true
	} else {
		dm.isSlice = true
	}

	// List 项如果是 struct ，是本编解码方案重点处理的情况
	if dm.itemKind == reflect.Struct && dm.itemType != cst.TypeTime {
		//dm.initListStructMeta()
		//dm.isStruct = true
		dm.ss = dts.SchemaAsReqByType(dm.itemType)
		dm.bindStructDec()
	} else {
		dm.bindListItemDec()
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

// Dest is just a base type value
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanJustBaseValue(d *subDecode) {
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

func (dm *decMeta) bindStructDec() {
	dm.kvPairDec = scanStructValue

	fLen := len(dm.ss.FieldsAttr)
	dm.fieldsDec = make([]decValFunc, fLen)

	i := -1
nextField:
	i++
	if i >= fLen {
		return
	}

	// 字段不是指针类型
	if dm.ss.FieldsAttr[i].PtrLevel == 0 {
		switch fa := dm.ss.FieldsAttr[i]; fa.Kind {
		case reflect.Int:
			dm.fieldsDec[i] = scanObjIntValue
		case reflect.Int8:
			dm.fieldsDec[i] = scanObjInt8Value
		case reflect.Int16:
			dm.fieldsDec[i] = scanObjInt16Value
		case reflect.Int32:
			dm.fieldsDec[i] = scanObjInt32Value
		case reflect.Int64:
			dm.fieldsDec[i] = scanObjInt64Value
		case reflect.Uint:
			dm.fieldsDec[i] = scanObjUintValue
		case reflect.Uint8:
			dm.fieldsDec[i] = scanObjUint8Value
		case reflect.Uint16:
			dm.fieldsDec[i] = scanObjUint16Value
		case reflect.Uint32:
			dm.fieldsDec[i] = scanObjUint32Value
		case reflect.Uint64:
			dm.fieldsDec[i] = scanObjUint64Value
		case reflect.Float32:
			dm.fieldsDec[i] = scanObjF32Value
		case reflect.Float64:
			dm.fieldsDec[i] = scanObjF64Value
		case reflect.String:
			dm.fieldsDec[i] = scanObjStrValue
		case reflect.Bool:
			dm.fieldsDec[i] = scanObjBoolValue
		case reflect.Interface:
			dm.fieldsDec[i] = scanObjAnyValue

		case reflect.Struct:
			if dm.ss.FieldsAttr[i].Type == cst.TypeTime {
				dm.fieldsDec[i] = scanObjTimeValue
			} else {
				dm.fieldsDec[i] = scanObjMixValue
			}
		case reflect.Slice:
			// 分情况，如果是字节切片，单独处理
			if fa.Type == cst.TypeBytes {
				// TODO: 字节切片的解析，把字符串当做base64编码看待
				dm.fieldsDec[i] = scanObjBytesValue
			} else {
				dm.fieldsDec[i] = scanObjMixValue
			}
		case reflect.Map, reflect.Array:
			dm.fieldsDec[i] = scanObjMixValue

		default:
			panic(errValueType)
		}
		goto nextField
	}

	// 字段是指针类型，我们需要判断的是真实的数据类型
	switch fa := dm.ss.FieldsAttr[i]; fa.Kind {
	case reflect.Int:
		dm.fieldsDec[i] = scanObjPtrIntValue
	case reflect.Int8:
		dm.fieldsDec[i] = scanObjPtrInt8Value
	case reflect.Int16:
		dm.fieldsDec[i] = scanObjPtrInt16Value
	case reflect.Int32:
		dm.fieldsDec[i] = scanObjPtrInt32Value
	case reflect.Int64:
		dm.fieldsDec[i] = scanObjPtrInt64Value
	case reflect.Uint:
		dm.fieldsDec[i] = scanObjPtrUintValue
	case reflect.Uint8:
		dm.fieldsDec[i] = scanObjPtrUint8Value
	case reflect.Uint16:
		dm.fieldsDec[i] = scanObjPtrUint16Value
	case reflect.Uint32:
		dm.fieldsDec[i] = scanObjPtrUint32Value
	case reflect.Uint64:
		dm.fieldsDec[i] = scanObjPtrUint64Value
	case reflect.Float32:
		dm.fieldsDec[i] = scanObjPtrF32Value
	case reflect.Float64:
		dm.fieldsDec[i] = scanObjPtrF64Value
	case reflect.String:
		dm.fieldsDec[i] = scanObjPtrStrValue
	case reflect.Bool:
		dm.fieldsDec[i] = scanObjPtrBoolValue
	case reflect.Interface:
		dm.fieldsDec[i] = scanObjPtrAnyValue
	case reflect.Struct:
		if dm.ss.FieldsAttr[i].Type == cst.TypeTime {
			dm.fieldsDec[i] = scanObjPtrTimeValue
		} else {
			dm.fieldsDec[i] = scanObjPtrMixValue
		}
	case reflect.Slice:
		// 分情况，如果是字节切片，单独处理
		if fa.Type == cst.TypeBytes {
			// TODO: 字节切片的解析，把字符串当做base64编码看待
			dm.fieldsDec[i] = scanObjPtrStrValue
		} else {
			dm.fieldsDec[i] = scanObjPtrMixValue
		}
	case reflect.Map, reflect.Array:
		dm.fieldsDec[i] = scanObjPtrMixValue
	default:
		panic(errValueType)
	}
	goto nextField
}

// List
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (dm *decMeta) bindListItemDec() {
	if dm.isArrBind {
		switch dm.itemKind {
		default:
			return
		case reflect.Interface:
			dm.itemDec = scanArrAnyValue
		case reflect.Struct:
			// Note: 有个特殊情况，当处理GsonRows解析时候，要特殊处理
			if dm.itemType == cst.TypeTime {
				dm.itemDec = scanArrTimeValue
			} else {
				dm.itemDec = scanArrMixValue
			}
		case reflect.Map, reflect.Array, reflect.Slice:
			dm.itemDec = scanArrMixValue
		}
	}

	switch dm.itemKind {
	default:
		return
	case reflect.Interface:
		dm.itemDec = scanListAnyValue
	case reflect.Struct:
		if dm.isArray {
			dm.itemDec = scanArrPtrMixValue
		} else {
			dm.itemDec = scanListMixValue // 这里只能是Slice
		}
		dm.isArrBind = true // Note：这些情况，无需用到缓冲池
	case reflect.Map, reflect.Array, reflect.Slice:
		if dm.isArray {
			dm.itemDec = scanArrPtrMixValue
		} else {
			dm.itemDec = scanListMixValue // 这里只能是Slice
		}
		dm.isArrBind = true // Note：这些情况，无需用到缓冲池
	}
}

func (dm *decMeta) bindListDec() {
	if !dm.isPtr {
		switch dm.itemKind {
		default:
			panic(errValueType)

		case reflect.Int:
			dm.listDec = decIntList
		case reflect.Int8:
			dm.listDec = decInt8List
		case reflect.Int16:
			dm.listDec = decInt16List
		case reflect.Int32:
			dm.listDec = decInt32List
		case reflect.Int64:
			dm.listDec = decInt64List
		case reflect.Uint:
			dm.listDec = decUintList
		case reflect.Uint8:
			dm.listDec = decUint8List
		case reflect.Uint16:
			dm.listDec = decUint16List
		case reflect.Uint32:
			dm.listDec = decUint32List
		case reflect.Uint64:
			dm.listDec = decUint64List

		case reflect.Float32:
			dm.listDec = decF32List
		case reflect.Float64:
			dm.listDec = decF64List

		case reflect.String:
			dm.listDec = decStrList
		case reflect.Bool:
			dm.listDec = decBoolList

		case reflect.Interface:
			dm.listDec = decListBaseType
		case reflect.Struct:
			if dm.itemType == cst.TypeTime {
				dm.listDec = decListBaseType
			} else {
				dm.listDec = decStructList
			}
		case reflect.Map, reflect.Array, reflect.Slice:
			dm.listDec = decListBaseType
		}
		return
	}

	// 数据项 为 指针类型
	switch dm.itemKind {
	default:
		panic(errValueType)

	case reflect.Int:
		dm.listDec = decIntListPtr
	case reflect.Int8:
		dm.listDec = decInt8ListPtr
	case reflect.Int16:
		dm.listDec = decInt16ListPtr
	case reflect.Int32:
		dm.listDec = decInt32ListPtr
	case reflect.Int64:
		dm.listDec = decInt64ListPtr
	case reflect.Uint:
		dm.listDec = decUintListPtr
	case reflect.Uint8:
		dm.listDec = decUint8ListPtr
	case reflect.Uint16:
		dm.listDec = decUint16ListPtr
	case reflect.Uint32:
		dm.listDec = decUint32ListPtr
	case reflect.Uint64:
		dm.listDec = decUint64ListPtr

	case reflect.Float32:
		dm.listDec = decF32ListPtr
	case reflect.Float64:
		dm.listDec = decF64ListPtr

	case reflect.String:
		dm.listDec = decStrListPtr
	case reflect.Bool:
		dm.listDec = decBoolListPtr

	case reflect.Interface:
		dm.listDec = decListBaseType
	case reflect.Struct:
		if dm.itemType == cst.TypeTime {
			dm.listDec = decListBaseType
		} else {
			dm.listDec = decStructList
		}
	case reflect.Map, reflect.Array, reflect.Slice:
		dm.listDec = decListBaseType
	}
}
