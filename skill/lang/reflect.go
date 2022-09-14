package lang

import (
	"reflect"
	"runtime"
)

func NameOfFunc(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
