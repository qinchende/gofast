package decode

//
//import (
//	"github.com/qinchende/gofast/store/gson"
//	"github.com/qinchende/gofast/store/jde"
//)
//
//type gsonDecode struct {
//	gr *gson.GsonRow
//	//keyIdx    int  // key index
//	//skipValue bool // 跳过当前要解析的值
//	//skipTotal bool // 跳过所有项目
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
////go:inline
//func (sd *gsonDecode) bindString(val string) (err int) {
//	//sd.obj.ss.BindString(sd.obj.objPtr, sd.keyIdx, val)
//	// sd.obj.setStringByIndex(sd.keyIdx, val)
//	return jde.noErr
//}
//
////go:inline
//func (sd *gsonDecode) bindBool(val bool) (err int) {
//	//sd.obj.ss.BindBool(sd.obj.objPtr, sd.keyIdx, val)
//	// sd.obj.setBoolByIndex(sd.keyIdx, val)
//	return jde.noErr
//}
//
////go:inline
//func (sd *gsonDecode) bindNumber(val string, hasDot bool) (err int) {
//	if _, err1 := jde.parseInt(val); err < 0 {
//		return err1
//	} else {
//		//sd.obj.ss.BindInt(sd.obj.objPtr, sd.keyIdx, num)
//		// sd.obj.setIntByIndex(sd.keyIdx, num)
//	}
//	return jde.noErr
//}
