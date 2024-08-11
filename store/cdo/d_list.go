package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"time"
	"unsafe"
)

func (d *decoder) decList() {
	off1, tLen := decListTypeU24(d.str[d.scan:])
	d.scan += off1

	if d.dm.isSlice {
		d.decSlice(int(tLen))
	} else {
		d.decArray(int(tLen))
	}
}

func (d *decoder) decSlice(tLen int) {
	// 含有指针字段的struct，需要注意gc标记的问题
	if d.dm.isUnsafe {
		d.dstPtr = rt.SliceToArraySafe(d.dstPtr, tLen, d.dm.listType)
	} else {
		d.dstPtr = rt.SliceToArray(d.dstPtr, d.dm.itemMemSize, tLen)
	}
	d.slice = rt.SliceHeader{DataPtr: d.dstPtr, Len: tLen, Cap: tLen}
	d.dm.listDec(d, tLen)
}

func (d *decoder) decArray(tLen int) {
	if d.dm.arrLen != tLen {
		panic(errListSize) // 目标和数据源数量不匹配，直接异常
	}
	d.slice = rt.SliceHeader{DataPtr: d.dstPtr, Len: tLen, Cap: tLen}
	d.dm.listDec(d, tLen)
}

// 检查 List item type 是否符合预期
func validListItemType(d *decoder, typ byte) int {
	pos := d.scan
	if d.str[pos] != typ {
		panic(errListType)
	}
	pos++
	return pos
}

func decListAll(d *decoder, tLen int) {
	for i := 0; i < tLen; i++ {
		d.dstPtr = unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
		if d.dm.isPtr {
			if d.str[d.scan] == FixNilMixed {
				d.scan++
				continue
			}
			d.dstPtr = getPtrValAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemTypeAbi)
		}
		d.dm.itemDec(d)
	}
}

func decListFuncPtr(d *decoder, tLen int, typ byte, fn func(iPtr unsafe.Pointer, s string) int) {
	pos := validListItemType(d, typ)
	for i := 0; i < tLen; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
		if d.str[pos] == FixNil {
			*(*unsafe.Pointer)(iPtr) = nil
			pos += 1
			continue
		}
		iPtr = getPtrValAddr(iPtr, d.dm.ptrLevel, d.dm.itemTypeAbi)
		pos += fn(iPtr, d.str[pos:])
	}
	d.scan = pos
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++ int +++
func decListInt(d *decoder, tLen int) {
	list := *(*[]int)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x63 {
			off, v = decListVarIntPart1(d.str[pos:], v)
		} else {
			off, v = decListVarIntPart2(d.str[pos:], v)
		}
		list[i] = toInt(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListIntPtr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindInt(iPtr, sym, v)
		return off
	})
}

// +++ int8 +++
func decListInt8(d *decoder, tLen int) {
	list := *(*[]int8)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x3F {
			off, v = decListVarInt8(d.str[pos:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toInt8(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListInt8Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindInt8(iPtr, sym, v)
		return off
	})
}

// +++ int16 +++
func decListInt16(d *decoder, tLen int) {
	list := *(*[]int16)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x40 {
			off, v = decListVarInt16(d.str[pos:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toInt16(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListInt16Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindInt16(iPtr, sym, v)
		return off
	})
}

// +++ int32 +++
func decListInt32(d *decoder, tLen int) {
	list := *(*[]int32)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x63 {
			off, v = decListVarIntPart1(d.str[pos:], v)
		} else {
			off, v = decListVarIntPart2(d.str[pos:], v)
		}
		list[i] = toInt32(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListInt32Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindInt32(iPtr, sym, v)
		return off
	})
}

// +++ int64 +++
func decListInt64(d *decoder, tLen int) {
	list := *(*[]int64)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x63 {
			off, v = decListVarIntPart1(d.str[pos:], v)
		} else {
			off, v = decListVarIntPart2(d.str[pos:], v)
		}
		list[i] = toInt64(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListInt64Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindInt64(iPtr, sym, v)
		return off
	})
}

// +++ uint +++
func decListUint(d *decoder, tLen int) {
	list := *(*[]uint)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x63 {
			off, v = decListVarIntPart1(d.str[pos:], v)
		} else {
			off, v = decListVarIntPart2(d.str[pos:], v)
		}
		list[i] = toUint(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListUintPtr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindUint(iPtr, sym, v)
		return off
	})
}

// +++ uint8 +++
func decListUint8(d *decoder, tLen int) {
	list := *(*[]uint8)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x3F {
			off, v = decListVarInt8(d.str[pos:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toUint8(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListUint8Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindUint8(iPtr, sym, v)
		return off
	})
}

// +++ uint16 +++
func decListUint16(d *decoder, tLen int) {
	list := *(*[]uint16)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x40 {
			off, v = decListVarInt16(d.str[pos:], v)
		} else {
			panic(errInfinity)
		}
		list[i] = toUint16(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListUint16Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindUint16(iPtr, sym, v)
		return off
	})
}

// +++ uint32 +++
func decListUint32(d *decoder, tLen int) {
	list := *(*[]uint32)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x63 {
			off, v = decListVarIntPart1(d.str[pos:], v)
		} else {
			off, v = decListVarIntPart2(d.str[pos:], v)
		}
		list[i] = toUint32(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListUint32Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindUint32(iPtr, sym, v)
		return off
	})
}

// +++ uint64 +++
func decListUint64(d *decoder, tLen int) {
	list := *(*[]uint64)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		sym, v := listVarIntHead(d.str[pos])
		var off int
		if v <= 0x63 {
			off, v = decListVarIntPart1(d.str[pos:], v)
		} else {
			off, v = decListVarIntPart2(d.str[pos:], v)
		}
		list[i] = toUint64(sym, v)
		pos += off
	}
	d.scan = pos
}

func decListUint64Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListVarInt, func(iPtr unsafe.Pointer, s string) int {
		sym, off, v := decListVarInt(s)
		bindUint64(iPtr, sym, v)
		return off
	})
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++ F32 +++
func decListF32(d *decoder, tLen int) {
	list := *(*[]float32)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListF32)
	for i := 0; i < len(list); i++ {
		list[i] = scanF32Val(d.str[pos:])
		pos += 4
	}
	d.scan = pos
}

func decListF32Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListF32, func(iPtr unsafe.Pointer, s string) int {
		bindF32(iPtr, scanF32Val(s))
		return 4
	})
}

// +++ F64 +++
func decListF64(d *decoder, tLen int) {
	list := *(*[]float64)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListF64)
	for i := 0; i < len(list); i++ {
		list[i] = scanF64Val(d.str[pos:])
		pos += 8
	}
	d.scan = pos
}

func decListF64Ptr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListF64, func(iPtr unsafe.Pointer, s string) int {
		bindF64(iPtr, scanF64Val(s))
		return 8
	})
}

// +++ Bool +++
func decListBool(d *decoder, tLen int) {
	list := *(*[]bool)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListBool)
	for i := 0; i < len(list); i++ {
		list[i] = scanBoolVal(d.str[pos:])
		pos += 1
	}
	d.scan = pos
}

func decListBoolPtr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListBool, func(iPtr unsafe.Pointer, s string) int {
		bindBool(iPtr, scanBoolVal(s))
		return 1
	})
}

// +++ String +++
func decListStr(d *decoder, tLen int) {
	list := *(*[]string)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListStr)
	for i := 0; i < len(list); i++ {
		off, str := scanString(d.str[pos:])
		pos += off
		list[i] = str
	}
	d.scan = pos
}

func decListStrPtr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListStr, func(iPtr unsafe.Pointer, s string) int {
		off, str := scanString(s)
		bindString(iPtr, str)
		return off
	})
}

// +++ time +++
func decListTime(d *decoder, tLen int) {
	list := *(*[]time.Time)(unsafe.Pointer(&d.slice))
	pos := validListItemType(d, ListTime)
	for i := 0; i < len(list); i++ {
		list[i] = scanTimeVal(d.str[pos:])
		pos += 8
	}
	d.scan = pos
}

func decListTimePtr(d *decoder, tLen int) {
	decListFuncPtr(d, tLen, ListTime, func(iPtr unsafe.Pointer, s string) int {
		bindTime(iPtr, scanTimeVal(s))
		return 8
	})
}

// []struct & []*struct
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decListStruct(d *decoder, tLen int) {
	pos := d.scan

	// 1. struct fields
	off, typ, fSize := decListSubtypeU16(d.str[pos:])
	if typ != ListObjFields {
		panic(errListType)
	}
	pos += off

	for i := 0; i < int(fSize); i++ {
		off1, fName := scanString(d.str[pos:])
		pos += off1
		d.fIdxes[i] = int16(d.dm.ss.ColumnIndex(fName))
	}
	d.scan = pos

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
			d.dstPtr = getPtrValAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemTypeAbi)
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
