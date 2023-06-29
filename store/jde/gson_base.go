package jde

import (
	"github.com/qinchende/gofast/store/gson"
)

//func DecodeGsonRowsFromBytes(v any, source []byte) error {
//	//return decodeFromString(v, source)
//	return nil
//}

func DecodeGsonRowsFromString(v any, str string) gson.RowsRet {
	ret := decGsonRows(v, str)
	if ret.Err == nil && ret.Scan != len(str) {
		ret.Err = errJsonRowsStr
	}
	return ret
}
