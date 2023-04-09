package jsonx

import (
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/gson"
)

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) parseJson() int {
	switch c := sd.sub[sd.scan]; c {
	case '{':
		// 先检查最后一个字符是否符合预期
		if sd.sub[len(sd.sub)-1] != '}' {
			sd.errPos = len(sd.sub) - 1
			return errChar
		}
		sd.sub = sd.sub[:len(sd.sub)-1]
		sd.scan++
		return sd.parseObject()
	case '[':
		if sd.sub[len(sd.sub)-1] != ']' {
			sd.errPos = len(sd.sub) - 1
			return errChar
		}
		sd.sub = sd.sub[:len(sd.sub)-1]
		sd.scan++
		return sd.parseArray()
	case 'n':
		if sd.sub == "null" {
			return noErr
		}
	}
	return errJson
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：dd.sub 肯定是 { 字符后面的字符串
// 返回 } 后面字符的 index
func (sd *subDecode) scanObject(pos int) int {
	var hasKV bool
	for pos < len(sd.sub) {
		c := sd.sub[pos]
		pos++

		if isSpace(c) {
			continue
		}
		if c == '}' {
			return pos
		}

		// 上面都不是，只能是k:v
		if hasKV {
			// 不是第一个kv, 先跳过一个逗号
			pos = sd.scanComma(pos - 1)
			if pos < 0 {
				return pos
			}
			pos += 1
		}
		pos = sd.scanKVItem(pos - 1)
		if pos < 0 {
			return pos
		}
		hasKV = true
	}
	return -1
}

// 必须是k:v, ...形式
func (sd *subDecode) scanKVItem(pos int) int {
	// A: 找冒号 +++
	colonPos := sd.scanColon(pos)
	// 1. 没有冒号，可能就是一段空白字符
	if colonPos < 0 {
		return colonPos
	}

	// 冒号找到了，看看后面的值
	pos = colonPos + 1
	for pos < len(sd.sub) {
		c := sd.sub[pos]
		pos++

		if isSpace(c) {
			continue
		}
		if c == '{' {
			return sd.scanObject(pos)
		}
		if c == '[' {
			return sd.scanArray(pos)
		}

		// one value
		if c == '"' {
			return sd.scanQuoteValue(pos)
		}
		return sd.scanNoQuoteValue(pos - 1)
	}
	return -1
}

// 剩下的全部应该是 k:v,k:v,k:{},k:[]
func (sd *subDecode) parseObject() (err int) {
loopKVItem:
	// TODO: 不应该先找逗号，而应该先找冒号，看冒号后面的第一个非空字符是否是{[，如果是就需要先跳过所有{}和[]的匹配对，再找后面的逗号
	// 注意不是所有{ } [ ] 字符都算，本身key 或者 value 是有可能包含这些特殊字符的。

	// A: 找冒号 +++
	colonPos := sd.scanColon(sd.scan)
	// 1. 没有找到冒号，可能就是一段空白字符
	if colonPos == errNotFound {
		tmp := trim(sd.sub[sd.scan:])
		if len(tmp) != 0 {
			return errChar
		}
		return noErr
	} else if colonPos < 0 {
		return colonPos
	}

	key := sd.sub[sd.scan:colonPos]
	sd.scan = colonPos
	// 2. TODO：这里冒号前面的Key其实就可以得到了，可以先判断目标对象是否有这个key，没有value都不用解析了，直接解析下一个

	// B: 找逗号 +++
	commaPos := sd.scanComma(colonPos + 1) // 从冒号后面开始查找第一个匹配的逗号
	val := ""
	sd.isMixedVal = false
	// 1. 没找到逗号，这是最后一个k:v了
	if commaPos == errNotFound {
		sd.scan = len(sd.sub)
		val = sd.sub[colonPos+1 : sd.scan]
	} else if commaPos <= 0 {
		return commaPos
	}
	// 2. 找到一个“,” 其前面部分当做一个k:v来解
	if commaPos > 0 {
		val = sd.sub[colonPos+1 : commaPos]
		sd.scan = commaPos + 1
	}

	if gr, ok := sd.dst.(*gson.GsonRow); ok {
		err = sd.setGsonValue(gr, key, val)
	} else {
		err = sd.setMapValue(key, val)
	}
	if err < 0 {
		return err
	}

	goto loopKVItem // 再找下一个
}

func (sd *subDecode) parseArray() int {
	return noErr

}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 找冒号："k":"v" | "k":[ |  "k":{
func (sd *subDecode) scanColon(pos int) int {
	quoteCt := 0
	for pos < len(sd.sub) {
		c := sd.sub[pos]
		pos++

		if c == '"' {
			// 非第一个"，如果后面的"前面有\，是允许的
			if quoteCt == 1 {
				if sd.sub[pos-1] == '\\' {
					pos++
					continue
				}
			}
			quoteCt++
			if quoteCt > 2 {
				sd.errPos = pos - 1
				return errChar
			}
		}
		if quoteCt == 2 {
			if c == ':' {
				return pos
			}
		}
	}
	return errNotFound
}

// 找逗号："k":"v",
func (sd *subDecode) scanComma(pos int) int {
	ckSpace := true
	for pos < len(sd.sub) {
		c := sd.sub[pos]
		if c == ',' {
			return pos
		}
		pos++

		if ckSpace {
			if isSpace(c) {
				continue
			}

			ckSpace = false
			if c == '{' {
				pos = sd.scanObject(pos)
				if pos < 0 {
					return pos
				}
				sd.isMixedVal = true
				return sd.scanComma(pos)
			}
			if c == '[' {
				pos = sd.scanArray(pos)
				if pos < 0 {
					return pos
				}
				sd.isMixedVal = true
				return sd.scanComma(pos)
			}
		}
	}
	return -1
}

// 前提：dd.sub 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanArray(pos int) int {
	for pos < len(sd.sub) {
		c := sd.sub[pos]
		pos++

		if c == ']' {
			return pos
		}
		if isSpace(c) {
			continue
		}

		pos = sd.scanArrItem(pos)
	}
	return -1
}

func (sd *subDecode) scanArrItem(pos int) int {

	return -1
}

func (sd *subDecode) scanQuoteValue(pos int) int {
	for pos < len(sd.sub) {
		c := sd.sub[pos]
		pos++

		if c == '"' {
			if sd.sub[pos-1] == '\\' {
				continue
			}
			return pos
		}
	}
	return -1
}

//// 跳过一个逗号
//func (dd *subDecode) scanComma(pos int) int {
//	for pos < len(dd.sub) {
//		c := dd.sub[pos]
//		pos++
//
//		if isSpace(c) {
//			continue
//		}
//		if c == ',' {
//			return pos
//		}
//		return -1
//	}
//	return -1
//}

//
//// 可能的情况 "k":v | "k":{} | "k":[]
//func (dd *subDecode) parseKV(colonPos, commaPos int) int {
//	key := dd.sub[dd.scan:colonPos]
//	val := dd.sub[colonPos+1 : commaPos]
//
//	if dd.gr != nil {
//		return dd.setGsonValue(key, val)
//	}
//	return dd.setMapValue(key, val)
//}

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

// 核查 key +++++++++
func (sd *subDecode) checkKey(key string, needCopy bool) (ret string, err int) {
	var inShare bool
	if key, err = cutKeyQuote(key); err < 0 {
		return
	}
	if key, inShare, err = sd.getStringLiteral(key); err < 0 {
		return
	}
	if inShare && needCopy {
		key = cloneString(key)
	}
	return key, noErr
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

// 检查是否为一个合法的 key
func cutKeyQuote(key string) (ret string, err int) {
	key = trim(key)

	if len(key) < 2 {
		return "", errChar
	}

	if key[0] != '"' || key[len(key)-1] != '"' {
		return "", errChar
	}
	return key[1 : len(key)-1], noErr
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
func (sd *subDecode) scanMatch(pos int, match string) int {
	end := pos + len(match)
	if end > len(sd.sub) {
		return -1
	}
	if sd.sub[pos:end] == match {
		return end
	}
	return -1
}

// 匹配一个数值，不带有正负号前缀。
// 0.234 | 234.23 | 23424 | 3.8e+07
func (sd *subDecode) scanNumber(pos int) int {
	if pos >= len(sd.sub) {
		return -1
	}

	var startZero, hasDot, needNum bool

	c := sd.sub[pos]
	if c < '0' || c > '9' {
		return -1
	}

	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
	if c == '0' {
		startZero = true
	}
	pos++

	for pos < len(sd.sub) {
		c = sd.sub[pos]
		pos++

		if startZero {
			switch c {
			case '.', 'e', 'E':
				startZero = false
			default:
				return pos - 1
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
			return sd.scanScientificNumberTail(pos)
		} else if c < '0' || c > '9' {
			pos--
			break
		} else {
			needNum = false
		}
	}

	if needNum {
		return -1
	}
	return pos
}

// 检查科学计数法（e|E）后面的字符串合法性
func (sd *subDecode) scanScientificNumberTail(pos int) int {
	if pos >= len(sd.sub) {
		return errEOF
	}

	c := sd.sub[pos]
	if c == '-' || c == '+' {
		pos++
	}

	// TODO: 加减号后面没有任何数字匹配，会如何
	for pos < len(sd.sub) {
		if sd.sub[pos] < '0' || sd.sub[pos] > '9' {
			return pos
		}
		pos++
	}
	return pos
}

func (sd *subDecode) scanNoQuoteValue(pos int) int {
	if pos >= len(sd.sub) {
		return -1
	}

	switch c := sd.sub[pos]; {
	case c >= '0' && c <= '9':
		// 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7
		return sd.scanNumber(pos)
	case c == '-':
		// -0.3 | -3.7E-7
		pos++
		return sd.scanNumber(pos)
	case c == 'f':
		return sd.scanMatch(pos, "false")
	case c == 't':
		return sd.scanMatch(pos, "true")
	case c == 'n':
		return sd.scanMatch(pos, "null")
	}
	return -1
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
