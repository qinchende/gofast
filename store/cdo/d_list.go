package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"unsafe"
)

func (d *subDecode) scanSlice() {
	// 1. source item length
	off1, itemLen := scanListTypeU24(d.str[d.scan:])
	d.scan += off1

	// 2. mem ready
	d.dstPtr = rt.SliceToArray(d.dstPtr, d.dm.itemMemSize, int(itemLen))

	d.dm.listDec(d, int(itemLen))
}

func (d *subDecode) scanArray() {
	// 1. source item length
	off1, typ, itemLen := scanTypeLen4(d.str[d.scan:])
	if typ != TypeList {
		panic(errChar)
	}
	d.scan += off1

	// 2. mem ready
	d.dstPtr = rt.SliceToArray(d.dstPtr, d.dm.itemMemSize, int(itemLen))

	d.dm.listDec(d, int(itemLen))

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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decIntList(d *subDecode, listSize int) {
	offS := d.scan
	if d.str[offS] != ListVarInt {
		panic(errChar)
	}
	offS++
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)

		// Part1
		typ, v := typeValue(d.str[offS])
		var off int
		if v <= 59 {
			off, v = scanU64Par1(d.str[offS:], v)
		} else {
			off, v = scanU64Par2(d.str[offS:], v)
		}

		// Part2
		//typ, v := typeValue(d.str[offS])
		//var off int
		//
		//s := d.str[offS:]
		//switch v {
		//default:
		//	off = 1
		//case 56:
		//	off, v = 2, uint64(s[1])
		//case 57:
		//	off, v = 3, uint64(s[1])|uint64(s[2])<<8
		//case 58:
		//	off, v = 4, uint64(s[1])|uint64(s[2])<<8|uint64(s[3])<<16
		//case 59:
		//	off, v = 5, uint64(s[1])|uint64(s[2])<<8|uint64(s[3])<<16|uint64(s[4])<<24
		//case 60:
		//	off, v = 6, uint64(s[1])|uint64(s[2])<<8|uint64(s[3])<<16|uint64(s[4])<<24|uint64(s[5])<<32
		//case 61:
		//	off, v = 7, uint64(s[1])|uint64(s[2])<<8|uint64(s[3])<<16|uint64(s[4])<<24|uint64(s[5])<<32|uint64(s[6])<<40
		//case 62:
		//	off, v = 8, uint64(s[1])|uint64(s[2])<<8|uint64(s[3])<<16|uint64(s[4])<<24|uint64(s[5])<<32|uint64(s[6])<<40|uint64(s[7])<<48
		//case 63:
		//	off, v = 9, uint64(s[1])|uint64(s[2])<<8|uint64(s[3])<<16|uint64(s[4])<<24|uint64(s[5])<<32|uint64(s[6])<<40|uint64(s[7])<<48|uint64(s[8])<<56
		//}

		// Part3
		//off, typ, v := scanTypeLen8(d.str[offS:])

		offS += off

		if typ == TypePosInt {
			bindInt(iPtr, int64(v))
		} else if typ == TypeNegInt {
			bindInt(iPtr, int64(-v))
		} else {
			panic(errChar)
		}
	}
	d.scan = offS
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decF32List(d *subDecode, listSize int) {
	offS := d.scan
	if d.str[offS] != ListF32 {
		panic(errChar)
	}
	offS++
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
		off, f32 := scanF32Val(d.str[offS:])
		offS += off
		bindF32(iPtr, f32)
	}
	d.scan = offS
}

func decF64List(d *subDecode, listSize int) {
	offS := d.scan
	if d.str[offS] != ListF64 {
		panic(errChar)
	}
	offS++
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
		off, f64 := scanF64Val(d.str[offS:])
		offS += off
		bindF64(iPtr, f64)
	}
	d.scan = offS
}

func decBoolList(d *subDecode, listSize int) {
	offS := d.scan
	if d.str[offS] != ListBool {
		panic(errChar)
	}
	offS++
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)
		off, b := scanBool(d.str[offS:])
		offS += off
		bindBool(iPtr, b)
	}
	d.scan = offS
}

func decStringList(d *subDecode, listSize int) {
	offS := d.scan
	if d.str[offS] != ListStr {
		panic(errChar)
	}
	offS++
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(d.dstPtr, i*d.dm.itemMemSize)

		off, str := scanString(d.str[offS:])
		offS += off
		bindString(iPtr, str)
	}
	d.scan = offS
}
