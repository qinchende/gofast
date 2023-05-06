package jde

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"reflect"
)

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) scanJson() (err errType) {
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
		sd.scanJsonEnd('}')
		return
	case '[':
		sd.scanJsonEnd(']')
		return
	case 'n':
		sd.scanJsonEnd('l')
		return
	}
	return errJson
}

// 只支持 } ] l 三个字符判断
func (sd *subDecode) scanJsonEnd(ch byte) {
	// 去掉尾部的空白字符
	for i := len(sd.str) - 1; i > 0; i-- {
		if !isBlankChar[sd.str[i]] {
			if sd.str[i] != ch {
				sd.scan = i
				panic(errChar)
			}
			sd.str = sd.str[:i+1]
			break
		}
	}

	if ch == '}' {
		sd.scan++
		sd.scanObject()
	} else if ch == ']' {
		sd.resetListPool()
		sd.scan++
		sd.scanList()
		sd.flushListPool()
	} else {
		sd.skipNull()
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 { 字符后面的字符串
// 返回 } 后面一个字符的 index
func (sd *subDecode) scanObject() {
	first := true
	for {
		for isBlankChar[sd.str[sd.scan]] {
			sd.scan++
		}

		switch c := sd.str[sd.scan]; c {
		case '}':
			sd.scan++
			return
		case ',':
			sd.scan++
			for isBlankChar[sd.str[sd.scan]] {
				sd.scan++
			}
			goto scanKVPair
		default:
			if first {
				first = false
				goto scanKVPair
			}
			panic(errChar)
		}

	scanKVPair:
		sd.scanKVItem()
	}
}

// 必须是k:v, ...形式。不能为空，而且前面空字符已跳过，否则错误
func (sd *subDecode) scanKVItem() {
	// A: 找 key 字符串
	start := sd.scan
	slash := sd.scanQuoteString()
	if slash {
		//sd.key = sd.unescapeString(start, sd.scan)
		sd.key = sd.str[start+1 : sd.unescapeEnd()]
	} else {
		sd.key = sd.str[start+1 : sd.scan-1]
	}

	// B: 跳过冒号
	for isBlankChar[sd.str[sd.scan]] {
		sd.scan++
	}
	if sd.str[sd.scan] == ':' {
		sd.scan++
		for isBlankChar[sd.str[sd.scan]] {
			sd.scan++
		}
	} else {
		panic(errChar)
	}

	// C: 找 value string，然后绑定
	sd.checkSkip()
	sd.scanObjValue()
}

func (sd *subDecode) scanSubObject() {
	sub := subDecode{
		str:       sd.str,
		scan:      sd.scan,
		skipTotal: sd.skipValue,
	}

	if sd.gr != nil {
		// TODO: 无法为子对象提供目标值，只能返回字符串
		sub.skipTotal = true
	} else {
		sd.skipValue = true
		*sub.mp = make(cst.KV)
		sd.mp.Set(sd.key, sub.mp)
	}

	sub.scanObject()
	//if err < 0 {
	//	sd.scan = sub.scan
	//	return
	//}

	if sd.gr != nil && sd.skipValue == false {
		val := sd.str[sd.scan-1 : sub.scan]
		// TODO: 这里要重新规划一下
		sd.gr.SetString(sd.key, val)
	}
	sd.scan = sub.scan
	return
}

//func (sd *subDecode) scanSubArray(key string) (val string, err int) {
//	sub := subDecode{
//		str:       sd.str,
//		scan:      sd.scan,
//		skipTotal: sd.skipValue,
//	}
//
//	if sd.gr != nil {
//		// TODO: 无法为子对象提供目标值，只能返回字符串
//		sub.skipTotal = true
//	} else {
//		sd.skipValue = true
//	}
//
//	err = sub.scanList()
//	if err < 0 {
//		sd.scan = sub.scan
//		return
//	}
//
//	if sd.gr != nil {
//		if sd.skipValue == false {
//			val = sd.str[sd.scan-1 : sub.scan]
//		}
//	} else {
//		//sd.mp.Set(key, sub.list)
//	}
//	sd.scan = sub.scan
//	return
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanList() {
	if !sd.isList {
		panic(errList)
	}

	// 根据目标值类型，直接匹配，提高效率
	switch {
	case isNumKind(sd.dm.itemKind) == true:
		sd.scanArrItems(sd.scanNumValue)
	case sd.dm.itemKind == reflect.String:
		sd.scanArrItems(sd.scanStrVal)
	case sd.dm.itemKind == reflect.Bool:
		sd.scanArrItems(sd.scanBoolVal)
	default:
		sd.scanArrItems(sd.scanObjValue)
	}
}

func (sd *subDecode) scanArrItems(scanValue func()) {
	first := true
	for {
		pos := sd.scan
		for isBlankChar[sd.str[pos]] {
			pos++
		}

		switch c := sd.str[pos]; c {
		case ']':
			sd.scan = pos + 1
			return
		case ',':
			pos++
			for isBlankChar[sd.str[pos]] {
				pos++
			}
		default:
			if first {
				first = false
			} else {
				sd.scan = pos
				panic(errChar)
			}
		}

		sd.scan = pos
		scanValue()
	}
}

func (sd *subDecode) scanStrVal() {
	if sd.skipValue {
		sd.skipQuoteString()
		return
	}

	start := sd.scan
	slash := sd.scanQuoteString()
	if slash {
		//sd.bindString(sd.unescapeString(start, sd.scan))
		sd.bindString(sd.str[start+1 : sd.unescapeEnd()])
	} else {
		sd.bindString(sd.str[start+1 : sd.scan-1])
	}
}

func (sd *subDecode) scanBoolVal() {
	switch sd.str[sd.scan] {
	case 't':
		sd.skipTrue()
		if sd.skipValue {
			return
		}
		sd.bindBool(true)
	case 'f':
		sd.skipFalse()
		if sd.skipValue {
			return
		}
		sd.bindBool(false)
	default:
		panic(errBool)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) scanQuoteString() (slash bool) {
	pos := sd.scan
	if sd.str[pos] != '"' {
		panic(errChar)
	}

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
				sd.pl.escPos = sd.pl.escPos[0:0]
			}
			sd.pl.escPos = append(sd.pl.escPos, pos)
			pos++
			//c = sd.str[pos]
			//if c < ' ' {
			//	sd.scan = pos
			//	panic(errChar)
			//}
		}
	}
}

func (sd *subDecode) skipQuoteString() {
	pos := sd.scan
	if sd.str[pos] != '"' {
		panic(errChar)
	}

	for {
		pos++
		switch c := sd.str[pos]; {
		case c == '"':
			sd.scan = pos + 1
			return
		case c == '\\':
			pos++
		}
	}
}

func (sd *subDecode) scanObjValue() {
	switch sd.str[sd.scan] {
	case '{':
		sd.scan++
		sd.scanSubObject()
	case '[':
		sd.scan++
		//err = sd.scanSubArray()
	case '"':
		sd.scanStrVal()
	default:
		sd.scanNoQuoteValue()
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 匹配一个数值，不带有正负号前缀。
// 0.234 | 234.23 | 23424 | 3.8e+07
func (sd *subDecode) scanNumValue() {
	pos := sd.scan
	start := pos
	var hasDot, needNum bool

	c := sd.str[pos]
	if c == '-' {
		pos++
		c = sd.str[pos]
	}

	//if c < '0' || c > '9' {
	//	panic(errNumberFmt)
	//}

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
	//pos++
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
			needNum = false
		}
	}

	if needNum {
		panic(errNumberFmt)
	}

over:
	sd.scan = pos
	if sd.skipValue {
		return
	}
	sd.bindNumber(sd.str[start:pos])
}

func (sd *subDecode) scanNoQuoteValue() {
	switch c := sd.str[sd.scan]; {
	case (c >= '0' && c <= '9') || c == '-':
		sd.scanNumValue() // 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7 | -0.3 | -3.7E-7
	case c == 'f':
		sd.skipFalse()
		if sd.skipValue {
			return
		}
		sd.bindBool(false)
	case c == 't':
		sd.skipTrue()
		if sd.skipValue {
			return
		}
		sd.bindBool(true)
	case c == 'n':
		sd.skipNull()
		if sd.skipValue {
			return
		}
		sd.bindNull()
	default:
		panic(errValue)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) skipNull() {
	s := sd.scan + 1
	if sd.str[s:s+3] == "ull" {
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
