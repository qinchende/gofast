// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"math"
	"net/http"
)

func init() {
	checkRuntimeVer()
}

type (
	AppHandler  func(gft *GoFast)
	HttpHandler func(http.HandlerFunc) http.HandlerFunc
	CtxHandler  func(ctx *Context)

	// 抽取出一些常用函数原型
	injectFunc func(*GoFast) *GoFast
)

const (
	gftSupportMinGoVer float64 = 1.18           // 支持的最小GO版本是 1.18 and later
	maxCtxCarryLen     int     = 8              // 请求上下文能携带的最大扩展数据项
	maxHttpHandlers    uint8   = math.MaxUint8  // 最多多少个全局拦截器
	maxRouteHandlers   int8    = math.MaxInt8   // 单路由最多中间件函数数量
	maxAllHandlers     uint16  = math.MaxUint16 // 全局所有路由节点的所有中间件函数最大总和
)

var (
	spf = fmt.Sprintf
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主动抛异常
func GFPanicIf(yes bool, val any) {
	if !yes {
		return
	}
	GFPanic(val)
}

func GFPanicErr(err error) {
	if err != nil {
		panic(cst.GFError(err))
	}
}

func GFPanic(val any) {
	if val == nil {
		return
	}

	switch val.(type) {
	case string:
		str := val.(string)
		if len(str) != 0 {
			panic(cst.GFError(errors.New(str)))
		}
	case error:
		panic(cst.GFError(val.(error)))
	default:
		str := lang.ToString(val)
		if len(str) != 0 {
			panic(cst.GFError(errors.New(str)))
		}
	}
}
