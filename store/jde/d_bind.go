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

//go:inline
func (sd *subDecode) isSkip() bool {
	return sd.skipValue || sd.skipTotal
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (sd *subDecode) bindString(val string) (err int) {
	// 如果是 struct
	if sd.isStruct {
		sd.dm.ss.BindString(sd.dstPtr, sd.keyIdx, val)
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
	if sd.isList {
		// 如果是定长的数组，而且值不是指针类型，可以直接设置值
		if sd.isArray && !sd.isPtr {
			if sd.arrIdx >= sd.dm.arrLen {
				sd.skipValue = true
				return noErr
			}
			if sd.isAny {
				bindArrValue[any](sd, val)
			} else {
				bindArrValue[string](sd, val)
			}
			return noErr
		}
		if sd.isAny {
			sd.pl.bufAny = append(sd.pl.bufAny, val)
		} else {
			sd.pl.bufStr = append(sd.pl.bufStr, val)
		}
		return noErr
	}

	return noErr
}

func (sd *subDecode) bindBool(val bool) (err int) {
	// 如果是 struct
	if sd.isStruct {
		sd.dm.ss.BindBool(sd.dstPtr, sd.keyIdx, val)
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
	if sd.isList {
		if sd.isArray && !sd.isPtr {
			if sd.arrIdx >= sd.dm.arrLen {
				sd.skipValue = true
				return noErr
			}
			if sd.isAny {
				bindArrValue[any](sd, val)
			} else {
				bindArrValue[bool](sd, val)
			}
			return noErr
		}
		if sd.isAny {
			sd.pl.bufAny = append(sd.pl.bufAny, val)
		} else {
			sd.pl.bufBol = append(sd.pl.bufBol, val)
		}
		return noErr
	}

	return noErr
}

func (sd *subDecode) bindNumber(val string, hasDot bool) (err int) {
	// 如果是 struct
	if sd.isStruct {
		if num, err1 := parseInt(val); err < 0 {
			return err1
		} else {
			sd.dm.ss.BindInt(sd.dstPtr, sd.keyIdx, num)
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
	if sd.isList {
		if allowInt(sd.dm.itemKind) {
			if sd.isArray && !sd.isPtr {
				if sd.arrIdx >= sd.dm.arrLen {
					sd.skipValue = true
					return noErr
				}

				if i64, err1 := parseInt(val); err1 < 0 {
					return errNumberFmt
				} else {
					if sd.isAny {
						bindArrValue[any](sd, i64)
					} else {
						sd.dm.arrSetInt(sd, i64)
					}
				}
				return noErr
			}

			if num, err1 := parseInt(val); err < 0 {
				return err1
			} else {
				if sd.isAny {
					sd.pl.bufAny = append(sd.pl.bufAny, num)
				} else {
					sd.pl.bufI64 = append(sd.pl.bufI64, num)
				}
			}

		} else if allowFloat(sd.dm.itemKind) {
			if sd.isArray && !sd.isPtr {
				if sd.arrIdx >= sd.dm.arrLen {
					sd.skipValue = true
					return noErr
				}

				if f64, err1 := strconv.ParseFloat(val, 64); err1 != nil {
					return errNumberFmt
				} else {
					if sd.isAny {
						bindArrValue[any](sd, f64)
					} else {
						sd.dm.arrSetFloat(sd, f64)
					}
				}
				return noErr
			}

			if num, err1 := strconv.ParseFloat(val, 64); err1 != nil {
				return errNumberFmt
			} else {
				if sd.isAny {
					sd.pl.bufAny = append(sd.pl.bufAny, num)
				} else {
					sd.pl.bufF64 = append(sd.pl.bufF64, num)
				}
			}

		} else {
			return errList
		}

		return noErr
	}

	return noErr
}

func (sd *subDecode) bindNull() (err int) {
	// 如果是数组
	if sd.isList {
		return noErr
	}

	// 如果是 struct

	return noErr
}
