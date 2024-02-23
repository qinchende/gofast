package jde

import (
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/store/gson"
)

// Decoder
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GsonRow ++++++
func DecodeGsonRowFromValueBytes(obj any, bs []byte) error {
	return DecodeGsonRowFromValueString(obj, lang.BTS(bs))
}

// 这里解析的 str字符串 只包含GsonRow的 values，而不包含 cls
// 所以API的名称叫：xxxValueString
func DecodeGsonRowFromValueString(obj any, str string) error {
	return decGsonRowOnlyValues(obj, str)
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
// GsonRow ++++++
func EncodeGsonRowOnlyValuesBytes(obj any) ([]byte, error) {
	return encGsonRowOnlyValues(obj)
}

func EncodeGsonRowOnlyValuesFromList(bf *[]byte, values []any) {
	encGsonRowFromValues(bf, values)
}

// GsonRows ++++++
func EncodeGsonRowsBytes(objs any) ([]byte, error) {
	return encGsonRows(gson.RowsEncPet{
		List: objs,
	})
}

func EncodeGsonRows2Bytes(objs any, cls []string) ([]byte, error) {
	return encGsonRows(gson.RowsEncPet{
		List: objs,
		Cls:  cls,
	})
}

func EncodeGsonRowsPetBytes(pet *gson.RowsEncPet) ([]byte, error) {
	return encGsonRows(*pet)
}

func EncodeGsonRowsPetString(pet *gson.RowsEncPet) (string, error) {
	ret, err := encGsonRows(*pet)
	return lang.BTS(ret), err
}
