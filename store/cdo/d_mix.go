package cdo

import (
	"github.com/qinchende/gofast/core/rt"
)

// array & slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

//func (d *subDecode) scanList() {
//	// A. 可能需要用到缓冲池记录临时数据
//	d.resetListPool()
//	// B. 根据目标值类型，直接匹配，提高性能
//	d.scanListItems()
//	// C. 将解析好的数据一次性绑定到对象上
//	d.flushListPool()
//}

func (d *subDecode) scanListItems() {
	off1, typ, size := scanTypeLen4(d.str[d.scan:])
	if typ != TypeList {
		panic(errChar)
	}

	d.scan += off1
	for i := 0; i < int(size); i++ {
		if d.skipValue {
			d.skipOneValue()
		} else {
			d.dm.itemDec(d)
			if d.dm.isArray {
				d.arrIdx++
				if d.arrIdx >= d.dm.arrLen {
					d.skipValue = true
				}
			}
		}
	}
	// 数组多余的部分需要重置成类型零值
	if d.arrIdx < d.dm.arrLen {
		d.resetArrLeftItems()
	}
	// 清理变量
	if d.share != nil {
		d.resetShareDecode()
	}
}

func (d *subDecode) skipList() {
	off1, typ, size := scanTypeLen4(d.str[d.scan:])
	if typ != TypeList {
		panic(errChar)
	}

	d.scan += off1
	for i := 0; i < int(size); i++ {
		d.skipOneValue()
	}
}

// struct & map
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (d *subDecode) scanKVS() {
	off, _, size := scanTypeLen2(d.str[d.scan:])
	d.scan += off

	//if typ != TypeMap {
	//	panic(errKV)
	//}

	for i := 0; i < int(size); i++ {
		off, key := scanString(d.str[d.scan:])
		d.scan += off
		d.dm.kvPairDec(d, key)
	}
}

func (d *subDecode) skipKVS() {
	off1, _, size := scanTypeLen2(d.str[d.scan:])
	//if typ != TypeMap {
	//	panic(errKV)
	//}

	for i := 0; i < int(size); i++ {
		off2 := skipString(d.str[off1:])
		d.scan += off2 + off1
		d.skipOneValue()
	}
}

// skip items
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (d *subDecode) skipOneValue() {
	c := d.str[0]
	typ := c & TypeMask
	val := c & TypeValMask
	off := 0

	switch typ {
	case TypeFixed:
		if val <= TypeList {
			d.skipList()
			break
		}
		switch val {
		default:
			panic(errChar)
		case FixNil, FixMixedNil, FixTrue, FixFalse:
			off = 1
		case FixF32:
			off = 5
		case FixF64:
			off = 9
		case FixTime:
			off = 5
		}
	case TypePosInt, TypeNegInt:
		if val <= 23 {
			off = 1
		} else {
			off = int(1 + val - 23)
		}
	case TypeStr:
		if val <= 27 {
			off = 1
		} else {
			off = int(1 + val - 27)
		}
	}

	// 解析标记往前走
	d.scan += off
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Scan Advanced mixed type, such as map | gson | struct
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// map +++++
// 目前只支持 map[string]any，并不支持其它map
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanCstKVValue(d *subDecode, k string) {

}

// map WebKV +++++
// 只支持 map[string]string
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanWebKVValue(d *subDecode, k string) {

}

// struct +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanStructValue(d *subDecode, key string) {
	// TODO: 此处 d.keyIdx 可以继续被优化
	if d.keyIdx = d.dm.ss.ColumnIndex(key); d.keyIdx < 0 {
		d.skipValue = true
		d.skipOneValue()
	} else {
		d.dm.fieldsDec[d.keyIdx](d) // 根据目标值类型来解析
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Scan Advanced mixed type, such as map | struct | array | slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// sash as map | struct
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjMixValue(d *subDecode) {
	d.startSubDecode(d.dm.ss.FieldsAttr[d.keyIdx].Type, fieldMixPtr(d))
}

func scanObjPtrMixValue(d *subDecode) {
	if d.str[d.scan] == FixMixedNil {
		fieldSetNil(d)
		d.scan++
	} else {
		d.startSubDecode(d.dm.ss.FieldsAttr[d.keyIdx].Type, fieldPtrDeep(d))
	}
}

// sash as array | slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// array and item not ptr
func scanArrMixValue(d *subDecode) {
	typ := itemType(d.str[d.scan:])

	switch typ {
	default:
		panic(errValueType)
	//case TypeMap:
	//	d.initShareDecode(arrMixItemPtr(d))
	//	d.share.scanKVS()
	//	d.scan = d.share.scan
	case TypeList:
		d.initShareDecode(arrMixItemPtr(d))
		scanListBaseType(d.share, 0)
		d.scan = d.share.scan
	}
}

// array and item is ptr
func scanArrPtrMixValue(d *subDecode) {
	if d.str[d.scan] == FixMixedNil {
		fieldSetNil(d)
		d.scan++
	} else {
		scanArrMixValue(d)
	}
}

// slice 中可能是实体对象，也可能是对象指针
func scanListMixValue(d *subDecode) {
	sh := (*rt.SliceHeader)(d.dstPtr)
	ptr := rt.SliceNextItem(sh, d.dm.itemMemSize)

	if d.dm.isPtr {
		ptr = sliceMixItemPtr(d, ptr)
	}
	d.initShareDecode(ptr)
	if d.share.dm.isList {
		scanListBaseType(d.share, 0)
	} else {
		d.share.scanKVS()
	}
	d.scan = d.share.scan

	//d.arrIdx++
	////sh.Len = d.arrIdx
}

// pointer +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanPointerValue(d *subDecode) {
	ptr := getPtrValueAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemKind, d.dm.itemType)
	d.startSubDecode(d.dm.itemType, ptr)
}
