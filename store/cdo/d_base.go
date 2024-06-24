package cdo

import (
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/core/rt"
	"math"
	"reflect"
	"unsafe"
)

func itemType(s string) (int uint8) {
	return s[0] & TypeMask
}

func listType(s string) (int uint8) {
	return 0
}

func typeValue(c byte) (uint8, uint64) {
	return c & TypeMask, uint64(c & TypeValMask)
}

func scanTypeLen2L2(s string) (int, uint8, uint16) {
	c := s[0]
	typ := c & ListMask
	size := uint16(c & ListValMask)

	switch size {
	default:
		return 1, typ, size // size <= 61
	case 62:
		return 2, typ, uint16(s[1])
	case 63:
		return 3, typ, uint16(s[1]) | uint16(s[2])<<8
	}
}

//func scanTypeU16(s string) (int, uint8, uint16) {
//	c := s[0]
//	typ := c & TypeMask
//	size := uint16(c & TypeValMask)
//
//	switch size {
//	default:
//		return 1, typ, size // size <= 29
//	case 30:
//		return 2, typ, uint16(s[1])
//	case 31:
//		return 3, typ, uint16(s[1]) | uint16(s[2])<<8
//
//	}
//}

func scanListTypeU24(s string) (int, uint32) {
	c := s[0]
	if c&TypeListMask != TypeList {
		panic(errChar)
	}
	size := uint32(c & TypeListValMask)

	switch size {
	default:
		return 1, size // size <= 28
	case 29:
		return 2, uint32(s[1])
	case 30:
		_ = s[2]
		return 3, uint32(s[1]) | uint32(s[2])<<8
	case 31:
		_ = s[3]
		return 4, uint32(s[1]) | uint32(s[2])<<8 | uint32(s[3])<<16
	}
}

func scanTypeU32By6(s string) (int, uint8, uint32) {
	c := s[0]
	typ := c & TypeMask
	size := uint32(c & TypeValMask)

	switch size {
	default:
		return 1, typ, size // size <= 59
	case 60:
		return 2, typ, uint32(s[1])
	case 61:
		_ = s[2]
		return 3, typ, uint32(s[1]) | uint32(s[2])<<8
	case 62:
		_ = s[3]
		return 4, typ, uint32(s[1]) | uint32(s[2])<<8 | uint32(s[3])<<16
	case 63:
		_ = s[4]
		return 5, typ, uint32(s[1]) | uint32(s[2])<<8 | uint32(s[3])<<16 | uint32(s[4])<<24
	}
}

func scanTypeU64By6(s string) (int, uint8, uint64) {
	c := s[0]
	typ := c & TypeMask
	size := uint64(c & TypeValMask)

	switch size {
	default:
		return 1, typ, size // size <= 55
	case 56:
		return 2, typ, uint64(s[1])
	case 57:
		_ = s[2]
		return 3, typ, uint64(s[1]) | uint64(s[2])<<8
	case 58:
		_ = s[3]
		return 4, typ, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16
	case 59:
		_ = s[4]
		return 5, typ, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24
	case 60:
		_ = s[5]
		return 6, typ, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32
	case 61:
		_ = s[6]
		return 7, typ, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40
	case 62:
		_ = s[7]
		return 8, typ, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48
	case 63:
		_ = s[8]
		return 9, typ, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48 | uint64(s[8])<<56
	}
}

// 用第一位标记符号的 VarInt
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ListVarInt
// 每一项第一个字节，第1位为符号位，后面7位是值
func listVarIntHead(c byte) (byte, uint64) {
	return c, uint64(c & 0x7F)
}

func scanListVarInt8(s string, v uint64) (int, uint64) {
	switch byte(v) >> 5 {
	case 0:
		return 1, uint64(s[0] & 0x1F)
	default:
		return 2, uint64(s[0]&0x1F)<<8 | uint64(s[1])
	}
}

func scanListVarInt16(s string, v uint64) (int, uint64) {
	switch byte(v) >> 5 {
	case 0:
		return 1, uint64(s[0] & 0x1F)
	case 1:
		_ = s[1]
		return 2, uint64(s[0]&0x1F)<<8 | uint64(s[1])
	default:
		_ = s[2]
		return 3, uint64(s[0]&0x1F)<<16 | uint64(s[1])<<8 | uint64(s[2])
	}
}

func scanListVarInt(s string) (byte, int, uint64) {
	typ, v := listVarIntHead(s[0])
	var off int
	if v <= 0x63 {
		off, v = scanListVarIntPart1(s, v)
	} else {
		off, v = scanListVarIntPart2(s, v)
	}
	return typ, off, v
}

func scanListVarIntPart1(s string, v uint64) (int, uint64) {
	switch byte(v) >> 5 {
	case 0:
		return 1, uint64(s[0] & 0x1F)
	case 1:
		return 2, uint64(s[0]&0x1F)<<8 | uint64(s[1])
	case 2:
		return 3, uint64(s[0]&0x1F)<<16 | uint64(s[1])<<8 | uint64(s[2])
	default:
		return 4, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16
	}
}

func scanListVarIntPart2(s string, v uint64) (int, uint64) {
	switch byte(v) & 0x0F {
	default:
		panic(errChar)
	case 4:
		_ = s[4]
		return 5, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24
	case 5:
		_ = s[5]
		return 6, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32
	case 6:
		_ = s[6]
		return 7, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40
	case 7:
		_ = s[7]
		return 8, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48
	case 8:
		_ = s[8]
		return 9, uint64(s[1]) | uint64(s[2])<<8 | uint64(s[3])<<16 | uint64(s[4])<<24 | uint64(s[5])<<32 | uint64(s[6])<<40 | uint64(s[7])<<48 | uint64(s[8])<<56
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func scanUint64(str string) (int, uint64) {
//	off, typ, val := scanTypeU64By6(str)
//	if typ != TypeVarIntPos {
//		panic(errChar)
//	}
//	return off, val
//}

func scanVarInt(str string) (int, byte, uint64) {
	return scanTypeU64By6(str)
}

func scanF32Val(s string) float32 {
	_ = s[3]
	return math.Float32frombits(uint32(s[0]) | uint32(s[1])<<8 | uint32(s[2])<<16 | uint32(s[3])<<24)
}

func scanF64Val(s string) float64 {
	_ = s[7]
	return math.Float64frombits(uint64(s[0]) | uint64(s[1])<<8 | uint64(s[2])<<16 | uint64(s[3])<<24 |
		uint64(s[4])<<32 | uint64(s[5])<<40 | uint64(s[6])<<48 | uint64(s[7])<<56)
}

func scanBoolVal(s string) bool {
	if s[0] == FixTrue {
		return true
	} else {
		return false
	}
}

func scanBool(s string) bool {
	switch s[0] {
	default:
		panic(errChar)
	case FixTrue:
		return true
	case FixFalse:
		return false
	}
}

func scanString(str string) (int, string) {
	off, typ, size := scanTypeU32By6(str)
	if typ != TypeStr {
		panic(errChar)
	}
	size += uint32(off)
	return int(size), str[off:size]
}

func skipString(str string) int {
	off, typ, size := scanTypeU32By6(str)
	if typ != TypeStr {
		panic(errChar)
	}
	return off + int(size)
}

func scanBytes(str string) (int, []byte) {
	off, v := scanString(str)
	return off, lang.STB(v)
}

// memory plan
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
			newPtr = unsafe.Pointer(&rt.SliceHeader{})
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
func arrItemPtr(d *subDecode) unsafe.Pointer {
	return unsafe.Add(d.dstPtr, d.arrIdx*d.dm.itemMemSize)
}

func arrMixItemPtr(d *subDecode) unsafe.Pointer {
	ptr := unsafe.Add(d.dstPtr, d.arrIdx*d.dm.itemMemSize)

	// 只有field字段为map或者slice的时候，值才可能是nil
	if d.dm.itemKind == reflect.Map {
		if *(*unsafe.Pointer)(ptr) == nil {
			*(*unsafe.Pointer)(ptr) = reflect.MakeMap(d.dm.itemType).UnsafePointer()
		}
	}
	return ptr
}

func arrMixItemPtrDeep(d *subDecode) unsafe.Pointer {
	ptr := unsafe.Add(d.dstPtr, d.arrIdx*d.dm.itemMemSize)
	return getPtrValueAddr(ptr, d.dm.ptrLevel, d.dm.itemKind, d.dm.itemType)
}

func sliceMixItemPtr(d *subDecode, ptr unsafe.Pointer) unsafe.Pointer {
	return getPtrValueAddr(ptr, d.dm.ptrLevel, d.dm.itemKind, d.dm.itemType)
}

// reset array left item
func (d *subDecode) resetArrLeftItems() {
	var dfValue unsafe.Pointer
	if !d.dm.isPtr {
		dfValue = zeroValues[d.dm.itemKind]
	}
	for i := d.arrIdx; i < d.dm.arrLen; i++ {
		*(*unsafe.Pointer)(unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)) = dfValue
	}
}

// struct & map & gson
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//

func fieldPtr(d *subDecode) unsafe.Pointer {
	return d.dm.ss.FieldsAttr[d.keyIdx].MyPtr(d.dstPtr)
}

func fieldMixPtr(d *subDecode) unsafe.Pointer {
	fa := &d.dm.ss.FieldsAttr[d.keyIdx]
	ptr := fa.MyPtr(d.dstPtr)

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

func fieldPtrDeep(d *subDecode) unsafe.Pointer {
	fa := &d.dm.ss.FieldsAttr[d.keyIdx]
	ptr := fa.MyPtr(d.dstPtr)
	return getPtrValueAddr(ptr, fa.PtrLevel, fa.Kind, fa.Type)
}

func fieldSetNil(d *subDecode) {
	*(*unsafe.Pointer)(fieldPtr(d)) = nil
}
