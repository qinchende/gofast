package jde

import "unicode/utf8"

const hex = "0123456789abcdef"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// NOTE：这里的String不会在前后添加""字符
func addStrNoQuotes(bs []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		if !noEscapeTable[s[i]] {
			return appendComplexStr(bs, s, i)
		}
	}
	return append(bs, s...)
}

func addStrQuotes(bs []byte, s string, c byte) []byte {
	bs = append(bs, '"')
	for i := 0; i < len(s); i++ {
		if !noEscapeTable[s[i]] {
			bs = appendComplexStr(bs, s, i)
			return append(bs, '"', c)
		}
	}
	bs = append(bs, s...)
	return append(bs, '"', c)
}

func appendComplexStr(bs []byte, s string, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				if start < i {
					bs = append(bs, s[start:i]...)
				}
				bs = append(bs, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if noEscapeTable[b] {
			i++
			continue
		}
		if start < i {
			bs = append(bs, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			bs = append(bs, '\\', b)
		case '\b':
			bs = append(bs, '\\', 'b')
		case '\f':
			bs = append(bs, '\\', 'f')
		case '\n':
			bs = append(bs, '\\', 'n')
		case '\r':
			bs = append(bs, '\\', 'r')
		case '\t':
			bs = append(bs, '\\', 't')
		default:
			bs = append(bs, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		bs = append(bs, s[start:]...)
	}
	return bs
}
