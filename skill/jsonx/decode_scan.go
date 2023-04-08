package jsonx

import "github.com/qinchende/gofast/skill/lang"

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (dd *fastDecode) parseJson() int {
	dd.skipHeadSpace()
	dd.skipTailSpace()

	if dd.head >= dd.tail {
		return errEOF
	}

	switch c := dd.src[dd.head]; c {
	case '{':
		if dd.src[dd.tail] != '}' {
			dd.flag = dd.tail + 1
			return errChar
		}
		dd.changeSub(dd.src[dd.head+1 : dd.tail])
		return dd.parseObject()
	case '[':
		if dd.src[dd.tail] != ']' {
			dd.flag = dd.tail + 1
			return errChar
		}
		dd.changeSub(dd.src[dd.head+1 : dd.tail])
		return dd.parseArray()
	case 'n':
		if dd.src[dd.head:dd.tail+1] == "null" {
			return noErr
		}
	}
	return errJson
}

// 剩下的全部应该是 k:v,k:v,k:{},k:[]
func (dd *fastDecode) parseObject() int {
loopKVItem:
	// TODO: 不应该先找逗号，而应该先找冒号，看冒号后面的第一个非空字符是否是{[，如果是就需要先跳过所有{}和[]的匹配对，再找后面的逗号
	// 注意不是所有{ } [ ] 字符都算，本身key 或者 value 是有可能包含这些特殊字符的。

	// A: 找冒号 +++
	colonIdx := dd.nextColon(int(dd.scan))
	// 1. 没有冒号，可能就是一段空白字符
	if colonIdx == -1 {
		tmp := trim(dd.sub[dd.scan:])
		if len(tmp) != 0 {
			return errChar
		}
		return noErr
	}
	// 2. TODO：这里冒号前面的Key其实就可以得到了，可以先判断目标对象是否有这个key，没有value都不用解析了，直接解析下一个

	// B: 找逗号 +++
	commaIdx := dd.nextComma(colonIdx + 1) // 从冒号后面开始查找
	// 1. 没找到逗号，这是最后一个k:v了
	if commaIdx == -1 {
		end := uint32(len(dd.sub))
		err := dd.parseKV(uint32(colonIdx), end)
		dd.isMixedVal = false
		dd.scan = end
		return err
	}
	// 2. 找到一个“,” 其前面部分当做一个k:v来解
	if commaIdx > 0 {
		err := dd.parseKV(uint32(colonIdx), uint32(commaIdx))
		dd.isMixedVal = false
		dd.scan = uint32(commaIdx) + 1
		if err != noErr {
			return err
		}
		goto loopKVItem // 再找下一个
	}
	return noErr
}

func (dd *fastDecode) parseArray() int {
	return noErr

}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 找冒号："k":"v" | "k":[ |  "k":{
func (dd *fastDecode) nextColon(idx int) int {
	quoteCt := 0
	for idx < len(dd.sub) {
		if dd.sub[idx] == '"' {
			// 非第一个"，如果后面的"前面有\，是允许的
			if quoteCt == 1 {
				if dd.sub[idx-1] == '\\' {
					continue
				}
			}
			quoteCt++
			if quoteCt > 2 {
				return -1
			}
		}
		if quoteCt == 2 {
			if dd.sub[idx] == ':' {
				return idx
			}
		}
		idx++
	}
	return -1
}

// 找逗号："k":"v",
func (dd *fastDecode) nextComma(pos int) int {
	ckSpace := true
	for pos < len(dd.sub) {
		c := dd.sub[pos]
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
				pos = dd.skipObject(pos)
				if pos < 0 {
					return pos
				}
				dd.isMixedVal = true
				return dd.skipComma(pos)
			}
			if c == '[' {
				pos = dd.skipArray(pos)
				if pos < 0 {
					return pos
				}
				dd.isMixedVal = true
				return dd.skipComma(pos)
			}
		}
	}
	return -1
}

// 前提：dd.sub 肯定是 { 字符后面的字符串
// 返回 } 后面字符的 index
func (dd *fastDecode) skipObject(pos int) int {
	var hasKV bool
	for pos < len(dd.sub) {
		c := dd.sub[pos]
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
			pos = dd.skipComma(pos - 1)
			if pos < 0 {
				return pos
			}
			pos += 1
		}
		pos = dd.skipKVItem(pos - 1)
		if pos < 0 {
			return pos
		}
		hasKV = true
	}
	return -1
}

// 前提：dd.sub 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (dd *fastDecode) skipArray(pos int) int {
	for pos < len(dd.sub) {
		c := dd.sub[pos]
		pos++

		if c == ']' {
			return pos
		}
		if isSpace(c) {
			continue
		}

		pos = dd.skipArrItem(pos)
	}
	return -1
}

// 必须是k:v, ...形式
func (dd *fastDecode) skipKVItem(pos int) int {
	// A: 找冒号 +++
	colonIdx := dd.nextColon(pos)
	// 1. 没有冒号，可能就是一段空白字符
	if colonIdx < 0 {
		return colonIdx
	}

	// 冒号找到了，看看后面的值
	pos = colonIdx + 1
	for pos < len(dd.sub) {
		c := dd.sub[pos]
		pos++

		if isSpace(c) {
			continue
		}
		if c == '{' {
			return dd.skipObject(pos)
		}
		if c == '[' {
			return dd.skipArray(pos)
		}
		if c == '"' {
			return dd.skipQuoteValue(pos)
		}
		return dd.skipNoQuoteValue(pos - 1)
	}
	return -1
}

func (dd *fastDecode) skipArrItem(pos int) int {

	return -1
}

func (dd *fastDecode) skipQuoteValue(pos int) int {
	for pos < len(dd.sub) {
		c := dd.sub[pos]
		pos++

		if c == '"' {
			if dd.sub[pos-1] == '\\' {
				continue
			}
			return pos
		}
	}
	return -1
}

// 跳过一个逗号
func (dd *fastDecode) skipComma(pos int) int {
	for pos < len(dd.sub) {
		c := dd.sub[pos]
		pos++

		if isSpace(c) {
			continue
		}
		if c == ',' {
			return pos
		}
		return -1
	}
	return -1
}

// 可能的情况 "k":v | "k":{} | "k":[]
func (dd *fastDecode) parseKV(colonIdx, commaIdx uint32) int {
	key := dd.sub[dd.scan:colonIdx]
	val := dd.sub[colonIdx+1 : commaIdx]

	if dd.gr != nil {
		return dd.setGsonValue(key, val)
	}
	return dd.setMapValue(key, val)
}

// 当目标为 cst.KV 类型时候，用此方法设置
func (dd *fastDecode) setMapValue(key, val string) (err int) {
	if key, err = dd.checkKey(key, true); err < 0 {
		return
	}

	if !dd.isMixedVal {
		if val, err = dd.checkValue(val); err < 0 {
			return
		}
	}

	// set k = v
	dd.dst.Set(key, val)
	return noErr
}

// 当目标为 gson.GsonRow 类型时候，用此方法设置
func (dd *fastDecode) setGsonValue(key, val string) (err int) {
	if key, err = dd.checkKey(key, false); err < 0 {
		return
	}
	keyIdx := dd.gr.KeyIndex(key)
	// 没有这个字段，直接返回了(此时再去解析后面的value是没有意义的)
	if keyIdx < 0 {
		return noErr
	}

	if !dd.isMixedVal {
		if val, err = dd.checkValue(val); err < 0 {
			return
		}
	}

	// set k = v
	dd.gr.SetStringByIndex(keyIdx, val)
	return noErr
}

// 核查 key +++++++++
func (dd *fastDecode) checkKey(key string, needCopy bool) (ret string, err int) {
	var inShare bool
	if key, err = cutKeyQuote(key); err < 0 {
		return
	}
	if key, inShare, err = dd.getStringLiteral(key); err < 0 {
		return
	}
	if inShare && needCopy {
		key = cloneString(key)
	}
	return key, noErr
}

// 核查 value ++++++++++
func (dd *fastDecode) checkValue(val string) (ret string, err int) {
	var hasQuote, inShare bool
	if val, hasQuote, err = cutValueQuote(val); err < 0 {
		return
	}
	if hasQuote {
		if val, inShare, err = dd.getStringLiteral(val); err < 0 {
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
func (dd *fastDecode) getStringLiteral(str string) (ret string, inShare bool, err int) {
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
				if dd.share == nil {
					defSize := len(str)
					if defSize > tempByteStackSize {
						defSize = tempByteStackSize
					}
					dd.share = make([]byte, defSize)
				}

				if len(str) <= len(dd.share) {
					newStr = dd.share[:]
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
func (dd *fastDecode) skipMatch(pos int, match string) int {
	end := pos + len(match)
	if end > len(dd.sub) {
		return -1
	}
	if dd.sub[pos:end] == match {
		return end
	}
	return -1
}

// 匹配一个数值，不带有正负号前缀。
// 0.234 | 234.23 | 23424 | 3.8e+07
func (dd *fastDecode) skipNumber(pos int) int {
	if pos >= len(dd.sub) {
		return -1
	}

	var startZero, hasDot, needNum bool

	c := dd.sub[pos]
	if c < '0' || c > '9' {
		return -1
	}

	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
	if c == '0' {
		startZero = true
	}
	pos++

	for pos < len(dd.sub) {
		c = dd.sub[pos]
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
			return dd.skipScientificNumberTail(pos)
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
func (dd *fastDecode) skipScientificNumberTail(pos int) int {
	if pos >= len(dd.sub) {
		return -1
	}

	c := dd.sub[pos]
	if c == '-' || c == '+' {
		pos++
	}

	// TODO: 加减号后面没有任何数字匹配，会如何
	for pos < len(dd.sub) {
		if dd.sub[pos] < '0' || dd.sub[pos] > '9' {
			return pos
		}
		pos++
	}
	return pos
}

func (dd *fastDecode) skipNoQuoteValue(pos int) int {
	if pos >= len(dd.sub) {
		return -1
	}

	switch c := dd.sub[pos]; {
	case c >= '0' && c <= '9':
		// 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7
		return dd.skipNumber(pos)
	case c == '-':
		// -0.3 | -3.7E-7
		pos++
		return dd.skipNumber(pos)
	case c == 'f':
		return dd.skipMatch(pos, "false")
	case c == 't':
		return dd.skipMatch(pos, "true")
	case c == 'n':
		return dd.skipMatch(pos, "null")
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
