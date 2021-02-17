// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/render"
	"github.com/qinchende/gofast/logx"
)

type AppConfig struct {
	// FuncMap          		template.FuncMap
	RunMode               string // 当前模式[debug|test|product]
	SecureJsonPrefix      string
	HTMLRender            render.HTMLRender
	MaxMultipartMemory    int64
	SecondsBeforeShutdown int64 // 退出server之前等待的seconds，等待清理释放资源
	RedirectTrailingSlash bool  // 重定向URL结尾的`/`符号
	// RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	DisableDefNotAllowed   bool
	DisableDefNoRoute      bool
	ForwardedByClientIP    bool
	// UseRawPath             bool
	// UnescapePathValues     bool
	RemoveExtraSlash    bool  // 规范请求的URL
	PrintRouteTrees     bool  // 是否打印出当前路由数
	modeType            int8  // 运行模式，整形方便比较，提高性能
	FitTimeout          int64 // 每次请求的超时时间（单位：毫秒）
	FitMaxReqContentLen int64 // 最大请求字节数
}

func (gft *GoFast) initServerEnv() {
	//FuncMap: 				template.FuncMap{},
	//RedirectTrailingSlash:  true,
	//RedirectFixedPath:      false,
	//HandleMethodNotAllowed: false,
	//ForwardedByClientIP:    true,
	//AppEngine:              defaultAppEngine,
	//UseRawPath:             false,
	//RemoveExtraSlash:       false,
	//UnescapePathValues:     true,
	//MaxMultipartMemory:     defaultMultipartMemory,
	//trees:                  make(methodTrees, 0, 9),
	//delims:                 render.Delims{Left: "{{", Right: "}}"},
	//secureJsonPrefix:       "while(1);",
	if gft.SecureJsonPrefix == "" {
		gft.SecureJsonPrefix = "while(1);"
	}
	if gft.MaxMultipartMemory == 0 {
		gft.MaxMultipartMemory = defMultipartMemory
	}
	gft.RedirectTrailingSlash = true
	gft.ForwardedByClientIP = true
	//gft.UnescapePathValues = true

	gft.SetMode(gft.RunMode)
	logx.SetDebugStatus(gft.modeType == modeDebug)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++
// 当前运行处于啥模式：
const (
	modeDebug   int8 = iota // 0
	modeTest                // 1
	modeProduct             // 2
)

const (
	DebugMode   = "debug"
	TestMode    = "test"
	ProductMode = "product"
)

//func IsDebugMode() bool {
//	return appCfg.modeType == modeDebug
//}

func (gft *GoFast) SetMode(mode string) {
	switch mode {
	case DebugMode, "":
		gft.RunMode = DebugMode
		gft.modeType = modeDebug
	case ProductMode:
		gft.RunMode = ProductMode
		gft.modeType = modeProduct
	case TestMode:
		gft.RunMode = TestMode
		gft.modeType = modeTest
	default:
		panic("GoFast mode unknown: " + mode)
	}
}
