package jde

import (
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/gson"
)

// Decoder
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GsonRow ++++++
func DecodeGsonRowFromValueBytes(obj any, bs []byte) gson.RowsDecRet {
	return DecodeGsonRowFromValueString(obj, lang.BTS(bs))
}

func DecodeGsonRowFromValueString(obj any, str string) gson.RowsDecRet {
	ret := decGsonRow(obj, str)
	if ret.Err == nil && ret.Scan != len(str) {
		ret.Err = errJsonRowStr
	}
	return ret
}

// GsonRows ++++++
func DecodeGsonRowsFromString(objs any, str string) gson.RowsDecRet {
	ret := decGsonRows(objs, str)
	if ret.Err == nil && ret.Scan != len(str) {
		ret.Err = errJsonRowsStr
	}
	return ret
}

// Encoder +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func EncodeGsonRows(objs any) ([]byte, error) {
	return encGsonRows(gson.RowsEncPet{
		Target: objs,
	})
}

func EncodeGsonRows2(objs any, fls string) ([]byte, error) {
	return encGsonRows(gson.RowsEncPet{
		Target: objs,
		FlsStr: fls,
	})
}

func EncodeGsonRowsPet(pet gson.RowsEncPet) ([]byte, error) {
	return encGsonRows(pet)
}
