package jde

func listItemPtr(sd *subDecode) uintptr {
	return sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (sd *subDecode) scanSubArray(key string) (val string, err int) {
//	sub := subDecode{
//		str:       sd.str,
//		scan:      sd.scan,
//		skipTotal: sd.skipValue,
//	}
//
//	if sd.gr != nil {
//		// TODO: 无法为子对象提供目标值，只能返回字符串
//		sub.skipTotal = true
//	} else {
//		sd.skipValue = true
//	}
//
//	err = sub.scanList()
//	if err < 0 {
//		sd.scan = sub.scan
//		return
//	}
//
//	if sd.gr != nil {
//		if sd.skipValue == false {
//			val = sd.str[sd.scan-1 : sub.scan]
//		}
//	} else {
//		//sd.mp.Set(key, sub.list)
//	}
//	sd.scan = sub.scan
//	return
//}

// 前提：sd.str 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanList() {
	if !sd.isList {
		panic(errList)
	}

	// A. 可能需要用到缓冲池记录临时数据
	sd.resetListPool()

	// B. 根据目标值类型，直接匹配，提高性能
	//switch sd.dm.itemKind {
	//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//	//sd.scanArrItems(sd.scanIntValue)
	//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//	//sd.scanArrItems(sd.scanUintValue)
	//case reflect.Float32, reflect.Float64:
	//	//sd.scanArrItems(sd.scanNumValue)
	//case reflect.String:
	//	sd.scanArrItems(sd.scanStrKindValue)
	//case reflect.Slice, reflect.Array:
	//	sd.scanSubObject()
	//case reflect.Struct:
	//	sd.scanSubObject()
	//case reflect.Interface:
	//	sd.scanArrItems(sd.scanOneValue)
	//default:
	//	panic(errSupport)
	//}
	sd.scanArrItems()

	// C. 将解析好的数据一次性绑定到对象上
	sd.flushListPool()
}

func (sd *subDecode) scanArrItems() {
	pos := sd.scan

	for isBlankChar[sd.str[pos]] {
		pos++
	}
	c := sd.str[pos]
	if c == ',' {
		goto errChar
	}

	for {
		// 不用switch, 比较顺序相对比较明确
		if c == ',' {
			pos++
		} else if c == ']' {
			// 数组多余的部分需要重置成类型零值
			if sd.arrIdx < sd.dm.arrLen {
				sd.resetArrLeftItems()
			}
			sd.scan = pos + 1
			return
		} else if sd.arrIdx > 0 {
			goto errChar
		}

		for isBlankChar[sd.str[pos]] {
			pos++
		}

		sd.scan = pos
		//scanValue()
		sd.dm.itemFunc(sd)
		sd.arrIdx++
		pos = sd.scan

		for isBlankChar[sd.str[pos]] {
			pos++
		}
		c = sd.str[pos]
	}

errChar:
	sd.scan = pos
	panic(errChar)
}

func (sd *subDecode) scanOneValue() {
	switch c := sd.str[sd.scan]; {
	case c == '{':
		sd.scan++
		sd.scanSubObject()
	case c == '[':
		sd.scan++
		//err = sd.scanSubArray()
	case c == '"':
		sd.scanQuoteStrValue()
	case c >= '0' && c <= '9', c == '-':
		//sd.scanNumValue()
	case c == 't':
		sd.skipTrue()
		if sd.skipValue {
			return
		}
		sd.bindBool(true)
	case c == 'f':
		sd.skipFalse()
		if sd.skipValue {
			return
		}
		sd.bindBool(false)
	default:
		sd.skipNull()
		if sd.skipValue {
			return
		}
		sd.bindBoolNull()
	}
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 匹配一个数值，对应于float类型
//// 0.234 | 234.23 | 23424 | 3.8e+07 | 3.7E-7 | -0.3 | -3.7E-7
//func (sd *subDecode) scanNumValue() {
//	pos := sd.scan
//	start := pos
//	var hasDot, needNum bool
//
//	c := sd.str[pos]
//	if c == '-' {
//		pos++
//		c = sd.str[pos]
//	}
//	// 0开头的数字，只能是：0 | 0.x | 0e | 0E
//	if c == '0' {
//		pos++
//		c = sd.str[pos]
//
//		switch c {
//		case '.', 'e', 'E':
//			goto loopNum
//		default:
//			goto over
//		}
//	}
//	needNum = true
//
//loopNum:
//	for {
//		c = sd.str[pos]
//		pos++
//
//		if c == '.' {
//			if hasDot == true {
//				panic(errNumberFmt)
//			}
//			hasDot = true
//			needNum = true
//		} else if c == 'e' || c == 'E' {
//			if needNum {
//				panic(errNumberFmt)
//			}
//			needNum = true
//
//			c := sd.str[pos]
//			if c == '-' || c == '+' {
//				pos++
//			}
//			for {
//				if c = sd.str[pos]; c < '0' || c > '9' {
//					break loopNum
//				} else {
//					needNum = false
//				}
//				pos++
//			}
//		} else if c < '0' || c > '9' {
//			pos--
//			break
//		} else {
//			needNum = false // 到这里，字符肯定是数字
//		}
//	}
//
//	if needNum {
//		panic(errNumberFmt)
//	}
//
//over:
//	sd.scan = pos
//	// 还剩下最后一种可能：null
//	if start == pos {
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindNumberNull()
//		return
//	}
//	if sd.skipValue {
//		return
//	}
//	sd.bindNumber(sd.str[start:pos])
//}

//func (sd *subDecode) scanIntValue() {
//	pos := sd.scan
//	start := pos
//
//	c := sd.str[pos]
//	if c == '-' {
//		pos++
//		c = sd.str[pos]
//	}
//	if c == '0' {
//		pos++
//		goto over
//	}
//	for {
//		if c < '0' || c > '9' {
//			break
//		}
//		pos++
//		c = sd.str[pos]
//	}
//over:
//	sd.scan = pos
//	// 还剩下最后一种可能：null +++
//	if start == pos {
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindIntNull()
//		return
//	}
//	if sd.skipValue {
//		return
//	}
//	sd.bindIntList(sd.str[start:pos])
//}

//func (sd *subDecode) scanUintValue() {
//	pos := sd.scan
//	start := pos
//
//	c := sd.str[pos]
//	if c == '0' {
//		pos++
//		goto over
//	}
//	for {
//		if c < '0' || c > '9' {
//			break
//		}
//		pos++
//		c = sd.str[pos]
//	}
//over:
//	sd.scan = pos
//	// 还剩下最后一种可能：null
//	if start == pos {
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindUintNull()
//		return
//	}
//	if sd.skipValue {
//		return
//	}
//	sd.bindUintList(sd.str[start:pos])
//}
