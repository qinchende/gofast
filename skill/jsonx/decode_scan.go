package jsonx

import (
	"github.com/qinchende/gofast/skill/lang"
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
	start := sd.scan + 1
	if ret = sd.scanQuoteString(); ret < 0 {
		return
	}
	key := sd.sub[start : sd.scan-1]

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

// 核查 key +++++++++
func (sd *subDecode) checkKey(key string, needCopy bool) (ret string, err int) {
	var inShare bool
	//if key, err = cutKeyQuote(key); err < 0 {
	//	return
	//}
	if key, inShare, err = sd.getStringLiteral(key); err < 0 {
		return
	}
	if inShare && needCopy {
		key = cloneString(key)
	}
	return key, noErr
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) scanQuoteString() (ret int) {
	if sd.sub[sd.scan] != '"' {
		return errChar
	}

	sd.scan++
	for sd.scan < len(sd.sub) {
		if sd.sub[sd.scan] == '"' {
			sd.scan++
			if sd.sub[sd.scan-2] == '\\' {
				continue
			}
			return noErr
		}
		sd.scan++
	}
	return scanEOF
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
			return "{}", sd.scanObject()
		case '[':
			return "[]", sd.scanArray()
		case '"':
			ret = sd.scanQuoteString()
			if ret >= 0 {
				val = sd.sub[start+1 : sd.scan-1]
			}
			return
		default:
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

// 核查 value ++++++++++
func (sd *subDecode) checkValue(val string) (ret string, err int) {
	var hasQuote, inShare bool
	if val, hasQuote, err = cutValueQuote(val); err < 0 {
		return
	}
	if hasQuote {
		if val, inShare, err = sd.getStringLiteral(val); err < 0 {
			return
		}
		// 证明此时用的是临时栈空间保存转义之后的Value，需要申请新的内存空间放置
		if inShare {
			val = cloneString(val)
		}
	}
	return val, noErr
}

// 当目标为 cst.KV 类型时候，用此方法设置
func (sd *subDecode) setMapValue(key, val string) (err int) {
	if key, err = sd.checkKey(key, true); err < 0 {
		return
	}

	if !sd.isMixedVal {
		if val, err = sd.checkValue(val); err < 0 {
			return
		}
	}

	// set k = v
	sd.dst.Set(key, val)
	return noErr
}

// 当目标为 gson.GsonRow 类型时候，用此方法设置
func (sd *subDecode) setGsonValue(gr *gson.GsonRow, key, val string) (err int) {
	if key, err = sd.checkKey(key, false); err < 0 {
		return
	}
	keyIdx := gr.KeyIndex(key)
	// 没有这个字段，直接返回了(此时再去解析后面的value是没有意义的)
	if keyIdx < 0 {
		return noErr
	}

	if !sd.isMixedVal {
		if val, err = sd.checkValue(val); err < 0 {
			return
		}
	}

	// set k = v
	gr.SetStringByIndex(keyIdx, val)
	return noErr
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ++++++++++++++++++++++++++++++++++++++需要转义的字符
// \\ 反斜杠
// \" 双引号
// \' 单引号 （没有这个）
// \/ 正斜杠
// \b 退格符
// \f 换页符
// \t 制表符
// \n 换行符
// \r 回车符
// \u 后面跟十六进制字符 （比如笑脸表情 \u263A）
// +++++++++++++++++++++++++++++++++++++++++++++++++++
// 一个合法的 Key，或者Value 字符串
func (sd *subDecode) getStringLiteral(str string) (ret string, inShare bool, err int) {
	var newStr []byte
	var step int
	var hasSlash bool

	for i := 0; i < len(str); i++ {
		c := str[i]
		// 不支持非可见字符
		if c < 32 {
			return "", false, errChar
		}
		if c == '\\' {
			// 第一次检索到有 \
			if hasSlash == false {
				hasSlash = true
				// TODO：这里发生了逃逸，需要用sync.Pool的方式，共享内存空间
				// 或者别的黑魔法操作内存
				// add by sdx 20230404 动态初始化 share 内存
				if sd.share == nil {
					defSize := len(str)
					if defSize > tempByteStackSize {
						defSize = tempByteStackSize
					}
					sd.share = make([]byte, defSize)
				}

				if len(str) <= len(sd.share) {
					newStr = sd.share[:]
					inShare = true
				} else {
					newStr = make([]byte, len(str))
				}
				for ; step < i; step++ {
					newStr[step] = str[step]
				}
			}
			i++
			c = str[i]
			// 判断 \ 后面的字符
			switch c {
			case '"', '/', '\\':
				newStr[step] = c
				step++
			//case '\'': // 这种情况认为是错误
			case 'b':
				newStr[step] = '\b'
				step++
			case 'f':
				newStr[step] = '\f'
				step++
			case 't':
				newStr[step] = '\t'
				step++
			case 'n':
				newStr[step] = '\n'
				step++
			case 'r':
				newStr[step] = '\r'
				step++
			case 'u': // TODO: uft8编码字符有待转换
			default:
				return "", false, errJson
			}
			continue
		}
		if hasSlash {
			newStr[step] = c
			step++
		}
		//	// ASCII
		//case c < utf8.RuneSelf:
		//	b[w] = c
		//	r++
		//	w++
		//
		//	// Coerce to well-formed UTF-8.
		//	default:
		//	rr, size := utf8.DecodeRune(s[r:])
		//	r += size
		//	w += utf8.EncodeRune(b[w:], rr)
	}
	if hasSlash {
		return lang.BTS(newStr[:step]), inShare, noErr
	}
	return str, false, noErr
}

// 检查是否是一个合法的 value 值
func cutValueQuote(val string) (ret string, hasQuote bool, err int) {
	val = trim(val)

	if len(val) < 1 {
		return "", false, errChar
	}

	// 如果 value 没有 双引号，可能是数值、true、false、null四种情况
	if val[0] != '"' {
		if err = checkNoQuoteValue(val); err < 0 {
			return "", false, err
		}
		return val, false, noErr
	}

	// 有双引号
	if len(val) == 1 {
		return "", false, errChar
	} else if val[len(val)-1] != '"' {
		return "", false, errChar
	}
	return val[1 : len(val)-1], true, noErr
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

// 没有双引号围起来的值，只可能是：数值、true、false、null
func checkNoQuoteValue(str string) int {
	if len(str) == 0 {
		return errJson
	}

	// 只有这几种首字符的可能
isNumber:
	c := str[0]
	if c >= '0' && c <= '9' {
		// 0 | 0.234 | 234.23 | 23424 | 3.8e+07
		dot := false
		if c == '0' {
			if len(str) < 3 || str[1] != '.' {
				return errChar
			}
		}
		for i := 1; i < len(str); i++ {
			c = str[i]
			if c == '.' {
				if dot == true {
					return errNumberFmt
				} else {
					dot = true
					continue
				}
			} else if c == 'e' || c == 'E' {
				return checkScientificNumberTail(str[i+1:])
			} else if c < '0' || c > '9' {
				return errNumberFmt
			}
		}
	} else if c == '-' {
		// -0 | -0.3 | +13.33 | -3.7E-7
		if len(str) >= 2 && str[1] >= '0' && str[1] <= '9' {
			str = str[1:]
			goto isNumber
		} else {
			return errNumberFmt
		}
	} else if c == 'f' {
		// false
		if str != "false" {
			return errChar
		}
	} else if c == 't' {
		// true
		if str != "true" {
			return errChar
		}
	} else if c == 'n' {
		if str != "null" {
			return errChar
		}
	} else {
		return errJson
	}
	return noErr
}

// 检查科学计数法（e|E）后面的字符串合法性
func checkScientificNumberTail(str string) int {
	if len(str) == 0 {
		return errNumberFmt
	}
	c := str[0]
	if c == '-' || c == '+' {
		str = str[1:]
	}

	if len(str) == 0 {
		return errNumberFmt
	}
	for i := range str {
		if str[i] < '0' || str[i] > '9' {
			return errNumberFmt
		}
	}
	return noErr
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

//// 找逗号："k":"v",
//func (sd *subDecode) scanComma() (ret int) {
//	ckBlank := true
//	for sd.scan < len(sd.sub) {
//		c := sd.sub[sd.scan]
//		if c == ',' {
//			return noErr
//		}
//		sd.scan++
//
//		if ckBlank {
//			if isSpace(c) {
//				continue
//			}
//
//			ckBlank = false
//			if c == '{' {
//				ret = sd.scanObject()
//				if ret < 0 {
//					return ret
//				}
//				sd.isMixedVal = true
//				return sd.scanComma()
//			}
//			if c == '[' {
//				ret = sd.scanArray()
//				if ret < 0 {
//					return ret
//				}
//				sd.isMixedVal = true
//				return sd.scanComma()
//			}
//		}
//	}
//	return scanEOF
//}

//
//// 剩下的全部应该是 k:v,k:v,k:{},k:[]
//func (sd *subDecode) parseObject() (err int) {
//loopKVItem:
//	// TODO: 不应该先找逗号，而应该先找冒号，看冒号后面的第一个非空字符是否是{[，如果是就需要先跳过所有{}和[]的匹配对，再找后面的逗号
//	// 注意不是所有{ } [ ] 字符都算，本身key 或者 value 是有可能包含这些特殊字符的。
//
//	// A: 找冒号 +++
//	colonPos := sd.scanColon(sd.scan)
//	// 1. 没有找到冒号，可能就是一段空白字符
//	if colonPos == errNotFound {
//		tmp := trim(sd.sub[sd.scan:])
//		if len(tmp) != 0 {
//			return errChar
//		}
//		return noErr
//	} else if colonPos < 0 {
//		return colonPos
//	}
//
//	key := sd.sub[sd.scan:colonPos]
//	sd.scan = colonPos
//	// 2. TODO：这里冒号前面的Key其实就可以得到了，可以先判断目标对象是否有这个key，没有value都不用解析了，直接解析下一个
//
//	// B: 找逗号 +++
//	commaPos := sd.scanComma(colonPos + 1) // 从冒号后面开始查找第一个匹配的逗号
//	val := ""
//	sd.isMixedVal = false
//	// 1. 没找到逗号，这是最后一个k:v了
//	if commaPos == errNotFound {
//		sd.scan = len(sd.sub)
//		val = sd.sub[colonPos+1 : sd.scan]
//	} else if commaPos <= 0 {
//		return commaPos
//	}
//	// 2. 找到一个“,” 其前面部分当做一个k:v来解
//	if commaPos > 0 {
//		val = sd.sub[colonPos+1 : commaPos]
//		sd.scan = commaPos + 1
//	}
//
//	if gr, ok := sd.dst.(*gson.GsonRow); ok {
//		err = sd.setGsonValue(gr, key, val)
//	} else {
//		err = sd.setMapValue(key, val)
//	}
//	if err < 0 {
//		return err
//	}
//
//	goto loopKVItem // 再找下一个
//}

//func (sd *subDecode) scanArray() int {
//	return noErr
//
//}

//
//// 可能的情况 "k":v | "k":{} | "k":[]
//func (sd *subDecode) parseKV(colonPos, commaPos int) int {
//	key := sd.sub[sd.scan:colonPos]
//	val := sd.sub[colonPos+1 : commaPos]
//
//	if sd.gr != nil {
//		return sd.setGsonValue(key, val)
//	}
//	return sd.setMapValue(key, val)
//}

//func (sd *subDecode) scanQuoteValue(pos int) int {
//	for pos < len(sd.sub) {
//		c := sd.sub[pos]
//		pos++
//
//		if c == '"' {
//			if sd.sub[pos-1] == '\\' {
//				continue
//			}
//			return pos
//		}
//	}
//	sd.pos = pos
//	return scanEOF
//}

//
//// 检查是否为一个合法的 key
//func cutKeyQuote(key string) (ret string, err int) {
//	key = trim(key)
//
//	if len(key) < 2 {
//		return "", errChar
//	}
//
//	if key[0] != '"' || key[len(key)-1] != '"' {
//		return "", errChar
//	}
//	return key[1 : len(key)-1], noErr
//}
