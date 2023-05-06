package decode

//
//import (
//	"fmt"
//	"reflect"
//)
//
//// 采用尽最大努力解析出正确结果的策略
//// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
//func (sd *subDecode) scanJson(pos int) (err errType) {
//	// 解析过程中异常，这里统一截获处理，返回解析错误编号
//	defer func() {
//		if pic := recover(); pic != nil {
//			if code, ok := pic.(errType); ok {
//				err = code
//			} else {
//				// 调试的时候打印错误信息
//				fmt.Println(pic)
//				err = errJson
//			}
//		}
//	}()
//
//	for isBlankChar[sd.str[pos]] {
//		pos++
//	}
//
//	switch sd.str[pos] {
//	case '{':
//		sd.scanJsonEnd(pos, '}')
//		return
//	case '[':
//		sd.scanJsonEnd(pos, ']')
//		return
//	case 'n':
//		sd.scanJsonEnd(pos, 'l')
//		return
//	}
//	return errJson
//}
//
//// 只支持 } ] l 三个字符判断
//func (sd *subDecode) scanJsonEnd(pos int, ch byte) {
//	// 去掉尾部的空白字符
//	for i := len(sd.str) - 1; i > 0; i-- {
//		if !isBlankChar[sd.str[i]] {
//			if sd.str[i] != ch {
//				sd.scan = i
//				panic(errChar)
//			}
//			sd.str = sd.str[:i+1]
//			break
//		}
//	}
//
//	if ch == '}' {
//		pos++
//		sd.scanObject(pos)
//	} else if ch == ']' {
//		sd.resetListPool()
//		pos++
//		sd.scanList(pos)
//		sd.flushListPool()
//	} else {
//		sd.skipNull(pos)
//	}
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 前提：sd.str 肯定是 { 字符后面的字符串
//// 返回 } 后面一个字符的 index
//func (sd *subDecode) scanObject(pos int) int {
//	first := true
//	for {
//		for isBlankChar[sd.str[pos]] {
//			pos++
//		}
//
//		switch c := sd.str[pos]; c {
//		case '}':
//			pos++
//			return pos
//		case ',':
//			pos++
//			for isBlankChar[sd.str[pos]] {
//				pos++
//			}
//			goto scanKVPair
//		default:
//			if first {
//				first = false
//				goto scanKVPair
//			}
//			sd.scan = pos
//			panic(errChar)
//		}
//
//	scanKVPair:
//		pos = sd.scanKVItem(pos)
//	}
//}
//
//// 必须是k:v, ...形式。不能为空，而且前面空字符已跳过，否则错误
//func (sd *subDecode) scanKVItem(pos int) int {
//	// A: 找 key 字符串
//	start := pos
//	var slash bool
//	pos, slash = sd.scanQuoteString(pos)
//	if slash {
//		sd.key = sd.unescapeString(start, pos)
//	} else {
//		sd.key = sd.str[start+1 : pos-1]
//	}
//
//	// B: 跳过冒号
//	for isBlankChar[sd.str[pos]] {
//		pos++
//	}
//	if sd.str[pos] == ':' {
//		pos++
//		for isBlankChar[sd.str[pos]] {
//			pos++
//		}
//	} else {
//		sd.scan = pos
//		panic(errChar)
//	}
//
//	// C: 找 value string，然后绑定
//	sd.checkSkip()
//	return sd.scanObjValue(pos)
//}
//
////func (sd *subDecode) scanSubObject() {
////	sub := subDecode{
////		str:       sd.str,
////		scan:      sd.scan,
////		skipTotal: sd.skipValue,
////	}
////
////	if sd.gr != nil {
////		// TODO: 无法为子对象提供目标值，只能返回字符串
////		sub.skipTotal = true
////	} else {
////		sd.skipValue = true
////		*sub.mp = make(cst.KV)
////		sd.mp.Set(sd.key, sub.mp)
////	}
////
////	sub.scanObject()
////	//if err < 0 {
////	//	sd.scan = sub.scan
////	//	return
////	//}
////
////	if sd.gr != nil && sd.skipValue == false {
////		val := sd.str[sd.scan-1 : sub.scan]
////		// TODO: 这里要重新规划一下
////		sd.gr.SetString(sd.key, val)
////	}
////	sd.scan = sub.scan
////	return
////}
//
////func (sd *subDecode) scanSubArray(key string) (val string, err int) {
////	sub := subDecode{
////		str:       sd.str,
////		scan:      sd.scan,
////		skipTotal: sd.skipValue,
////	}
////
////	if sd.gr != nil {
////		// TODO: 无法为子对象提供目标值，只能返回字符串
////		sub.skipTotal = true
////	} else {
////		sd.skipValue = true
////	}
////
////	err = sub.scanList()
////	if err < 0 {
////		sd.scan = sub.scan
////		return
////	}
////
////	if sd.gr != nil {
////		if sd.skipValue == false {
////			val = sd.str[sd.scan-1 : sub.scan]
////		}
////	} else {
////		//sd.mp.Set(key, sub.list)
////	}
////	sd.scan = sub.scan
////	return
////}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 前提：sd.str 肯定是 [ 字符后面的字符串
//// 返回 ] 后面字符的 index
//func (sd *subDecode) scanList(pos int) int {
//	if !sd.isList {
//		panic(errList)
//	}
//
//	// 根据目标值类型，直接匹配，提高效率
//	switch {
//	case isNumKind(sd.dm.itemKind) == true:
//		return sd.scanArrItems(pos, sd.scanNumValue)
//	case sd.dm.itemKind == reflect.String:
//		return sd.scanArrItems(pos, sd.scanStrVal)
//	case sd.dm.itemKind == reflect.Bool:
//		return sd.scanArrItems(pos, sd.scanBoolVal)
//	default:
//		return sd.scanArrItems(pos, sd.scanObjValue)
//	}
//}
//
//func (sd *subDecode) scanArrItems(pos int, scanValue func(int) int) int {
//	//pos := sd.scan
//	first := true
//	for {
//		for isBlankChar[sd.str[pos]] {
//			pos++
//		}
//
//		switch c := sd.str[pos]; c {
//		case ']':
//			pos++
//			return pos
//		case ',':
//			pos++
//			for isBlankChar[sd.str[pos]] {
//				pos++
//			}
//		default:
//			if first {
//				first = false
//			} else {
//				sd.scan = pos
//				panic(errChar)
//			}
//		}
//
//		pos = scanValue(pos)
//	}
//}
//
//func (sd *subDecode) scanStrVal(pos int) int {
//	start := pos
//
//	var slash bool
//	pos, slash = sd.scanQuoteString(pos)
//	if sd.skipValue {
//		return pos
//	}
//
//	if slash {
//		sd.bindString(sd.unescapeString(start, pos))
//	} else {
//		sd.bindString(sd.str[start+1 : pos-1])
//	}
//	return pos
//}
//
//func (sd *subDecode) scanBoolVal(pos int) int {
//	switch sd.str[pos] {
//	case 't':
//		pos = sd.skipTrue(pos)
//		if sd.skipValue {
//			return pos
//		}
//		sd.bindBool(true)
//	case 'f':
//		pos = sd.skipFalse(pos)
//		if sd.skipValue {
//			return pos
//		}
//		sd.bindBool(false)
//	default:
//		sd.scan = pos
//		panic(errBool)
//	}
//	return pos
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (sd *subDecode) scanQuoteString(pos int) (nPos int, slash bool) {
//	//pos := sd.scan
//	if sd.str[pos] != '"' {
//		sd.scan = pos
//		panic(errChar)
//	}
//
//	for {
//		pos++
//
//		switch c := sd.str[pos]; {
//		//case c < ' ':
//		//	sd.scan = pos
//		//	panic(errChar)
//		case c == '"':
//			pos++
//			sd.scan = pos + 1
//			return
//		case c == '\\':
//			slash = true
//			pos++
//			c = sd.str[pos]
//			if c < ' ' {
//				sd.scan = pos
//				panic(errChar)
//			}
//		}
//	}
//}
//
//func (sd *subDecode) scanObjValue(pos int) int {
//	switch sd.str[pos] {
//	case '{':
//		pos++
//		//return sd.scanSubObject(pos)
//		return pos
//	case '[':
//		pos++
//		//err = sd.scanSubArray()
//		return pos
//	case '"':
//		return sd.scanStrVal(pos)
//	default:
//		return sd.scanNoQuoteValue(pos)
//	}
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 匹配一个数值，不带有正负号前缀。
//// 0.234 | 234.23 | 23424 | 3.8e+07
//func (sd *subDecode) scanNumValue(pos int) int {
//	//pos := sd.scan
//	start := pos
//	var hasDot, needNum bool
//
//	c := sd.str[pos]
//	if c == '-' {
//		pos++
//		c = sd.str[pos]
//	}
//
//	//if c < '0' || c > '9' {
//	//	panic(errNumberFmt)
//	//}
//
//	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
//	if c == '0' {
//		pos++
//		c = sd.str[pos]
//
//		switch c {
//		case '.', 'e', 'E':
//			goto loopNum
//		default:
//			goto over
//		}
//	}
//	//pos++
//	needNum = true
//
//loopNum:
//	for {
//		c = sd.str[pos]
//		pos++
//
//		if c == '.' {
//			if hasDot == true {
//				panic(errNumberFmt)
//			}
//			hasDot = true
//			needNum = true
//		} else if c == 'e' || c == 'E' {
//			if needNum {
//				panic(errNumberFmt)
//			}
//			needNum = true
//
//			c := sd.str[pos]
//			if c == '-' || c == '+' {
//				pos++
//			}
//			for {
//				if c = sd.str[pos]; c < '0' || c > '9' {
//					break loopNum
//				} else {
//					needNum = false
//				}
//				pos++
//			}
//		} else if c < '0' || c > '9' {
//			pos--
//			break
//		} else {
//			needNum = false
//		}
//	}
//
//	if needNum {
//		panic(errNumberFmt)
//	}
//
//over:
//	//sd.scan = pos
//	if sd.skipValue {
//		return pos
//	}
//	sd.bindNumber(sd.str[start:pos])
//	return pos
//}
//
//func (sd *subDecode) scanNoQuoteValue(pos int) int {
//	switch c := sd.str[sd.scan]; {
//	case (c >= '0' && c <= '9') || c == '-':
//		//sd.scanNumValue() // 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7 | -0.3 | -3.7E-7
//	case c == 'f':
//		pos = sd.skipFalse(pos)
//		if sd.skipValue {
//			return pos
//		}
//		sd.bindBool(false)
//	case c == 't':
//		pos = sd.skipTrue(pos)
//		if sd.skipValue {
//			return pos
//		}
//		sd.bindBool(true)
//	case c == 'n':
//		pos = sd.skipNull(pos)
//		if sd.skipValue {
//			return pos
//		}
//		sd.bindNull()
//	default:
//		sd.scan = pos
//		panic(errValue)
//	}
//	return pos
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (sd *subDecode) skipNull(pos int) int {
//	pos += 4
//	if sd.str[pos-3:pos] == "ull" {
//		return pos
//	}
//	panic(errNull)
//}
//
//func (sd *subDecode) skipTrue(pos int) int {
//	pos += 4
//	if sd.str[pos-3:pos] == "rue" {
//		return pos
//	}
//	panic(errBool)
//}
//
//func (sd *subDecode) skipFalse(pos int) int {
//	pos += 5
//	if sd.str[pos-4:pos] == "alse" {
//		return pos
//	}
//	panic(errBool)
//}
