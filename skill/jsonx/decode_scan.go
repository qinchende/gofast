package jsonx

import "github.com/qinchende/gofast/skill/lang"

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (dd *fastDecode) parseJson() error {
	if ok := dd.skipHeadSpace(); !ok {
		return sErr
	}

	c := dd.src[dd.head]
	switch c {
	case '{':
		if ok := dd.skipTailSpace(); !ok {
			return sErr
		}
		if dd.src[dd.tail] != '}' {
			return sErr
		}
		return dd.parseObject()
	case '[':
		if exist := dd.skipTailSpace(); !exist {
			return sErr
		}
		if dd.src[dd.tail] != ']' {
			return sErr
		}
		return dd.parseArray()
	default:
		return sErr
	}
}

func (dd *fastDecode) parseObject() error {
	dd.head++
	dd.tail--

	// 剩下的全部应该是 k:v,k:v,k:{},k:[]
loopComma:
	// TODO: 不应该先找逗号，而应该先找冒号，看冒号后面的第一个非空字符是否是{[，如果是就需要先跳过所有{}和[]的匹配对，再找后面的逗号
	// 注意不是所有{ } [ ] 字符都算，本身key 或者 value 是有可能包含这些特殊字符的。
	off := dd.nextComma()
	if off < 0 {
		nh := dd.tail + 1
		span := dd.src[dd.head:nh]
		err := dd.parseKV(span)
		dd.head = nh
		return err
	} else if off == 0 {
		return sErr
	} else if off > 0 {
		nh := dd.head + off + 1
		span := dd.src[dd.head : nh-1]
		if err := dd.parseKV(span); err != nil {
			return err
		}
		dd.head = nh
		goto loopComma
	}
	return nil
}

func (dd *fastDecode) parseArray() error {
	return nil

}

//func (dd *fastDecode) parseLiteral() error {
//	return nil
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 下一个逗号
func (dd *fastDecode) nextComma() int {
	off := dd.head
	for off < dd.tail {
		c := dd.src[off]
		if c == ',' {
			return off - dd.head
		}
		off++
	}
	return -1
}

// 可能的情况 "k":v | "k":{} | "k":[]
func (dd *fastDecode) parseKV(kvPair string) error {
	// 冒号不能是key中的冒号
	colonIdx := colonInKVString(kvPair)
	// 可能就是一段空白字符
	if colonIdx == -1 {
		kvPair = trim(kvPair)
		if len(kvPair) != 0 {
			return sErr
		}
	}
	// 长度小于2的key是不可能的，等于2时可能是空字符串: ""
	if colonIdx < 2 {
		return sErr
	}
	key := kvPair[:colonIdx]
	val := kvPair[colonIdx+1:]

	if dd.gr != nil {
		return dd.setGsonValue(key, val)
	}
	return dd.setMapValue(key, val)
}

// 核查 key +++++++++
func (dd *fastDecode) checkKey(key string, needCopy bool) (ret string, err error) {
	var inShare bool
	if key, err = cutKeyQuote(key); err != nil {
		return
	}
	if key, inShare, err = dd.getStringLiteral(key); err != nil {
		return
	}
	if inShare && needCopy {
		key = copyString(key)
	}
	return key, nil
}

// 核查 value ++++++++++
func (dd *fastDecode) checkValue(val string) (ret string, err error) {
	var hasQuote, inShare bool
	if val, hasQuote, err = cutValueQuote(val); err != nil {
		return
	}
	if hasQuote {
		if val, inShare, err = dd.getStringLiteral(val); err != nil {
			return
		}
		// 证明此时用的是临时栈空间保存转义之后的Value，需要申请新的内存空间放置
		if inShare {
			val = copyString(val)
		}
	}
	return val, nil
}

// 当目标为 cst.KV 类型时候，用此方法设置
func (dd *fastDecode) setMapValue(key, val string) (err error) {
	if key, err = dd.checkKey(key, true); err != nil {
		return
	}
	if val, err = dd.checkValue(val); err != nil {
		return
	}

	// set k = v
	dd.dst.Set(key, val)
	return nil
}

// 当目标为 gson.GsonRow 类型时候，用此方法设置
func (dd *fastDecode) setGsonValue(key, val string) (err error) {
	if key, err = dd.checkKey(key, false); err != nil {
		return
	}
	keyIdx := dd.gr.KeyIndex(key)
	// 没有这个字段，直接返回了(此时再去解析后面的value是没有意义的)
	if keyIdx < 0 {
		return nil
	}

	if val, err = dd.checkValue(val); err != nil {
		return
	}

	// set k = v
	dd.gr.SetStringByIndex(keyIdx, val)
	return nil
}

func copyString(src string) string {
	tmp := make([]byte, len(src))
	copy(tmp, src)
	return lang.BTS(tmp)
}

func (dd *fastDecode) skipHeadSpace() bool {
	for dd.head < dd.tail {
		c := dd.src[dd.head]
		if !isSpace(c) {
			return true
		}
		dd.head++
	}
	return false
}

func (dd *fastDecode) skipTailSpace() bool {
	for dd.head < dd.tail {
		c := dd.src[dd.tail]
		if !isSpace(c) {
			return true
		}
		dd.tail--
	}
	return false
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\n' || c == '\r' || c == '\t')
}

// 在"k":v形式的字符串中找到冒号
func colonInKVString(str string) int {
	quoteCt := 0
	for i := range str {
		if str[i] == '"' {
			// 非第一个"，如果后面的"前面有\，是允许的
			if quoteCt == 1 {
				if str[i-1] == '\\' {
					continue
				}
			}
			quoteCt++
			if quoteCt > 2 {
				return -1
			}
		}
		if quoteCt == 2 {
			if str[i] == ':' {
				return i
			}
		}
	}
	return -1
}

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
// 一个合法的 Key，或者Value 字符串
func (dd *fastDecode) getStringLiteral(str string) (ret string, inShare bool, err error) {
	var newStr []byte
	var step int
	var hasSlash bool

	for i := 0; i < len(str); i++ {
		c := str[i]
		// 不支持非可见字符
		if c < 32 {
			return "", false, sErr
		}
		if c == '\\' {
			// 第一次检索到有 \
			if hasSlash == false {
				hasSlash = true
				// TODO：这里发生了逃逸，需要用sync.Pool的方式，共享内存空间
				// 或者别的黑魔法操作内存
				// add by sdx 20230404 动态初始化 share 内存
				if dd.share == nil {
					defSize := len(str)
					if defSize > tempByteStackSize {
						defSize = tempByteStackSize
					}
					dd.share = make([]byte, defSize)
				}

				if len(str) <= len(dd.share) {
					newStr = dd.share[:]
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
				return "", false, sErr
			}
			continue
		}
		if hasSlash {
			newStr[step] = c
			step++
		}
		//	// ASCII
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
		return lang.BTS(newStr[:step]), inShare, nil
	}
	return str, false, nil
}

// 检查是否为一个合法的 key
func cutKeyQuote(key string) (ret string, err error) {
	key = trim(key)

	if len(key) < 2 {
		return "", sErr
	}

	if key[0] != '"' || key[len(key)-1] != '"' {
		return "", sErr
	}
	return key[1 : len(key)-1], nil
}

// 检查是否是一个合法的 value 值
func cutValueQuote(val string) (ret string, hasQuote bool, err error) {
	val = trim(val)

	if len(val) < 1 {
		return "", false, sErr
	}

	// 如果 value 没有 双引号，可能是数值、true、false、null四种情况
	if val[0] != '"' {
		if err = checkNoQuoteValue(val); err != nil {
			return "", false, err
		}
		return val, false, nil
	}

	// 有双引号
	if len(val) == 1 {
		return "", false, sErr
	} else if val[len(val)-1] != '"' {
		return "", false, sErr
	}
	return val[1 : len(val)-1], true, nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 没有双引号围起来的值，只可能是：数值、true、false、null
func checkNoQuoteValue(str string) error {
	if len(str) == 0 {
		return sErr
	}

	// 只有这几种首字符的可能
isNumber:
	c := str[0]
	if c >= '0' && c <= '9' {
		// 0.234 | 234.23 | 23424 | 3.8e+07
		dot := false
		if c == '0' {
			if len(str) < 3 || str[1] != '.' {
				return sErr
			}
		}
		for i := 1; i < len(str); i++ {
			c = str[i]
			if c == '.' {
				if dot == true {
					return sErr
				} else {
					dot = true
					continue
				}
			} else if c == 'e' || c == 'E' {
				return checkScientificNumberTail(str[i+1:])
			} else if c < '0' || c > '9' {
				return sErr
			}
		}
	} else if c == '-' {
		// -0.3 | +13.33 | -3.7E-7
		if len(str) >= 2 && str[1] >= '0' && str[1] <= '9' {
			str = str[1:]
			goto isNumber
		} else {
			return sErr
		}
	} else if c == 'f' {
		// false
		if str != "false" {
			return sErr
		}
	} else if c == 't' {
		// true
		if str != "true" {
			return sErr
		}
	} else if c == 'n' {
		if str != "null" {
			return sErr
		}
	} else {
		return sErr
	}
	return nil
}

// 检查科学计数法（e|E）后面的字符串合法性
func checkScientificNumberTail(str string) error {
	if len(str) == 0 {
		return sErr
	}
	c := str[0]
	if c == '-' || c == '+' {
		str = str[1:]
	}

	if len(str) == 0 {
		return sErr
	}
	for i := range str {
		if str[i] < '0' || str[i] > '9' {
			return sErr
		}
	}
	return nil
}
