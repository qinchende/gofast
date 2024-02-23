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
	TypeWebKV     reflect.Type
	TypeStrAnyMap reflect.Type
	TypeTime      reflect.Type
	TypeDuration  reflect.Type
	TypeBytes     reflect.Type
)

func init() {
	TypeCstKV = reflect.TypeOf(KV{})
	TypeWebKV = reflect.TypeOf(WebKV{})
	TypeStrAnyMap = reflect.TypeOf(map[string]any{})
	TypeTime = reflect.TypeOf(time.Time{})
	TypeDuration = reflect.TypeOf(time.Duration(0))
	TypeBytes = reflect.TypeOf(make([]byte, 0))
}
