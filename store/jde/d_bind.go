package jde

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

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
		if sd.isArrBind {
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
		if sd.isArrBind {
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// numbers
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

	if sd.isList {
		sd.bindFloatList(val)
	}
}

func (sd *subDecode) bindIntList(val string) {
	if sd.isArrBind {
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
}

func (sd *subDecode) bindUintList(val string) {
	if sd.isArrBind {
		if sd.arrIdx >= sd.dm.arrLen {
			sd.skipValue = true
			return
		}
		sd.bindUintArr(parseUint(val))
		return
	}

	if sd.isAny {
		sd.pl.bufAny = append(sd.pl.bufAny, parseUint(val))
	} else {
		sd.pl.bufU64 = append(sd.pl.bufU64, parseUint(val))
	}
}

func (sd *subDecode) bindFloatList(val string) {
	if sd.isArrBind {
		if sd.arrIdx >= sd.dm.arrLen {
			sd.skipValue = true
			return
		}
		sd.bindFloatArr(parseFloat(val))
		return
	}

	if sd.isAny {
		sd.pl.bufAny = append(sd.pl.bufAny, parseFloat(val))
	} else {
		sd.pl.bufF64 = append(sd.pl.bufF64, parseFloat(val))
	}
}

// null value
// 分不同函数，避免分支条件判断影响性能
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) bindIntNull() {
	if sd.isArrBind {
		if sd.arrIdx >= sd.dm.arrLen {
			sd.skipValue = true
			return
		}
		sd.arrIdx++
		return
	}
	if sd.isAny {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufAny))
		}
		sd.pl.bufAny = append(sd.pl.bufAny, nil)
	} else {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufI64))
		}
		sd.pl.bufI64 = append(sd.pl.bufI64, 0)
	}
}

func (sd *subDecode) bindUintNull() {
	if sd.isArrBind {
		if sd.arrIdx >= sd.dm.arrLen {
			sd.skipValue = true
			return
		}
		sd.arrIdx++
		return
	}
	if sd.isAny {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufAny))
		}
		sd.pl.bufAny = append(sd.pl.bufAny, nil)
	} else {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufU64))
		}
		sd.pl.bufU64 = append(sd.pl.bufU64, 0)
	}
}

func (sd *subDecode) bindFloatNull() {
	if sd.isArrBind {
		if sd.arrIdx >= sd.dm.arrLen {
			sd.skipValue = true
			return
		}
		sd.arrIdx++
		return
	}

	if sd.isAny {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufAny))
		}
		sd.pl.bufAny = append(sd.pl.bufAny, nil)
	} else {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufF64))
		}
		sd.pl.bufF64 = append(sd.pl.bufF64, 0)
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) bindNumberNull() {
	sd.bindNull(func() {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufF64))
		}
		sd.pl.bufF64 = append(sd.pl.bufF64, 0)
	})
}

func (sd *subDecode) bindStringNull() {
	sd.bindNull(func() {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufStr))
		}
		sd.pl.bufStr = append(sd.pl.bufStr, "")
	})
}

func (sd *subDecode) bindBoolNull() {
	sd.bindNull(func() {
		if sd.isPtr {
			sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufBol))
		}
		sd.pl.bufBol = append(sd.pl.bufBol, false)
	})
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) bindNull(fn func()) {
	if sd.isList {
		if sd.isArrBind {
			if sd.arrIdx >= sd.dm.arrLen {
				sd.skipValue = true
				return
			}
			sd.arrIdx++
			return
		}
		if sd.isAny {
			if sd.isPtr {
				sd.pl.nilPos = append(sd.pl.nilPos, len(sd.pl.bufAny))
			}
			sd.pl.bufAny = append(sd.pl.bufAny, nil)
		} else {
			fn()
		}
		return
	}

	if sd.isSuperKV {
		if sd.gr != nil {
			sd.gr.SetByIndex(sd.keyIdx, nil)
			return
		}
		sd.mp.Set(sd.key, nil)
		return
	}

	//// 如果是 struct，字段是默认值，相当于啥也不做
	//if sd.isStruct {
	//	return
	//}
}
