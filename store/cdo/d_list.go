package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"unsafe"
)

func (d *subDecode) scanList() {
	off1, tLen := scanListTypeU24(d.str[d.scan:])
	d.scan += off1

	if d.dm.isSlice {
		d.scanSlice(int(tLen))
	} else {
		d.scanArray(int(tLen))
	}
}

func (d *subDecode) scanSlice(size int) {
	d.dstPtr = rt.SliceToArray(d.dstPtr, d.dm.itemMemSize, size)
	d.slice = rt.SliceHeader{DataPtr: d.dstPtr, Len: size, Cap: size}
	d.dm.listDec(d, size)
}

func (d *subDecode) scanArray(size int) {
	// 数据源中的数据量和数组大小不匹配的时候直接异常
	if d.dm.arrLen != size {
		panic(errListSize)
	}
	d.slice = rt.SliceHeader{DataPtr: d.dstPtr, Len: size, Cap: size}
	d.dm.listDec(d, size)

	//// 清理变量
	//if d.share != nil {
	//	d.resetShareDecode()
	//}
}

// 检查 List item type 是否符合预期
func validListItemType(d *subDecode, typ byte) int {
	offS := d.scan
	if d.str[offS] != typ {
		panic(errListType)
	}
	offS++
	return offS
}

func decListBaseType(d *subDecode, tLen int) {
	skipValue := false

	// 循环记录
	for i := 0; i < tLen; i++ {
		d.dstPtr = unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)

		if d.dm.isPtr {
			// 本项值为nil，直接跳过本条记录解析
			if d.str[d.scan] == FixNil {
				d.scan++
				continue
			}
			d.dstPtr = getPtrValueAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemKind, d.dm.itemType)
		}

		if skipValue {
			d.skipOneValue()
		} else {
			d.dm.itemDec(d)
			if d.dm.isArray {
				d.arrIdx++
				if d.arrIdx >= d.dm.arrLen {
					skipValue = true
				}
			}
		}
	}
}

func decListItemPtr(d *subDecode, tLen int, typ byte, fn func(iPtr unsafe.Pointer, s string) int) {
	offS := validListItemType(d, typ)
	ptrS := d.dm.itemMemSize
	for i := 0; i < tLen; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*ptrS)
		if d.str[offS] == FixNil {
			*(*unsafe.Pointer)(iPtr) = nil
			offS += 1
			continue
		}
		iPtr = getPtrValueAddr(iPtr, d.dm.ptrLevel, d.dm.itemKind, d.dm.itemType)
		offS += fn(iPtr, d.str[offS:])
	}
	d.scan = offS
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++ int +++
func decIntList(d *subDecode, tLen int) {
	list := *(*[]int)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x63 {
			off, v = scanListVarIntPart1(d.str[offS:], v)
		} else {
			off, v = scanListVarIntPart2(d.str[offS:], v)
		}
		list[i] = toInt(sym, v)
		offS += off
	}
	d.scan = offS
}

func decIntListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindInt(iPtr, sym, v)
		return off
	})
}

// +++ int8 +++
func decInt8List(d *subDecode, tLen int) {
	list := *(*[]int8)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x3F {
			off, v = scanListVarInt8(d.str[offS:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toInt8(sym, v)
		offS += off
	}
	d.scan = offS
}

func decInt8ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindInt8(iPtr, sym, v)
		return off
	})
}

// +++ int16 +++
func decInt16List(d *subDecode, tLen int) {
	list := *(*[]int16)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x40 {
			off, v = scanListVarInt16(d.str[offS:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toInt16(sym, v)
		offS += off
	}
	d.scan = offS
}

func decInt16ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindInt16(iPtr, sym, v)
		return off
	})
}

// +++ int32 +++
func decInt32List(d *subDecode, tLen int) {
	list := *(*[]int32)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x63 {
			off, v = scanListVarIntPart1(d.str[offS:], v)
		} else {
			off, v = scanListVarIntPart2(d.str[offS:], v)
		}
		list[i] = toInt32(sym, v)
		offS += off
	}
	d.scan = offS
}

func decInt32ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindInt32(iPtr, sym, v)
		return off
	})
}

// +++ int64 +++
func decInt64List(d *subDecode, tLen int) {
	list := *(*[]int64)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x63 {
			off, v = scanListVarIntPart1(d.str[offS:], v)
		} else {
			off, v = scanListVarIntPart2(d.str[offS:], v)
		}
		list[i] = toInt64(sym, v)
		offS += off
	}
	d.scan = offS
}

func decInt64ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindInt64(iPtr, sym, v)
		return off
	})
}

// +++ uint +++
func decUintList(d *subDecode, tLen int) {
	list := *(*[]uint)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x63 {
			off, v = scanListVarIntPart1(d.str[offS:], v)
		} else {
			off, v = scanListVarIntPart2(d.str[offS:], v)
		}
		list[i] = toUint(sym, v)
		offS += off
	}
	d.scan = offS
}

func decUintListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindUint(iPtr, sym, v)
		return off
	})
}

// +++ uint8 +++
func decUint8List(d *subDecode, tLen int) {
	list := *(*[]uint8)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x3F {
			off, v = scanListVarInt8(d.str[offS:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toUint8(sym, v)
		offS += off
	}
	d.scan = offS
}

func decUint8ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindUint8(iPtr, sym, v)
		return off
	})
}

// +++ uint16 +++
func decUint16List(d *subDecode, tLen int) {
	list := *(*[]uint16)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x40 {
			off, v = scanListVarInt16(d.str[offS:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toUint16(sym, v)
		offS += off
	}
	d.scan = offS
}

func decUint16ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindUint16(iPtr, sym, v)
		return off
	})
}

// +++ uint32 +++
func decUint32List(d *subDecode, tLen int) {
	list := *(*[]uint32)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x63 {
			off, v = scanListVarIntPart1(d.str[offS:], v)
		} else {
			off, v = scanListVarIntPart2(d.str[offS:], v)
		}
		list[i] = toUint32(sym, v)
		offS += off
	}
	d.scan = offS
}

func decUint32ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindUint32(iPtr, sym, v)
		return off
	})
}

// +++ uint64 +++
func decUint64List(d *subDecode, tLen int) {
	list := *(*[]uint64)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[offS])
		var off int
		if v <= 0x63 {
			off, v = scanListVarIntPart1(d.str[offS:], v)
		} else {
			off, v = scanListVarIntPart2(d.str[offS:], v)
		}
		list[i] = toUint64(sym, v)
		offS += off
	}
	d.scan = offS
}

func decUint64ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := scanListVarInt(s)
		bindUint64(iPtr, sym, v)
		return off
	})
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++ F32 +++
func decF32List(d *subDecode, tLen int) {
	list := *(*[]float32)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListF32)
	for i := 0; i < len(list); i++ {
		list[i] = scanF32Val(d.str[offS:])
		offS += 4
	}
	d.scan = offS
}

func decF32ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListF32, func(iPtr unsafe.Pointer, s string) int {
		bindF32(iPtr, scanF32Val(s))
		return 4
	})
}

// +++ F64 +++
func decF64List(d *subDecode, tLen int) {
	list := *(*[]float64)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListF64)
	for i := 0; i < len(list); i++ {
		list[i] = scanF64Val(d.str[offS:])
		offS += 8
	}
	d.scan = offS
}

func decF64ListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListF64, func(iPtr unsafe.Pointer, s string) int {
		bindF64(iPtr, scanF64Val(s))
		return 8
	})
}

// +++ Bool +++
func decBoolList(d *subDecode, tLen int) {
	list := *(*[]bool)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListBool)
	for i := 0; i < len(list); i++ {
		list[i] = scanBoolVal(d.str[offS:])
		offS += 1
	}
	d.scan = offS
}

func decBoolListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListBool, func(iPtr unsafe.Pointer, s string) int {
		bindBool(iPtr, scanBoolVal(s))
		return 1
	})
}

// +++ String +++
func decStrList(d *subDecode, tLen int) {
	list := *(*[]string)(unsafe.Pointer(&d.slice))
	offS := validListItemType(d, ListStr)
	for i := 0; i < len(list); i++ {
		off, str := scanString(d.str[offS:])
		list[i] = str
		offS += off
	}
	d.scan = offS
}

func decStrListPtr(d *subDecode, tLen int) {
	decListItemPtr(d, tLen, ListStr, func(iPtr unsafe.Pointer, s string) int {
		off, str := scanString(s)
		bindString(iPtr, str)
		return off
	})
}

// []struct & []*struct
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decStructList(d *subDecode, tLen int) {
	offS := d.scan

	// 1. struct fields
	off, typ, fSize := scanListSubtypeU16(d.str[offS:])
	if typ != ListObjFields {
		panic(errListType)
	}
	offS += off

	for i := 0; i < int(fSize); i++ {
		off1, fName := scanString(d.str[offS:])
		offS += off1
		d.fIdxes[i] = int16(d.dm.ss.ColumnIndex(fName))
	}
	d.scan = offS

	// 2. records values
	listPtr := d.dstPtr
	for i := 0; i < tLen; i++ {
		d.dstPtr = unsafe.Add(listPtr, i*d.dm.itemMemSize)

		// 如果是指针比如：[]*struct，需要分配空间
		if d.dm.isPtr {
			if d.str[d.scan] == FixNilMixed {
				d.scan++
				continue
			}
			d.dstPtr = getPtrValueAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemKind, d.dm.itemType)
		}

		// 循环字段
		for j := 0; j < int(fSize); j++ {
			d.fIdx = int(d.fIdxes[j])
			if d.fIdx >= 0 {
				d.dm.fieldsDec[d.fIdx](d)
			} else {
				d.skipOneValue()
			}
		}
	}
}
