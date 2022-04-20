// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/sysx"
)

// GoFast WEB框架的配置参数
type AppConfig struct {
	LogConfig logx.LogConfig
	// FuncMap          	template.FuncMap
	// RedirectFixedPath    bool // 此项特性无多大必要，不兼容Gin
	Name                   string `cnf:",def=GoFastSite"`
	Addr                   string `cnf:",def=0.0.0.0:8099"`
	RunMode                string `cnf:",def=debug,enum=debug|test|product"` // 当前模式[debug|test|product]
	SecureJsonPrefix       string `cnf:",def=while(1);"`
	MaxMultipartMemory     int64  `cnf:",def=33554432"` // 最大上传文件的大小，默认32MB
	SecondsBeforeShutdown  int64  `cnf:",def=1000"`     // 退出server之前等待的毫秒，等待清理释放资源
	RedirectTrailingSlash  bool   `cnf:",def=false"`    // 探测url后面加减'/'之后是否能匹配路由（这个时代默认不需要了）
	HandleMethodNotAllowed bool   `cnf:",def=false"`
	DisableDefNotAllowed   bool   `cnf:",def=false"`
	DisableDefNoRoute      bool   `cnf:",def=false"`
	ForwardedByClientIP    bool   `cnf:",def=true"`
	RemoveExtraSlash       bool   `cnf:",def=false"`                       // 规范请求的URL
	UseRawPath             bool   `cnf:",def=false"`                       // 默认取原始的Path，不需要自动转义
	UnescapePathValues     bool   `cnf:",def=true"`                        // 默认把URL中的参数值做转义
	PrintRouteTrees        bool   `cnf:",def=false"`                       // 是否打印出当前路由数
	FitDisableTimeout      bool   `cnf:",def=false"`                       // 是否支持超时自动返回
	NeedSysCheck           bool   `cnf:",def=true"`                        // 是否启动CPU使用情况的定时检查工作
	NeedSysPrint           bool   `cnf:",def=true"`                        // 定时打印系统检查日志
	FitReqTimeout          int64  `cnf:",def=3000"`                        // 每次请求的超时时间（单位：毫秒）
	FitMaxContentLength    int64  `cnf:",def=0"`                           // 最大请求字节数，32MB（def=33554432）
	FitMaxConnections      int32  `cnf:",def=1000000,range=[0:100000000]"` // 最大请求处理数，默认100万个请求同时进入
	FitJwtSecret           string `cnf:",NA"`                              // JWT认证的秘钥
	FitLogType             string `cnf:",def=json,enum=json|sdx"`
	modeType               int8   `cnf:",NA"` // 内部记录状态
	//HTMLRender             render.HTMLRender `cnf:",NA"`
	//EnableRouteMonitor bool `cnf:",def=true"` // 是否统计路由的访问处理情况，为单个路由的熔断降载做储备
}

func (gft *GoFast) initServerConfig() {
	//if gft.MaxMultipartMemory == 0 {
	//	gft.MaxMultipartMemory = defMultipartMemory
	//}
	//if gft.FitMaxReqContentLen == 0 {
	//	gft.FitMaxReqContentLen = defMultipartMemory
	//}

	// 是否启动CPU检查
	if gft.NeedSysCheck {
		sysx.StartSysCheck(gft.NeedSysPrint)
	}
	gft.SetMode(gft.RunMode)
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
	logx.SetDebugStatus(gft.modeType == modeDebug)
}

func (gft *GoFast) IsDebugging() bool {
	return gft.modeType == modeDebug
}

func (gft *GoFast) AppName() (name string) {
	name = gft.Name
	if len(name) <= 0 {
		name = gft.Addr
	}
	return
}

// 日志文件的目标系统
const (
	LogTypeConsole    = "console"
	LogTypeELK        = "elk"
	LogTypePrometheus = "prometheus"
)
