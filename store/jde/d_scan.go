package jde

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"reflect"
)

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) scanJson() (err int) {
	// 万一解析过程中异常，这里统一截获处理，返回解析错误
	defer func() {
		if pic := recover(); pic != nil {
			fmt.Println(pic) // 调试的时候打印错误信息
			err = errJson
		}
	}()

	sd.skipBlank()

	switch sd.str[sd.scan] {
	case '{':
		return sd.scanJsonEnd('}')
	case '[':
		return sd.scanJsonEnd(']')
	case 'n':
		return sd.scanJsonEnd('l')
	}
	return errJson
}

// 只支持 } ] l 三个字符判断
func (sd *subDecode) scanJsonEnd(ch byte) (err int) {
	// 去掉尾部的空字符
	for i := len(sd.str) - 1; i > 0; i-- {
		if !isBlankChar[sd.str[i]] {
			if sd.str[i] != ch {
				sd.scan = i
				return errChar
			}
			sd.str = sd.str[:i+1]
			break
		}
	}

	if ch == '}' {
		sd.scan++
		err = sd.scanObject()
	} else if ch == ']' {
		sd.resetListPool()
		sd.scan++
		if err = sd.scanList(); err < 0 {
			return err
		}
		sd.flushListPool()
	} else {
		err = sd.skipMatch(bytesNull)
	}
	if err == scanEOF {
		return noErr
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 { 字符后面的字符串
// 返回 } 后面一个字符的 index
func (sd *subDecode) scanObject() (err int) {
	first := true
	for {
		sd.skipBlank()

		switch c := sd.str[sd.scan]; c {
		case '}':
			sd.scan++
			return noErr
		case ',':
			sd.scan++
			sd.skipBlank()
			goto scanKVPair
		default:
			if first {
				first = false
				goto scanKVPair
			}
			return errChar
		}

	scanKVPair:
		if err = sd.scanKVItem(); err < 0 {
			return
		}
	}
}

// 必须是k:v, ...形式。不能为空，而且前面空字符已跳过，否则错误
func (sd *subDecode) scanKVItem() (err int) {
	// A: 找 key 字符串
	var slash bool
	start := sd.scan
	if slash, err = sd.scanQuoteString(); err < 0 {
		return
	}
	if slash {
		if sd.key, err = sd.unescapeString(start, sd.scan); err < 0 {
			return
		}
	} else {
		sd.key = sd.str[start+1 : sd.scan-1]
	}

	// B: 跳过冒号
	sd.skipBlank()
	if sd.str[sd.scan] == ':' {
		sd.scan++
		sd.skipBlank()
	} else {
		return errChar
	}

	// C: 找 value string，然后绑定
	sd.setSkip()
	err = sd.scanObjValue()
	sd.skipValue = false
	return
}

func (sd *subDecode) scanSubObject() (err int) {
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

	err = sub.scanObject()
	if err < 0 {
		sd.scan = sub.scan
		return
	}

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
func (sd *subDecode) scanList() (err int) {
	if !sd.isList {
		return errList
	}

	// 根据目标值类型，直接匹配，提高效率
	switch {
	case isNumKind(sd.arr.itemKind) == true:
		err = sd.scanArrItems(sd.scanNumValue)
	case sd.arr.itemKind == reflect.String:
		err = sd.scanArrItems(sd.scanStrVal)
	case sd.arr.itemKind == reflect.Bool:
		err = sd.scanArrItems(sd.scanBoolVal)
	default:
		err = sd.scanArrItems(sd.scanObjValue)
	}
	sd.skipValue = false
	return
}

func (sd *subDecode) scanArrItems(scanValue func() int) (err int) {
	first := true
	for {
		sd.skipBlank()

		switch c := sd.str[sd.scan]; c {
		case ']':
			sd.scan++
			return noErr
		case ',':
			sd.scan++
			sd.skipBlank()
		default:
			if first {
				first = false
			} else {
				return errChar
			}
		}

		if err = scanValue(); err < 0 {
			return
		}
	}
}

func (sd *subDecode) scanStrVal() (err int) {
	start := sd.scan

	var slash bool
	if slash, err = sd.scanQuoteString(); err < 0 {
		return
	} else if sd.isSkip() {
		return noErr
	}

	var val string
	if slash {
		if val, err = sd.unescapeString(start, sd.scan); err < 0 {
			return
		}
	} else {
		val = sd.str[start+1 : sd.scan-1]
	}
	return sd.bindString(val)
}

func (sd *subDecode) scanBoolVal() (err int) {
	if sd.isSkip() {
		return noErr
	}
	return noErr
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) scanQuoteString() (slash bool, err int) {
	if sd.str[sd.scan] != '"' {
		return false, errChar
	}

	for {
		sd.scan++

		switch c := sd.str[sd.scan]; {
		//case c < ' ':
		//	return false, errChar
		case c == '"':
			sd.scan++
			return slash, noErr
		case c == '\\':
			slash = true
			sd.scan++
			c = sd.str[sd.scan]
			if c < ' ' {
				return false, errChar
			}
		}
	}
}

func (sd *subDecode) scanObjValue() (err int) {
	switch sd.str[sd.scan] {
	case '{':
		sd.scan++
		err = sd.scanSubObject()
	case '[':
		sd.scan++
		//err = sd.scanSubArray()
	case '"':
		err = sd.scanStrVal()
	default:
		err = sd.scanNoQuoteValue()
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 匹配一个数值，不带有正负号前缀。
// 0.234 | 234.23 | 23424 | 3.8e+07
func (sd *subDecode) scanNumValue() int {
	start := sd.scan
	var startZero, hasDot, needNum bool

	c := sd.str[sd.scan]
	if c == '-' {
		sd.scan++
		c = sd.str[sd.scan]
	}

	if c < '0' || c > '9' {
		return errNumberFmt
	}

	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
	if c == '0' {
		startZero = true
	}
	sd.scan++

loopNum:
	for {
		c = sd.str[sd.scan]
		sd.scan++

		if startZero {
			switch c {
			case '.', 'e', 'E':
				startZero = false
			default:
				return errNumberFmt
			}
		}

		if c == '.' {
			if hasDot == true {
				return errNumberFmt
			}
			hasDot = true
			needNum = true
		} else if c == 'e' || c == 'E' {
			if needNum {
				return errNumberFmt
			}
			needNum = true

			c := sd.str[sd.scan]
			if c == '-' || c == '+' {
				sd.scan++
			}
			for {
				if c = sd.str[sd.scan]; c < '0' || c > '9' {
					break loopNum
				} else {
					needNum = false
				}
				sd.scan++
			}
		} else if c < '0' || c > '9' {
			sd.scan--
			break
		} else {
			needNum = false
		}
	}

	if needNum {
		return errNumberFmt
	}

	if sd.isSkip() {
		return noErr
	}
	return sd.bindNumber(sd.str[start:sd.scan], hasDot)
}

func (sd *subDecode) scanNoQuoteValue() (err int) {
	switch c := sd.str[sd.scan]; {
	case (c >= '0' && c <= '9') || c == '-':
		return sd.scanNumValue() // 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7 | -0.3 | -3.7E-7
	case c == 'f':
		if err = sd.skipMatch(bytesFalse); err < 0 {
			return
		} else if sd.isSkip() {
			return noErr
		}
		return sd.bindBool(false)
	case c == 't':
		if err = sd.skipMatch(bytesTrue); err < 0 {
			return
		} else if sd.isSkip() {
			return noErr
		}
		return sd.bindBool(true)
	case c == 'n':
		if err = sd.skipMatch(bytesNull); err < 0 {
			return
		} else if sd.isSkip() {
			return noErr
		}
		return sd.bindNull()
	}

	return errValue
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//go:inline
func (sd *subDecode) skipBlank() {
	for isBlankChar[sd.str[sd.scan]] {
		sd.scan++
	}
}

func (sd *subDecode) skipMatch(match string) int {
	pos := sd.scan + len(match)
	if pos > len(sd.str) {
		return errMismatch
	}
	if sd.str[sd.scan:pos] == match {
		sd.scan = pos
		return noErr
	}
	return errMismatch
}
