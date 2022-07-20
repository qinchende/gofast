// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/sysx"
)

// GoFast WEB框架的配置参数
type GfConfig struct {
	LogConfig logx.LogConfig
	// FuncMap          	template.FuncMap
	// RedirectFixedPath    bool // 此项特性无多大必要，不兼容Gin
	Name                  string `v:"def=GoFastSite"`
	Addr                  string `v:"def=0.0.0.0:8099,match=ipv4:port"`
	RunMode               string `v:"def=debug,enum=debug|test|product"` // 当前模式[debug|test|product]
	SecureJsonPrefix      string `v:"def=while(1);"`
	MaxMultipartMemory    int64  `v:"def=33554432"` // 最大上传文件的大小，默认32MB
	SecondsBeforeShutdown int64  `v:"def=1000"`     // 退出server之前等待的毫秒，等待清理释放资源
	RedirectTrailingSlash bool   `v:"def=false"`    // 探测url后面加减'/'之后是否能匹配路由（这个时代默认不需要了）
	CheckOtherMethodRoute bool   `v:"def=false"`    // 检查其它Method下，是否有对应的路由
	DefNotAllowedHandler  bool   `v:"def=true"`     // 是否采用默认的NotAllowed处理函数
	DefNoRouteHandler     bool   `v:"def=true"`     // 是否采用默认的NoRoute匹配函数
	ForwardedByClientIP   bool   `v:"def=true"`
	RemoveExtraSlash      bool   `v:"def=false"`                       // 规范请求的URL
	UseRawPath            bool   `v:"def=false"`                       // 默认取原始的Path，不需要自动转义
	UnescapePathValues    bool   `v:"def=true"`                        // 默认把URL中的参数值做转义
	PrintRouteTrees       bool   `v:"def=false"`                       // 是否打印出当前路由数
	NeedSysCheck          bool   `v:"def=true"`                        // 是否启动CPU使用情况的定时检查工作
	NeedSysPrint          bool   `v:"def=true"`                        // 定时打印系统检查日志
	SdxEnableTimeout      bool   `v:"def=true"`                        // 默认启动超时拦截
	SdxDefTimeout         int64  `v:"def=3000"`                        // 每次请求的超时时间（单位：毫秒）
	FitMaxContentLength   int64  `v:"def=33554432"`                    // 最大请求字节数，32MB（33554432），传0不限制
	FitMaxConnections     int32  `v:"def=1000000,range=[0:100000000]"` // 最大同时请求数，默认100万同时进入，传0不限制
	FitJwtSecret          string `v:""`                                // JWT认证的秘钥
	FitLogType            string `v:"def=json,enum=json|sdx"`          // 日志类型
	modeType              int8   `v:""`                                // 内部记录状态
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
