package jde

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"unsafe"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// scan string
func (sd *subDecode) scanQuoteStr() (slash bool) {
	pos := sd.scan
	for {
		pos++

		switch c := sd.str[pos]; {
		//case c < ' ':
		//	sd.scan = pos
		//	panic(errChar)
		case c == '"':
			sd.scan = pos + 1
			return
		case c == '\\':
			if !slash {
				slash = true
				sd.escPos = sd.escPos[0:0]
			}
			sd.escPos = append(sd.escPos, pos)
			pos++
			//c = sd.str[pos]
			//if c < ' ' {
			//	sd.scan = pos
			//	panic(errChar)
			//}
		}
	}
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
	}
	sd.pl.bufBol = append(sd.pl.bufBol, v)
}

// any +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanObjAnyValue(sd *subDecode) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubObject(rfTypeOfKV, unsafe.Pointer(&newMap))
		bindAny(fieldPtr(sd), newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubList(rfTypeOfList, unsafe.Pointer(&newList))
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
		//sd.scanSubObject()
	case c == '[':
		//err = sd.scanSubArray()
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
		//sd.scanSubObject()
	case c == '[':
		//err = sd.scanSubArray()
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
		//nsd := newSubDecode(sd)
		//nsd.scanObject()
	case c == '[':
		//err = sd.scanSubArray()
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

//// struct +++++
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func scanObjStructValue(sd *subDecode) {
//	switch c := sd.str[sd.scan]; {
//	case c == '{':
//		sd.scanSubObject()
//	default:
//		sd.skipNull()
//	}
//}

// map +++++
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanMapAnyValue(sd *subDecode, k string) {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		newMap := make(cst.KV)
		sd.scanSubObject(rfTypeOfKV, unsafe.Pointer(&newMap))
		sd.mp.Set(k, newMap)
	case c == '[':
		newList := make([]any, 0)
		sd.scanSubList(rfTypeOfList, unsafe.Pointer(&newList))
		sd.mp.Set(k, newList)
	case c == '"':
		start := sd.scan + 1
		slash := sd.scanQuoteStr()
		if slash {
			sd.mp.Set(k, sd.str[start:sd.unescapeEnd()])
		} else {
			sd.mp.Set(k, sd.str[start:sd.scan-1])
		}
	case c >= '0' && c <= '9', c == '-':
		if start := sd.scanNumValue(); start > 0 {
			sd.mp.Set(k, lang.ParseFloat(sd.str[start:sd.scan]))
		}
	case c == 't':
		sd.skipTrue()
		sd.mp.Set(k, true)
	case c == 'f':
		sd.skipFalse()
		sd.mp.Set(k, false)
	default:
		sd.skipNull()
		sd.mp.Set(k, nil)
	}
}
