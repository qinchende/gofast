// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/aid/logx"
	"time"
)

// GoFast WEB框架的配置参数
type AppConfig struct {
	AppName        string        `v:"must"`                                // 应用名称
	ServerName     string        `v:"must"`                                // 运行终端编号
	ListenAddr     string        `v:"def=0.0.0.0:8099,match=ipv4:port"`    // 监听ip:port
	RunMode        string        `v:"def=product,enum=debug|test|product"` // 当前模式[debug|test|product]
	BeforeShutdown time.Duration `v:"def=1000ms"`                          // 退出APP之前等待的毫秒，等待清理释放资源

	LogConfig logx.LogConfig
	WebConfig WebConfig
	SdxConfig SdxConfig // sdx middleware configs

	runModeInt8 int8 // 运行模式
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

// 闪电侠实现的中间件控制参数
type SdxConfig struct {
	//SysStateMonitor  bool  `v:"def=true"`                  // 是否启动系统资源使用情况的定时检查工作
	PrintSysState    bool  `v:"def=true"`                  // 定时打印系统资源状态检查日志
	PrintRouteState  bool  `v:"def=true"`                  // 定时打印路由访问统计数据
	MaxContentLength int64 `v:"def=33554432"`              // 最大请求字节数，32MB（33554432），传0不限制
	MaxConnections   int32 `v:"def=0,range=[0:100000000]"` // 最大同时请求数，0不限制

	EnableSpecialHandlers bool          `v:"def=true"`   // 是否启用默认的特殊路由中间件
	EnableTrack           bool          `v:"def=false"`  // 启动链路追踪
	EnableGunzip          bool          `v:"def=false"`  // 启动gunzip
	EnableShedding        bool          `v:"def=true"`   // 启动降载限制访问
	EnableTimeout         bool          `v:"def=true"`   // 启动超时拦截
	DefaultTimeout        time.Duration `v:"def=3000ms"` // 默认请求超时时间（单位：毫秒）
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
	return gft.runModeInt8 == modeDebug
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
		gft.runModeInt8 = modeDebug
	case DevMode:
		gft.RunMode = DevMode
		gft.runModeInt8 = modeDev
	case ProductMode:
		gft.RunMode = ProductMode
		gft.runModeInt8 = modeProduct
	case TestMode:
		gft.RunMode = TestMode
		gft.runModeInt8 = modeTest
	default:
		panic("GoFast mode unknown: " + mode)
	}
}
