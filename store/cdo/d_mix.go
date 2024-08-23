package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"time"
	"unsafe"
)

// Basic type value
func (d *decoder) decBasic() {
	d.dm.itemDec(d)
}

// pointer +++++
func (d *decoder) decPointer() {
	d.dstPtr = getPtrValAddr(d.dstPtr, d.dm.ptrLevel, d.dm.itemTypeAbi)
	d.dm.itemDec(d)
}

// any value +++++
func (d *decoder) decAny() any {
	c := d.str[d.scan]
	typ := c & TypeMask
	val := c & TypeValMask

	switch typ {
	case TypeFixed:
		if val >= TypeList {
			d.decList()
			break
		}
		switch val {
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
	case TypeStr:
	}
	return nil
}

// +++ struct & map +++
func (d *decoder) kvLenPos() (int, int) {
	pos := d.scan
	off, tLen := decListTypeU24(d.str[pos:])
	pos += off

	if d.str[pos] != ListKV {
		panic(errListType)
	}
	pos++
	return int(tLen), pos
}

func (d *decoder) decStruct() {
	tLen, pos := d.kvLenPos()
	d.scan = pos
	for i := 0; i < tLen; i++ {
		d.dm.kvPairDec(d, decStrVal(d))
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// list[map | array | slice]
func decListMixItem(d *decoder) {
	sh := (*rt.SliceHeader)(d.dstPtr)
	ptr := rt.SliceNextItemSafe(sh, d.dm.itemMemSize, d.dm.itemType)
	if d.dm.isPtr {
		ptr = getPtrValAddr(ptr, d.dm.ptrLevel, d.dm.itemTypeAbi)
	}
	d.runSub(d.dm.itemType, ptr)
}

// Struct Field
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decField(d *decoder, key string) {
	if d.fIdx = d.dm.ss.ColumnIndex(key); d.fIdx >= 0 {
		d.dm.fieldsDec[d.fIdx](d)
	} else {
		d.skipOneValue()
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

func decVarIntFieldPtr(binder intBinder) decValFunc {
	return func(d *decoder) {
		if fieldNotNil(d) {
			typ, v := decVarInt(d)
			binder(fieldPtrDeep(d), typ, v)
		}
	}
}

func decInt(d *decoder) {
	typ, v := decVarInt(d)
	bindInt(d.dstPtr, typ, v)
}

func decIntField(d *decoder) {
	typ, v := decVarInt(d)
	bindInt(fieldPtr(d), typ, v)
}

func decInt8(d *decoder) {
	typ, v := decVarInt(d)
	bindInt8(d.dstPtr, typ, v)
}

func decInt8Field(d *decoder) {
	typ, v := decVarInt(d)
	bindInt8(fieldPtr(d), typ, v)
}

func decInt16(d *decoder) {
	typ, v := decVarInt(d)
	bindInt16(d.dstPtr, typ, v)
}

func decInt16Field(d *decoder) {
	typ, v := decVarInt(d)
	bindInt16(fieldPtr(d), typ, v)
}

func decInt32(d *decoder) {
	typ, v := decVarInt(d)
	bindInt32(d.dstPtr, typ, v)
}

func decInt32Field(d *decoder) {
	typ, v := decVarInt(d)
	bindInt32(fieldPtr(d), typ, v)
}

func decInt64(d *decoder) {
	typ, v := decVarInt(d)
	bindInt64(d.dstPtr, typ, v)
}

func decInt64Field(d *decoder) {
	typ, v := decVarInt(d)
	bindInt64(fieldPtr(d), typ, v)
}

func decUint(d *decoder) {
	typ, v := decVarInt(d)
	bindUint(d.dstPtr, typ, v)
}

func decUintField(d *decoder) {
	typ, v := decVarInt(d)
	bindUint(fieldPtr(d), typ, v)
}

func decUint8(d *decoder) {
	typ, v := decVarInt(d)
	bindUint8(d.dstPtr, typ, v)
}

func decUint8Field(d *decoder) {
	typ, v := decVarInt(d)
	bindUint8(fieldPtr(d), typ, v)
}

func decUint16(d *decoder) {
	typ, v := decVarInt(d)
	bindUint16(d.dstPtr, typ, v)
}

func decUint16Field(d *decoder) {
	typ, v := decVarInt(d)
	bindUint16(fieldPtr(d), typ, v)
}

func decUint32(d *decoder) {
	typ, v := decVarInt(d)
	bindUint32(d.dstPtr, typ, v)
}

func decUint32Field(d *decoder) {
	typ, v := decVarInt(d)
	bindUint32(fieldPtr(d), typ, v)
}

func decUint64(d *decoder) {
	typ, v := decVarInt(d)
	bindUint64(d.dstPtr, typ, v)
}

func decUint64Field(d *decoder) {
	typ, v := decVarInt(d)
	bindUint64(fieldPtr(d), typ, v)
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

func decF32(d *decoder) {
	bindF32(d.dstPtr, decFixF32(d))
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

func decF64(d *decoder) {
	bindF64(d.dstPtr, decFixF64(d))
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

func decTime(d *decoder) {
	bindTime(d.dstPtr, decFixTime(d))
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
func decStrVal(d *decoder) string {
	off, str := scanString(d.str[d.scan:])
	d.scan += off
	return str
}

func decStr(d *decoder) {
	bindString(d.dstPtr, decStrVal(d))
}

func decStrField(d *decoder) {
	bindString(fieldPtr(d), decStrVal(d))
}

func decStrFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindString(fieldPtrDeep(d), decStrVal(d))
	}
}

// []byte +++++
func decBytesVal(d *decoder) []byte {
	off, bs := scanBytes(d.str[d.scan:])
	d.scan += off
	return bs
}

func decBytes(d *decoder) {
	bindBytes(d.dstPtr, decBytesVal(d))
}

func decBytesField(d *decoder) {
	bindBytes(fieldPtr(d), decBytesVal(d))
}

func decBytesFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindBytes(fieldPtrDeep(d), decBytesVal(d))
	}
}

// bool +++++
func decBoolVal(d *decoder) bool {
	v := scanBool(d.str[d.scan:])
	d.scan += 1
	return v
}

func decBool(d *decoder) {
	bindBool(d.dstPtr, decBoolVal(d))
}

func decBoolField(d *decoder) {
	bindBool(fieldPtr(d), decBoolVal(d))
}

func decBoolFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindBool(fieldPtrDeep(d), decBoolVal(d))
	}
}

// any +++++
func decAnyVal(d *decoder, iPtr unsafe.Pointer) {
	v := *((*any)(iPtr))
	vTyp := reflect.TypeOf(v)
	if v == nil || vTyp.Kind() != reflect.Pointer {
		bindAny(iPtr, d.decAny())
		return
	}
	vTyp = vTyp.Elem()
	d.runSub(vTyp, (*rt.AFace)(iPtr).DataPtr)
}

func decAny(d *decoder) {
	decAnyVal(d, d.dstPtr)
}

func decAnyField(d *decoder) {
	decAnyVal(d, fieldPtr(d))
}

func decAnyFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		bindAny(fieldPtrDeep(d), d.decAny())
	}
}

// mixed field +++++  such as map|struct
func decMix(d *decoder) {
	d.runSub(d.dm.itemType, d.dstPtr)
}

func decMixField(d *decoder) {
	d.runSub(d.dm.ss.FieldsAttr[d.fIdx].Type, fieldMixPtr(d))
}

func decMixFieldPtr(d *decoder) {
	if fieldNotNil(d) {
		d.runSub(d.dm.ss.FieldsAttr[d.fIdx].Type, fieldPtrDeep(d))
	}
}
