// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/sysx"
)

// GoFast WEB框架的配置参数
type GfConfig struct {
	LogConfig   logx.LogConfig
	AppName     string `v:"required"`
	ListenAddr  string `v:"def=0.0.0.0:8099,route=ipv4:port"`    // 监听ip:port
	RunningMode string `v:"def=product,enum=debug|test|product"` // 当前模式[debug|test|product]
	//LogType     string `v:"def=json,enum=json|sdx"`              // 日志类型

	// 配置主体Web框架控制参数
	BeforeShutdownMS      int64  `v:"def=1000"`      // 退出server之前等待的毫秒，等待清理释放资源
	RedirectTrailingSlash bool   `v:"def=false"`     // 探测url后面加减'/'之后是否能匹配路由（这个时代默认不需要了）
	CheckOtherMethodRoute bool   `v:"def=false"`     // 检查其它Method下，是否有对应的路由
	RemoveExtraSlash      bool   `v:"def=false"`     // 规范请求的URL
	UseRawPath            bool   `v:"def=false"`     // 默认取原始的Path，不需要自动转义
	UnescapePathValues    bool   `v:"def=true"`      // 是否把URL中的参数值做转义
	DefNotAllowedHandler  bool   `v:"def=true"`      // 是否采用默认的NotAllowed处理函数
	DefNoRouteHandler     bool   `v:"def=true"`      // 是否采用默认的NoRoute匹配函数
	ForwardedByClientIP   bool   `v:"def=true"`      // 是否从"X-Forwarded-For"的header中提取请求IP地址
	SecureJsonPrefix      string `v:"def=while(1);"` // JsonP安全前缀
	MaxMultipartBytes     int64  `v:"def=33554432"`  // 最大上传文件的大小，默认32MB
	ApplyUrlParamsToPms   bool   `v:"def=true"`      // 将UrlParams解析的参数自动加入Pms
	PrintRouteTrees       bool   `v:"def=false"`     // 是否打印出当前路由数

	SdxConfig cst.SdxConfig // middleware configs

	modeType int8 `v:""` // 内部记录状态
	//EnableRouteMonitor bool `cnf:",def=true"` // 是否统计路由的访问处理情况，为单个路由的熔断降载做储备
}

func (gft *GoFast) initServerConfig() {
	// 是否启动CPU检查
	if gft.SdxConfig.NeedSysCheck {
		sysx.StartSysCheck(gft.SdxConfig.NeedSysPrint)
	}
	gft.SetMode(gft.RunningMode)
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

func (gft *GoFast) SetMode(mode string) {
	switch mode {
	case DebugMode, "":
		gft.RunningMode = DebugMode
		gft.modeType = modeDebug
	case ProductMode:
		gft.RunningMode = ProductMode
		gft.modeType = modeProduct
	case TestMode:
		gft.RunningMode = TestMode
		gft.modeType = modeTest
	default:
		panic("GoFast mode unknown: " + mode)
	}
}

func (gft *GoFast) IsDebugging() bool {
	return gft.modeType == modeDebug
}

func (gft *GoFast) ProjectName() (name string) {
	name = gft.AppName
	if len(name) <= 0 {
		name = gft.ListenAddr
	}
	return
}
