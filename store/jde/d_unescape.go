package jde

import (
	"reflect"
	"unsafe"
)

// TODO: 这里都还没有考虑unicode字符的情况
// 这是一种不安全的处理方式，直接修改了原始字节数组。
func (sd *subDecode) unescapeEnd() int {
	pos := sd.escPos[0]

	s := []byte{}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sh.Data = (*(*reflect.StringHeader)(unsafe.Pointer(&sd.str))).Data
	sh.Len, sh.Cap = pos, sd.scan

	for i := 0; i < len(sd.escPos); i++ {
		if pos < sd.escPos[i] {
			s = append(s, sd.str[pos:sd.escPos[i]]...)
		}

		pos = sd.escPos[i] + 1
		c := unescapeChar[sd.str[pos]]
		if c == 0 {
			sd.scan = pos
			panic(errChar)
		}
		//if c == 'u' {
		//	// unicode字符
		//}
		//// ASCII
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

		s = append(s, c)

		pos++
	}

	if pos < sd.scan {
		s = append(s, sd.str[pos:sd.scan]...)
	}

	return len(s) - 1
}
