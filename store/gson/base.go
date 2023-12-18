package gson

import (
	"reflect"
)

//const (
//	StrTypeOfGsonRow  = "gson.GsonRow"
//	StrTypeOfGsonRows = "gson.GsonRows"
//	StrTypeOfRowsRet  = "gson.RowsDecRet"
//	StrTypeOfRowsPet  = "gson.RowsDecPet"
//)

var (
	TypeGsonRow    reflect.Type
	TypeRowsDecPet reflect.Type
)

func init() {
	TypeGsonRow = reflect.TypeOf(GsonRow{})
	TypeRowsDecPet = reflect.TypeOf(RowsDecPet{})
}
