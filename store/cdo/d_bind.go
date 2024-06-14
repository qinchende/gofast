package cdo

import (
	"github.com/qinchende/gofast/core/cst"
	"math"
	"time"
	"unsafe"
)

// int
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//

func bindInt(p unsafe.Pointer, v int64) {
	*(*int)(p) = int(v)
}

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
	if v > MaxUint08 {
		panic(errInfinity)
	}
	*(*uint8)(p) = uint8(v)
}

func bindUint16(p unsafe.Pointer, v uint64) {
	if v > MaxUint16 {
		panic(errInfinity)
	}
	*(*uint16)(p) = uint16(v)
}

func bindUint32(p unsafe.Pointer, v uint64) {
	if v > MaxUint32 {
		panic(errInfinity)
	}
	*(*uint32)(p) = uint32(v)
}

func bindUint64(p unsafe.Pointer, v uint64) {
	*(*uint64)(p) = v
}

// float
func bindF32(p unsafe.Pointer, v float32) {
	*(*float32)(p) = v
}

func bindF64(p unsafe.Pointer, v float64) {
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

// float32
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjPtrF32Value(d *subDecode) {
	if d.str[d.scan] == FixNil {
		fieldSetNil(d)
		d.scan++
	} else {
		v := scanF32Val(d.str[d.scan:])
		d.scan += 4
		bindF32(fieldPtrDeep(d), v)
	}
}

func scanObjF32Value(d *subDecode) {
	v := scanF32Val(d.str[d.scan:])
	d.scan += 4
	bindF32(fieldPtr(d), v)
}

// float64
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

//func scanListFloat64Value(d *subDecode) {
//	off, v := scanF64Val(d.str[d.scan:])
//	d.scan += off
//	d.pl.bufF64 = append(d.pl.bufF64, v)
//}

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

//func scanListTimeValue(d *subDecode) {
//	v := false
//	switch d.str[d.scan] {
//	case 't':
//		//d.skipTrue()
//		v = true
//	case 'f':
//		//d.skipFalse()
//	default:
//		//d.skipNull()
//		d.pl.nulPos = append(d.pl.nulPos, len(d.pl.bufBol))
//	}
//	d.pl.bufBol = append(d.pl.bufBol, v)
//}

// bool +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

//func scanListBoolValue(d *subDecode) {
//	off, v := scanBoolVal(d.str[d.scan:])
//	d.scan += off
//	d.pl.bufBol = append(d.pl.bufBol, v)
//}

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
