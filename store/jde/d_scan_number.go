package jde

import (
	"github.com/qinchende/gofast/skill/lang"
	"math"
)

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

// 匹配一个数值，对应于float类型
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

			c := sd.str[pos]
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

// int +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjIntValue(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
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
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int8 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt8Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt8(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
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
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int16 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt16Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt16(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
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
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int32 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt32Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt32(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
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
	}
	sd.pl.bufI64 = append(sd.pl.bufI64, v)
}

// int64 +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjInt64Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt64(fieldPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
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
	}
	sd.pl.bufF64 = append(sd.pl.bufF64, v)
}