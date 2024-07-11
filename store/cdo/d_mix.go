package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"time"
	"unsafe"
)

func (d *decoder) skipList() {
	off1, typ, size := scanTypeU32By6(d.str[d.scan:])
	if typ != TypeList {
		panic(errCdoChar)
	}

	d.scan += off1
	for i := 0; i < int(size); i++ {
		d.skipOneValue()
	}
}

func (d *decoder) skipOneValue() {
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
			panic(errCdoChar)
		case FixNil, FixNilMixed, FixTrue, FixFalse:
			off = 1
		case FixF32:
			off = 5
		case FixF64:
			off = 9
		case FixTime:
			off = 5
		}
	case TypeVarIntPos, TypeVarIntNeg:
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

func decAny(d *decoder) any {
	c := d.str[0]
	typ := c & TypeMask
	v := c & TypeValMask

	switch typ {
	case TypeFixed:
		if v >= TypeList {
			d.skipList()
			break
		}
		switch byte(v) {
		default:
			panic(errCdoChar)
		case FixNil, FixNilMixed, FixTrue, FixFalse:
		case FixF32:
			return decFixF32(d)
		case FixF64:
			return decFixF64(d)
		case FixTime:
			return decFixTime(d)
		}
	case TypeVarIntPos, TypeVarIntNeg:
		if v <= 23 {
			//off = 1
		} else {
			//off = int(1 + v - 23)
		}
	case TypeStr:
		if v <= 27 {
			//off = 1
		} else {
			//off = int(1 + v - 27)
		}
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// pointer +++++
func scanPointerValue(d *decoder) {
	ptr := getPtrValAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemType)
	d.runSub(d.dm.itemType, ptr)
}

// map & array & slice
func decListMixItem(d *decoder) {
	sh := (*rt.SliceHeader)(d.dstPtr)
	ptr := rt.SliceNextItem(sh, d.dm.itemMemSize)

	if d.dm.isPtr {
		ptr = sliceMixItemPtr(d, ptr)
	}
	d.initSub(ptr)
	if d.sub.dm.isList {
		decListAll(d.sub, 0)
	} else {
		d.sub.scanKVS()
	}
	d.scan = d.sub.scan
}

// map +++++
// 目前只支持 map[string]any，并不支持其它map
func scanCstKVValue(d *decoder, k string) {

}

// map WebKV +++++
// 只支持 map[string]string
func scanWebKVValue(d *decoder, k string) {

}

// +++ struct & map +++
func (d *decoder) scanKVS() {
	off1, tLen := decListTypeU24(d.str[d.scan:])
	d.scan += off1

	if d.str[d.scan] != ListKV {
		panic(errListType)
	}
	d.scan++

	for i := 0; i < int(tLen); i++ {
		off, key := scanString(d.str[d.scan:])
		d.scan += off
		d.dm.kvPairDec(d, key)
	}
}

func (d *decoder) skipKVS() {
	//off1, _, size := scanTypeU16(d.str[d.scan:])
	//if typ != TypeMap {
	//	panic(errKV)
	//}

	//for i := 0; i < int(size); i++ {
	//	off2 := skipString(d.str[off1:])
	//	d.scan += off2 + off1
	//	d.skipOneValue()
	//}
}

// Struct Field
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decField(d *decoder, key string) {
	if d.fIdx = d.dm.ss.ColumnIndex(key); d.fIdx < 0 {
		d.skipValue = true
		d.skipOneValue()
	} else {
		d.dm.fieldsDec[d.fIdx](d)
	}
}

func fieldNotNil(d *decoder) bool {
	if d.str[d.scan] == FixNil {
		*(*unsafe.Pointer)(fieldPtr(d)) = nil
		d.scan++
		return false
	}
	return true
}

// VarInt +++++
func decVarInt(d *decoder) (byte, uint64) {
	off, typ, v := scanVarIntVal(d.str[d.scan:])
	d.scan += off
	if typ == TypeVarIntPos || typ == TypeVarIntNeg {
		return typ, v
	}
	panic(errValType)
}

func decVarIntField(d *decoder) (unsafe.Pointer, byte, uint64) {
	off, typ, v := scanVarIntVal(d.str[d.scan:])
	d.scan += off
	if typ == TypeVarIntPos || typ == TypeVarIntNeg {
		return fieldPtr(d), typ, v
	}
	panic(errValType)
}

func decVarIntFieldPtr(binder intBinder) decValFunc {
	return func(d *decoder) {
		if fieldNotNil(d) {
			typ, v := decVarInt(d)
			binder(fieldPtrDeep(d), typ, v)
		}
	}
}

// float32 +++++
func decFixF32(d *decoder) float32 {
	pos := d.scan
	if d.str[pos] != FixF32 {
		panic(errValType)
	}
	pos++
	v := scanF32Val(d.str[pos:])
	d.scan = pos + 4
	return v
}

func decF32Field(d *decoder) {
	bindF32(fieldPtr(d), decFixF32(d))
}

func decF32FieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindF32(fieldPtrDeep(d), decFixF32(d))
	}
}

// float64 +++++
func decFixF64(d *decoder) float64 {
	pos := d.scan
	if d.str[pos] != FixF64 {
		panic(errValType)
	}
	pos++
	v := scanF64Val(d.str[pos:])
	d.scan = pos + 8
	return v
}

func decF64Field(d *decoder) {
	bindF64(fieldPtr(d), decFixF64(d))
}

func decF64FieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindF64(fieldPtrDeep(d), decFixF64(d))
	}
}

// time.Time +++++
func decFixTime(d *decoder) time.Time {
	pos := d.scan
	if d.str[pos] != FixTime {
		panic(errValType)
	}
	pos++
	v := scanTimeVal(d.str[pos:])
	d.scan = pos + 8
	return v
}

func decTimeField(d *decoder) {
	bindTime(fieldPtr(d), decFixTime(d))
}

func decTimeFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindTime(fieldPtrDeep(d), decFixTime(d))
	}
}

// string +++++
func decStrField(d *decoder) {
	off, str := scanString(d.str[d.scan:])
	d.scan += off
	bindString(fieldPtr(d), str)
}

func decStrFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		off, str := scanString(d.str[d.scan:])
		d.scan += off
		bindString(fieldPtrDeep(d), str)
	}
}

// []byte +++++
func decBytesField(d *decoder) {
	off, bs := scanBytes(d.str[d.scan:])
	d.scan += off
	bindBytes(fieldPtr(d), bs)
}

func decBytesFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		off, bs := scanBytes(d.str[d.scan:])
		d.scan += off
		bindBytes(fieldPtrDeep(d), bs)
	}
}

// bool +++++
func decBoolField(d *decoder) {
	v := scanBool(d.str[d.scan:])
	d.scan += 1
	bindBool(fieldPtr(d), v)
}

func decBoolFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		v := scanBool(d.str[d.scan:])
		d.scan += 1
		bindBool(fieldPtrDeep(d), v)
	}
}

// any +++++
func decAnyField(d *decoder) {
	bindAny(fieldPtr(d), decAny(d))
}

func decAnyFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindAny(fieldPtrDeep(d), decAny(d))
	}
}

// mixed field +++++  such as map|struct
func decMixField(d *decoder) {
	d.runSub(d.dm.ss.FieldsAttr[d.fIdx].Type, fieldMixPtr(d))
}

func decMixFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		d.runSub(d.dm.ss.FieldsAttr[d.fIdx].Type, fieldPtrDeep(d))
	}
}
