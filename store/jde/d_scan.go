package jde

import (
	"fmt"
)

// 采用尽最大努力解析出正确结果的策略
// 可能解析过程中出现错误，所有最终需要通过判断返回的error来确定解析是否成功，发生错误时已经解析的结果不可信，请不要使用
func (sd *subDecode) scanJson() (err errType) {
	// 解析过程中异常，这里统一截获处理，返回解析错误编号
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				err = code
			} else {
				// 调试的时候打印错误信息
				fmt.Println(pic)
				err = errJson
			}
		}
	}()

	for isBlankChar[sd.str[sd.scan]] {
		sd.scan++
	}

	switch sd.str[sd.scan] {
	case '{':
		sd.scanJsonEnd('}')
		return
	case '[':
		sd.scanJsonEnd(']')
		return
	case 'n':
		sd.scanJsonEnd('l')
		return
	}
	return errJson
}

// 只支持 } ] l 三个字符判断
func (sd *subDecode) scanJsonEnd(ch byte) {
	// 去掉尾部的空白字符
	for i := len(sd.str) - 1; i > 0; i-- {
		if !isBlankChar[sd.str[i]] {
			if sd.str[i] != ch {
				sd.scan = i
				panic(errChar)
			}
			sd.str = sd.str[:i+1]
			break
		}
	}

	if ch == '}' {
		sd.scan++
		sd.scanObject()
	} else if ch == ']' {
		sd.scan++
		sd.scanList()
	} else {
		sd.skipNull()
	}
}

// skip items
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) skipOneValue() {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		sd.scan++
		sd.scanSubObject()
	case c == '[':
		sd.scan++
		//err = sd.scanSubArray()
	case c == '"':
		sd.skipQuoteStr()
	case c >= '0' && c <= '9', c == '-':
		sd.scanNumValue()
	case c == 't':
		sd.skipTrue()
	case c == 'f':
		sd.skipFalse()
	default:
		sd.skipNull()
	}
}

func (sd *subDecode) skipQuoteStr() {
	pos := sd.scan
	for {
		pos++
		switch c := sd.str[pos]; {
		case c == '"':
			sd.scan = pos + 1
			return
		case c == '\\':
			pos++ // 跳过 '\' 后面的一个字符
		}
	}
}

func (sd *subDecode) skipNull() {
	s := sd.scan
	if sd.str[s:s+4] == "null" {
		sd.scan += 4
		return
	}
	panic(errNull)
}

func (sd *subDecode) skipTrue() {
	s := sd.scan + 1
	if sd.str[s:s+3] == "rue" {
		sd.scan += 4
		return
	}
	panic(errBool)
}

func (sd *subDecode) skipFalse() {
	s := sd.scan + 1
	if sd.str[s:s+4] == "alse" {
		sd.scan += 5
		return
	}
}
