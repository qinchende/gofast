package jde

import (
	"fmt"
)

// start decode json
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) scanStart() (err errType) {
	// 解析过程中异常，这里统一截获处理，返回解析错误编号
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				err = code
			} else {
				// 调试的时候打印错误信息
				fmt.Println(pic)
				err = errJson
			}
		}
	}()

	for isBlankChar[sd.str[sd.scan]] {
		sd.scan++
	}

	switch sd.str[sd.scan] {
	case '{':
		if sd.dm.isSuperKV || sd.dm.isStruct {
			sd.scanObject()
		} else {
			return errObject
		}
	case '[':
		if sd.dm.isList {
			sd.scanList()
		} else {
			return errList
		}
	default:
		sd.skipNull()
	}
	return
}

// array & slice
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (sd *subDecode) scanOneValue() {
//	switch c := sd.str[sd.scan]; {
//	case c == '{':
//		sd.scanSubObject()
//	case c == '[':
//		//err = sd.scanSubArray()
//	case c == '"':
//		sd.scanQuoteStrValue()
//	case c >= '0' && c <= '9', c == '-':
//		//sd.scanNumValue()
//	case c == 't':
//		sd.skipTrue()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBool(true)
//	case c == 'f':
//		sd.skipFalse()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBool(false)
//	default:
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBoolNull()
//	}
//}

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

// skip items
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) skipOneValue() {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		//sd.scanSubObject()
	case c == '[':
		//err = sd.scanSubArray()
	case c == '"':
		sd.skipQuoteStr()
	case c >= '0' && c <= '9', c == '-':
		sd.scanNumValue()
	case c == 't':
		sd.skipTrue()
	case c == 'f':
		sd.skipFalse()
	default:
		sd.skipNull()
	}
}

func (sd *subDecode) skipQuoteStr() {
	pos := sd.scan
	for {
		pos++
		switch c := sd.str[pos]; {
		case c == '"':
			sd.scan = pos + 1
			return
		case c == '\\':
			pos++ // 跳过 '\' 后面的一个字符
		}
	}
}

func (sd *subDecode) skipNull() {
	s := sd.scan
	if sd.str[s:s+4] == "null" {
		sd.scan += 4
		return
	}
	panic(errNull)
}

func (sd *subDecode) skipTrue() {
	s := sd.scan + 1
	if sd.str[s:s+3] == "rue" {
		sd.scan += 4
		return
	}
	panic(errBool)
}

func (sd *subDecode) skipFalse() {
	s := sd.scan + 1
	if sd.str[s:s+4] == "alse" {
		sd.scan += 5
		return
	}
	panic(errBool)
}
