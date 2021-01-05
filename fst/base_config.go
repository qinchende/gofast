// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

import (
	"gofast/fst/render"
	"gofast/skill"
)

type AppConfig struct {
	//FuncMap          		template.FuncMap
	RunMode                string // 当前模式，字符串
	SecureJsonPrefix       string
	HTMLRender             render.HTMLRender
	MaxMultipartMemory     int64
	SecondsBeforeShutdown  int64
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	ForwardedByClientIP    bool
	UseRawPath             bool
	UnescapePathValues     bool
	RemoveExtraSlash       bool
	PrintRouteTrees        bool // 是否打印出当前路由数
	modeType               int8 // 运行模式，整形方便比较，提高性能

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

	gft.SetMode(gft.RunMode)
	skill.SetDebugStatus(gft.modeType == modeDebug)
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
		gft.modeType = modeDebug
	case ProductMode:
		gft.modeType = modeProduct
	case TestMode:
		gft.modeType = modeTest
	default:
		panic("GoFast mode unknown: " + mode)
	}
}
