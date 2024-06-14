package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"unsafe"
)

func (d *subDecode) scanList() {
	// 1. source item length
	off1, itemLen := scanListTypeU24(d.str[d.scan:])
	d.scan += off1

	if d.dm.isList {
		d.scanSlice(int(itemLen))
	} else {
		d.scanArray(int(itemLen))
	}
}

func (d *subDecode) scanSlice(size int) {
	d.dstPtr = rt.SliceToArray(d.dstPtr, d.dm.itemMemSize, size)
	d.slice = rt.SliceHeader{DataPtr: d.dstPtr, Len: size, Cap: size}
	d.dm.listDec(d, size)
}

func (d *subDecode) scanArray(size int) {
	d.slice = rt.SliceHeader{DataPtr: d.dstPtr, Len: size, Cap: size}
	d.dm.listDec(d, size)

	//// 数组多余的部分需要重置成类型零值
	//if d.arrIdx < d.dm.arrLen {
	//	d.resetArrLeftItems()
	//}
	//// 清理变量
	//if d.share != nil {
	//	d.resetShareDecode()
	//}
}

// 检查 List item type 是否符合预期
func checkListItemType(d *subDecode, typ byte) int {
	offS := d.scan
	if d.str[offS] != typ {
		panic(errListType)
	}
	offS++
	return offS
}

func decListBaseType(d *subDecode, listSize int) {
	skipValue := false

	// 循环记录
	for i := 0; i < listSize; i++ {
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

func decListItemPtr(d *subDecode, listSize int, typ byte, fn func(iPtr unsafe.Pointer, s string) int) {
	offS := checkListItemType(d, typ)
	ptrS := d.dm.itemMemSize
	for i := 0; i < listSize; i++ {
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
func decIntList(d *subDecode, listSize int) {
	list := *(*[]int)(unsafe.Pointer(&d.slice))
	offS := checkListItemType(d, ListVarInt)
	for i := 0; i < len(list); i++ {
		c := d.str[offS]
		typ := c & 0x80
		v := uint64(c & 0x7F)

		var off int
		if v <= 122 {
			off, v = scanU64ValBy7Part1(d.str[offS:], v)
		} else {
			off, v = scanU64ValBy7Part2(d.str[offS:], v)
		}
		if typ == ListVarIntPos {
			list[i] = int(v)
		} else {
			list[i] = int(-v)
		}
		offS += off
	}
	d.scan = offS
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++ F32 +++
func decF32List(d *subDecode, listSize int) {
	list := *(*[]float32)(unsafe.Pointer(&d.slice))
	offS := checkListItemType(d, ListF32)
	for i := 0; i < len(list); i++ {
		list[i] = scanF32Val(d.str[offS:])
		offS += 4
	}
	d.scan = offS
}

func decF32ListPtr(d *subDecode, listSize int) {
	decListItemPtr(d, listSize, ListF32, func(iPtr unsafe.Pointer, s string) int {
		bindF32(iPtr, scanF32Val(s))
		return 4
	})
}

// +++ F64 +++
func decF64List(d *subDecode, listSize int) {
	list := *(*[]float64)(unsafe.Pointer(&d.slice))
	offS := checkListItemType(d, ListF64)
	for i := 0; i < len(list); i++ {
		list[i] = scanF64Val(d.str[offS:])
		offS += 8
	}
	d.scan = offS
}

func decF64ListPtr(d *subDecode, listSize int) {
	decListItemPtr(d, listSize, ListF64, func(iPtr unsafe.Pointer, s string) int {
		bindF64(iPtr, scanF64Val(s))
		return 8
	})
}

// +++ Bool +++
func decBoolList(d *subDecode, listSize int) {
	list := *(*[]bool)(unsafe.Pointer(&d.slice))
	offS := checkListItemType(d, ListBool)
	for i := 0; i < len(list); i++ {
		list[i] = scanBoolVal(d.str[offS:])
		offS += 1
	}
	d.scan = offS
}

func decBoolListPtr(d *subDecode, listSize int) {
	decListItemPtr(d, listSize, ListBool, func(iPtr unsafe.Pointer, s string) int {
		bindBool(iPtr, scanBoolVal(s))
		return 1
	})
}

// +++ String +++
func decStrList(d *subDecode, listSize int) {
	list := *(*[]string)(unsafe.Pointer(&d.slice))
	offS := checkListItemType(d, ListStr)
	for i := 0; i < len(list); i++ {
		off, str := scanString(d.str[offS:])
		list[i] = str
		offS += off
	}
	d.scan = offS
}

func decStrListPtr(d *subDecode, listSize int) {
	decListItemPtr(d, listSize, ListStr, func(iPtr unsafe.Pointer, s string) int {
		off, str := scanString(s)
		bindString(iPtr, str)
		return off
	})
}
