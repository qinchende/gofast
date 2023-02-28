package cst

import "github.com/qinchende/gofast/skill/lang"

// GoFast框架主动抛异常
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 为了和Runtime抛异常区别开来，GoFast主动抛出的异常都是自定义数据类型
func Panic(val any) {
	if val == nil {
		return
	}

	switch val.(type) {
	case string:
		panic(TypeString(val.(string)))
	case error:
		panic(TypeError(val.(error)))
	case int:
		panic(TypeInt(val.(int)))
	default:
		panic(TypeString(lang.ToString(val)))
	}
}

func PanicIf(ifTrue bool, val any) {
	if ifTrue {
		Panic(val)
	}
}

func PanicIfErr(err error) {
	if err != nil {
		panic(TypeError(err))
	}
}

func PanicString(str string) {
	panic(TypeString(str))
}
