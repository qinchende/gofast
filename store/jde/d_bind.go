package jde

import (
	"reflect"
	"strconv"
)

// skip some items
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//go:inline
func (sd *subDecode) checkSkip() {
	// 如果是 struct ，就找找是否支持这个字段
	if sd.dm.isStruct {
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

//go:inline
func (sd *subDecode) isSkip() bool {
	return sd.skipValue || sd.skipTotal
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//go:inline
func (sd *subDecode) bindString(val string) (err int) {
	// 如果是 struct
	if sd.dm.isStruct {
		sd.dm.ss.BindString(sd.dstPtr, sd.keyIdx, val)
		//sd.obj.setStringByIndex(sd.keyIdx, val)
		return noErr
	}

	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetStringByIndex(sd.keyIdx, val)
			return noErr
		}
		sd.mp.Set(sd.key, val)
		return noErr
	}

	// 如果是数组
	if sd.dm.isList {
		if !allowStr(sd.dm.itemKind) {
			return errList
		}
		if sd.dm.isArray {
			if len(sd.pl.bufStr) >= sd.arrLen {
				sd.skipValue = true
				return noErr
			}
			//if !sd.arr.isPtr {
			//	bindArrValue[string](sd.arr, val)
			//	return noErr
			//}
		}
		sd.pl.bufStr = append(sd.pl.bufStr, val)
		return noErr
	}

	return noErr
}

//go:inline
func (sd *subDecode) bindBool(val bool) (err int) {
	// 如果是 struct
	if sd.dm.isStruct {
		sd.dm.ss.BindBool(sd.dstPtr, sd.keyIdx, val)
		//sd.dm.setBoolByIndex(sd.keyIdx, val)
		return noErr
	}

	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetByIndex(sd.keyIdx, val)
			return noErr
		}
		sd.mp.Set(sd.key, val)
		return noErr
	}

	// 如果是数组
	if sd.dm.isList {
		if !allowBool(sd.dm.itemKind) {
			return errList
		}
		if sd.dm.isArray {
			if len(sd.pl.bufStr) >= sd.arrLen {
				sd.skipValue = true
				return noErr
			}
		}
		sd.pl.bufBol = append(sd.pl.bufBol, val)
		return noErr
	}

	return noErr
}

//go:inline
func (sd *subDecode) bindNumber(val string, hasDot bool) (err int) {
	// 如果是 struct
	if sd.dm.isStruct {
		if num, err1 := parseInt(val); err < 0 {
			return err1
		} else {
			sd.dm.ss.BindInt(sd.dstPtr, sd.keyIdx, num)
			//sd.dm.setIntByIndex(sd.keyIdx, num)
		}
		return noErr
	}

	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetStringByIndex(sd.keyIdx, val)
			return noErr
		}
		sd.mp.Set(sd.key, val)
		return noErr
	}

	// 如果是数组
	if sd.dm.isList {
		// 如果目标是 any 值
		if sd.dm.itemKind == reflect.Interface {
			if sd.dm.isArray {
				if len(sd.pl.bufStr) >= sd.arrLen {
					sd.skipValue = true
					return noErr
				}
			}
			sd.pl.bufStr = append(sd.pl.bufStr, val)
			return noErr
		}

		if allowInt(sd.dm.itemKind) {
			//if sd.isArray && !sd.arr.isPtr {
			//	if sd.arr.arrIdx >= sd.arr.arrLen {
			//		sd.skipValue = true
			//		return noErr
			//	}
			//
			//	if num, err1 := parseInt(val); err < 0 {
			//		return err1
			//	} else {
			//		sd.arr.arrIntFunc(sd.arr, num)
			//	}
			//	return noErr
			//}

			if len(sd.pl.bufI64) >= sd.arrLen {
				sd.skipValue = true
				return noErr
			}

			if num, err1 := parseInt(val); err < 0 {
				return err1
			} else {
				sd.pl.bufI64 = append(sd.pl.bufI64, num)
			}

		} else if allowFloat(sd.dm.itemKind) {
			if sd.dm.isArray && len(sd.pl.bufF64) >= sd.arrLen {
				sd.skipValue = true
				return noErr
			}
			if num, err1 := strconv.ParseFloat(val, 64); err1 != nil {
				return errNumberFmt
			} else {
				sd.pl.bufF64 = append(sd.pl.bufF64, num)
			}

		} else {
			return errList
		}
		return noErr
	}

	// 如果是 struct

	return noErr
}

func (sd *subDecode) bindNull() (err int) {
	// 如果是数组
	if sd.dm.isList {
		return noErr
	}

	// 如果是 struct

	return noErr
}
