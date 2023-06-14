package jde

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
	"unsafe"
)

// Note: 此函数只适合 object 的 field，List 的 item 为 指针类型 的情形。非指针不能调用此方法
func getPtrValueAddr(ptr unsafe.Pointer, ptrLevel uint8, kd reflect.Kind, rfType reflect.Type) unsafe.Pointer {
	for ptrLevel > 1 {
		if *(*unsafe.Pointer)(ptr) == nil {
			tpPtr := unsafe.Pointer(new(unsafe.Pointer))
			*(*unsafe.Pointer)(ptr) = tpPtr
			ptr = tpPtr
		} else {
			ptr = *(*unsafe.Pointer)(ptr)
		}

		ptrLevel--
	}

	if *(*unsafe.Pointer)(ptr) == nil {
		var newPtr unsafe.Pointer

		switch kd {
		case reflect.Map:
			newPtr = unsafe.Pointer(new(unsafe.Pointer))
			*(*unsafe.Pointer)(newPtr) = reflect.MakeMap(rfType).UnsafePointer()
		case reflect.Slice:
			newPtr = unsafe.Pointer(&reflect.SliceHeader{})
			*(*unsafe.Pointer)(newPtr) = reflect.MakeSlice(rfType, 0, 0).UnsafePointer()
		default:
			newPtr = reflect.New(rfType).UnsafePointer()
		}

		*(*unsafe.Pointer)(ptr) = newPtr
		return newPtr
	}
	return *(*unsafe.Pointer)(ptr)
}

// array & slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func arrItemPtr(sd *subDecode) unsafe.Pointer {
	return unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(sd.arrIdx*sd.dm.itemRawSize))
}

func arrMixItemPtr(sd *subDecode) unsafe.Pointer {
	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(sd.arrIdx*sd.dm.itemRawSize))

	// 只有field字段为map或者slice的时候，值才可能是nil
	if sd.dm.itemKind == reflect.Map {
		if *(*unsafe.Pointer)(ptr) == nil {
			*(*unsafe.Pointer)(ptr) = reflect.MakeMap(sd.dm.itemType).UnsafePointer()
		}
	}
	return ptr
}

func arrMixItemPtrDeep(sd *subDecode) unsafe.Pointer {
	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(sd.arrIdx*sd.dm.itemRawSize))
	return getPtrValueAddr(ptr, sd.dm.ptrLevel, sd.dm.itemKind, sd.dm.itemType)
}

func sliceMixItemPtr(sd *subDecode, ptr unsafe.Pointer) unsafe.Pointer {
	return getPtrValueAddr(ptr, sd.dm.ptrLevel, sd.dm.itemKind, sd.dm.itemType)
}

// reset array left item
func (sd *subDecode) resetArrLeftItems() {
	var dfValue unsafe.Pointer
	if !sd.dm.isPtr {
		dfValue = zeroValues[sd.dm.itemKind]
	}
	for i := sd.arrIdx; i < sd.dm.arrLen; i++ {
		*(*unsafe.Pointer)(unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(i*sd.dm.itemRawSize))) = dfValue
	}
}

// struct & map & gson
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func fieldPtr(sd *subDecode) unsafe.Pointer {
	return unsafe.Pointer(uintptr(sd.dstPtr) + sd.dm.ss.FieldsAttr[sd.keyIdx].Offset)
}

func fieldMixPtr(sd *subDecode) unsafe.Pointer {
	fa := &sd.dm.ss.FieldsAttr[sd.keyIdx]
	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + fa.Offset)

	if fa.Kind == reflect.Map {
		if *(*unsafe.Pointer)(ptr) == nil {
			*(*unsafe.Pointer)(ptr) = reflect.MakeMap(fa.Type).UnsafePointer()
		}
	}

	//// 只有field字段为map或者slice的时候，值才可能是nil
	//if *(*unsafe.Pointer)(ptr) == nil {
	//	switch fa.Kind {
	//	// Note: 当 array & slice & struct 的时候，相当于是值类型，直接返回首地址即可
	//	//default:
	//	//	panic(errSupport)
	//	case reflect.Map:
	//		*(*unsafe.Pointer)(ptr) = reflect.MakeMap(fa.Type).UnsafePointer()
	//		//case reflect.Slice:
	//		// Note: fa.Kind == reflect.Slice，
	//		// 此时可能申请slice对象没有意义，因为解析程序会自己创建临时空间，完成之后替换旧内存
	//		// 但如果slice中的项还是 mix 类型，可能又不一样了，这种情况解析程序不会申请临时空间
	//		//newPtr := reflect.MakeSlice(fa.Type, 0, 4).UnsafePointer()	// 默认给4个值的空间，避免扩容
	//		//*(*unsafe.Pointer)(ptr) = *(*unsafe.Pointer)(newPtr)
	//	}
	//}
	return ptr
}

func fieldPtrDeep(sd *subDecode) unsafe.Pointer {
	fa := &sd.dm.ss.FieldsAttr[sd.keyIdx]
	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + fa.Offset)
	return getPtrValueAddr(ptr, fa.PtrLevel, fa.Kind, fa.Type)
}

func fieldSetNil(sd *subDecode) {
	*(*unsafe.Pointer)(fieldPtr(sd)) = nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Scan Advanced mixed type, such as map | gson | struct
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// map +++++
// 目前只支持 map[string]any，并不支持其它map
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanMapAnyValue(sd *subDecode, k string) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubDecode(rfTypeOfKV, unsafe.Pointer(&newMap))
		sd.mp.Set(k, newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubDecode(rfTypeOfList, unsafe.Pointer(&newList))
		sd.mp.Set(k, newList)
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			sd.mp.Set(k, sd.str[start:sd.unescapeEnd()])
		} else {
			sd.mp.Set(k, sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		if start := sd.scanNumValue(); start > 0 {
			// 可以选项，不解析，直接返回字符串
			sd.mp.Set(k, lang.ParseFloat(sd.str[start:sd.scan]))
		}
	case c == 't':
		sd.skipTrue()
		sd.mp.Set(k, true)
	case c == 'f':
		sd.skipFalse()
		sd.mp.Set(k, false)
	default:
		sd.skipNull()
		sd.mp.Set(k, nil)
	}
}

// GsonRow +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanGsonValue(sd *subDecode, k string) {
	kIdx := 0
	if kIdx = sd.gr.KeyIndex(k); kIdx < 0 {
		sd.skipValue = true
		sd.skipOneValue()
		return
	}

	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubDecode(rfTypeOfKV, unsafe.Pointer(&newMap))
		sd.gr.SetByIndex(kIdx, newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubDecode(rfTypeOfList, unsafe.Pointer(&newList))
		sd.gr.SetByIndex(kIdx, newList)
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			sd.gr.SetStringByIndex(kIdx, sd.str[start:sd.unescapeEnd()])
		} else {
			sd.gr.SetStringByIndex(kIdx, sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		if start := sd.scanNumValue(); start > 0 {
			sd.gr.SetByIndex(kIdx, lang.ParseFloat(sd.str[start:sd.scan]))
		}
	case c == 't':
		sd.skipTrue()
		sd.gr.SetByIndex(kIdx, true)
	case c == 'f':
		sd.skipFalse()
		sd.gr.SetByIndex(kIdx, false)
	default:
		sd.skipNull()
		sd.gr.SetByIndex(kIdx, nil)
	}
}

// struct +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanStructValue(sd *subDecode, key string) {
	// TODO: 此处 sd.keyIdx 可以继续被优化
	if sd.keyIdx = sd.dm.ss.ColumnIndex(key); sd.keyIdx < 0 {
		sd.skipValue = true
		sd.skipOneValue()
	} else {
		sd.dm.fieldsDec[sd.keyIdx](sd) // 根据目标值类型来解析
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Scan Advanced mixed type, such as map | struct | array | slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// sash as map | struct
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjMixValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case '{', '[':
		sd.scanSubDecode(sd.dm.ss.FieldsAttr[sd.keyIdx].Type, fieldMixPtr(sd))
	default:
		sd.skipNull()
	}
}

func scanObjPtrMixValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case '{', '[':
		sd.scanSubDecode(sd.dm.ss.FieldsAttr[sd.keyIdx].Type, fieldPtrDeep(sd))
	default:
		sd.skipNull()
		fieldSetNil(sd)
	}
}

// sash as array | slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// array and item not ptr
func scanArrMixValue(sd *subDecode) {
	// TODO：在这里循环处理
	switch c := sd.str[sd.scan]; {
	case c == '{':
		sd.initShareDecode(arrMixItemPtr(sd))
		sd.share.scanObject()
		sd.scan = sd.share.scan
	case c == '[':
		sd.initShareDecode(arrMixItemPtr(sd))
		sd.share.scanList()
		sd.scan = sd.share.scan
	default:
		sd.skipNull()
	}
}

// array and item is ptr
func scanArrPtrMixValue(sd *subDecode) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		sd.initShareDecode(arrMixItemPtrDeep(sd))
		sd.share.scanObject()
		sd.scan = sd.share.scan
	case c == '[':
		sd.initShareDecode(arrMixItemPtrDeep(sd))
		sd.share.scanList()
		sd.scan = sd.share.scan
	default:
		sd.skipNull()
		fieldSetNil(sd)
	}
}

// slice 中可能是实体对象，也可能是对象指针
func scanListMixValue(sd *subDecode) {
	sh := (*reflect.SliceHeader)(sd.dstPtr)
	if sd.arrIdx == sh.Cap {
		var oldMem = []byte{}
		old := (*reflect.SliceHeader)(unsafe.Pointer(&oldMem))
		old.Data, old.Len, old.Cap = sh.Data, sh.Len*sd.dm.itemRawSize, sh.Cap*sd.dm.itemRawSize

		// TODO: 需要有更高效的扩容算法
		newLen := int(float64(sh.Cap)*1.6) + 4 // 一种简易的动态扩容算法
		//fmt.Printf("growing len: %d, cap: %d \n\n", sh.Len*sd.dm.itemRawSize, newLen*sd.dm.itemRawSize)
		*(*[]byte)(sd.dstPtr) = make([]byte, sh.Len*sd.dm.itemRawSize, newLen*sd.dm.itemRawSize)

		copy(*(*[]byte)(sd.dstPtr), oldMem)
		sh.Len, sh.Cap = sd.arrIdx, newLen
	}
	ptr := unsafe.Pointer(sh.Data + uintptr(sd.arrIdx*sd.dm.itemRawSize))

	switch sd.str[sd.scan] {
	case '{', '[':
		if sd.dm.isPtr {
			ptr = sliceMixItemPtr(sd, ptr)
		}
		sd.initShareDecode(ptr)
		if sd.share.dm.isList {
			sd.share.scanList()
		} else {
			sd.share.scanObject()
		}
		sd.scan = sd.share.scan
	default:
		sd.skipNull()
		if sd.dm.isPtr {
			*(*unsafe.Pointer)(ptr) = nil
		}
	}
	sd.arrIdx++
	sh.Len = sd.arrIdx
}

// pointer +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanPointerValue(sd *subDecode) {
	ptr := getPtrValueAddr(sd.dstPtr, sd.dm.ptrLevel, sd.dm.itemKind, sd.dm.itemType)
	sd.scanSubDecode(sd.dm.itemType, ptr)
}
