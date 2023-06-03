package decode

import (
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

//
//// 肯定是 "*?" 的字符串
//// 这是零新增内存方案，直接移动原始[]byte。还可以用共享内存方案实现
//func (sd *subDecode) unescapeString(start, end int) (val string) {
//	str := sd.str[start:end]
//
//	var bs []byte
//	var slash bool
//	var pos, ct int
//	end = 2 // 只是复用end变量
//
//	for i := 0; i < len(str); i++ {
//		c := str[i]
//
//		if c == '\\' {
//			if slash == false {
//				bs = lang.STB(str[i:])[0:0]
//				end = i
//			}
//			slash = true
//
//			if ct > 0 {
//				bs = append(bs, str[pos:pos+ct]...)
//				pos = 0
//				ct = 0
//			}
//
//			i++
//			switch c = str[i]; c {
//			case '"', '\\', '/':
//				bs = append(bs, c)
//			case 'b':
//				bs = append(bs, '\b')
//			case 'f':
//				bs = append(bs, '\f')
//			case 't':
//				bs = append(bs, '\t')
//			case 'n':
//				bs = append(bs, '\n')
//			case 'r':
//				bs = append(bs, '\r')
//			case 'u':
//				// TODO: uft8编码字符有待转换
//				//bs = append(bs, '\u0233')
//			default:
//				//case '\'': // 这种情况认为是错误
//				sd.scan = start + i
//				panic(errChar)
//			}
//
//			//// Quote, control characters are invalid.
//			//case c == '"', c < ' ':
//			//	return
//			//
//			//	// ASCII
//			//case c < utf8.RuneSelf:
//			//	b[w] = c
//			//	r++
//			//	w++
//			//
//			//// Coerce to well-formed UTF-8.
//			//default:
//			//	rr, size := utf8.DecodeRune(s[r:])
//			//	r += size
//			//	w += utf8.EncodeRune(b[w:], rr)
//
//			continue
//		}
//
//		if slash {
//			if ct == 0 {
//				pos = i
//			}
//			ct++
//		}
//	}
//
//	if ct > 0 {
//		bs = append(bs, str[pos:pos+ct]...)
//		pos = 0
//		ct = 0
//	}
//
//	end += len(bs)
//	for i := end; i < len(str); i++ {
//		bs = append(bs, ' ') // 填充空格
//	}
//	return str[1 : end-1]
//}
//
//func (sd *subDecode) unescapeCopy(str string) (ret string, inShare bool) {
//	var newStr []byte
//	var step int
//	var hasSlash bool
//
//	for i := 0; i < len(str); i++ {
//		c := str[i]
//		if c == '\\' {
//			// 第一次检索到有 \
//			if hasSlash == false {
//				//hasSlash = true
//				//// TODO：这里发生了逃逸，需要用sync.Pool的方式，共享内存空间
//				//// 或者别的黑魔法操作内存
//				//// add by sdx 20230404 动态初始化 share 内存
//				//if sd.share == nil {
//				//	defSize := len(str)
//				//	if defSize > tempByteStackSize {
//				//		defSize = tempByteStackSize
//				//	}
//				//	sd.share = make([]byte, defSize)
//				//}
//				//
//				//if len(str) <= len(sd.share) {
//				//	newStr = sd.share[:]
//				//	inShare = true
//				//} else {
//				//	newStr = make([]byte, len(str))
//				//}
//				//for ; step < i; step++ {
//				//	newStr[step] = str[step]
//				//}
//			}
//			i++
//			c = str[i]
//			// 判断 \ 后面的字符
//			switch c {
//			case '"', '/', '\\':
//				newStr[step] = c
//				step++
//			//case '\'': // 这种情况认为是错误
//			case 'b':
//				newStr[step] = '\b'
//				step++
//			case 'f':
//				newStr[step] = '\f'
//				step++
//			case 't':
//				newStr[step] = '\t'
//				step++
//			case 'n':
//				newStr[step] = '\n'
//				step++
//			case 'r':
//				newStr[step] = '\r'
//				step++
//			case 'u': // TODO: uft8编码字符有待转换
//			default:
//				panic(errJson)
//			}
//			continue
//		}
//		if hasSlash {
//			newStr[step] = c
//			step++
//		}
//		// ASCII
//		//case c < utf8.RuneSelf:
//		//	b[w] = c
//		//	r++
//		//	w++
//		//
//		//	// Coerce to well-formed UTF-8.
//		//	default:
//		//	rr, size := utf8.DecodeRune(s[r:])
//		//	r += size
//		//	w += utf8.EncodeRune(b[w:], rr)
//	}
//	if hasSlash {
//		return lang.BTS(newStr[:step]), inShare
//	}
//	return str, false
//}
//
//func cloneString(src string) string {
//	tmp := make([]byte, len(src))
//	copy(tmp, src)
//	return lang.BTS(tmp)
//}

//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 字符串处理
//func sliceSetString(val []string, sd *subDecode) {
//	ptrLevel := sd.dm.ptrLevel
//
//	// 如果绑定对象是字符串切片
//	newArr := make([]string, len(val))
//	copy(newArr, val)
//	if ptrLevel <= 0 {
//		*(sd.dst.(*[]string)) = newArr
//		return
//	}
//
//	// 一级指针
//	ptrLevel--
//	ret1 := copySlice[string](sd, ptrLevel, newArr)
//	if ret1 == nil {
//		return
//	}
//
//	// 二级指针
//	ptrLevel--
//	ret2 := copySlice[*string](sd, ptrLevel, ret1)
//	if ret2 == nil {
//		return
//	}
//
//	// 三级指针
//	ptrLevel--
//	_ = copySlice[**string](sd, ptrLevel, ret2)
//	return
//}
//
//// Bool处理
//func sliceSetBool(val []bool, sd *subDecode) {
//	ptrLevel := sd.dm.ptrLevel
//
//	// 如果绑定对象是字符串切片
//	newArr := make([]bool, len(val))
//	copy(newArr, val)
//	if ptrLevel <= 0 {
//		*(sd.dst.(*[]bool)) = newArr
//		return
//	}
//
//	// 一级指针
//	ptrLevel--
//	ret1 := copySlice[bool](sd, ptrLevel, newArr)
//	if ret1 == nil {
//		return
//	}
//
//	// 二级指针
//	ptrLevel--
//	ret2 := copySlice[*bool](sd, ptrLevel, ret1)
//	if ret2 == nil {
//		return
//	}
//
//	// 三级指针
//	ptrLevel--
//	_ = copySlice[**bool](sd, ptrLevel, ret2)
//	return
//}
//
//func sliceSetAny(val []any, sd *subDecode) {
//	ptrLevel := sd.dm.ptrLevel
//
//	// 如果绑定对象是字符串切片
//	newArr := make([]any, len(val))
//	copy(newArr, val)
//	if ptrLevel <= 0 {
//		*(sd.dst.(*[]any)) = newArr
//		return
//	}
//
//	// 一级指针
//	ptrLevel--
//	ret1 := copySlice[any](sd, ptrLevel, newArr)
//	if ret1 == nil {
//		return
//	}
//
//	// 二级指针
//	ptrLevel--
//	ret2 := copySlice[*any](sd, ptrLevel, ret1)
//	if ret2 == nil {
//		return
//	}
//
//	// 三级指针
//	ptrLevel--
//	_ = copySlice[**any](sd, ptrLevel, ret2)
//	return
//}

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

// 来自标准库 json/decode_xxx.go的函数
func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}
