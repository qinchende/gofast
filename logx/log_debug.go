// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const gftSupportMinGoVer = 14

var isDebug = false
var _once sync.Once

// 每个程序只设置一次debug标志位，后面的设置都失效
func SetDebugStatus(yn bool) {
	_once.Do(func() {
		isDebug = yn
	})
}

func IsDebugging() bool {
	return isDebug
}

func DebugPrint(format string, v ...any) {
	if isDebug {
		Info("[Debug] ", fmt.Sprintf(format, v...))
	}
}

func DebugPrintError(err error) {
	if err != nil && isDebug {
		Info("[Debug] ", fmt.Sprintf("[ERROR] %v\n", err))
	}
}

func GetMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func DebugPrintWarningDefault() {
	if v, e := GetMinVer(runtime.Version()); e == nil && v <= gftSupportMinGoVer {
		DebugPrint("[WARN] Now GoFast requires Go 1.14 or later and Go 1.16 will be required soon.\n")
	}
	DebugPrint("[WARN] Creating an Engine instance with the Logger and Recovery middleware already attached.\n")
}

//func DebugPrintWARNINGNew() {
//	DebugPrint(`[WARN] Running in "debug" mode. Switch to "release" mode in production.
// - using env:	export GoFast_MODE=release
// - using code:	fst.SetMode(fst.ReleaseMode)
//
//`)
//}
//
//func DebugPrintWARNINGSetHTMLTemplate() {
//	DebugPrint(`[WARN] Since SetHTMLTemplate() is NOT thread-safe. It should only be called
//at initialization. ie. before any route is registered or the router is listening in a socket:
//
//	router := fst.Default()
//	router.SetHTMLTemplate(template) // << good place
//
//`)
//}
