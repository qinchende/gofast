package jde

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"math"
	"reflect"
	"unsafe"
)

// 解析成 int64 ，可以是正负整数
func (sd *subDecode) scanIntValue() int {
	pos := sd.scan
	start := pos

	c := sd.str[pos]
	if c == '-' {
		pos++
		c = sd.str[pos]
	}
	if c == '0' {
		pos++
		goto over
	}
	for {
		if c < '0' || c > '9' {
			break
		}
		pos++
		c = sd.str[pos]
	}
over:
	sd.scan = pos
	// 还剩下最后一种可能：null +++
	if start == pos {
		sd.skipNull()
		return -1
	}
	return start
}

func (sd *subDecode) scanIntMust() int {
	max := len(sd.str)
	pos := sd.scan
	start := pos

	c := sd.str[pos]
	if c == '-' {
		pos++
		if pos == max {
			goto over
		}
		c = sd.str[pos]
	}
	if c == '0' {
		pos++
		goto over
	}
	for {
		if c < '0' || c > '9' {
			break
		}
		pos++
		if pos == max {
			goto over
		}
		c = sd.str[pos]
	}
over:
	sd.scan = pos
	return start
}

// 解析成 uint64 ，只能是正整数
func (sd *subDecode) scanUintValue() int {
	pos := sd.scan
	start := pos

	c := sd.str[pos]
	if c == '0' {
		pos++
		goto over
	}
	for {
		if c < '0' || c > '9' {
			break
		}
		pos++
		c = sd.str[pos]
	}
over:
	sd.scan = pos
	// 还剩下最后一种可能：null
	if start == pos {
		sd.skipNull()
		return -1
	}
	return start
}

func (sd *subDecode) scanUintMust() int {
	max := len(sd.str)
	pos := sd.scan
	start := pos

	c := sd.str[pos]
	if c == '0' {
		pos++
		goto over
	}
	for {
		if c < '0' || c > '9' {
			break
		}
		pos++
		if pos == max {
			goto over
		}
		c = sd.str[pos]
	}
over:
	sd.scan = pos
	return start
}

// 匹配一个数值
// 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7 | -0.3 | -3.7E-7
func (sd *subDecode) scanNumValue() int {
	pos := sd.scan
	start := pos
	var hasDot, needNum bool

	c := sd.str[pos]
	if c == '-' {
		pos++
		c = sd.str[pos]
	}
	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
	if c == '0' {
		pos++
		c = sd.str[pos]

		switch c {
		case '.', 'e', 'E':
			goto loopNum
		default:
			goto over
		}
	}
	needNum = true

loopNum:
	for {
		c = sd.str[pos]
		pos++

		if c == '.' {
			if hasDot == true {
				panic(errNumberFmt)
			}
			hasDot = true
			needNum = true
		} else if c == 'e' || c == 'E' {
			if needNum {
				panic(errNumberFmt)
			}
			needNum = true

			c = sd.str[pos]
			if c == '-' || c == '+' {
				pos++
			}
			for {
				if c = sd.str[pos]; c < '0' || c > '9' {
					break loopNum
				} else {
					needNum = false
				}
				pos++
			}
		} else if c < '0' || c > '9' {
			pos--
			break
		} else {
			needNum = false // 到这里，字符肯定是数字
		}
	}

	if needNum {
		panic(errNumberFmt)
	}

over:
	sd.scan = pos
	// 还剩下最后一种可能：null
	if start == pos {
		sd.skipNull()
		return -1
	}
	return start
}

func (sd *subDecode) scanNumMust() int {
	max := len(sd.str)
	pos := sd.scan
	start := pos
	var hasDot, needNum bool

	c := sd.str[pos]
	if c == '-' {
		pos++
		if pos == max {
			goto over
		}
		c = sd.str[pos]
	}
	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
	if c == '0' {
		pos++
		if pos == max {
			goto over
		}
		c = sd.str[pos]

		switch c {
		case '.', 'e', 'E':
			goto loopNum
		default:
			goto over
		}
	}
	needNum = true

loopNum:
	for {
		if pos == max {
			goto over
		}
		c = sd.str[pos]
		pos++

		if c == '.' {
			if hasDot == true {
				panic(errNumberFmt)
			}
			hasDot = true
			needNum = true
		} else if c == 'e' || c == 'E' {
			if needNum {
				panic(errNumberFmt)
			}
			needNum = true

			if pos == max {
				goto over
			}
			c = sd.str[pos]
			if c == '-' || c == '+' {
				pos++
			}
			for {
				if pos == max {
					goto over
				}
				if c = sd.str[pos]; c < '0' || c > '9' {
					break loopNum
				} else {
					needNum = false
				}
				pos++
			}
		} else if c < '0' || c > '9' {
			pos--
			break
		} else {
			needNum = false // 到这里，字符肯定是数字
		}
	}

	if needNum {
		panic(errNumberFmt)
	}

over:
	sd.scan = pos
	return start
}

// int
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
func bindFloat32(p unsafe.Pointer, v float64) {
	if v < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
		panic(errInfinity)
	}
	*(*float32)(p) = float32(v)
}

func bindFloat64(p unsafe.Pointer, v float64) {
	*(*float64)(p) = v
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// int +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjIntValue(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanObjPtrIntValue(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt(fieldPtrDeep(sd), lang.ParseInt(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrIntValue(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
	}
	bindInt(arrItemPtr(sd), v)
}

func scanListIntValue(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufI64))
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int8 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt8Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt8(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanObjPtrInt8Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt8(fieldPtrDeep(sd), lang.ParseInt(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrInt8Value(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
	}
	bindInt8(arrItemPtr(sd), v)
}

func scanListInt8Value(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
		if v < math.MinInt8 || v > math.MaxInt8 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufI64))
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int16 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt16Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt16(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanObjPtrInt16Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt16(fieldPtrDeep(sd), lang.ParseInt(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrInt16Value(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
	}
	bindInt16(arrItemPtr(sd), v)
}

func scanListInt16Value(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
		if v < math.MinInt16 || v > math.MaxInt16 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufI64))
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int32 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt32Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt32(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanObjPtrInt32Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt32(fieldPtrDeep(sd), lang.ParseInt(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrInt32Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt32(arrItemPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanListInt32Value(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
		if v < math.MinInt32 || v > math.MaxInt32 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufI64))
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int64 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt64Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt64(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanObjPtrInt64Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt64(fieldPtrDeep(sd), lang.ParseInt(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrInt64Value(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
	}
	bindInt64(arrItemPtr(sd), v)
}

func scanListInt64Value(sd *subDecode) {
	v := int64(0)
	if start := sd.scanIntValue(); start > 0 {
		v = lang.ParseInt(sd.str[start:sd.scan])
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufI64))
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// uint +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjUintValue(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint(fieldPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjPtrUintValue(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint(fieldPtrDeep(sd), lang.ParseUint(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrUintValue(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
	}
	bindUint(arrItemPtr(sd), v)
}

func scanListUintValue(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufU64))
	}
	sd.pl.bufU64 = append(sd.pl.bufU64, v)
}

// uint8 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjUint8Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint8(fieldPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjPtrUint8Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint8(fieldPtrDeep(sd), lang.ParseUint(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrUint8Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
	}
	bindUint8(arrItemPtr(sd), v)
}

func scanListUint8Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
		if v > math.MaxUint8 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufU64))
	}
	sd.pl.bufU64 = append(sd.pl.bufU64, v)
}

// uint16 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjUint16Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint16(fieldPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjPtrUint16Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint16(fieldPtrDeep(sd), lang.ParseUint(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrUint16Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
	}
	bindUint16(arrItemPtr(sd), v)
}

func scanListUint16Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
		if v > math.MaxUint16 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufU64))
	}
	sd.pl.bufU64 = append(sd.pl.bufU64, v)
}

// uint32 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjUint32Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint32(fieldPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjPtrUint32Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint32(fieldPtrDeep(sd), lang.ParseUint(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrUint32Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
	}
	bindUint32(arrItemPtr(sd), v)
}

func scanListUint32Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
		if v > math.MaxUint32 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufU64))
	}
	sd.pl.bufU64 = append(sd.pl.bufU64, v)
}

// uint64 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjUint64Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint64(fieldPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjPtrUint64Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint64(fieldPtrDeep(sd), lang.ParseUint(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrUint64Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
	}
	bindUint64(arrItemPtr(sd), v)
}

func scanListUint64Value(sd *subDecode) {
	v := uint64(0)
	if start := sd.scanUintValue(); start > 0 {
		v = lang.ParseUint(sd.str[start:sd.scan])
		if v > math.MaxUint64 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufU64))
	}
	sd.pl.bufU64 = append(sd.pl.bufU64, v)
}

// float32
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjFloat32Value(sd *subDecode) {
	if start := sd.scanNumValue(); start > 0 {
		bindFloat32(fieldPtr(sd), lang.ParseFloat(sd.str[start:sd.scan]))
	}
}

func scanObjPtrFloat32Value(sd *subDecode) {
	if start := sd.scanNumValue(); start > 0 {
		bindFloat32(fieldPtrDeep(sd), lang.ParseFloat(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrFloat32Value(sd *subDecode) {
	v := float64(0)
	if start := sd.scanNumValue(); start > 0 {
		v = lang.ParseFloat(sd.str[start:sd.scan])
	}
	bindFloat32(arrItemPtr(sd), v)
}

func scanListFloat32Value(sd *subDecode) {
	v := float64(0)
	if start := sd.scanNumValue(); start > 0 {
		v = lang.ParseFloat(sd.str[start:sd.scan])
		if v < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
			panic(errInfinity)
		}
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufF64))
	}
	sd.pl.bufF64 = append(sd.pl.bufF64, v)
}

// float64
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjFloat64Value(sd *subDecode) {
	if start := sd.scanNumValue(); start > 0 {
		bindFloat64(fieldPtr(sd), lang.ParseFloat(sd.str[start:sd.scan]))
	}
}

func scanObjPtrFloat64Value(sd *subDecode) {
	if start := sd.scanNumValue(); start > 0 {
		bindFloat64(fieldPtrDeep(sd), lang.ParseFloat(sd.str[start:sd.scan]))
	} else {
		fieldSetNil(sd)
	}
}

func scanArrFloat64Value(sd *subDecode) {
	v := float64(0)
	if start := sd.scanNumValue(); start > 0 {
		v = lang.ParseFloat(sd.str[start:sd.scan])
	}
	bindFloat64(arrItemPtr(sd), v)
}

func scanListFloat64Value(sd *subDecode) {
	v := float64(0)
	if start := sd.scanNumValue(); start > 0 {
		v = lang.ParseFloat(sd.str[start:sd.scan])
	} else {
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufF64))
	}
	sd.pl.bufF64 = append(sd.pl.bufF64, v)
}

// string +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjStrValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			bindString(fieldPtr(sd), sd.str[start:sd.unescapeEnd()])
		} else {
			bindString(fieldPtr(sd), sd.str[start:sd.scan-1])
		}
	default:
		sd.skipNull()
	}
}

func scanObjPtrStrValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			bindString(fieldPtrDeep(sd), sd.str[start:sd.unescapeEnd()])
		} else {
			bindString(fieldPtrDeep(sd), sd.str[start:sd.scan-1])
		}
	default:
		sd.skipNull()
		fieldSetNil(sd)
	}
}

func scanArrStrValue(sd *subDecode) {
	v := ""
	switch sd.str[sd.scan] {
	case '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			v = sd.str[start:sd.unescapeEnd()]
		} else {
			v = sd.str[start : sd.scan-1]
		}
	default:
		sd.skipNull()
	}
	bindString(arrItemPtr(sd), v)
}

func scanListStrValue(sd *subDecode) {
	v := ""
	switch sd.str[sd.scan] {
	case '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			v = sd.str[start:sd.unescapeEnd()]
		} else {
			v = sd.str[start : sd.scan-1]
		}
	default:
		sd.skipNull()
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufStr))
	}
	sd.pl.bufStr = append(sd.pl.bufStr, v)
}

// bool +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjBoolValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case 't':
		sd.skipTrue()
		bindBool(fieldPtr(sd), true)
	case 'f':
		sd.skipFalse()
		bindBool(fieldPtr(sd), false)
	default:
		sd.skipNull()
	}
}

func scanObjPtrBoolValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case 't':
		sd.skipTrue()
		bindBool(fieldPtrDeep(sd), true)
	case 'f':
		sd.skipFalse()
		bindBool(fieldPtrDeep(sd), false)
	default:
		sd.skipNull()
		fieldSetNil(sd)
	}
}

func scanArrBoolValue(sd *subDecode) {
	v := false
	switch sd.str[sd.scan] {
	case 't':
		sd.skipTrue()
		v = true
	case 'f':
		sd.skipFalse()
	default:
		sd.skipNull()
	}
	bindBool(arrItemPtr(sd), v)
}

func scanListBoolValue(sd *subDecode) {
	v := false
	switch sd.str[sd.scan] {
	case 't':
		sd.skipTrue()
		v = true
	case 'f':
		sd.skipFalse()
	default:
		sd.skipNull()
		sd.pl.nulPos = append(sd.pl.nulPos, len(sd.pl.bufBol))
	}
	sd.pl.bufBol = append(sd.pl.bufBol, v)
}

// any +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjAnyValue(sd *subDecode) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubDecode(rfTypeOfKV, unsafe.Pointer(&newMap))
		bindAny(fieldPtr(sd), newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubDecode(rfTypeOfList, unsafe.Pointer(&newList))
		bindAny(fieldPtr(sd), newList)
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			bindAny(fieldPtr(sd), sd.str[start:sd.unescapeEnd()])
		} else {
			bindAny(fieldPtr(sd), sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		if start := sd.scanNumValue(); start > 0 {
			bindAny(fieldPtr(sd), lang.ParseFloat(sd.str[start:sd.scan]))
		}
	case c == 't':
		sd.skipTrue()
		bindAny(fieldPtr(sd), true)
	case c == 'f':
		sd.skipFalse()
		bindAny(fieldPtr(sd), false)
	default:
		sd.skipNull()
	}
}

func scanObjPtrAnyValue(sd *subDecode) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubDecode(rfTypeOfKV, unsafe.Pointer(&newMap))
		bindAny(fieldPtrDeep(sd), newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubDecode(rfTypeOfList, unsafe.Pointer(&newList))
		bindAny(fieldPtrDeep(sd), newList)
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			bindAny(fieldPtrDeep(sd), sd.str[start:sd.unescapeEnd()])
		} else {
			bindAny(fieldPtrDeep(sd), sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		if start := sd.scanNumValue(); start > 0 {
			bindAny(fieldPtrDeep(sd), lang.ParseFloat(sd.str[start:sd.scan]))
		}
	case c == 't':
		sd.skipTrue()
		bindAny(fieldPtrDeep(sd), true)
	case c == 'f':
		sd.skipFalse()
		bindAny(fieldPtrDeep(sd), false)
	default:
		sd.skipNull()
		fieldSetNil(sd)
	}
}

func scanArrAnyValue(sd *subDecode) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubDecode(rfTypeOfKV, unsafe.Pointer(&newMap))
		bindAny(arrItemPtr(sd), newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubDecode(rfTypeOfList, unsafe.Pointer(&newList))
		bindAny(arrItemPtr(sd), newList)
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			bindAny(arrItemPtr(sd), sd.str[start:sd.unescapeEnd()])
		} else {
			bindAny(arrItemPtr(sd), sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		if start := sd.scanNumValue(); start > 0 {
			bindAny(arrItemPtr(sd), lang.ParseFloat(sd.str[start:sd.scan]))
		}
	case c == 't':
		sd.skipTrue()
		bindAny(arrItemPtr(sd), true)
	case c == 'f':
		sd.skipFalse()
		bindAny(arrItemPtr(sd), false)
	default:
		sd.skipNull()
		bindAny(arrItemPtr(sd), nil)
	}
}

func scanListAnyValue(sd *subDecode) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubDecode(rfTypeOfKV, unsafe.Pointer(&newMap))
		sd.pl.bufAny = append(sd.pl.bufAny, newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubDecode(rfTypeOfList, unsafe.Pointer(&newList))
		sd.pl.bufAny = append(sd.pl.bufAny, newList)
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			sd.pl.bufAny = append(sd.pl.bufAny, sd.str[start:sd.unescapeEnd()])
		} else {
			sd.pl.bufAny = append(sd.pl.bufAny, sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		if start := sd.scanNumValue(); start > 0 {
			sd.pl.bufAny = append(sd.pl.bufAny, lang.ParseFloat(sd.str[start:sd.scan]))
		}
	case c == 't':
		sd.skipTrue()
		sd.pl.bufAny = append(sd.pl.bufAny, true)
	case c == 'f':
		sd.skipFalse()
		sd.pl.bufAny = append(sd.pl.bufAny, false)
	default:
		sd.skipNull()
		sd.pl.bufAny = append(sd.pl.bufAny, nil)
	}
}

// Dest is just a base type value
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanJustBaseValue(sd *subDecode) {
	switch c := sd.str[sd.scan]; {
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			bindString(sd.dstPtr, sd.str[start:sd.unescapeEnd()])
		} else {
			bindString(sd.dstPtr, sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		// NOTE：只能是数值类型
		switch sd.dm.itemKind {
		case reflect.Int:
			bindInt(sd.dstPtr, lang.ParseInt(sd.str[sd.scanIntMust():sd.scan]))
		case reflect.Int8:
			bindInt8(sd.dstPtr, lang.ParseInt(sd.str[sd.scanIntMust():sd.scan]))
		case reflect.Int16:
			bindInt16(sd.dstPtr, lang.ParseInt(sd.str[sd.scanIntMust():sd.scan]))
		case reflect.Int32:
			bindInt32(sd.dstPtr, lang.ParseInt(sd.str[sd.scanIntMust():sd.scan]))
		case reflect.Int64:
			bindInt64(sd.dstPtr, lang.ParseInt(sd.str[sd.scanIntMust():sd.scan]))
		case reflect.Uint:
			bindUint(sd.dstPtr, lang.ParseUint(sd.str[sd.scanUintMust():sd.scan]))
		case reflect.Uint8:
			bindUint8(sd.dstPtr, lang.ParseUint(sd.str[sd.scanUintMust():sd.scan]))
		case reflect.Uint16:
			bindUint16(sd.dstPtr, lang.ParseUint(sd.str[sd.scanUintMust():sd.scan]))
		case reflect.Uint32:
			bindUint32(sd.dstPtr, lang.ParseUint(sd.str[sd.scanUintMust():sd.scan]))
		case reflect.Uint64:
			bindUint64(sd.dstPtr, lang.ParseUint(sd.str[sd.scanUintMust():sd.scan]))
		case reflect.Float32:
			bindFloat32(sd.dstPtr, lang.ParseFloat(sd.str[sd.scanNumMust():sd.scan]))
		case reflect.Float64:
			bindFloat64(sd.dstPtr, lang.ParseFloat(sd.str[sd.scanNumMust():sd.scan]))
		default:
			panic(errValueType)
		}
	case c == 't':
		sd.skipTrue()
		bindBool(sd.dstPtr, true)
	case c == 'f':
		sd.skipFalse()
		bindBool(sd.dstPtr, false)
	default:
		sd.skipNull()
	}
}
