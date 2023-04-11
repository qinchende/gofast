package jsonx

import (
	"github.com/qinchende/gofast/store/gson"
)

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) parseJson() (ret int) {
	if ret = sd.skipBlank(); ret < 0 {
		return errJson
	}

	switch sd.sub[sd.scan] {
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
		if sd.scan < len(sd.sub) {
			return errChar
		}
	}
	if ret == scanEOF {
		return noErr
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.sub 肯定是 { 字符后面的字符串
// 返回 } 后面字符的 index
func (sd *subDecode) scanObject() (ret int) {
	var hasKV bool
	for sd.scan < len(sd.sub) {
		switch c := sd.sub[sd.scan]; {
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
	start := sd.scan
	if slash, ret = sd.scanQuoteString(); ret < 0 {
		return
	}
	var key string
	if slash {
		if key, ret = sd.unescapeString(start, sd.scan); ret < 0 {
			return
		}
	} else {
		key = sd.sub[start+1 : sd.scan-1]
	}

	// B: 跳过冒号
	if ret = sd.skipSeparator(':'); ret < 0 {
		return
	}
	if ret = sd.skipBlank(); ret < 0 {
		return
	}

	// 2. TODO：这里冒号前面的Key其实就可以得到了，可以先判断目标对象是否有这个key，没有value就跳过值，然后解析下一个
	//keyIdx := -1
	//gr, _ := sd.dst.(*gson.GsonRow)
	//if gr != nil {
	//	keyIdx = gr.KeyIndex(key)
	//}

	// C: 找Value
	value := ""
	if value, ret = sd.scanValue(); ret < 0 {
		return
	}

	// D: 记录解析的 KV
	if gr, ok := sd.dst.(*gson.GsonRow); ok {
		ret = sd.setGsonValue(gr, key, value)
	} else {
		ret = sd.setMapValue(key, value)
	}
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) scanQuoteString() (slash bool, ret int) {
	if sd.sub[sd.scan] != '"' {
		return false, errChar
	}
	sd.scan++
	for sd.scan < len(sd.sub) {
		c := sd.sub[sd.scan]
		sd.scan++

		if c < 32 {
			sd.scan--
			return false, errChar
		} else if c == '\\' {
			slash = true
			if sd.scan < len(sd.sub) {
				c = sd.sub[sd.scan]
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

func (sd *subDecode) scanValue() (val string, ret int) {
	start := sd.scan
	for sd.scan < len(sd.sub) {
		c := sd.sub[sd.scan]
		sd.scan++

		if isSpace(c) {
			continue
		}

		// TODO：这里需要完善Object,Array
		switch c {
		case '{':
			sd.directString = true
			return "{}", sd.scanObject()
		case '[':
			sd.directString = true
			return "[]", sd.scanArray()
		case '"':
			sd.scan--
			var slash bool
			slash, ret = sd.scanQuoteString()
			if ret < 0 {
				return
			}
			if slash {
				val, ret = sd.unescapeString(start, sd.scan)
			} else {
				val = sd.sub[start+1 : sd.scan-1]
			}
			return
		default:
			sd.scan--
			ret = sd.scanNoQuoteValue()
			if ret >= 0 {
				val = sd.sub[start:sd.scan]
			}
			return
		}
	}
	return "", scanEOF
}

// 跳过一个分割符号，前面可以是空字符
// 比如 , | :
func (sd *subDecode) skipSeparator(ch byte) (ret int) {
	for sd.scan < len(sd.sub) {
		c := sd.sub[sd.scan]
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

// 当目标为 cst.KV 类型时候，用此方法设置
func (sd *subDecode) setMapValue(key, val string) (err int) {
	// set k = v
	sd.dst.Set(key, val)
	return noErr
}

// 当目标为 gson.GsonRow 类型时候，用此方法设置
func (sd *subDecode) setGsonValue(gr *gson.GsonRow, key, val string) (err int) {
	keyIdx := gr.KeyIndex(key)
	// 没有这个字段，直接返回了(此时再去解析后面的value是没有意义的)
	if keyIdx < 0 {
		return noErr
	}

	// set k = v
	gr.SetStringByIndex(keyIdx, val)
	return noErr
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 匹配一个数值，不带有正负号前缀。
// 0.234 | 234.23 | 23424 | 3.8e+07
func (sd *subDecode) scanNumber() int {
	if sd.scan >= len(sd.sub) {
		return scanEOF
	}

	var startZero, hasDot, needNum bool

	c := sd.sub[sd.scan]
	if c < '0' || c > '9' {
		return -1
	}

	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
	if c == '0' {
		startZero = true
	}
	sd.scan++

	for sd.scan < len(sd.sub) {
		c = sd.sub[sd.scan]
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
	if sd.scan >= len(sd.sub) {
		return scanEOF
	}

	c := sd.sub[sd.scan]
	if c == '-' || c == '+' {
		sd.scan++
	}

	// TODO: 加减号后面没有任何数字匹配，会如何
	for sd.scan < len(sd.sub) {
		if sd.sub[sd.scan] < '0' || sd.sub[sd.scan] > '9' {
			return noErr
		}
		sd.scan++
	}
	return noErr
}

func (sd *subDecode) scanNoQuoteValue() (ret int) {
	if sd.scan >= len(sd.sub) {
		return scanEOF
	}

	switch c := sd.sub[sd.scan]; {
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.sub 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanArray() (ret int) {
	for sd.scan < len(sd.sub) {
		c := sd.sub[sd.scan]
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

func (sd *subDecode) scanArrItem() int {

	return -1
}
