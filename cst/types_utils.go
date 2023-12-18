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
	TypeCstKV     reflect.Type
	TypeStrAnyMap reflect.Type
	TypeTime      reflect.Type
	TypeBytes     reflect.Type
)

func init() {
	TypeCstKV = reflect.TypeOf(KV{})
	TypeStrAnyMap = reflect.TypeOf(map[string]any{})
	TypeTime = reflect.TypeOf(time.Time{})
	TypeBytes = reflect.TypeOf(make([]byte, 0))
}
