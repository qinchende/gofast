package cdo

import (
	"unsafe"
)

//func (d *subDecode) resetForListStruct() {
//	//d.columns = d.columns[0:0]
//}

func decListStruct(d *subDecode, tLen int) {
	//d.resetForListStruct()
	offS := d.scan

	//// 1. List length ++++++++++++++++++++++++++++++++++++++++++
	//off1, typ, size1 := scanTypeU32By6(d.str[offS:])
	//if typ != TypeArrSame {
	//	panic(errChar)
	//}
	//offS += off1

	// 2. Struct fields ++++++++++++++++++++++++++++++++++++++++
	off2, typ2, fSize := scanTypeLen2L2(d.str[offS:])
	if typ2 != ListObjFields {
		panic(errChar)
	}
	offS += off2

	for i := 0; i < int(fSize); i++ {
		off, fName := scanString(d.str[offS:])
		offS += off
		d.clsIdx[i] = int8(d.dm.ss.ColumnIndex(fName))
	}
	//d.clsCt = int(fSize) // 多少个有效字段
	d.scan = offS

	// 3. Records value ++++++++++++++++++++++++++++++++++++++++
	//tLen := int(size1)
	//ptr := rt.SliceToArray(d.dstPtr, d.dm.itemMemSize, tLen) // 当切片类型时 []struct

	// 循环记录
	for i := 0; i < tLen; i++ {
		d.dstPtr = unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
		//itemPtr := d.dstPtr

		// 如果是指针比如：[]*struct，需要分配空间
		if d.dm.isPtr {
			// 本项值为nil，直接跳过本条记录解析
			if d.str[d.scan] == FixNilMixed {
				d.scan++
				continue
			}
			d.dstPtr = getPtrValueAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemKind, d.dm.itemType)
		}

		// 循环字段
		for j := 0; j < int(fSize); j++ {
			d.keyIdx = int(d.clsIdx[j])
			if d.keyIdx >= 0 {
				d.dm.fieldsDec[d.keyIdx](d)
			} else {
				d.skipOneValue()
			}
		}
	}
}
