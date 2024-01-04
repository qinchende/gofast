package jde

import (
	"github.com/qinchende/gofast/core/rt"
	"unicode/utf8"
	"unsafe"
)

// TODO: 这里都还没有考虑unicode字符的情况
// 这是一种不安全的处理方式，直接修改了原始字节数组。
func (sd *subDecode) unescapeEnd() int {
	pos := sd.escPos[0]

	var bs []byte
	sh := (*rt.SliceHeader)(unsafe.Pointer(&bs))
	sh.DataPtr = (*(*rt.StringHeader)(unsafe.Pointer(&sd.str))).DataPtr
	sh.Len, sh.Cap = pos, sd.scan

	for i := 0; i < len(sd.escPos); i++ {
		if pos < sd.escPos[i] {
			bs = append(bs, sd.str[pos:sd.escPos[i]]...)
		}

		pos = sd.escPos[i] + 1
		c := unescapeChar[sd.str[pos]]
		// 如果\ 后面跟着的不是特定字符，肯定有问题的
		if c == 0 {
			sd.scan = pos
			panic(errEscape)
		}

		// ++++++++++++++++++++++++++++++++++++++++
		// 开始特殊字符处理，比如：\u0233
		if c == 'u' {
			if pos+6 > sd.scan {
				sd.scan = pos
				panic(errEscape)
			}

			v1 := hexToInt[sd.str[pos+1]]
			v2 := hexToInt[sd.str[pos+2]]
			v3 := hexToInt[sd.str[pos+3]]
			v4 := hexToInt[sd.str[pos+4]]
			code := rune((v1 << 12) | (v2 << 8) | (v3 << 4) | v4)
			pos += 4

			if code >= 0xd800 && code < 0xdc00 {
				if pos+6 > sd.scan {
					sd.scan = pos
					panic(errEscape)
				}
				if sd.str[pos+1] == '\\' && sd.str[pos+2] == 'u' {
					v1 = hexToInt[sd.str[pos+3]]
					v2 = hexToInt[sd.str[pos+4]]
					v3 = hexToInt[sd.str[pos+5]]
					v4 = hexToInt[sd.str[pos+6]]
					lo := rune((v1 << 12) | (v2 << 8) | (v3 << 4) | v4)
					pos += 6

					if lo >= 0xdc00 && lo < 0xe000 {
						code = (code-0xd800)<<10 | (lo - 0xdc00) + 0x10000
					}
				}
			}

			var b [utf8.UTFMax]byte
			n := utf8.EncodeRune(b[:], code)

			bs = append(bs, b[:n]...)
			pos++

			continue
		}
		// ++++++++++++++++++++++++++++++++++++++++

		bs = append(bs, c)
		pos++
	}

	if pos < sd.scan {
		bs = append(bs, sd.str[pos:sd.scan]...)
	}

	return len(bs) - 1
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
//
//// 来自标准库 json/decode_xxx.go的函数 ，放在这里供参考
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
