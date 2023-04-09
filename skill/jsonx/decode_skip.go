package jsonx

func (sd *subDecode) skipHeadSpace() {
	for sd.scan < len(sd.sub) {
		if !isSpace(sd.sub[sd.scan]) {
			return
		}
		sd.scan++
	}
}

func (sd *subDecode) skipTailSpace() {
	tail := len(sd.sub) - 1
	for tail >= sd.scan {
		if !isSpace(sd.sub[tail]) {
			break
		}
		tail--
	}
	sd.sub = sd.sub[:tail+1]
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：dd.sub 肯定是 { 字符后面的字符串
// 返回 } 后面字符的 index
func (sd *subDecode) skipObject(pos int) int {
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
			pos = sd.skipComma(pos - 1)
			if pos < 0 {
				return pos
			}
			pos += 1
		}
		pos = sd.skipKVItem(pos - 1)
		if pos < 0 {
			return pos
		}
		hasKV = true
	}
	return -1
}

// 必须是k:v, ...形式
func (sd *subDecode) skipKVItem(pos int) int {
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
			return sd.skipObject(pos)
		}
		if c == '[' {
			return sd.skipArray(pos)
		}

		// one value
		if c == '"' {
			return sd.skipQuoteValue(pos)
		}
		return sd.skipNoQuoteValue(pos - 1)
	}
	return -1
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 找冒号："k":"v" | "k":[ |  "k":{
func (sd *subDecode) skipColon(pos int) int {
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
func (sd *subDecode) skipComma(pos int) int {
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
				pos = sd.skipObject(pos)
				if pos < 0 {
					return pos
				}
				sd.isMixedVal = true
				return sd.skipComma(pos)
			}
			if c == '[' {
				pos = sd.skipArray(pos)
				if pos < 0 {
					return pos
				}
				sd.isMixedVal = true
				return sd.skipComma(pos)
			}
		}
	}
	return -1
}

// 前提：dd.sub 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) skipArray(pos int) int {
	for pos < len(sd.sub) {
		c := sd.sub[pos]
		pos++

		if c == ']' {
			return pos
		}
		if isSpace(c) {
			continue
		}

		pos = sd.skipArrItem(pos)
	}
	return -1
}

func (sd *subDecode) skipArrItem(pos int) int {

	return -1
}

func (sd *subDecode) skipQuoteValue(pos int) int {
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
//func (dd *subDecode) skipComma(pos int) int {
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) skipMatch(pos int, match string) int {
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
func (sd *subDecode) skipNumber(pos int) int {
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
			return sd.skipScientificNumberTail(pos)
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
func (sd *subDecode) skipScientificNumberTail(pos int) int {
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

func (sd *subDecode) skipNoQuoteValue(pos int) int {
	if pos >= len(sd.sub) {
		return -1
	}

	switch c := sd.sub[pos]; {
	case c >= '0' && c <= '9':
		// 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7
		return sd.skipNumber(pos)
	case c == '-':
		// -0.3 | -3.7E-7
		pos++
		return sd.skipNumber(pos)
	case c == 'f':
		return sd.skipMatch(pos, "false")
	case c == 't':
		return sd.skipMatch(pos, "true")
	case c == 'n':
		return sd.skipMatch(pos, "null")
	}
	return -1
}
