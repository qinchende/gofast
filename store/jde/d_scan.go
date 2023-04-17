package jde

import (
	"github.com/qinchende/gofast/cst"
)

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) parseJson() (err int) {
	if err = sd.skipBlank(); err < 0 {
		return errJson
	}

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
	if c, pos := sd.lastNotBlank(); c != ch {
		sd.scan = pos
		return errChar
	}

	if ch == '}' {
		sd.scan++
		err = sd.scanObject()
	} else if ch == ']' {
		sd.startListPool()
		sd.scan++
		err = sd.scanArray()
		if err < 0 {
			return err
		}
		sd.endListPool()
	} else {
		err = sd.skipMatch(bytesNull)
		if sd.scan < len(sd.str) {
			return errChar
		}
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
	if sd.isList {
		return errObject
	}

	var hasKV bool
	for sd.scan < len(sd.str) {
		switch c := sd.str[sd.scan]; {
		case isSpace(c):
			sd.scan++
			continue
		case c == '}':
			sd.scan++
			return noErr
		}

		// 只能是k:v,k:v 不是第一个k:v, 先跳过一个逗号
		if hasKV {
			if err = sd.skipSeparator(','); err < 0 {
				return
			}
			if err = sd.skipBlank(); err < 0 {
				return
			}
		} else {
			hasKV = true
		}

		if err = sd.scanKVItem(); err < 0 {
			return
		}
	}
	return scanEOF
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
	if err = sd.skipSeparator(':'); err < 0 {
		return
	}
	if err = sd.skipBlank(); err < 0 {
		return
	}

	// C: 找Value，然后直接赋值对象
	return sd.scanValue()
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
		sub.mp = make(cst.KV)
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
//	err = sub.scanArray()
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
func (sd *subDecode) scanArray() (err int) {
	if !sd.isList {
		return errArray
	}

	var hasItem bool
	for sd.scan < len(sd.str) {
		c := sd.str[sd.scan]

		if isSpace(c) {
			continue
		}
		if c == ']' {
			sd.scan++
			return noErr
		}

		// A. 如果已有Item，跳过','
		if hasItem {
			if err = sd.skipSeparator(','); err < 0 {
				return
			}
			if err = sd.skipBlank(); err < 0 {
				return
			}
		} else {
			hasItem = true
		}

		// B: 找Value
		if err = sd.scanValue(); err < 0 {
			return
		}
	}
	return scanEOF
}

func (sd *subDecode) scanArrItem() int {

	return -1
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) scanQuoteString() (slash bool, err int) {
	if sd.str[sd.scan] != '"' {
		return false, errChar
	}
	sd.scan++
	for sd.scan < len(sd.str) {
		c := sd.str[sd.scan]
		sd.scan++

		if c < 32 {
			sd.scan--
			return false, errChar
		} else if c == '\\' {
			slash = true
			if sd.scan < len(sd.str) {
				c = sd.str[sd.scan]
				if c == '"' || c == '\\' {
					sd.scan++
				}
			}
		} else if c == '"' {
			return slash, noErr
		}
	}
	return slash, scanEOF
}

func (sd *subDecode) scanValue() (err int) {
	sd.setSkipFlag()

	var val string
	start := sd.scan
	for sd.scan < len(sd.str) {
		c := sd.str[sd.scan]
		sd.scan++

		if isSpace(c) {
			continue
		}

		switch c {
		case '{':
			err = sd.scanSubObject()
		case '[':
			//err = sd.scanSubArray()
		case '"':
			sd.scan--
			var slash bool
			if slash, err = sd.scanQuoteString(); err < 0 {
				return
			} else if sd.isSkip() {
				return noErr
			}
			if slash {
				if val, err = sd.unescapeString(start, sd.scan); err < 0 {
					return
				}
			} else {
				val = sd.str[start+1 : sd.scan-1]
			}
			err = sd.bindString(val)
		default:
			sd.scan--
			err = sd.scanNoQuoteValue()
		}
		sd.skipValue = false
		return
	}
	return scanEOF
}

// 跳过一个分割符号，前面可以是空字符
// 比如 ',' 或者 ':'
func (sd *subDecode) skipSeparator(ch byte) (err int) {
	for sd.scan < len(sd.str) {
		c := sd.str[sd.scan]
		if isSpace(c) {
			sd.scan++
			continue
		}
		if c == ch {
			sd.scan++
			return noErr
		}
		return errChar
	}
	return scanEOF
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 匹配一个数值，不带有正负号前缀。
// 0.234 | 234.23 | 23424 | 3.8e+07
func (sd *subDecode) scanNumberValue() int {
	if sd.scan >= len(sd.str) {
		return scanEOF
	}

	start := sd.scan
	var startZero, hasDot, needNum bool

	c := sd.str[sd.scan]
	if c < '0' || c > '9' {
		return errNumberFmt
	}

	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
	if c == '0' {
		startZero = true
	}
	sd.scan++

	for sd.scan < len(sd.str) {
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
			return sd.scanScientificNumberTail()
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
	return sd.bindNumber(sd.str[start:sd.scan])
}

// 检查科学计数法（e|E）后面的字符串合法性
func (sd *subDecode) scanScientificNumberTail() int {
	if sd.scan >= len(sd.str) {
		return scanEOF
	}

	c := sd.str[sd.scan]
	if c == '-' || c == '+' {
		sd.scan++
	}

	// TODO: 加减号后面没有任何数字匹配，会如何
	for sd.scan < len(sd.str) {
		if sd.str[sd.scan] < '0' || sd.str[sd.scan] > '9' {
			return noErr
		}
		sd.scan++
	}
	return noErr
}

func (sd *subDecode) scanNoQuoteValue() (err int) {
	if sd.scan >= len(sd.str) {
		return scanEOF
	}

	switch c := sd.str[sd.scan]; {
	case c >= '0' && c <= '9':
		return sd.scanNumberValue() // 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7
	case c == '-':
		sd.scan++
		return sd.scanNumberValue() // -0.3 | -3.7E-7
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
func (sd *subDecode) lastNotBlank() (byte, int) {
	for i := len(sd.str) - 1; i < len(sd.str); i-- {
		if !isSpace(sd.str[i]) {
			sd.str = sd.str[:i+1] // cut 最后的空字符
			return sd.str[i], i
		}
	}
	return 0, 0
}

func (sd *subDecode) skipBlank() int {
	for sd.scan < len(sd.str) {
		if !isSpace(sd.str[sd.scan]) {
			return noErr
		}
		sd.scan++
	}
	return scanEOF
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
