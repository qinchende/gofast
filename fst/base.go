// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast"
	"github.com/qinchende/gofast/core/dts"
	"log"
	"math"
	"net/http"
)

type (
	AppHandler  func(gft *GoFast)
	HttpHandler func(http.HandlerFunc) http.HandlerFunc
	CtxHandler  func(ctx *Context)

	appSelfFunc func(*GoFast) *GoFast
)

const (
	gftSupportMinGoVer float64 = 1.20           // 支持的最小GO版本是 1.20 and later
	maxHttpHandlers    uint8   = math.MaxUint8  // 最多多少个全局拦截器
	maxRouteHandlers   int8    = math.MaxInt8   // 单路由最多中间件函数数量
	maxAllHandlers     uint16  = math.MaxUint16 // 全局所有路由节点的所有中间件函数最大总和
)

var (
	pBindOptions         = dts.AsOptions(dts.AsDef)
	pBindAndValidOptions = dts.AsOptions(dts.AsReq)
)

func init() {
	log.Println("Current GoFast version: " + gofast.Version)

	// 检查Go版本是否符合要求
	checkRuntimeVer()

	// 初始化包共享变量
	pBindOptions.CacheSchema = true
}
