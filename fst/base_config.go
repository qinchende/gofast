// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/render"
	"github.com/qinchende/gofast/logx"
)

// GoFast WEB框架的配置参数
type AppConfig struct {
	// FuncMap          	template.FuncMap
	// UseRawPath           bool
	// UnescapePathValues   bool
	// RedirectFixedPath    bool
	Addr                   string `json:",default=127.0.0.1:8099"`
	Name                   string `json:",optional,default=GoFastSite"`
	RunMode                string `json:",default=debug,options=debug|test|product"` // 当前模式[debug|test|product]
	SecureJsonPrefix       string `json:",optional,default=while(1);"`
	MaxMultipartMemory     int64  `json:",default=33554432"` // 最大上传文件的大小
	SecondsBeforeShutdown  int64  `json:",default=1000"`     // 退出server之前等待的毫秒，等待清理释放资源
	RedirectTrailingSlash  bool   `json:",default=true"`     // 重定向URL结尾的`/`符号
	HandleMethodNotAllowed bool   `json:",default=false"`
	DisableDefNotAllowed   bool   `json:",default=false"`
	DisableDefNoRoute      bool   `json:",default=false"`
	ForwardedByClientIP    bool   `json:",default=true"`
	RemoveExtraSlash       bool   `json:",default=false"`                       // 规范请求的URL
	PrintRouteTrees        bool   `json:",default=true"`                        // 是否打印出当前路由数
	FitReqTimeout          int64  `json:",default=3000"`                        // 每次请求的超时时间（单位：毫秒）
	FitMaxReqContentLen    int64  `json:",default=33554432"`                    // 最大请求字节数
	FitMaxReqCount         int32  `json:",default=1000000,range=[0:100000000]"` // 最大请求处理数
	FitJwtSecret           string `json:",optional"`                            // JWT认证的秘钥
	FitLogType             string `json:",default=json,options=json|sdx"`

	HTMLRender render.HTMLRender `json:",optional"`

	// 内部记录状态
	modeType int8
}

func (gft *GoFast) initServerEnv() {
	//if gft.MaxMultipartMemory == 0 {
	//	gft.MaxMultipartMemory = defMultipartMemory
	//}
	//if gft.FitMaxReqContentLen == 0 {
	//	gft.FitMaxReqContentLen = defMultipartMemory
	//}

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

// 日志文件的目标系统
const (
	LogTypeConsole    = "console"
	LogTypeELK        = "elk"
	LogTypePrometheus = "prometheus"
)
