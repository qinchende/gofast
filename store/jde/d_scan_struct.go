package jde

import (
	"github.com/qinchende/gofast/skill/lang"
)

func fieldPtr(sd *subDecode) uintptr {
	return sd.dstPtr + sd.dm.ss.FieldsOffset[sd.keyIdx]
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 { 字符后面的字符串
// 返回 } 后面一个字符的 index
func (sd *subDecode) scanObject() {
	first := true
	pos := sd.scan
	for {
		if isBlankChar[sd.str[pos]] {
			pos++
			continue
		}

		switch c := sd.str[pos]; c {
		case '}':
			sd.scan = pos + 1
			return
		case ',':
			pos++
			for isBlankChar[sd.str[pos]] {
				pos++
			}
			goto scanKVPair
		default:
			if first {
				first = false
				goto scanKVPair
			}
			sd.scan = pos
			panic(errChar)
		}

	scanKVPair:
		// A: 找 key 字符串
		start := pos
		if sd.str[start] != '"' {
			sd.scan = pos
			panic(errChar)
		}

		sd.scan = pos
		slash := sd.scanQuoteString()
		pos = sd.scan

		if slash {
			sd.key = sd.str[start+1 : sd.unescapeEnd()]
		} else {
			sd.key = sd.str[start+1 : pos-1]
		}

		// B: 跳过冒号
		for isBlankChar[sd.str[pos]] {
			pos++
		}
		if sd.str[pos] == ':' {
			pos++
			for isBlankChar[sd.str[pos]] {
				pos++
			}
		} else {
			sd.scan = pos
			panic(errChar)
		}

		sd.scan = pos
		// C: 找 value string，然后绑定
		// sd.checkSkip() // 确定key是否存在，以及索引位置
		if sd.keyIdx = sd.dm.ss.ColumnIndex(sd.key); sd.keyIdx < 0 {
			sd.skipOneValue()
		} else {
			// TODO: 要根据目标值类型，来解析 ++++
			//sd.scanOneValue()
			sd.dm.ssFunc[sd.keyIdx](sd)
			// +++++++++++++++++++++++++++++++++++
		}
		pos = sd.scan
	}
}

// int64
// +++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjIntValue(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt(sd.dm.nextPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanObjInt8Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt8(sd.dm.nextPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

func scanObjInt16Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt16(sd.dm.nextPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}
func scanObjInt32Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt32(sd.dm.nextPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}
func scanObjInt64Value(sd *subDecode) {
	if start := sd.scanIntValue(); start > 0 {
		bindInt64(sd.dm.nextPtr(sd), lang.ParseInt(sd.str[start:sd.scan]))
	}
}

// uint64
// +++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjUintValue(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint(sd.dm.nextPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjUint8Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint8(sd.dm.nextPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjUint16Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint16(sd.dm.nextPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjUint32Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint32(sd.dm.nextPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

func scanObjUint64Value(sd *subDecode) {
	if start := sd.scanUintValue(); start > 0 {
		bindUint64(sd.dm.nextPtr(sd), lang.ParseUint(sd.str[start:sd.scan]))
	}
}

// float64
// +++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjFloat32Value(sd *subDecode) {
	if start := sd.scanNumValue(); start > 0 {
		bindFloat32(sd.dm.nextPtr(sd), lang.ParseFloat(sd.str[start:sd.scan]))
	}
}

func scanObjFloat64Value(sd *subDecode) {
	if start := sd.scanNumValue(); start > 0 {
		bindFloat64(sd.dm.nextPtr(sd), lang.ParseFloat(sd.str[start:sd.scan]))
	}
}

// float64
// +++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjStrValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case '"':
		start := sd.scan + 1
		slash := sd.scanQuoteString()
		if slash {
			bindString(sd.dm.nextPtr(sd), sd.str[start:sd.unescapeEnd()])
		} else {
			bindString(sd.dm.nextPtr(sd), sd.str[start:sd.scan-1])
		}
	default:
		sd.skipNull()
	}
}

func scanObjBolValue(sd *subDecode) {
	switch sd.str[sd.scan] {
	case 't':
		sd.skipTrue()
		bindBool(sd.dm.nextPtr(sd), true)
	case 'f':
		sd.skipFalse()
		bindBool(sd.dm.nextPtr(sd), false)
	default:
		sd.skipNull()
	}
}

func scanObjAnyValue(sd *subDecode) {
	sd.scanBoolValue()
}
