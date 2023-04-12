package jsonx

//func (sd *subDecode) skipHeadSpace() {
//	for sd.scan < len(sd.str) {
//		if !isSpace(sd.str[sd.scan]) {
//			return
//		}
//		sd.scan++
//	}
//}
//
//func (sd *subDecode) skipTailSpace() {
//	tail := len(sd.str) - 1
//	for tail >= sd.scan {
//		if !isSpace(sd.str[tail]) {
//			break
//		}
//		tail--
//	}
//	sd.str = sd.str[:tail+1]
//}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 前提：dd.str 肯定是 { 字符后面的字符串
//// 返回 } 后面字符的 index
//func (sd *subDecode) skipObject(pos int) int {
//	var hasKV bool
//	for pos < len(sd.str) {
//		c := sd.str[pos]
//		pos++
//
//		if isSpace(c) {
//			continue
//		}
//		if c == '}' {
//			return pos
//		}
//
//		// 上面都不是，只能是k:v
//		if hasKV {
//			// 不是第一个kv, 先跳过一个逗号
//			pos = sd.skipSeparator(',')
//			if pos < 0 {
//				return pos
//			}
//			pos += 1
//		}
//		pos = sd.skipKVItem(pos - 1)
//		if pos < 0 {
//			return pos
//		}
//		hasKV = true
//	}
//	return -1
//}
//
//// 必须是k:v, ...形式
//func (sd *subDecode) skipKVItem(pos int) int {
//	// A: 找冒号 +++
//	colonPos := sd.skipSeparator(':')
//	// 1. 没有冒号，可能就是一段空白字符
//	if colonPos < 0 {
//		return colonPos
//	}
//
//	// 冒号找到了，看看后面的值
//	pos = colonPos + 1
//	for pos < len(sd.str) {
//		c := sd.str[pos]
//		pos++
//
//		if isSpace(c) {
//			continue
//		}
//		if c == '{' {
//			return sd.skipObject(pos)
//		}
//		if c == '[' {
//			return sd.skipArray(pos)
//		}
//
//		// one value
//		if c == '"' {
//			return sd.skipQuoteValue(pos)
//		}
//		return sd.skipNoQuoteValue(pos - 1)
//	}
//	return -1
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 找冒号："k":"v" | "k":[ |  "k":{
//func (sd *subDecode) skipColon(pos int) int {
//	quoteCt := 0
//	for pos < len(sd.str) {
//		c := sd.str[pos]
//		pos++
//
//		if c == '"' {
//			// 非第一个"，如果后面的"前面有\，是允许的
//			if quoteCt == 1 {
//				if sd.str[pos-1] == '\\' {
//					pos++
//					continue
//				}
//			}
//			quoteCt++
//			if quoteCt > 2 {
//				return errChar
//			}
//		}
//		if quoteCt == 2 {
//			if c == ':' {
//				return pos
//			}
//		}
//	}
//	return noErr
//}

//// 找逗号："k":"v",
//func (sd *subDecode) skipComma(pos int) int {
//	ckSpace := true
//	for pos < len(sd.str) {
//		c := sd.str[pos]
//		if c == ',' {
//			return pos
//		}
//		pos++
//
//		if ckSpace {
//			if isSpace(c) {
//				continue
//			}
//
//			ckSpace = false
//			if c == '{' {
//				pos = sd.skipObject(pos)
//				if pos < 0 {
//					return pos
//				}
//				sd.directString = true
//				return sd.skipComma(pos)
//			}
//			if c == '[' {
//				pos = sd.skipArray(pos)
//				if pos < 0 {
//					return pos
//				}
//				sd.directString = true
//				return sd.skipComma(pos)
//			}
//		}
//	}
//	return -1
//}
//
//// 前提：dd.str 肯定是 [ 字符后面的字符串
//// 返回 ] 后面字符的 index
//func (sd *subDecode) skipArray(pos int) int {
//	for pos < len(sd.str) {
//		c := sd.str[pos]
//		pos++
//
//		if c == ']' {
//			return pos
//		}
//		if isSpace(c) {
//			continue
//		}
//
//		pos = sd.skipArrItem(pos)
//	}
//	return -1
//}
//
//func (sd *subDecode) skipArrItem(pos int) int {
//
//	return -1
//}
//
//func (sd *subDecode) skipQuoteValue(pos int) int {
//	for pos < len(sd.str) {
//		c := sd.str[pos]
//		pos++
//
//		if c == '"' {
//			if sd.str[pos-1] == '\\' {
//				continue
//			}
//			return pos
//		}
//	}
//	return -1
//}

//// 跳过一个逗号
//func (dd *subDecode) skipComma(pos int) int {
//	for pos < len(dd.str) {
//		c := dd.str[pos]
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

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (sd *subDecode) skipMatch(pos int, match string) int {
//	end := pos + len(match)
//	if end > len(sd.str) {
//		return -1
//	}
//	if sd.str[pos:end] == match {
//		return end
//	}
//	return -1
//}
//
//// 匹配一个数值，不带有正负号前缀。
//// 0.234 | 234.23 | 23424 | 3.8e+07
//func (sd *subDecode) skipNumber(pos int) int {
//	if pos >= len(sd.str) {
//		return -1
//	}
//
//	var startZero, hasDot, needNum bool
//
//	c := sd.str[pos]
//	if c < '0' || c > '9' {
//		return -1
//	}
//
//	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
//	if c == '0' {
//		startZero = true
//	}
//	pos++
//
//	for pos < len(sd.str) {
//		c = sd.str[pos]
//		pos++
//
//		if startZero {
//			switch c {
//			case '.', 'e', 'E':
//				startZero = false
//			default:
//				return pos - 1
//			}
//		}
//
//		if c == '.' {
//			if hasDot == true {
//				return -1
//			}
//			hasDot = true
//			needNum = true
//		} else if c == 'e' || c == 'E' {
//			if needNum {
//				return -1
//			}
//			return sd.skipScientificNumberTail(pos)
//		} else if c < '0' || c > '9' {
//			pos--
//			break
//		} else {
//			needNum = false
//		}
//	}
//
//	if needNum {
//		return -1
//	}
//	return pos
//}
//
//// 检查科学计数法（e|E）后面的字符串合法性
//func (sd *subDecode) skipScientificNumberTail(pos int) int {
//	if pos >= len(sd.str) {
//		return scanEOF
//	}
//
//	c := sd.str[pos]
//	if c == '-' || c == '+' {
//		pos++
//	}
//
//	// TODO: 加减号后面没有任何数字匹配，会如何
//	for pos < len(sd.str) {
//		if sd.str[pos] < '0' || sd.str[pos] > '9' {
//			return pos
//		}
//		pos++
//	}
//	return pos
//}
//
//func (sd *subDecode) skipNoQuoteValue(pos int) int {
//	if pos >= len(sd.str) {
//		return -1
//	}
//
//	switch c := sd.str[pos]; {
//	case c >= '0' && c <= '9':
//		// 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7
//		return sd.skipNumber(pos)
//	case c == '-':
//		// -0.3 | -3.7E-7
//		pos++
//		return sd.skipNumber(pos)
//	case c == 'f':
//		return sd.skipMatch("false")
//	case c == 't':
//		return sd.skipMatch("true")
//	case c == 'n':
//		return sd.skipMatch("null")
//	}
//	return -1
//}

//// 没有双引号围起来的值，只可能是：数值、true、false、null
//func checkNoQuoteValue(str string) int {
//	if len(str) == 0 {
//		return errJson
//	}
//
//	// 只有这几种首字符的可能
//isNumber:
//	c := str[0]
//	if c >= '0' && c <= '9' {
//		// 0 | 0.234 | 234.23 | 23424 | 3.8e+07
//		dot := false
//		if c == '0' {
//			if len(str) < 3 || str[1] != '.' {
//				return errChar
//			}
//		}
//		for i := 1; i < len(str); i++ {
//			c = str[i]
//			if c == '.' {
//				if dot == true {
//					return errNumberFmt
//				} else {
//					dot = true
//					continue
//				}
//			} else if c == 'e' || c == 'E' {
//				return checkScientificNumberTail(str[i+1:])
//			} else if c < '0' || c > '9' {
//				return errNumberFmt
//			}
//		}
//	} else if c == '-' {
//		// -0 | -0.3 | +13.33 | -3.7E-7
//		if len(str) >= 2 && str[1] >= '0' && str[1] <= '9' {
//			str = str[1:]
//			goto isNumber
//		} else {
//			return errNumberFmt
//		}
//	} else if c == 'f' {
//		// false
//		if str != "false" {
//			return errChar
//		}
//	} else if c == 't' {
//		// true
//		if str != "true" {
//			return errChar
//		}
//	} else if c == 'n' {
//		if str != "null" {
//			return errChar
//		}
//	} else {
//		return errJson
//	}
//	return noErr
//}

//// 检查科学计数法（e|E）后面的字符串合法性
//func checkScientificNumberTail(str string) int {
//	if len(str) == 0 {
//		return errNumberFmt
//	}
//	c := str[0]
//	if c == '-' || c == '+' {
//		str = str[1:]
//	}
//
//	if len(str) == 0 {
//		return errNumberFmt
//	}
//	for i := range str {
//		if str[i] < '0' || str[i] > '9' {
//			return errNumberFmt
//		}
//	}
//	return noErr
//}

//// 找逗号："k":"v",
//func (sd *subDecode) scanComma() (ret int) {
//	ckBlank := true
//	for sd.scan < len(sd.str) {
//		c := sd.str[sd.scan]
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
//				sd.directString = true
//				return sd.scanComma()
//			}
//			if c == '[' {
//				ret = sd.scanArray()
//				if ret < 0 {
//					return ret
//				}
//				sd.directString = true
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
//		tmp := trim(sd.str[sd.scan:])
//		if len(tmp) != 0 {
//			return errChar
//		}
//		return noErr
//	} else if colonPos < 0 {
//		return colonPos
//	}
//
//	key := sd.str[sd.scan:colonPos]
//	sd.scan = colonPos
//	// 2. TODO：这里冒号前面的Key其实就可以得到了，可以先判断目标对象是否有这个key，没有value都不用解析了，直接解析下一个
//
//	// B: 找逗号 +++
//	commaPos := sd.scanComma(colonPos + 1) // 从冒号后面开始查找第一个匹配的逗号
//	val := ""
//	sd.directString = false
//	// 1. 没找到逗号，这是最后一个k:v了
//	if commaPos == errNotFound {
//		sd.scan = len(sd.str)
//		val = sd.str[colonPos+1 : sd.scan]
//	} else if commaPos <= 0 {
//		return commaPos
//	}
//	// 2. 找到一个“,” 其前面部分当做一个k:v来解
//	if commaPos > 0 {
//		val = sd.str[colonPos+1 : commaPos]
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
//	key := sd.str[sd.scan:colonPos]
//	val := sd.str[colonPos+1 : commaPos]
//
//	if sd.gr != nil {
//		return sd.setGsonValue(key, val)
//	}
//	return sd.setMapValue(key, val)
//}

//func (sd *subDecode) scanQuoteValue(pos int) int {
//	for pos < len(sd.str) {
//		c := sd.str[pos]
//		pos++
//
//		if c == '"' {
//			if sd.str[pos-1] == '\\' {
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

//
//// 核查 key +++++++++
//func (sd *subDecode) checkKey(key string, needCopy bool) (ret string, err int) {
//	var inShare bool
//	//if key, err = cutKeyQuote(key); err < 0 {
//	//	return
//	//}
//	if key, inShare, err = sd.unescapeString(key); err < 0 {
//		return
//	}
//	if inShare && needCopy {
//		key = cloneString(key)
//	}
//	return key, noErr
//}

//// 检查是否是一个合法的 value 值
//func cutValueQuote(val string) (ret string, hasQuote bool, err int) {
//	val = trim(val)
//
//	if len(val) < 1 {
//		return "", false, errChar
//	}
//
//	// 如果 value 没有 双引号，可能是数值、true、false、null四种情况
//	if val[0] != '"' {
//		if err = checkNoQuoteValue(val); err < 0 {
//			return "", false, err
//		}
//		return val, false, noErr
//	}
//
//	// 有双引号
//	if len(val) == 1 {
//		return "", false, errChar
//	} else if val[len(val)-1] != '"' {
//		return "", false, errChar
//	}
//	return val[1 : len(val)-1], true, noErr
//}

//
//// 核查 value ++++++++++
//func (sd *subDecode) checkValue(val string) (ret string, err int) {
//	var hasQuote, inShare bool
//	if val, hasQuote, err = cutValueQuote(val); err < 0 {
//		return
//	}
//	if hasQuote {
//		if val, inShare, err = sd.unescapeString(val); err < 0 {
//			return
//		}
//		// 证明此时用的是临时栈空间保存转义之后的Value，需要申请新的内存空间放置
//		if inShare {
//			val = cloneString(val)
//		}
//	}
//	return val, noErr
//}

//// 当目标为 cst.KV 类型时候，用此方法设置
//func (sd *subDecode) setMapValue(key, val string) (err int) {
//	// set k = v
//	sd.dst.Set(key, val)
//	return noErr
//}

//// 当目标为 gson.GsonRow 类型时候，用此方法设置
//func (sd *subDecode) setGsonValue(key, val string) (err int) {
//	keyIdx := sd.gr.KeyIndex(key)
//	// 没有这个字段，直接返回了(此时再去解析后面的value是没有意义的)
//	if keyIdx < 0 {
//		return noErr
//	}
//
//	// set k = v
//	sd.gr.SetStringByIndex(keyIdx, val)
//	return noErr
//}
