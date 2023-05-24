package jde

import "unsafe"

func arrItemPtr(sd *subDecode) unsafe.Pointer {
	return unsafe.Pointer(sd.dstPtr + uintptr(sd.arrIdx*sd.dm.arrItemBytes))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 前提：sd.str 肯定是 [ 字符后面的字符串
// 返回 ] 后面字符的 index
func (sd *subDecode) scanList() {
	// A. 可能需要用到缓冲池记录临时数据
	sd.resetListPool()

	// B. 根据目标值类型，直接匹配，提高性能
	sd.scanListItems()

	// C. 将解析好的数据一次性绑定到对象上
	sd.flushListPool()
}

func (sd *subDecode) scanListItems() {
	pos := sd.scan

	pos++
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
			sd.dm.listItemDec(sd)
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
