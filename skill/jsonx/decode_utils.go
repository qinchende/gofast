package jsonx

import (
	"github.com/qinchende/gofast/skill/lang"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

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

const (
	spaceCharMask = (1 << ' ') | (1 << '\t') | (1 << '\r') | (1 << '\n')
)

func isSpace(c byte) bool {
	return spaceCharMask&(1<<c) != 0
}

//
//func trimHead(str string) int {
//	i := 0
//	for ; i < len(str); i++ {
//		if !isSpace(str[i]) {
//			break
//		}
//	}
//	return i
//}
//
//func trimTail(str string) int {
//	tail := len(str) - 1
//	for ; tail >= 0; tail-- {
//		if !isSpace(str[tail]) {
//			break
//		}
//	}
//	return tail
//}

func trim(str string) string {
	s := 0
	e := len(str) - 1
	for s < e {
		c := str[s]
		if !isSpace(c) {
			break
		}
		s++
	}
	for s < e {
		c := str[e]
		if !isSpace(c) {
			break
		}
		e--
	}
	return str[s : e+1]
}

func cloneString(src string) string {
	tmp := make([]byte, len(src))
	copy(tmp, src)
	return lang.BTS(tmp)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 一个合法的 Key，或者Value 字符串

func (sd *subDecode) unescapeCopy(str string) (ret string, inShare bool, err int) {
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
		// ASCII
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
