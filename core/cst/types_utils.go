package cst

import (
	"reflect"
	"time"
)

//const (
//	StrTypeOfKV        = "cst.KV"
//	StrTypeOfWebKV     = "cst.WebKV"
//	StrTypeOfStrAnyMap = "map[string]interface {}"
//	StrTypeOfStrStrMap = "map[string]string"
//	StrTypeOfTime      = "time.Time"
//	StrTypeOfBytes     = "[]byte"
//)

var (
	// common map
	TypeCstKV     reflect.Type
	TypeWebKV     reflect.Type
	TypeMapStrAny reflect.Type
	TypeMapStrStr reflect.Type

	TypeTime     reflect.Type
	TypeDuration reflect.Type
	TypeBytes    reflect.Type
)

func init() {
	TypeCstKV = reflect.TypeOf(KV{})
	TypeWebKV = reflect.TypeOf(WebKV{})
	TypeMapStrAny = reflect.TypeOf(map[string]any{})
	TypeMapStrStr = reflect.TypeOf(map[string]string{})

	TypeTime = reflect.TypeOf(time.Time{})
	TypeDuration = reflect.TypeOf(time.Duration(0))
	TypeBytes = reflect.TypeOf(make([]byte, 0))
}
