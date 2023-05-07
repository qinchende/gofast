package jde

import "strconv"

// fast number value parser
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var (
	pow10u64 = [...]uint64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	}
	pow10u64Len = len(pow10u64)
)

func parseUint(s string) uint64 {
	maxDigit := len(s)
	if maxDigit > pow10u64Len {
		panic(errNumberFmt)
	}
	sum := uint64(0)
	for i := 0; i < maxDigit; i++ {
		c := uint64(s[i]) - 48
		digitValue := pow10u64[maxDigit-i-1]
		sum += c * digitValue
	}
	return sum
}

var (
	pow10i64 = [...]int64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}
	pow10i64Len = len(pow10i64)
)

func parseInt(s string) int64 {
	isNegative := false
	if s[0] == '-' {
		s = s[1:]
		isNegative = true
	}
	maxDigit := len(s)
	if maxDigit > pow10i64Len {
		panic(errNumberFmt)
	}
	sum := int64(0)
	for i := 0; i < maxDigit; i++ {
		c := int64(s[i]) - 48
		digitValue := pow10i64[maxDigit-i-1]
		sum += c * digitValue
	}
	if isNegative {
		return -1 * sum
	}
	return sum
}

//go:inline
func parseFloat(s string) float64 {
	if f64, err := strconv.ParseFloat(s, 64); err != nil {
		panic(errNumberFmt)
	} else {
		return f64
	}
}

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
