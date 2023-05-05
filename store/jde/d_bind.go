package jde

import (
	"strconv"
)

// skip some items
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) checkSkip() {
	// 如果是 struct ，就找找是否支持这个字段
	if sd.isStruct {
		if sd.keyIdx = sd.dm.ss.ColumnIndex(sd.key); sd.keyIdx < 0 {
			sd.skipValue = true
		} else {
			sd.skipValue = false
		}
		return
	}
	// PS: 可以先判断目标对象是否有这个key，没有就跳过value，解析下一个kv
	if sd.gr != nil {
		if sd.keyIdx = sd.gr.KeyIndex(sd.key); sd.keyIdx < 0 {
			sd.skipValue = true
		} else {
			sd.skipValue = false
		}
		return
	}
}

////go:inline
//func (sd *subDecode) isSkip() bool {
//	return sd.skipValue || sd.skipTotal
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) bindString(val string) {
	// 如果是 struct
	if sd.isStruct {
		sd.dm.ss.BindString(sd.dstPtr, sd.keyIdx, val)
		return
	}

	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetStringByIndex(sd.keyIdx, val)
			return
		}
		sd.mp.Set(sd.key, val)
		return
	}

	// 如果是数组
	if sd.isList {
		// 如果是定长的数组，而且值不是指针类型，可以直接设置值
		if sd.isArray && !sd.isPtr {
			if sd.arrIdx >= sd.dm.arrLen {
				sd.skipValue = true
				return
			}
			sd.bindStringArr(val)
			return
		}
		if sd.isAny {
			sd.pl.bufAny = append(sd.pl.bufAny, val)
		} else {
			sd.pl.bufStr = append(sd.pl.bufStr, val)
		}
		return
	}
}

func (sd *subDecode) bindBool(val bool) {
	// 如果是 struct
	if sd.isStruct {
		sd.dm.ss.BindBool(sd.dstPtr, sd.keyIdx, val)
		return
	}

	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetByIndex(sd.keyIdx, val)
			return
		}
		sd.mp.Set(sd.key, val)
		return
	}

	// 如果是数组
	if sd.isList {
		if sd.isArray && !sd.isPtr {
			if sd.arrIdx >= sd.dm.arrLen {
				sd.skipValue = true
				return
			}
			sd.bindBoolArr(val)
			return
		}
		if sd.isAny {
			sd.pl.bufAny = append(sd.pl.bufAny, val)
		} else {
			sd.pl.bufBol = append(sd.pl.bufBol, val)
		}
		return
	}

	return
}

func (sd *subDecode) bindNumber(val string) {
	// 如果是 struct
	if sd.isStruct {
		sd.dm.ss.BindInt(sd.dstPtr, sd.keyIdx, parseInt(val))
		return
	}

	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetStringByIndex(sd.keyIdx, val)
			return
		}
		sd.mp.Set(sd.key, val)
		return
	}

	// 如果是数组
	if sd.isList {
		// 只能是整形
		if allowInt(sd.dm.itemKind) {
			if sd.isArray && !sd.isPtr {
				if sd.arrIdx >= sd.dm.arrLen {
					sd.skipValue = true
					return
				}
				sd.bindIntArr(parseInt(val))
				return
			}

			if sd.isAny {
				sd.pl.bufAny = append(sd.pl.bufAny, parseInt(val))
			} else {
				sd.pl.bufI64 = append(sd.pl.bufI64, parseInt(val))
			}
			return
		}

		// 只能是整形
		if allowFloat(sd.dm.itemKind) {
			if sd.isArray && !sd.isPtr {
				if sd.arrIdx >= sd.dm.arrLen {
					sd.skipValue = true
					return
				}

				if f64, err1 := strconv.ParseFloat(val, 64); err1 != nil {
					panic(errNumberFmt)
				} else {
					sd.bindFloatArr(f64)
				}
				return
			}

			if num, err1 := strconv.ParseFloat(val, 64); err1 != nil {
				panic(errNumberFmt)
			} else {
				if sd.isAny {
					sd.pl.bufAny = append(sd.pl.bufAny, num)
				} else {
					sd.pl.bufF64 = append(sd.pl.bufF64, num)
				}
			}
			return
		}

		// 错误
		panic(errList)
	}
}

func (sd *subDecode) bindNull() {
	// 如果是数组
	if sd.isList {
		return
	}

	// 如果是 struct

	return
}
