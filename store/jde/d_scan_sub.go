package jde

import (
	"reflect"
	"unsafe"
)

// array & slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func arrItemPtr(sd *subDecode) unsafe.Pointer {
	return unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.arrItemBytes))
}

// 前提：sd.str 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanList() {
	// A. 可能需要用到缓冲池记录临时数据
	sd.resetListPool()

	// B. 根据目标值类型，直接匹配，提高性能
	sd.scanListItems()

	// C. 将解析好的数据一次性绑定到对象上
	sd.flushListPool()
}

func (sd *subDecode) scanListItems() {
	pos := sd.scan

	pos++
	for isBlankChar[sd.str[pos]] {
		pos++
	}
	c := sd.str[pos]
	if c == ',' {
		goto errChar
	}

	for {
		// 不用switch, 比较顺序相对比较明确
		if c == ',' {
			pos++
		} else if c == ']' {
			// 数组多余的部分需要重置成类型零值
			if sd.arrIdx < sd.dm.arrLen {
				sd.resetArrLeftItems()
			}
			sd.scan = pos + 1
			return
		} else if sd.arrIdx > 0 {
			goto errChar
		}

		for isBlankChar[sd.str[pos]] {
			pos++
		}

		sd.scan = pos
		if sd.skipValue {
			sd.skipOneValue()
		} else {
			sd.dm.listItemDec(sd)
			if sd.dm.isArray {
				sd.arrIdx++
				if sd.arrIdx >= sd.dm.arrLen {
					sd.skipValue = true
				}
			}
		}
		pos = sd.scan

		for isBlankChar[sd.str[pos]] {
			pos++
		}
		c = sd.str[pos]
	}

errChar:
	sd.scan = pos
	panic(errChar)
}

// struct & map & gson
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func fieldPtr(sd *subDecode) unsafe.Pointer {
	return unsafe.Pointer(sd.dstPtr + sd.dm.ss.FieldsAttr[sd.keyIdx].Offset)
}

func fieldMixPtr(sd *subDecode) unsafe.Pointer {
	fa := &sd.dm.ss.FieldsAttr[sd.keyIdx]
	ptr := unsafe.Pointer(sd.dstPtr + fa.Offset)

	if *(*unsafe.Pointer)(ptr) == nil {
		var newPtr unsafe.Pointer

		switch fa.Kind {
		case reflect.Map:
			newPtr = reflect.MakeMap(fa.Type).UnsafePointer()
		case reflect.Slice:
			newPtr = reflect.MakeSlice(fa.Type, 0, 0).UnsafePointer()
		default:
			newPtr = reflect.New(fa.Type).UnsafePointer()
		}

		//newPtr := rVal.UnsafePointer()
		*(*unsafe.Pointer)(ptr) = newPtr
		return newPtr
	}
	return *(*unsafe.Pointer)(ptr)
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
		sd.dm.kvPairDec(sd, key)
		pos = sd.scan
	}

errChar:
	sd.scan = pos
	panic(errChar)
}
