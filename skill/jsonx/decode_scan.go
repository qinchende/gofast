package jsonx

// 采用尽最大努力解析出正确结果的策略
func (dd *fastDecode) parseJson() error {
	if ok := dd.nextHead(); !ok {
		return sErr
	}

	c := dd.src[dd.head]
	switch c {
	case '{':
		if ok := dd.nextTail(); !ok {
			return sErr
		}
		if dd.src[dd.tail] != '}' {
			return sErr
		}
		return dd.parseObject()
	case '[':
		if exist := dd.nextTail(); !exist {
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

	// 剩下的全部应该是 k:v,k:v
loopComma:
	off := dd.nextComma()
	if off < 0 {
		nh := dd.tail + 1
		span := dd.src[dd.head:nh]
		dd.head = nh
		return dd.parseKV(span)
	} else if off == 0 {
		return sErr
	} else if off > 0 {
		nh := dd.head + off + 1
		span := dd.src[dd.head : nh-1]
		dd.head = nh
		if err := dd.parseKV(span); err != nil {
			return err
		}
		goto loopComma
	}
	return nil
}

func (dd *fastDecode) parseArray() error {
	return nil

}

func (dd *fastDecode) parseLiteral() error {
	return nil
}

func (dd *fastDecode) parseKV(kvItem string) error {
	colon := nextColon(kvItem)
	if colon == -1 {
		kvItem = trim(kvItem)
		if len(kvItem) != 0 {
			return sErr
		}
	}
	if colon < 3 {
		return sErr
	}
	key := kvItem[:colon]
	val := kvItem[colon+1:]

	var err error
	if key, err = checkKey(key); err != nil {
		return err
	}

	if val, err = checkValue(val); err != nil {
		return err
	}

	// k = v
	dd.dst.Set(key, val)
	return nil
}

func (dd *fastDecode) nextHead() bool {
	for dd.head < dd.tail {
		c := dd.src[dd.head]
		if !isSpace(c) {
			return true
		}
		dd.head++
	}
	return false
}

func (dd *fastDecode) nextTail() bool {
	for dd.head < dd.tail {
		c := dd.src[dd.tail]
		if !isSpace(c) {
			return true
		}
		dd.tail--
	}
	return false
}

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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\n' || c == '\r' || c == '\t')
}

func nextColon(str string) int {
	for i := range str {
		if str[i] == ':' {
			return i
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

func checkString(str string) error {
	for i := range str {
		if str[i] < 32 || str[i] == '\\' {
			return sErr
		}
	}
	return nil
}

func checkKey(key string) (string, error) {
	key = trim(key)

	if len(key) < 3 {
		return "", sErr
	}
	if key[0] != '"' || key[len(key)-1] != '"' {
		return "", sErr
	}
	key = key[1 : len(key)-1]
	if err := checkString(key); err != nil {
		return "", err
	}
	return key, nil
}

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
			if str[i] == '.' {
				if dot == true {
					return sErr
				} else {
					dot = true
					continue
				}
			} else if str[i] == 'e' || str[i] == 'E' {
				return checkScientificNumberTail(str[i+1:])
			} else if str[i] < '0' || str[i] > '9' {
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

func checkValue(val string) (string, error) {
	val = trim(val)
	if len(val) < 1 {
		return "", sErr
	}

	if val[0] == '"' {
		if len(val) == 1 {
			return "", sErr
		} else if val[len(val)-1] != '"' {
			return "", sErr
		}
		val = val[1 : len(val)-1]

		if err := checkString(val); err != nil {
			return "", err
		}
	} else {
		if err := checkNoQuoteValue(val); err != nil {
			return "", err
		}
	}
	return val, nil
}
