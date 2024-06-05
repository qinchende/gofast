package cdo

import (
	"github.com/qinchende/gofast/core/cst"
	"math"
	"reflect"
	"time"
	"unsafe"
)

// int
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//go:inline
func bindInt(p unsafe.Pointer, v int64) {
	*(*int)(p) = int(v)
}

//go:inline
func bindInt8(p unsafe.Pointer, v int64) {
	if v < math.MinInt8 || v > math.MaxInt8 {
		panic(errInfinity)
	}
	*(*int8)(p) = int8(v)
}

func bindInt16(p unsafe.Pointer, v int64) {
	if v < math.MinInt16 || v > math.MaxInt16 {
		panic(errInfinity)
	}
	*(*int16)(p) = int16(v)
}

func bindInt32(p unsafe.Pointer, v int64) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		panic(errInfinity)
	}
	*(*int32)(p) = int32(v)
}

func bindInt64(p unsafe.Pointer, v int64) {
	*(*int64)(p) = v
}

// uint
func bindUint(p unsafe.Pointer, v uint64) {
	*(*uint)(p) = uint(v)
}

func bindUint8(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint8 {
		panic(errInfinity)
	}
	*(*uint8)(p) = uint8(v)
}

func bindUint16(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint16 {
		panic(errInfinity)
	}
	*(*uint16)(p) = uint16(v)
}

func bindUint32(p unsafe.Pointer, v uint64) {
	if v > math.MaxUint32 {
		panic(errInfinity)
	}
	*(*uint32)(p) = uint32(v)
}

func bindUint64(p unsafe.Pointer, v uint64) {
	*(*uint64)(p) = v
}

// float
func bindFloat32(p unsafe.Pointer, v float32) {
	*(*float32)(p) = v
}

func bindFloat64(p unsafe.Pointer, v float64) {
	*(*float64)(p) = v
}

// []byte
func bindBytes(p unsafe.Pointer, v []byte) {
	*(*[]byte)(p) = v
}

// string & bool & any
func bindString(p unsafe.Pointer, v string) {
	*(*string)(p) = v
}

func bindBool(p unsafe.Pointer, v bool) {
	*(*bool)(p) = v
}

func bindAny(p unsafe.Pointer, v any) {
	*(*any)(p) = v
}

// 时间默认都是按 RFC3339 格式存储并解析
func bindTime(p unsafe.Pointer, v string) {
	if tm, err := time.Parse(cst.TimeFmtRFC3339, v); err != nil {
		panic(err)
	} else {
		*(*time.Time)(p) = tm
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// int +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrIntValue(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt(fieldPtrDeep(d), v)
	}
}

func scanObjIntValue(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt(fieldPtr(d), v)
}

func scanArrIntValue(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt(arrItemPtr(d), v)
}

func scanListIntValue(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	d.pl.bufInt = append(d.pl.bufInt, int(v))
}

// int8 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt8Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt8(fieldPtrDeep(d), v)
	}
}

func scanObjInt8Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt8(fieldPtr(d), v)
}

func scanArrInt8Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt8(arrItemPtr(d), v)
}

func scanListInt8Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off

	if v < math.MinInt8 || v > math.MaxInt8 {
		panic(errInfinity)
	}
	d.pl.bufI8 = append(d.pl.bufI8, int8(v))
}

// int16 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt16Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt16(fieldPtrDeep(d), v)
	}
}

func scanObjInt16Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt16(fieldPtr(d), v)
}

func scanArrInt16Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt16(arrItemPtr(d), v)
}

func scanListInt16Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off

	if v < math.MinInt16 || v > math.MaxInt16 {
		panic(errInfinity)
	}
	d.pl.bufI16 = append(d.pl.bufI16, int16(v))
}

// int32 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt32Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt32(fieldPtrDeep(d), v)
	}
}

func scanObjInt32Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt32(fieldPtr(d), v)
}

func scanArrInt32Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt32(arrItemPtr(d), v)
}

func scanListInt32Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off

	if v < math.MinInt32 || v > math.MaxInt32 {
		panic(errInfinity)
	}
	d.pl.bufI32 = append(d.pl.bufI32, int32(v))
}

// int64 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrInt64Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt64(fieldPtrDeep(d), v)
	}
}

func scanObjInt64Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt64(fieldPtr(d), v)
}

func scanArrInt64Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	bindInt64(arrItemPtr(d), v)
}

func scanListInt64Value(d *subDecode) {
	off, v := scanInt64(d.str[d.scan:])
	d.scan += off
	d.pl.bufI64 = append(d.pl.bufI64, v)
}

// uint +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUintValue(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint(fieldPtrDeep(d), v)
	}
}

func scanObjUintValue(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint(fieldPtr(d), v)
}

func scanArrUintValue(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint(arrItemPtr(d), v)
}

func scanListUintValue(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	d.pl.bufUint = append(d.pl.bufUint, uint(v))
}

// uint8 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint8Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint8(fieldPtrDeep(d), v)
	}
}

func scanObjUint8Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint8(fieldPtr(d), v)
}

func scanArrUint8Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint8(arrItemPtr(d), v)
}

func scanListUint8Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	if v > math.MaxUint8 {
		panic(errInfinity)
	}
	d.pl.bufU8 = append(d.pl.bufU8, uint8(v))
}

// uint16 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint16Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint16(fieldPtrDeep(d), v)
	}
}

func scanObjUint16Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint16(fieldPtr(d), v)
}

func scanArrUint16Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint16(arrItemPtr(d), v)
}

func scanListUint16Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	if v > math.MaxUint16 {
		panic(errInfinity)
	}
	d.pl.bufU16 = append(d.pl.bufU16, uint16(v))
}

// uint32 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint32Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint32(fieldPtrDeep(d), v)
	}
}

func scanObjUint32Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint32(fieldPtr(d), v)
}

func scanArrUint32Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint32(arrItemPtr(d), v)
}

func scanListUint32Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	if v > math.MaxUint32 {
		panic(errInfinity)
	}
	d.pl.bufU32 = append(d.pl.bufU32, uint32(v))
}

// uint64 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrUint64Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint64(fieldPtrDeep(d), v)
	}
}

func scanObjUint64Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint64(fieldPtr(d), v)
}

func scanArrUint64Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	bindUint64(arrItemPtr(d), v)
}

func scanListUint64Value(d *subDecode) {
	off, v := scanUint64(d.str[d.scan:])
	d.scan += off
	d.pl.bufU64 = append(d.pl.bufU64, v)
}

// float32
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrFloat32Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanFloat32(d.str[d.scan:])
		d.scan += off
		bindFloat32(fieldPtrDeep(d), v)
	}
}

func scanObjFloat32Value(d *subDecode) {
	off, v := scanFloat32(d.str[d.scan:])
	d.scan += off
	bindFloat32(fieldPtr(d), v)
}

func scanArrFloat32Value(d *subDecode) {
	off, v := scanFloat32(d.str[d.scan:])
	d.scan += off
	bindFloat32(arrItemPtr(d), v)
}

func scanListFloat32Value(d *subDecode) {
	off, v := scanFloat32(d.str[d.scan:])
	d.scan += off
	d.pl.bufF32 = append(d.pl.bufF32, v)
}

// float64
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrFloat64Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanFloat64(d.str[d.scan:])
		d.scan += off
		bindFloat64(fieldPtrDeep(d), v)
	}
}

func scanObjFloat64Value(d *subDecode) {
	off, v := scanFloat64(d.str[d.scan:])
	d.scan += off
	bindFloat64(fieldPtr(d), v)
}

func scanArrFloat64Value(d *subDecode) {
	off, v := scanFloat64(d.str[d.scan:])
	d.scan += off
	bindFloat64(arrItemPtr(d), v)
}

func scanListFloat64Value(d *subDecode) {
	off, v := scanFloat64(d.str[d.scan:])
	d.scan += off
	d.pl.bufF64 = append(d.pl.bufF64, v)
}

// string +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrStrValue(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, str := scanString(d.str[d.scan:])
		d.scan += off
		bindString(fieldPtrDeep(d), str)
	}
}

func scanObjStrValue(d *subDecode) {
	off, str := scanString(d.str[d.scan:])
	d.scan += off
	bindString(fieldPtr(d), str)
}

func scanArrStrValue(d *subDecode) {
	off, str := scanString(d.str[d.scan:])
	d.scan += off
	bindString(arrItemPtr(d), str)
}

func scanListStrValue(d *subDecode) {
	off, str := scanString(d.str[d.scan:])
	d.scan += off
	d.pl.bufStr = append(d.pl.bufStr, str)
}

// []byte +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjBytesValue(d *subDecode) {
	off, bs := scanBytes(d.str[d.scan:])
	d.scan += off
	bindBytes(fieldPtr(d), bs)
}

// time.Time +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjTimeValue(d *subDecode) {
	switch d.str[d.scan] {
	case '"':
		start := d.scan + 1
		//d.scanQuoteStr()
		bindTime(fieldPtr(d), d.str[start:d.scan-1])
	default:
	}
}

func scanObjPtrTimeValue(d *subDecode) {
	switch d.str[d.scan] {
	case '"':
		start := d.scan + 1
		//d.scanQuoteStr()
		bindTime(fieldPtrDeep(d), d.str[start:d.scan-1])
	default:
		fieldSetNil(d)
	}
}

func scanArrTimeValue(d *subDecode) {
	v := ""
	switch d.str[d.scan] {
	case '"':
		start := d.scan + 1
		//d.scanQuoteStr()
		v = d.str[start : d.scan-1]
	default:
	}
	bindTime(arrItemPtr(d), v)
}

func scanListTimeValue(d *subDecode) {
	v := false
	switch d.str[d.scan] {
	case 't':
		//d.skipTrue()
		v = true
	case 'f':
		//d.skipFalse()
	default:
		//d.skipNull()
		d.pl.nulPos = append(d.pl.nulPos, len(d.pl.bufBol))
	}
	d.pl.bufBol = append(d.pl.bufBol, v)
}

// bool +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrBoolValue(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		off, v := scanBool(d.str[d.scan:])
		d.scan += off
		bindBool(fieldPtrDeep(d), v)
	}
}

func scanObjBoolValue(d *subDecode) {
	off, v := scanBool(d.str[d.scan:])
	d.scan += off
	bindBool(fieldPtr(d), v)
}

func scanArrBoolValue(d *subDecode) {
	off, v := scanBool(d.str[d.scan:])
	d.scan += off
	bindBool(arrItemPtr(d), v)
}

func scanListBoolValue(d *subDecode) {
	off, v := scanBool(d.str[d.scan:])
	d.scan += off
	d.pl.bufBol = append(d.pl.bufBol, v)
}

// any +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjAnyValue(d *subDecode) {
}

func scanObjPtrAnyValue(d *subDecode) {
}

func scanArrAnyValue(d *subDecode) {
}

func scanListAnyValue(d *subDecode) {
}

// Dest is just a base type value
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanJustBaseValue(d *subDecode) {
	// NOTE：只能是数值类型
	switch d.dm.itemKind {
	case reflect.Int:
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt(d.dstPtr, v)
	case reflect.Int8:
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt8(d.dstPtr, v)
	case reflect.Int16:
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt16(d.dstPtr, v)
	case reflect.Int32:
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt32(d.dstPtr, v)
	case reflect.Int64:
		off, v := scanInt64(d.str[d.scan:])
		d.scan += off
		bindInt64(d.dstPtr, v)
	case reflect.Uint:
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint(d.dstPtr, v)
	case reflect.Uint8:
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint8(d.dstPtr, v)
	case reflect.Uint16:
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint16(d.dstPtr, v)
	case reflect.Uint32:
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint32(d.dstPtr, v)
	case reflect.Uint64:
		off, v := scanUint64(d.str[d.scan:])
		d.scan += off
		bindUint64(d.dstPtr, v)
	case reflect.Float32:
		off, v := scanFloat32(d.str[d.scan:])
		d.scan += off
		bindFloat32(d.dstPtr, v)
	case reflect.Float64:
		off, v := scanFloat64(d.str[d.scan:])
		d.scan += off
		bindFloat64(d.dstPtr, v)
	default:
		panic(errValueType)
	}
}
