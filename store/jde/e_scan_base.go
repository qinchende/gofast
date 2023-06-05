package jde

import (
	"reflect"
	"unsafe"
)

//// Note: 此函数只适合 object 的 field，List 的 item 为 指针类型 的情形。非指针不能调用此方法
//func getPtrValueAddr(ptr unsafe.Pointer, ptrLevel uint8, kind reflect.Kind, rfType reflect.Type) unsafe.Pointer {
//	for ptrLevel > 1 {
//		if *(*unsafe.Pointer)(ptr) == nil {
//			tpPtr := unsafe.Pointer(new(unsafe.Pointer))
//			*(*unsafe.Pointer)(ptr) = tpPtr
//			ptr = tpPtr
//		} else {
//			ptr = *(*unsafe.Pointer)(ptr)
//		}
//
//		ptrLevel--
//	}
//
//	if *(*unsafe.Pointer)(ptr) == nil {
//		var newPtr unsafe.Pointer
//
//		switch kind {
//		case reflect.Map:
//			newPtr = unsafe.Pointer(new(unsafe.Pointer))
//			*(*unsafe.Pointer)(newPtr) = reflect.MakeMap(rfType).UnsafePointer()
//		case reflect.Slice:
//			newPtr = unsafe.Pointer(&reflect.SliceHeader{})
//			*(*unsafe.Pointer)(newPtr) = reflect.MakeSlice(rfType, 0, 0).UnsafePointer()
//		default:
//			newPtr = reflect.New(rfType).UnsafePointer()
//		}
//
//		*(*unsafe.Pointer)(ptr) = newPtr
//		return newPtr
//	}
//	return *(*unsafe.Pointer)(ptr)
//}

// array & slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encArrItemPtr(sd *subEncode, idx int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(idx*sd.dm.itemBytes))
}

func encArrMixItemPtr(sd *subEncode, idx int) unsafe.Pointer {
	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(idx*sd.dm.itemBytes))

	// 只有field字段为map或者slice的时候，值才可能是nil
	if sd.dm.itemBaseKind == reflect.Map {
		if *(*unsafe.Pointer)(ptr) == nil {
			*(*unsafe.Pointer)(ptr) = reflect.MakeMap(sd.dm.itemBaseType).UnsafePointer()
		}
	}
	return ptr
}

func encArrMixItemPtrDeep(sd *subEncode, idx int) unsafe.Pointer {
	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + uintptr(idx*sd.dm.itemBytes))
	return getPtrValueAddr(ptr, sd.dm.ptrLevel, sd.dm.itemBaseKind, sd.dm.itemBaseType)
}

func encSliceMixItemPtr(sd *subEncode, idx int, ptr unsafe.Pointer) unsafe.Pointer {
	return getPtrValueAddr(ptr, sd.dm.ptrLevel, sd.dm.itemBaseKind, sd.dm.itemBaseType)
}

//// struct & map & gson
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func fieldPtr(sd *subEncode) unsafe.Pointer {
//	return unsafe.Pointer(uintptr(sd.dstPtr) + sd.dm.ss.FieldsAttr[sd.keyIdx].Offset)
//}
//
//func fieldMixPtr(sd *subEncode) unsafe.Pointer {
//	fa := &sd.dm.ss.FieldsAttr[sd.keyIdx]
//	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + fa.Offset)
//
//	if fa.Kind == reflect.Map {
//		if *(*unsafe.Pointer)(ptr) == nil {
//			*(*unsafe.Pointer)(ptr) = reflect.MakeMap(fa.Type).UnsafePointer()
//		}
//	}
//
//	//// 只有field字段为map或者slice的时候，值才可能是nil
//	//if *(*unsafe.Pointer)(ptr) == nil {
//	//	switch fa.Kind {
//	//	// Note: 当 array & slice & struct 的时候，相当于是值类型，直接返回首地址即可
//	//	//default:
//	//	//	panic(errSupport)
//	//	case reflect.Map:
//	//		*(*unsafe.Pointer)(ptr) = reflect.MakeMap(fa.Type).UnsafePointer()
//	//		//case reflect.Slice:
//	//		// Note: fa.Kind == reflect.Slice，
//	//		// 此时可能申请slice对象没有意义，因为解析程序会自己创建临时空间，完成之后替换旧内存
//	//		// 但如果slice中的项还是 mix 类型，可能又不一样了，这种情况解析程序不会申请临时空间
//	//		//newPtr := reflect.MakeSlice(fa.Type, 0, 4).UnsafePointer()	// 默认给4个值的空间，避免扩容
//	//		//*(*unsafe.Pointer)(ptr) = *(*unsafe.Pointer)(newPtr)
//	//	}
//	//}
//	return ptr
//}
//
//func fieldPtrDeep(sd *subEncode) unsafe.Pointer {
//	fa := &sd.dm.ss.FieldsAttr[sd.keyIdx]
//	ptr := unsafe.Pointer(uintptr(sd.dstPtr) + fa.Offset)
//	return getPtrValueAddr(ptr, fa.PtrLevel, fa.Kind, fa.Type)
//}
//
//func fieldSetNil(sd *subEncode) {
//	*(*unsafe.Pointer)(fieldPtr(sd)) = nil
//}

func encObjPtrIntValue(se *subEncode, idx int) {
}

func encListIntValue(se *subEncode, idx int) {
}

func encListStrValue(se *subEncode, idx int) {
	//ptr := encArrItemPtr(se, idx)

	se.bs = append(se.bs, '"')
}