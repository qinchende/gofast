package jde

import (
	"strconv"
)

//go:inline
func (sd *subDecode) setSkipFlag() {
	// PS: 可以先判断目标对象是否有这个key，没有就跳过value，解析下一个kv
	if sd.gr != nil {
		if sd.keyIdx = sd.gr.KeyIndex(sd.key); sd.keyIdx < 0 {
			sd.skipValue = true
		}
	}
	// 如果是 struct ，就找找是否支持这个字段
}

//go:inline
func (sd *subDecode) isSkip() bool {
	if sd.skipValue || sd.skipTotal {
		sd.skipValue = false
		return true
	}
	return false
}

//go:inline
func (sd *subDecode) bindString(val string) (err int) {
	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetStringByIndex(sd.keyIdx, val)
			return noErr
		}
		sd.mp.Set(sd.key, val)
		return noErr
	}

	// 如果是数组
	if sd.isList {
		if !allowStr(sd.arr.itemKind) {
			return errArray
		}
		if sd.isArray && len(sd.pl.arrStr) >= sd.arr.arrSize {
			return noErr
		}
		sd.pl.arrStr = append(sd.pl.arrStr, val)
		return noErr
	}

	// 如果是 struct

	return noErr
}

//go:inline
func (sd *subDecode) bindBool(val bool) (err int) {
	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetByIndex(sd.keyIdx, val)
			return noErr
		}
		sd.mp.Set(sd.key, val)
		return noErr
	}

	// 如果是数组
	if sd.isList {
		if !allowBool(sd.arr.itemKind) {
			return errArray
		}
		if sd.isArray && len(sd.pl.arrStr) >= sd.arr.arrSize {
			return noErr
		}
		sd.pl.arrBool = append(sd.pl.arrBool, val)
		return noErr
	}

	// 如果是 struct

	return noErr
}

//go:inline
func (sd *subDecode) bindNumber(val string) (err int) {
	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetStringByIndex(sd.keyIdx, val)
			return noErr
		}
		sd.mp.Set(sd.key, val)
		return noErr
	}

	// 如果是数组
	if sd.isList {
		if allowInt(sd.arr.itemKind) {
			if sd.isArray && len(sd.pl.arrI64) >= sd.arr.arrSize {
				return noErr
			}
			if num, err1 := parseInt(val); err < 0 {
				return err1
			} else {
				sd.pl.arrI64 = append(sd.pl.arrI64, num)
			}
		} else if allowFloat(sd.arr.itemKind) {
			if sd.isArray && len(sd.pl.arrF64) >= sd.arr.arrSize {
				return noErr
			}
			if num, err1 := strconv.ParseFloat(val, 64); err1 != nil {
				return errNumberFmt
			} else {
				sd.pl.arrF64 = append(sd.pl.arrF64, num)
			}
		} else {
			return errArray
		}
		return noErr
	}

	// 如果是 struct

	return noErr
}

func (sd *subDecode) bindNull() (err int) {
	//if sd.gr != nil {
	//	sd.gr.SetByIndex(sd.keyIdx, nil)
	//} else if sd.mp != nil {
	//	sd.mp.Set(sd.key, nil)
	//}

	// 如果是数组
	if sd.isList {
		return noErr
	}

	// 如果是 struct

	return noErr
}
