package lang

import (
	"reflect"
	"runtime"
)

var Placeholder PlaceholderType

type (
	GenericType     = interface{}
	PlaceholderType = struct{}
)

func NameOfFunc(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
