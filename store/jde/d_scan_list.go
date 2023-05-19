package jde

func arrItemPtr(sd *subDecode) uintptr {
	return sd.dstPtr + uintptr(sd.arrIdx*sd.dm.itemSize)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanList() {
	if !sd.dm.isList {
		panic(errList)
	}

	// A. 可能需要用到缓冲池记录临时数据
	sd.resetListPool()

	// B. 根据目标值类型，直接匹配，提高性能
	sd.scanListItems()

	// C. 将解析好的数据一次性绑定到对象上
	sd.flushListPool()
}

func (sd *subDecode) scanListItems() {
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
		if sd.skipValue {
			sd.skipOneValue()
		} else {
			sd.dm.itemDec(sd)
			if sd.dm.isArray {
				sd.arrIdx++
				if sd.arrIdx >= sd.dm.arrLen {
					sd.skipValue = true
				}
			}
		}
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

//func (sd *subDecode) scanOneValue() {
//	switch c := sd.str[sd.scan]; {
//	case c == '{':
//		sd.scan++
//		sd.scanSubObject()
//	case c == '[':
//		sd.scan++
//		//err = sd.scanSubArray()
//	case c == '"':
//		sd.scanQuoteStrValue()
//	case c >= '0' && c <= '9', c == '-':
//		//sd.scanNumValue()
//	case c == 't':
//		sd.skipTrue()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBool(true)
//	case c == 'f':
//		sd.skipFalse()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBool(false)
//	default:
//		sd.skipNull()
//		if sd.skipValue {
//			return
//		}
//		sd.bindBoolNull()
//	}
//}
