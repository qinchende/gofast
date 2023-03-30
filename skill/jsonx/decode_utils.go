package jsonx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
)

const (
	// Continue.
	scanContinue     = iota // uninteresting byte
	scanBeginLiteral        // end implied by next result != scanContinue
	scanBeginObject         // begin object
	scanObjectKey           // just finished object key (string)
	scanObjectValue         // just finished non-last object value
	scanEndObject           // end object (implies scanObjectValue if possible)
	scanBeginArray          // begin array
	scanArrayValue          // just finished array value
	scanEndArray            // end array (implies scanArrayValue if possible)
	scanSkipSpace           // space byte; can skip; known to be last "continue" result

	// Stop.
	scanEnd   // top-level value ended *before* this byte; known to be first "stop" result
	scanError // hit an error, scanner.err.
)

var (
	syntaxError = errors.New("jsonx: json syntax error.")
)

type gsonDecode struct {
	dest   cst.SuperKV
	origin string
	src    string
	head   int
	tail   int
	//hSeek   int
	//tSeek   int
	braces  bracesMark
	squares squaresMark
}

func (dd *gsonDecode) init(dst cst.SuperKV, src string) error {
	//dd.searchBrackets()
	//// 左右大括号数量不一致，格式错误(object)
	//if len(dd.braces.left) != len(dd.braces.right) {
	//	return syntaxError
	//}
	//// 左右中括号数量不一致，格式错误(array)
	//if len(dd.squares.left) != len(dd.squares.right) {
	//	return syntaxError
	//}

	dd.dest = dst
	dd.origin = src
	dd.src = src
	dd.head = 0
	dd.tail = len(dd.src) - 1
	return nil
}

// 采用尽最大努力解析出正确结果的策略
func (dd *gsonDecode) parse() error {
	if ok := dd.nextHead(); !ok {
		return syntaxError
	}

	c := dd.src[dd.head]
	switch c {
	case '{':
		if ok := dd.nextTail(); !ok {
			return syntaxError
		}
		if dd.src[dd.tail] != '}' {
			return syntaxError
		}
		return dd.parseObject()
	case '[':
		if exist := dd.nextTail(); !exist {
			return syntaxError
		}
		if dd.src[dd.tail] != ']' {
			return syntaxError
		}
		return dd.parseArray()
	default:
		return syntaxError
	}
}

func (dd *gsonDecode) parseObject() error {
	dd.head++
	dd.tail--

loopComma:
	// 剩下的全部是 k:v,k:v
	off := dd.nextComma()
	if off < 0 {
		nh := dd.tail + 1
		span := dd.src[dd.head:nh]
		dd.head = nh
		return dd.parseKV(span)
	} else if off == 0 {
		return syntaxError
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

func (dd *gsonDecode) parseArray() error {
	return nil

}

func (dd *gsonDecode) parseLiteral() error {
	return nil
}

func (dd *gsonDecode) parseKV(kvItem string) error {
	colon := nextColon(kvItem)
	if colon == -1 {
		kvItem = trim(kvItem)
		if len(kvItem) != 0 {
			return syntaxError
		}
	}
	if colon < 3 {
		return syntaxError
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
	dd.dest.Set(key, val)
	return nil
}

func (dd *gsonDecode) nextHead() bool {
	for dd.head < dd.tail {
		c := dd.src[dd.head]
		if !isSpace(c) {
			return true
		}
		dd.head++
	}
	return false
}

func (dd *gsonDecode) nextTail() bool {
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
func (dd *gsonDecode) nextComma() int {
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

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
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

func checkKey(key string) (string, error) {
	key = trim(key)

	if len(key) < 3 {
		return "", syntaxError
	}
	if key[0] != '"' || key[len(key)-1] != '"' {
		return "", syntaxError
	}
	return key[1 : len(key)-1], nil
}

func checkValue(val string) (string, error) {
	val = trim(val)
	if len(val) < 1 {
		return "", syntaxError
	}

	if val[0] == '"' {
		if len(val) == 1 {
			return "", syntaxError
		} else if val[len(val)-1] != '"' {
			return "", syntaxError
		}
		return val[1 : len(val)-1], nil
	}
	return val, nil
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (d *decodeState) value(v reflect.Value) error {
//	switch d.opcode {
//	default:
//		panic(phasePanicMsg)
//
//	case scanBeginArray:
//		if v.IsValid() {
//			if err := d.array(v); err != nil {
//				return err
//			}
//		} else {
//			d.skip()
//		}
//		d.scanNext()
//
//	case scanBeginObject:
//		if v.IsValid() {
//			if err := d.object(v); err != nil {
//				return err
//			}
//		} else {
//			d.skip()
//		}
//		d.scanNext()
//
//	case scanBeginLiteral:
//		// All bytes inside literal return scanContinue op code.
//		start := d.readIndex()
//		d.rescanLiteral()
//
//		if v.IsValid() {
//			if err := d.literalStore(d.data[start:d.readIndex()], v, false); err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//// rescanLiteral is similar to scanWhile(scanContinue), but it specialises the
//// common case where we're decoding a literal. The decoder scans the input
//// twice, once for syntax errors and to check the length of the value, and the
//// second to perform the decoding.
////
//// Only in the second step do we use decodeState to tokenize literals, so we
//// know there aren't any syntax errors. We can take advantage of that knowledge,
//// and scan a literal's bytes much more quickly.
//func (d *decodeState) rescanLiteral() {
//	data, i := d.data, d.off
//Switch:
//	switch data[i-1] {
//	case '"': // string
//		for ; i < len(data); i++ {
//			switch data[i] {
//			case '\\':
//				i++ // escaped char
//			case '"':
//				i++ // tokenize the closing quote too
//				break Switch
//			}
//		}
//	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-': // number
//		for ; i < len(data); i++ {
//			switch data[i] {
//			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
//				'.', 'e', 'E', '+', '-':
//			default:
//				break Switch
//			}
//		}
//	case 't': // true
//		i += len("rue")
//	case 'f': // false
//		i += len("alse")
//	case 'n': // null
//		i += len("ull")
//	}
//	if i < len(data) {
//		d.opcode = stateEndValue(&d.scan, data[i])
//	} else {
//		d.opcode = scanEnd
//	}
//	d.off = i + 1
//}
//
//// stateEndValue is the state after completing a value,
//// such as after reading `{}` or `true` or `["x"`.
//func stateEndValue(s *scanner, c byte) int {
//	n := len(s.parseState)
//	if n == 0 {
//		// Completed top-level before the current byte.
//		s.step = stateEndTop
//		s.endTop = true
//		return stateEndTop(s, c)
//	}
//	if isSpace(c) {
//		s.step = stateEndValue
//		return scanSkipSpace
//	}
//	ps := s.parseState[n-1]
//	switch ps {
//	case parseObjectKey:
//		if c == ':' {
//			s.parseState[n-1] = parseObjectValue
//			s.step = stateBeginValue
//			return scanObjectKey
//		}
//		return s.error(c, "after object key")
//	case parseObjectValue:
//		if c == ',' {
//			s.parseState[n-1] = parseObjectKey
//			s.step = stateBeginString
//			return scanObjectValue
//		}
//		if c == '}' {
//			s.popParseState()
//			return scanEndObject
//		}
//		return s.error(c, "after object key:value pair")
//	case parseArrayValue:
//		if c == ',' {
//			s.step = stateBeginValue
//			return scanArrayValue
//		}
//		if c == ']' {
//			s.popParseState()
//			return scanEndArray
//		}
//		return s.error(c, "after array element")
//	}
//	return s.error(c, "")
//}
//
//// stateEndTop is the state after finishing the top-level value,
//// such as after reading `{}` or `[1,2,3]`.
//// Only space characters should be seen now.
//func stateEndTop(s *scanner, c byte) int {
//	if !isSpace(c) {
//		// Complain about non-space byte on next call.
//		s.error(c, "after top-level value")
//	}
//	return scanEnd
//}
//
//// stateBeginValue is the state at the beginning of the input.
//func stateBeginValue(s *scanner, c byte) int {
//	if isSpace(c) {
//		return scanSkipSpace
//	}
//	switch c {
//	case '{':
//		s.step = stateBeginStringOrEmpty
//		return s.pushParseState(c, parseObjectKey, scanBeginObject)
//	case '[':
//		s.step = stateBeginValueOrEmpty
//		return s.pushParseState(c, parseArrayValue, scanBeginArray)
//	case '"':
//		s.step = stateInString
//		return scanBeginLiteral
//	case '-':
//		s.step = stateNeg
//		return scanBeginLiteral
//	case '0': // beginning of 0.123
//		s.step = state0
//		return scanBeginLiteral
//	case 't': // beginning of true
//		s.step = stateT
//		return scanBeginLiteral
//	case 'f': // beginning of false
//		s.step = stateF
//		return scanBeginLiteral
//	case 'n': // beginning of null
//		s.step = stateN
//		return scanBeginLiteral
//	}
//	if '1' <= c && c <= '9' { // beginning of 1234.5
//		s.step = state1
//		return scanBeginLiteral
//	}
//	return s.error(c, "looking for beginning of value")
//}
//
//// stateBeginStringOrEmpty is the state after reading `{`.
//func stateBeginStringOrEmpty(s *scanner, c byte) int {
//	if isSpace(c) {
//		return scanSkipSpace
//	}
//	if c == '}' {
//		n := len(s.parseState)
//		s.parseState[n-1] = parseObjectValue
//		return stateEndValue(s, c)
//	}
//	return stateBeginString(s, c)
//}
//
//// stateBeginString is the state after reading `{"key": value,`.
//func stateBeginString(s *scanner, c byte) int {
//	if isSpace(c) {
//		return scanSkipSpace
//	}
//	if c == '"' {
//		s.step = stateInString
//		return scanBeginLiteral
//	}
//	return s.error(c, "looking for beginning of object key string")
//}
//
//// stateInString is the state after reading `"`.
//func stateInString(s *scanner, c byte) int {
//	if c == '"' {
//		s.step = stateEndValue
//		return scanContinue
//	}
//	if c == '\\' {
//		s.step = stateInStringEsc
//		return scanContinue
//	}
//	if c < 0x20 {
//		return s.error(c, "in string literal")
//	}
//	return scanContinue
//}
//
//// stateInStringEsc is the state after reading `"\` during a quoted string.
//func stateInStringEsc(s *scanner, c byte) int {
//	switch c {
//	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
//		s.step = stateInString
//		return scanContinue
//	case 'u':
//		s.step = stateInStringEscU
//		return scanContinue
//	}
//	return s.error(c, "in string escape code")
//}
//
//// stateInStringEscU is the state after reading `"\u` during a quoted string.
//func stateInStringEscU(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
//		s.step = stateInStringEscU1
//		return scanContinue
//	}
//	// numbers
//	return s.error(c, "in \\u hexadecimal character escape")
//}
//
//// stateInStringEscU1 is the state after reading `"\u1` during a quoted string.
//func stateInStringEscU1(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
//		s.step = stateInStringEscU12
//		return scanContinue
//	}
//	// numbers
//	return s.error(c, "in \\u hexadecimal character escape")
//}
//
//// stateInStringEscU12 is the state after reading `"\u12` during a quoted string.
//func stateInStringEscU12(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
//		s.step = stateInStringEscU123
//		return scanContinue
//	}
//	// numbers
//	return s.error(c, "in \\u hexadecimal character escape")
//}
//
//// stateInStringEscU123 is the state after reading `"\u123` during a quoted string.
//func stateInStringEscU123(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
//		s.step = stateInString
//		return scanContinue
//	}
//	// numbers
//	return s.error(c, "in \\u hexadecimal character escape")
//}
//
//// stateNeg is the state after reading `-` during a number.
//func stateNeg(s *scanner, c byte) int {
//	if c == '0' {
//		s.step = state0
//		return scanContinue
//	}
//	if '1' <= c && c <= '9' {
//		s.step = state1
//		return scanContinue
//	}
//	return s.error(c, "in numeric literal")
//}
//
//// state1 is the state after reading a non-zero integer during a number,
//// such as after reading `1` or `100` but not `0`.
//func state1(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' {
//		s.step = state1
//		return scanContinue
//	}
//	return state0(s, c)
//}
//
//// state0 is the state after reading `0` during a number.
//func state0(s *scanner, c byte) int {
//	if c == '.' {
//		s.step = stateDot
//		return scanContinue
//	}
//	if c == 'e' || c == 'E' {
//		s.step = stateE
//		return scanContinue
//	}
//	return stateEndValue(s, c)
//}
//
//// stateDot is the state after reading the integer and decimal point in a number,
//// such as after reading `1.`.
//func stateDot(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' {
//		s.step = stateDot0
//		return scanContinue
//	}
//	return s.error(c, "after decimal point in numeric literal")
//}
//
//// stateDot0 is the state after reading the integer, decimal point, and subsequent
//// digits of a number, such as after reading `3.14`.
//func stateDot0(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' {
//		return scanContinue
//	}
//	if c == 'e' || c == 'E' {
//		s.step = stateE
//		return scanContinue
//	}
//	return stateEndValue(s, c)
//}
//
//// stateE is the state after reading the mantissa and e in a number,
//// such as after reading `314e` or `0.314e`.
//func stateE(s *scanner, c byte) int {
//	if c == '+' || c == '-' {
//		s.step = stateESign
//		return scanContinue
//	}
//	return stateESign(s, c)
//}
//
//// stateESign is the state after reading the mantissa, e, and sign in a number,
//// such as after reading `314e-` or `0.314e+`.
//func stateESign(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' {
//		s.step = stateE0
//		return scanContinue
//	}
//	return s.error(c, "in exponent of numeric literal")
//}
//
//// stateE0 is the state after reading the mantissa, e, optional sign,
//// and at least one digit of the exponent in a number,
//// such as after reading `314e-2` or `0.314e+1` or `3.14e0`.
//func stateE0(s *scanner, c byte) int {
//	if '0' <= c && c <= '9' {
//		return scanContinue
//	}
//	return stateEndValue(s, c)
//}
//
//// stateT is the state after reading `t`.
//func stateT(s *scanner, c byte) int {
//	if c == 'r' {
//		s.step = stateTr
//		return scanContinue
//	}
//	return s.error(c, "in literal true (expecting 'r')")
//}
//
//// stateTr is the state after reading `tr`.
//func stateTr(s *scanner, c byte) int {
//	if c == 'u' {
//		s.step = stateTru
//		return scanContinue
//	}
//	return s.error(c, "in literal true (expecting 'u')")
//}
//
//// stateTru is the state after reading `tru`.
//func stateTru(s *scanner, c byte) int {
//	if c == 'e' {
//		s.step = stateEndValue
//		return scanContinue
//	}
//	return s.error(c, "in literal true (expecting 'e')")
//}
//
//// stateF is the state after reading `f`.
//func stateF(s *scanner, c byte) int {
//	if c == 'a' {
//		s.step = stateFa
//		return scanContinue
//	}
//	return s.error(c, "in literal false (expecting 'a')")
//}
//
//// stateFa is the state after reading `fa`.
//func stateFa(s *scanner, c byte) int {
//	if c == 'l' {
//		s.step = stateFal
//		return scanContinue
//	}
//	return s.error(c, "in literal false (expecting 'l')")
//}
//
//// stateFal is the state after reading `fal`.
//func stateFal(s *scanner, c byte) int {
//	if c == 's' {
//		s.step = stateFals
//		return scanContinue
//	}
//	return s.error(c, "in literal false (expecting 's')")
//}
//
//// stateFals is the state after reading `fals`.
//func stateFals(s *scanner, c byte) int {
//	if c == 'e' {
//		s.step = stateEndValue
//		return scanContinue
//	}
//	return s.error(c, "in literal false (expecting 'e')")
//}
//
//// stateN is the state after reading `n`.
//func stateN(s *scanner, c byte) int {
//	if c == 'u' {
//		s.step = stateNu
//		return scanContinue
//	}
//	return s.error(c, "in literal null (expecting 'u')")
//}
//
//// stateNu is the state after reading `nu`.
//func stateNu(s *scanner, c byte) int {
//	if c == 'l' {
//		s.step = stateNul
//		return scanContinue
//	}
//	return s.error(c, "in literal null (expecting 'l')")
//}
//
//// stateNul is the state after reading `nul`.
//func stateNul(s *scanner, c byte) int {
//	if c == 'l' {
//		s.step = stateEndValue
//		return scanContinue
//	}
//	return s.error(c, "in literal null (expecting 'l')")
//}
//
//// stateError is the state after reaching a syntax error,
//// such as after reading `[1}` or `5.1.2`.
//func stateError(s *scanner, c byte) int {
//	return scanError
//}
//
//// error records an error and switches to the error state.
//func (s *scanner) error(c byte, context string) int {
//	s.step = stateError
//	s.err = &SyntaxError{"invalid character " + quoteChar(c) + " " + context, s.bytes}
//	return scanError
//}
//
//// quoteChar formats c as a quoted character literal
//func quoteChar(c byte) string {
//	// special cases - different from quoted strings
//	if c == '\'' {
//		return `'\''`
//	}
//	if c == '"' {
//		return `'"'`
//	}
//
//	// use quoted string with different quotation marks
//	s := strconv.Quote(string(c))
//	return "'" + s[1:len(s)-1] + "'"
//}
