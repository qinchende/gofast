package jde

import (
	"github.com/qinchende/gofast/store/gson"
)

// Decoder +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DecodeGsonRowsFromString(v any, str string) gson.RowsDecRet {
	ret := decGsonRows(v, str)
	if ret.Err == nil && ret.Scan != len(str) {
		ret.Err = errJsonRowsStr
	}
	return ret
}

func DecodeGsonRowFromValueString(v any, str string) gson.RowsDecRet {
	ret := decGsonRows(v, str)
	if ret.Err == nil && ret.Scan != len(str) {
		ret.Err = errJsonRowsStr
	}
	return ret
}

// Encoder +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func EncodeGsonRows(v any) ([]byte, error) {
	return encGsonRows(gson.RowsEncPet{
		Target: v,
	})
}

func EncodeGsonRows2(v any, fls string) ([]byte, error) {
	return encGsonRows(gson.RowsEncPet{
		Target: v,
		FlsStr: fls,
	})
}

func EncodeGsonRowsPet(pet gson.RowsEncPet) ([]byte, error) {
	return encGsonRows(pet)
}
