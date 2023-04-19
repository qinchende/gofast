package jde

import (
	"reflect"
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
		if sd.arr.itemKind != reflect.String {
			return errArray
		}
		if sd.isArray {

		} else {
			sd.pl.arrStr = append(sd.pl.arrStr, val)
		}
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
		if !isNumKind(sd.arr.itemKind) {
			return errArray
		}
		if sd.isArray {

		} else {
			if num, err1 := parseInt(val); err < 0 {
				return err1
			} else {
				sd.pl.arrI64 = append(sd.pl.arrI64, num)
			}
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

// Set Values
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (ap *listDest) bindString(val string) (err int) {
	if ap.itemKind != reflect.String {
		return errArray
	}

	return noErr
}
