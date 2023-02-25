// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GoFast框架主动抛异常
func PanicIf(ifTrue bool, val any) {
	if ifTrue {
		Panic(val)
	}
}

// 为了和Runtime抛异常区别开来，GoFast主动抛出的异常都是自定义数据类型
func Panic(val any) {
	if val == nil {
		return
	}

	switch val.(type) {
	case string:
		panic(cst.TypeString(val.(string)))
	case error:
		panic(cst.TypeError(val.(error)))
	case int:
		panic(cst.TypeInt(val.(int)))
	default:
		panic(cst.TypeString(lang.ToString(val)))
	}
}

func PanicIfErr(err error) {
	if err != nil {
		panic(cst.TypeError(err))
	}
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 简易的抛出异常的方式，终止执行链，返回错误
//func (c *Context) PanicIf(yes bool, val any) {
//	if !yes {
//		return
//	}
//
//	switch val.(type) {
//	case string:
//		panic(cst.TypeString(val.(string)))
//	case error:
//		panic(cst.TypeError(val.(error)))
//	case int:
//		panic(cst.TypeInt(val.(int)))
//	default:
//		panic(cst.TypeString(lang.ToString(val)))
//	}
//}
//
//// 如果发现有错误信息，就抛异常终止handlers
//func (c *Context) PanicIfError(err error) {
//	c.PanicIf(err != nil, err)
//}
