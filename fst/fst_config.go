// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/aid/logx"
	"time"
)

// GoFast Framework Config
type ServerConfig struct {
	AppName        string        `v:"must"`                                // 应用名称
	HostName       string        `v:"must"`                                // 运行终端编号
	ListenAddr     string        `v:"def=0.0.0.0:8099,match=ipv4:port"`    // 监听ip:port
	RunMode        string        `v:"def=product,enum=debug|test|product"` // 当前模式[debug|test|product]
	BeforeShutdown time.Duration `v:"def=1000ms"`                          // 退出APP之前等待的毫秒，等待清理释放资源

	WebConfig WebConfig      // http config
	LogConfig logx.LogConfig // 日志配置

	iRunMode int8 // 运行模式
}

// 配置http相关设置，比如路由相关控制参数
type WebConfig struct {
	SecureJsonPrefix      string `v:"def=while(1);"` // JsonP安全前缀
	MaxMultipartBytes     int64  `v:"def=33554432"`  // 最大上传文件的大小，默认32MB
	RedirectTrailingSlash bool   `v:"def=false"`     // 探测url后面加减'/'之后是否能匹配路由（这个时代默认不需要了）
	CheckOtherMethodRoute bool   `v:"def=false"`     // 检查其它Method下，是否有对应的路由
	RemoveExtraSlash      bool   `v:"def=false"`     // 规范请求的URL
	UseRawPath            bool   `v:"def=false"`     // 默认取原始的Path，不需要自动转义
	UnescapePathValues    bool   `v:"def=true"`      // 是否把URL中的参数值做转义
	ForwardedByClientIP   bool   `v:"def=true"`      // 是否从"X-Forwarded-For"的header中提取请求IP地址
	ApplyUrlParamsToPms   bool   `v:"def=true"`      // 将UrlParams解析的参数自动加入Pms
	PrintRouteTrees       bool   `v:"def=false"`     // 是否打印出当前路由数
	CacheQueryValues      bool   `v:"def=true"`      // 存储解析后的QueryValues，方便下次访问

	//LogType     string `v:"def=json,enum=json|sdx"`              // 日志类型
	//EnableRouteMonitor bool `cnf:",def=true"` // 是否统计路由的访问处理情况，为单个路由的熔断降载做储备
	//DefNotAllowedHandler  bool   `v:"def=true"`      // 是否采用默认的NotAllowed处理函数
	//DefNoRouteHandler     bool   `v:"def=true"`      // 是否采用默认的NoRoute匹配函数
}

func (gft *GoFast) ProjectName() (name string) {
	name = gft.AppName
	if len(name) <= 0 {
		name = gft.ListenAddr
	}
	return
}

// 是否处于Debug模式
func (gft *GoFast) IsDebugging() bool {
	return gft.iRunMode == modeDebug
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 当前运行处于啥模式：
const (
	modeDebug   int8 = iota // 0
	modeDev                 // 1
	modeTest                // 2
	modeProduct             // 3
)

const (
	DebugMode   = "debug"
	DevMode     = "dev"
	TestMode    = "test"
	ProductMode = "product"
)

// 设置运行模式
func (gft *GoFast) SetMode(mode string) {
	switch mode {
	case DebugMode, "":
		gft.RunMode = DebugMode
		gft.iRunMode = modeDebug
	case DevMode:
		gft.RunMode = DevMode
		gft.iRunMode = modeDev
	case ProductMode:
		gft.RunMode = ProductMode
		gft.iRunMode = modeProduct
	case TestMode:
		gft.RunMode = TestMode
		gft.iRunMode = modeTest
	default:
		panic("GoFast mode unknown: " + mode)
	}
}
