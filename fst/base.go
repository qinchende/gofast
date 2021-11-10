// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"fmt"
	"github.com/qinchende/gofast/fst/cst"
	"math"
	"net/http"
)

type (
	KV          map[string]interface{}
	FitFunc     func(http.HandlerFunc) http.HandlerFunc
	CtxHandler  func(ctx *Context)
	CtxHandlers []CtxHandler
	AppHandler  func(gft *GoFast)
	AppHandlers []AppHandler
	GFPanic     error

	// 抽取出一些常用函数原型
	fitRegFunc func(*GoFast) *GoFast
)

const (
	//BodyBytesKey     = "_qinchende/gofast/bodybyteskey" // 记录POST提交时 body 的字节流，方便后期复用
	maxFits          uint8  = math.MaxUint8  // 最多多少个全局拦截器
	maxRouteHandlers int8   = math.MaxInt8   // 单路由最多中间件函数数量
	maxAllHandlers   uint16 = math.MaxUint16 // 全局所有路由节点的所有中间件函数最大总和
	//routePathMaxLen    uint8 = 255      // 路由字符串最长长度
	//routeMaxHandlers   uint8 = 255      // 路由 handlers 最大长度
	//defMultipartMemory int64 = 32 << 20 // 32 MB
)

var (
	spf            = fmt.Sprintf
	mimePlain      = []string{cst.MIMEPlain}
	default404Body = []byte("404 (page not found)")
	default405Body = []byte("405 (method not allowed)")
)
