package lang

import (
	"reflect"
	"runtime"
)

var Placeholder PlaceholderType

type (
	PlaceholderType = struct{}
)

func NameOfFunc(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
