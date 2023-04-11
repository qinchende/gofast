package jsonx

import (
	"github.com/pkg/profile"
	"github.com/qinchende/gofast/skill/lang"
)

func (sd *subDecode) unescapeString(start, end int) (val string, err int) {
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	// 或者直接生成文件

	str := sd.sub[start:end]

	//var bytes []byte
	bytes := lang.STB(str)

	var slash bool
	var off, pos, ct int

	for i := 0; i < len(bytes); i++ {
		c := bytes[i]

		if c == '\\' {
			if slash == false {
				//bs = bytes[i:][0:0]
				off = i
			}
			slash = true

			if ct > 0 {
				//copy(bytes[off:off+ct], bytes[pos:pos+ct])
				//bs = append(bs, bytes[pos:pos+ct]...)
				off += ct
				pos = 0
				ct = 0
			}

			i++
			switch c = bytes[i]; c {
			case '"', '\\', '/':
				bytes[off] = c
				off++
				//bs = append(bs, c)
			case 'b':
				bytes[off] = '\b'
				off++
				//bs = append(bs, '\b')
			case 'f':
				bytes[off] = '\f'
				off++
				//bs = append(bs, '\f')
			case 't':
				bytes[off] = '\t'
				off++
				//bs = append(bs, '\t')
			case 'n':
				bytes[off] = '\n'
				off++
				//bs = append(bs, '\n')
			case 'r':
				bytes[off] = '\r'
				off++
				//bs = append(bs, '\r')
			case 'u':
				// TODO: uft8编码字符有待转换
				// bs = append(bs, '\u0233')
			default:
				//case '\'': // 这种情况认为是错误
				sd.scan = start + i
				return "", errChar
			}

			//// Quote, control characters are invalid.
			//case c == '"', c < ' ':
			//	return
			//
			//	// ASCII
			//case c < utf8.RuneSelf:
			//	b[w] = c
			//	r++
			//	w++
			//
			//// Coerce to well-formed UTF-8.
			//default:
			//	rr, size := utf8.DecodeRune(s[r:])
			//	r += size
			//	w += utf8.EncodeRune(b[w:], rr)

			continue
		}

		if slash {
			//bs = append(bs, c)
			bytes[off] = c
			off++
			//if ct == 0 {
			//	pos = i
			//}
			//ct++
		}
	}

	if ct > 0 {
		//copy(bytes[off:off+ct], bytes[pos:pos+ct])
		//bs = append(bs, bytes[pos:pos+ct]...)
		off += ct
		pos = 0
		ct = 0
	}

	//if len(str) > 0 {
	//	return str, noErr
	//}

	//off = off + len(bs)
	//for i := len(bs); i < cap(bs); i++ {
	//	bs = append(bs, ' ')
	//}
	end = off - 1
	for ; off < len(bytes); off++ {
		//bs = append(bs, ' ')
		bytes[off] = ' '
	}

	if pos < 0 {

	}
	return str[1:end], noErr
}
