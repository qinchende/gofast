package jsonx

import "github.com/qinchende/gofast/cst"

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) parseJson() (ret int) {
	if ret = sd.skipBlank(); ret < 0 {
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
func (sd *subDecode) scanJsonEnd(ch byte) (ret int) {
	if c, pos := sd.lastNotBlank(); c != ch {
		sd.scan = pos
		return errChar
	}

	if ch == '}' {
		sd.scan++
		ret = sd.scanObject()
	} else if ch == ']' {
		sd.scan++
		ret = sd.scanArray()
	} else {
		ret = sd.skipMatch(bytesNull)
		if sd.scan < len(sd.str) {
			return errChar
		}
	}
	if ret == scanEOF {
		return noErr
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 { 字符后面的字符串
// 返回 } 后面一个字符的 index
func (sd *subDecode) scanObject() (ret int) {
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
			if ret = sd.skipSeparator(','); ret < 0 {
				return
			}
			if ret = sd.skipBlank(); ret < 0 {
				return
			}
		}
		if ret = sd.scanKVItem(); ret < 0 {
			return
		}
		hasKV = true
	}
	return scanEOF
}

// 必须是k:v, ...形式。不能为空，而且前面空字符已跳过，否则错误
func (sd *subDecode) scanKVItem() (ret int) {
	// A: 找 key 字符串
	var slash bool
	var key, value string

	idx := sd.scan
	if slash, ret = sd.scanQuoteString(); ret < 0 {
		return
	}
	if slash {
		if key, ret = sd.unescapeString(idx, sd.scan); ret < 0 {
			return
		}
	} else {
		key = sd.str[idx+1 : sd.scan-1]
	}

	// B: 跳过冒号
	if ret = sd.skipSeparator(':'); ret < 0 {
		return
	}
	if ret = sd.skipBlank(); ret < 0 {
		return
	}

	// PS: 可以先判断目标对象是否有这个key，没有就跳过value，解析下一个kv
	idx = -1 // Note: 这里只是复用了 前面的 idx 变量
	if sd.gr != nil {
		idx = sd.gr.KeyIndex(key)
		if idx < 0 {
			sd.skipValue = true
		}
	}

	// C: 找Value
	value, ret = sd.scanValue(key)

	// 不需要值时候，直接跳过去，找下一个KV
	if sd.skipValue || sd.skipTotal {
		sd.skipValue = false
		return
	}

	if ret < 0 {
		return
	}

	// D: 记录解析的 KV
	if sd.gr != nil {
		sd.gr.SetStringByIndex(idx, value)
	} else {
		sd.dst.Set(key, value)
	}
	return noErr
}

func (sd *subDecode) scanSubObject(key string) (val string, ret int) {
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
		newKV := make(cst.KV)
		sd.dst.Set(key, newKV)
		sub.dst = &newKV
	}

	ret = sub.scanObject()
	if ret < 0 {
		sd.scan = sub.scan
		return
	}

	if sd.gr != nil && sd.skipValue == false {
		val = sd.str[sd.scan-1 : sub.scan]
	}
	sd.scan = sub.scan
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanArray() (ret int) {
	for sd.scan < len(sd.str) {
		c := sd.str[sd.scan]
		sd.scan++

		if c == ']' {
			return noErr
		}
		if isSpace(c) {
			continue
		}

		ret = sd.scanArrItem()
	}
	return scanEOF
}

// TODO：需要实现解析List
func (sd *subDecode) scanArrItem() int {

	return -1
}

func (sd *subDecode) scanSubArray(key string) (val string, ret int) {
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
		sub.list = make([]any, 0)
	}

	ret = sub.scanArray()
	if ret < 0 {
		sd.scan = sub.scan
		return
	}

	if sd.gr != nil {
		if sd.skipValue == false {
			val = sd.str[sd.scan-1 : sub.scan]
		}
	} else {
		sd.dst.Set(key, sub.list)
	}
	sd.scan = sub.scan
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) scanQuoteString() (slash bool, ret int) {
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

func (sd *subDecode) scanValue(key string) (val string, ret int) {
	start := sd.scan
	for sd.scan < len(sd.str) {
		c := sd.str[sd.scan]
		sd.scan++

		if isSpace(c) {
			continue
		}

		// TODO：这里需要完善Object,Array
		switch c {
		case '{':
			return sd.scanSubObject(key)
		case '[':
			return sd.scanSubArray(key)
		case '"':
			sd.scan--
			var slash bool
			slash, ret = sd.scanQuoteString()
			if ret < 0 {
				return
			}
			if sd.skipValue {
				return "", noErr
			}
			if slash {
				val, ret = sd.unescapeString(start, sd.scan)
			} else {
				val = sd.str[start+1 : sd.scan-1]
			}
			return
		default:
			sd.scan--
			ret = sd.scanNoQuoteValue()
			if ret < 0 {
				return
			}
			if sd.skipValue {
				return "", noErr
			}
			val = sd.str[start:sd.scan]
			return
		}
	}
	return "", scanEOF
}

// 跳过一个分割符号，前面可以是空字符
// 比如 , | :
func (sd *subDecode) skipSeparator(ch byte) (ret int) {
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
func (sd *subDecode) scanNumber() int {
	if sd.scan >= len(sd.str) {
		return scanEOF
	}

	var startZero, hasDot, needNum bool

	c := sd.str[sd.scan]
	if c < '0' || c > '9' {
		return -1
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
				return sd.scan - 1
			}
		}

		if c == '.' {
			if hasDot == true {
				return -1
			}
			hasDot = true
			needNum = true
		} else if c == 'e' || c == 'E' {
			if needNum {
				return -1
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
	return noErr
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

func (sd *subDecode) scanNoQuoteValue() (ret int) {
	if sd.scan >= len(sd.str) {
		return scanEOF
	}

	switch c := sd.str[sd.scan]; {
	case c >= '0' && c <= '9':
		return sd.scanNumber() // 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7
	case c == '-':
		sd.scan++
		return sd.scanNumber() // -0.3 | -3.7E-7
	case c == 'f':
		return sd.skipMatch(bytesFalse)
	case c == 't':
		return sd.skipMatch(bytesTrue)
	case c == 'n':
		return sd.skipMatch(bytesNull)
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

func (sd *subDecode) skipMatch(match string) (ret int) {
	ret = sd.scan + len(match)
	if ret > len(sd.str) {
		return errMismatch
	}
	if sd.str[sd.scan:ret] == match {
		sd.scan = ret
		return noErr
	}
	return errMismatch
}

//func (sd *subDecode) skipNull() (ret int) {
//	ret = sd.scan + 4
//	if ret > len(sd.str) {
//		return errNull
//	}
//	if sd.str[sd.scan:ret] == bytesNull {
//		sd.scan = ret
//		return noErr
//	} else {
//		return errNull
//	}
//}
