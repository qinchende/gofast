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
	//I           interface{}
	KV          map[string]interface{}
	IncHandler  func(w *GFResponse, r *http.Request)
	IncHandlers []IncHandler
	CtxHandler  func(ctx *Context)
	CtxHandlers []CtxHandler
	AppHandler  func(gft *GoFast)
	AppHandlers []AppHandler
	GFPanic     error

	// 抽取出一些常用函数原型
	goFastRegFunc func(*GoFast) *GoFast
)

const (
	//BodyBytesKey     = "_qinchende/gofast/bodybyteskey" // 记录POST提交时 body 的字节流，方便后期复用
	maxFitLen int = math.MaxInt8 // 最多多少个中间件函数
	//routePathMaxLen    uint8 = 255      // 路由字符串最长长度
	//routeMaxHandlers   uint8 = 255      // 路由 handlers 最大长度
	//defMultipartMemory int64 = 32 << 20 // 32 MB
)

// Content-Type MIME of the most common data formats.
// 常量值，供外部访问调用
const (
	MIMEJSON          = cst.MIMEJSON
	MIMEHTML          = cst.MIMEHTML
	MIMEAppXML        = cst.MIMEAppXML
	MIMETextXML       = cst.MIMETextXML
	MIMEPlain         = cst.MIMEPlain
	MIMEPOSTForm      = cst.MIMEPOSTForm
	MIMEMultiPOSTForm = cst.MIMEMultiPOSTForm
	MIMEYaml          = cst.MIMEYAML
)

var (
	spf            = fmt.Sprintf
	mimePlain      = []string{MIMEPlain}
	default404Body = []byte("404 (page not found)")
	default405Body = []byte("405 (method not allowed)")
)
