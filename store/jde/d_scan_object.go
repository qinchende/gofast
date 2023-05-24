package jde

import (
	"reflect"
	"unsafe"
)

func fieldPtr(sd *subDecode) unsafe.Pointer {
	return unsafe.Pointer(sd.dstPtr + sd.dm.ss.FieldsAttr[sd.keyIdx].Offset)
}

func fieldPtrDeep(sd *subDecode) unsafe.Pointer {
	fa := &sd.dm.ss.FieldsAttr[sd.keyIdx]
	ptr := unsafe.Pointer(sd.dstPtr + fa.Offset)

	ptrLevel := fa.PtrLevel
	for ptrLevel > 1 {
		if *(*unsafe.Pointer)(ptr) == nil {
			tpPtr := unsafe.Pointer(new(unsafe.Pointer))
			*(*unsafe.Pointer)(ptr) = tpPtr
			ptr = tpPtr
		} else {
			ptr = *(*unsafe.Pointer)(ptr)
		}

		ptrLevel--
	}

	if *(*unsafe.Pointer)(ptr) == nil {
		rVal := reflect.New(fa.Type)

		newPtr := rVal.UnsafePointer()
		*(*unsafe.Pointer)(ptr) = newPtr
		return newPtr
	}
	return *(*unsafe.Pointer)(ptr)
}

func fieldSetNil(sd *subDecode) {
	*(*unsafe.Pointer)(fieldPtr(sd)) = nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 { 字符后面的字符串
// 返回 } 后面一个字符的 index
func (sd *subDecode) scanObject() {
	first := true
	pos := sd.scan

	pos++
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
			goto errChar
		}

	scanKVPair:
		// A: 找 key 字符串
		start := pos
		if sd.str[start] != '"' {
			goto errChar
		}

		sd.scan = pos
		slash := sd.scanQuoteStr()
		pos = sd.scan

		var key string
		if slash {
			key = sd.str[start+1 : sd.unescapeEnd()]
		} else {
			key = sd.str[start+1 : pos-1]
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
			goto errChar
		}

		// C: 找 value string，然后绑定
		sd.scan = pos
		sd.dm.objValDec(sd, key)
		pos = sd.scan
	}

errChar:
	sd.scan = pos
	panic(errChar)
}

func scanStructValue(sd *subDecode, key string) {
	if sd.keyIdx = sd.dm.ss.ColumnIndex(key); sd.keyIdx < 0 {
		sd.skipValue = true
		sd.skipOneValue()
	} else {
		sd.dm.fieldsDec[sd.keyIdx](sd) // TODO: 要根据目标值类型，来解析
	}
}

//func scanMapValue(sd *subDecode, key string) {
//	//sd.dm.mp.Set(key, sub.mp)
//}

func scanGsonValue(sd *subDecode, key string) {
	if sd.keyIdx = sd.gr.KeyIndex(key); sd.keyIdx < 0 {
		sd.skipValue = true
		sd.skipOneValue()
	} else {
		//sd.dm.gr.SetByIndex()
		//sd.dm.gr.SetStringByIndex()
	}
}
