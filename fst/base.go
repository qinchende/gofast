// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

import (
	"fmt"
	"github.com/qinchende/gofast/fst/binding"
)

type I interface{}
type KV map[string]interface{}
type CtxHandler func(*Context)
type CtxHandlers []CtxHandler
type AppHandler func(*GoFast)
type AppHandlers []AppHandler

const (
	// BodyBytesKey indicates a default body bytes key.
	BodyBytesKey = "_qinchende/gofast/bodybyteskey"

	//routePathMaxLen    uint8 = 255      // 路由字符串最长长度
	//routeMaxHandlers   uint8 = 255      // 路由 handlers 最大长度
	defMultipartMemory int64 = 32 << 20 // 32 MB
)

// Content-Type MIME of the most common data formats.
// 常量值，供外部访问调用
const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
	MIMEYAML              = binding.MIMEYAML
)

var (
	spf            = fmt.Sprintf
	mimePlain      = []string{MIMEPlain}
	default404Body = []byte("404 (PAGE NOT FOND)")
	default405Body = []byte("405 (METHOD NOT ALLOWED)")
)
