package jde

import "reflect"

func (sd *subDecode) mustList() (err int) {
	if sd.kind != reflect.Slice && sd.kind != reflect.Array {
		return errArray
	}
	return noErr
}

func (sd *subDecode) mustObject() (err int) {
	if sd.mp == nil && sd.gr == nil && sd.kind != reflect.Struct {
		return errObject
	}
	return noErr
}

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
	if sd.gr != nil {
		sd.gr.SetStringByIndex(sd.keyIdx, val)
	} else if sd.mp != nil {
		sd.mp.Set(sd.key, val)
	}
	return noErr
}

func (sd *subDecode) bindBool(val bool) (err int) {
	if sd.gr != nil {
		sd.gr.SetByIndex(sd.keyIdx, val)
	} else if sd.mp != nil {
		sd.mp.Set(sd.key, val)
	}
	return noErr
}

func (sd *subDecode) bindNumber(val string) (err int) {
	if sd.gr != nil {
		sd.gr.SetStringByIndex(sd.keyIdx, val)
	} else if sd.mp != nil {
		sd.mp.Set(sd.key, val)
	}
	return noErr
}

func (sd *subDecode) bindNull() (err int) {
	//if sd.gr != nil {
	//	sd.gr.SetByIndex(sd.keyIdx, nil)
	//} else if sd.mp != nil {
	//	sd.mp.Set(sd.key, nil)
	//}
	return noErr
}
