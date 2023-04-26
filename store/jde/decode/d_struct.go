package decode

//
//type structDecode struct {
//	//pl  *fastPool
//	obj *structPost // struct pet
//
//	skipValue bool // 跳过当前要解析的值
//	skipTotal bool // 跳过所有项目
//
//	//key    string // 当前KV对的Key值
//	keyIdx int // key index
//
//	//str    string // 本段字符串
//	//scan   int    // 自己的扫描进度，当解析错误时，这个就是定位
//	//key    string // 当前KV对的Key值
//	//keyIdx int    // key index
//
//}
//
////go:inline
//func (sd *structDecode) setSkip(key string) {
//	if sd.keyIdx = sd.obj.ss.ColumnIndex(key); sd.keyIdx < 0 {
//		sd.skipValue = true
//	} else {
//		sd.skipValue = false
//	}
//	return
//}
//
////go:inline
//func (sd *structDecode) isSkip() bool {
//	return sd.skipValue || sd.skipTotal
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (sd *structDecode) keyIndex(val string) int {
//	return sd.obj.ss.ColumnIndex(val)
//}
//
////go:inline
//func (sd *structDecode) bindString(val string) (err int) {
//	sd.obj.ss.BindString(sd.obj.objPtr, sd.keyIdx, val)
//	// sd.obj.setStringByIndex(sd.keyIdx, val)
//	return noErr
//}
//
////go:inline
//func (sd *structDecode) bindBool(val bool) (err int) {
//	sd.obj.ss.BindBool(sd.obj.objPtr, sd.keyIdx, val)
//	// sd.obj.setBoolByIndex(sd.keyIdx, val)
//	return noErr
//}
//
////go:inline
//func (sd *structDecode) bindNumber(val string, hasDot bool) (err int) {
//	if num, err1 := parseInt(val); err < 0 {
//		return err1
//	} else {
//		sd.obj.ss.BindInt(sd.obj.objPtr, sd.keyIdx, num)
//		// sd.obj.setIntByIndex(sd.keyIdx, num)
//	}
//	return noErr
//}
//
//func (sd *structDecode) bindNull() int {
//	//return sd.obj.ss.ColumnIndex(val)
//	return noErr
//}
