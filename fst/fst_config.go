// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
)

// GoFast WEB框架的配置参数
type GfConfig struct {
	AppName          string `v:"required"`                            // 应用名称
	ServerName       string `v:"required"`                            // 运行终端编号
	ListenAddr       string `v:"def=0.0.0.0:8099,match=ipv4:port"`    // 监听ip:port
	RunMode          string `v:"def=product,enum=debug|test|product"` // 当前模式[debug|test|product]
	BeforeShutdownMS int64  `v:"def=1000"`                            // 退出APP之前等待的毫秒，等待清理释放资源

	LogConfig logx.LogConfig
	WebConfig cst.WebConfig
	SdxConfig cst.SdxConfig // sdx middleware configs

	runModeInt8 int8 // 运行模式
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
