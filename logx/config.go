// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

var currConfig *LogConfig

type LogConfig struct {
	ServiceName        string `v:""`
	Mode               string `v:"def=console,enum=console|file|volume"`
	Level              string `v:"def=info,enum=info|error|severe"` // 记录日志的级别
	Path               string `v:"def=logs"`                        // 日志文件路径
	FilePrefix         string `v:""`                                // 日志文件名统一前缀
	FileNumber         int8   `v:"def=0,range=[0:3]"`               // 日志文件数量
	Compress           bool   `v:""`
	KeepDays           int    `v:""`
	StackArchiveMillis int    `v:"def=100"`
	StyleName          string `v:"def=sdx,enum=json|json-mini|sdx|sdx-mini"`
	style              int8   // 日志模板样式
}

const (
	fileAll   int8 = iota // 默认0：不同级别放入不同的日志文件
	fileOne               // 1：全部放在一个日志文件access中
	fileTwo               // 2：只分access和error两个文件
	fileThree             // 3：只分access和error和stat三个文件
)

// 日志样式名称
const (
	styleJsonStr     = "json"
	styleJsonMiniStr = "json-mini"
	styleSdxStr      = "sdx"
	styleSdxMiniStr  = "sdx-mini"
)

// 日志样式类型
const (
	StyleJson int8 = iota
	StyleJsonMini
	StyleSdx
	StyleSdxMini
)

func Style() int8 {
	return currConfig.style
}
