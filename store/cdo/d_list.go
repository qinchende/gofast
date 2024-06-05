package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"unsafe"
)

func (d *subDecode) scanList() {
	// 1. List length
	off1, typ, listSize := scanTypeLen4(d.str[d.scan:])
	if typ != TypeArray {
		panic(errChar)
	}
	d.scan += off1
	// 2. mem ready
	d.dstPtr = rt.SliceToArray(d.dstPtr, d.dm.itemMemSize, int(listSize))

	d.dm.listDec(d, int(listSize))

	//// 数组多余的部分需要重置成类型零值
	//if d.arrIdx < d.dm.arrLen {
	//	d.resetArrLeftItems()
	//}
	//// 清理变量
	//if d.share != nil {
	//	d.resetShareDecode()
	//}
}

func scanListBaseType(d *subDecode, listSize int) {
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

func decIntList[T constraints.Integer](d *subDecode, listSize int) {
	tpInt := d.scan
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
		off, typ, val := scanTypeLen8(d.str[tpInt:])
		tpInt += off

		if typ == TypePosInt {
			bindInt(iPtr, int64(val))
		} else if typ == TypeNegInt {
			bindInt(iPtr, int64(-val))
		} else {
			panic(errChar)
		}
	}
	d.scan = tpInt
}

func decStringList(d *subDecode, listSize int) {
	tpInt := d.scan
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)

		off, str := scanString(d.str[tpInt:])
		tpInt += off
		bindString(iPtr, str)
	}
	d.scan = tpInt
}
