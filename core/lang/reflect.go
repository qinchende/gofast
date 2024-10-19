package lang

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

func FuncFullName(fnc any) string {
	return runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name()
}

func FuncName(fnc any) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fnc).Pointer()).Name()
	return strings.TrimLeft(filepath.Ext(fullName), ".")
}
