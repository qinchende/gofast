package jde

import (
	"reflect"
)

func (sd *subDecode) setSkipFlag() {
	// PS: 可以先判断目标对象是否有这个key，没有就跳过value，解析下一个kv
	if sd.gr != nil {
		if sd.keyIdx = sd.gr.KeyIndex(sd.key); sd.keyIdx < 0 {
			sd.skipValue = true
		}
	}
	// 如果是 struct ，就找找是否支持这个字段
}

func (sd *subDecode) isSkip() bool {
	if sd.skipValue || sd.skipTotal {
		sd.skipValue = false
		return true
	}
	return false
}

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
		if sd.arr.recKind != reflect.String {
			return errArray
		}
		pl.arrStr = append(pl.arrStr, val)
		return noErr
	}

	// 如果是 struct

	return noErr
}

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
		err = sd.arr.bindString(val)
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

func (ap *arrPet) bindString(val string) (err int) {
	if ap.recKind != reflect.String {
		return errArray
	}

	return noErr
}
