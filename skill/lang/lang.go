package lang

import (
	"reflect"
	"runtime"
)

var Placeholder PlaceholderType

type (
	GenericType     = any
	PlaceholderType = struct{}
)

func NameOfFunc(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
