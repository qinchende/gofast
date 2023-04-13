package jde

import "github.com/qinchende/gofast/skill/lang"

// 肯定是 "*?" 的字符串
// 这是零新增内存方案，还可以用共享内存方案实现
func (sd *subDecode) unescapeString(start, end int) (val string, err int) {
	str := sd.str[start:end]

	var bs []byte
	var slash bool
	var pos, ct int
	end = 2 // 只是复用end变量

	for i := 0; i < len(str); i++ {
		c := str[i]

		if c == '\\' {
			if slash == false {
				bs = lang.STB(str[i:])[0:0]
				end = i
			}
			slash = true

			if ct > 0 {
				bs = append(bs, str[pos:pos+ct]...)
				pos = 0
				ct = 0
			}

			i++
			switch c = str[i]; c {
			case '"', '\\', '/':
				bs = append(bs, c)
			case 'b':
				bs = append(bs, '\b')
			case 'f':
				bs = append(bs, '\f')
			case 't':
				bs = append(bs, '\t')
			case 'n':
				bs = append(bs, '\n')
			case 'r':
				bs = append(bs, '\r')
			case 'u':
				// TODO: uft8编码字符有待转换
				//bs = append(bs, '\u0233')
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
			if ct == 0 {
				pos = i
			}
			ct++
		}
	}

	if ct > 0 {
		bs = append(bs, str[pos:pos+ct]...)
		pos = 0
		ct = 0
	}

	end += len(bs)
	for i := end; i < len(str); i++ {
		bs = append(bs, ' ') // 填充空格
	}
	return str[1 : end-1], noErr
}

func (sd *subDecode) unescapeCopy(str string) (ret string, inShare bool, err int) {
	var newStr []byte
	var step int
	var hasSlash bool

	for i := 0; i < len(str); i++ {
		c := str[i]
		if c == '\\' {
			// 第一次检索到有 \
			if hasSlash == false {
				//hasSlash = true
				//// TODO：这里发生了逃逸，需要用sync.Pool的方式，共享内存空间
				//// 或者别的黑魔法操作内存
				//// add by sdx 20230404 动态初始化 share 内存
				//if sd.share == nil {
				//	defSize := len(str)
				//	if defSize > tempByteStackSize {
				//		defSize = tempByteStackSize
				//	}
				//	sd.share = make([]byte, defSize)
				//}
				//
				//if len(str) <= len(sd.share) {
				//	newStr = sd.share[:]
				//	inShare = true
				//} else {
				//	newStr = make([]byte, len(str))
				//}
				//for ; step < i; step++ {
				//	newStr[step] = str[step]
				//}
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
func cloneString(src string) string {
	tmp := make([]byte, len(src))
	copy(tmp, src)
	return lang.BTS(tmp)
}
