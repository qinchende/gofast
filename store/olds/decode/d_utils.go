package decode

//func (sd *subDecode) scanQuoteStrValue() {
//	pos := sd.scan
//
//	if sd.skipValue {
//		for {
//			pos++
//			switch c := sd.str[pos]; {
//			case c == '"':
//				sd.scan = pos + 1
//				return
//			case c == '\\':
//				pos++ // 跳过 '\' 后面的一个字符
//			}
//		}
//	}
//
//	pos++
//	slash := sd.scanQuoteStr()
//	if slash {
//		sd.bindString(sd.str[pos:sd.unescapeEnd()])
//	} else {
//		sd.bindString(sd.str[pos : sd.scan-1])
//	}
//}
//
//func (sd *subDecode) scanStrKindValue() {
//	switch sd.str[sd.scan] {
//	case '"':
//		sd.scanQuoteStrValue()
//	default:
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindStringNull()
//	}
//}

//func (sd *subDecode) scanBoolValue() {
//	switch sd.str[sd.scan] {
//	case 't':
//		sd.skipTrue()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBool(true)
//		return
//	case 'f':
//		sd.skipFalse()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBool(false)
//	default:
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBoolNull()
//	}
//}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 匹配一个数值，对应于float类型
//// 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7 | -0.3 | -3.7E-7
//func (sd *subDecode) scanNumValue() {
//	pos := sd.scan
//	start := pos
//	var hasDot, needNum bool
//
//	c := sd.str[pos]
//	if c == '-' {
//		pos++
//		c = sd.str[pos]
//	}
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
//			needNum = false // 到这里，字符肯定是数字
//		}
//	}
//
//	if needNum {
//		panic(errNumberFmt)
//	}
//
//over:
//	sd.scan = pos
//	// 还剩下最后一种可能：null
//	if start == pos {
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindNumberNull()
//		return
//	}
//	if sd.skipValue {
//		return
//	}
//	sd.bindNumber(sd.str[start:pos])
//}

//func (sd *subDecode) scanIntValue() {
//	pos := sd.scan
//	start := pos
//
//	c := sd.str[pos]
//	if c == '-' {
//		pos++
//		c = sd.str[pos]
//	}
//	if c == '0' {
//		pos++
//		goto over
//	}
//	for {
//		if c < '0' || c > '9' {
//			break
//		}
//		pos++
//		c = sd.str[pos]
//	}
//over:
//	sd.scan = pos
//	// 还剩下最后一种可能：null +++
//	if start == pos {
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindIntNull()
//		return
//	}
//	if sd.skipValue {
//		return
//	}
//	sd.bindIntList(sd.str[start:pos])
//}

//func (sd *subDecode) scanUintValue() {
//	pos := sd.scan
//	start := pos
//
//	c := sd.str[pos]
//	if c == '0' {
//		pos++
//		goto over
//	}
//	for {
//		if c < '0' || c > '9' {
//			break
//		}
//		pos++
//		c = sd.str[pos]
//	}
//over:
//	sd.scan = pos
//	// 还剩下最后一种可能：null
//	if start == pos {
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindUintNull()
//		return
//	}
//	if sd.skipValue {
//		return
//	}
//	sd.bindUintList(sd.str[start:pos])
//}

//func (sd *subDecode) skipQuoteString() {
//	pos := sd.scan
//	//if sd.str[pos] != '"' {
//	//	panic(errChar)
//	//}
//	for {
//		pos++
//		switch c := sd.str[pos]; {
//		case c == '"':
//			sd.scan = pos + 1
//			return
//		case c == '\\':
//			pos++
//		}
//	}
//}

//// 检查科学计数法（e|E）后面的字符串合法性
//func (sd *subDecode) scanScientificNumberTail() int {
//	c := sd.str[sd.scan]
//	if c == '-' || c == '+' {
//		sd.scan++
//	}
//
//	for {
//		if c = sd.str[sd.scan]; c < '0' || c > '9' {
//			return noErr
//		}
//		sd.scan++
//	}
//}

//func trimHead(str string) int {
//	i := 0
//	for ; i < len(str); i++ {
//		if !isBlank(str[i]) {
//			break
//		}
//	}
//	return i
//}
//
//func trimTail(str string) int {
//	tail := len(str) - 1
//	for ; tail >= 0; tail-- {
//		if !isBlank(str[tail]) {
//			break
//		}
//	}
//	return tail
//}
//
//func trim(str string) string {
//	s := 0
//	e := len(str) - 1
//	for s < e {
//		c := str[s]
//		if !isBlank(c) {
//			break
//		}
//		s++
//	}
//	for s < e {
//		c := str[e]
//		if !isBlank(c) {
//			break
//		}
//		e--
//	}
//	return str[s : e+1]
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 一个合法的 Key，或者Value 字符串

//// 来自标准库 json/decode_xxx.go的函数
//func unquoteBytes(s []byte) (t []byte, ok bool) {
//	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
//		return
//	}
//	s = s[1 : len(s)-1]
//
//	// Check for unusual characters. If there are none,
//	// then no unquoting is needed, so return a slice of the
//	// original bytes.
//	r := 0
//	for r < len(s) {
//		c := s[r]
//		if c == '\\' || c == '"' || c < ' ' {
//			break
//		}
//		if c < utf8.RuneSelf {
//			r++
//			continue
//		}
//		rr, size := utf8.DecodeRune(s[r:])
//		if rr == utf8.RuneError && size == 1 {
//			break
//		}
//		r += size
//	}
//	if r == len(s) {
//		return s, true
//	}
//
//	b := make([]byte, len(s)+2*utf8.UTFMax)
//	w := copy(b, s[0:r])
//	for r < len(s) {
//		// Out of room? Can only happen if s is full of
//		// malformed UTF-8 and we're replacing each
//		// byte with RuneError.
//		if w >= len(b)-2*utf8.UTFMax {
//			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
//			copy(nb, b[0:w])
//			b = nb
//		}
//		switch c := s[r]; {
//		case c == '\\':
//			r++
//			if r >= len(s) {
//				return
//			}
//			switch s[r] {
//			default:
//				return
//			case '"', '\\', '/', '\'':
//				b[w] = s[r]
//				r++
//				w++
//			case 'b':
//				b[w] = '\b'
//				r++
//				w++
//			case 'f':
//				b[w] = '\f'
//				r++
//				w++
//			case 'n':
//				b[w] = '\n'
//				r++
//				w++
//			case 'r':
//				b[w] = '\r'
//				r++
//				w++
//			case 't':
//				b[w] = '\t'
//				r++
//				w++
//			case 'u':
//				r--
//				rr := getu4(s[r:])
//				if rr < 0 {
//					return
//				}
//				r += 6
//				if utf16.IsSurrogate(rr) {
//					rr1 := getu4(s[r:])
//					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
//						// A valid pair; consume.
//						r += 6
//						w += utf8.EncodeRune(b[w:], dec)
//						break
//					}
//					// Invalid surrogate; fall back to replacement rune.
//					rr = unicode.ReplacementChar
//				}
//				w += utf8.EncodeRune(b[w:], rr)
//			}
//
//		// Quote, control characters are invalid.
//		case c == '"', c < ' ':
//			return
//
//		// ASCII
//		case c < utf8.RuneSelf:
//			b[w] = c
//			r++
//			w++
//
//		// Coerce to well-formed UTF-8.
//		default:
//			rr, size := utf8.DecodeRune(s[r:])
//			r += size
//			w += utf8.EncodeRune(b[w:], rr)
//		}
//	}
//	return b[0:w], true
//}
//
//// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
//// or it returns -1.
//func getu4(s []byte) rune {
//	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
//		return -1
//	}
//	var r rune
//	for _, c := range s[2:6] {
//		switch {
//		case '0' <= c && c <= '9':
//			c = c - '0'
//		case 'a' <= c && c <= 'f':
//			c = c - 'a' + 10
//		case 'A' <= c && c <= 'F':
//			c = c - 'A' + 10
//		default:
//			return -1
//		}
//		r = r*16 + rune(c)
//	}
//	return r
//}
