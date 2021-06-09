// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

var currConfig *LogConfig

type LogConfig struct {
	ServiceName        string `json:",optional"`
	Mode               string `json:",default=console,options=console|file|volume"`
	Level              string `json:",default=info,options=info|error|severe"` // 记录日志的级别
	Path               string `json:",default=logs"`                           // 日志文件路径
	FilePrefix         string `json:",optional"`                               // 日志文件名统一前缀
	FileNumber         int8   `json:",default=0,range=[0:2]"`                  // 日志文件数量
	Compress           bool   `json:",optional"`
	KeepDays           int    `json:",optional"`
	StackArchiveMillis int    `json:",default=100"`
	NeedCpuMem         bool   `json:",default=true"`
	StyleName          string `json:",default=json"`
	style              int8   `inner:",optional"` // 日志模板样式
}

const (
	fileAll int8 = iota // 默认0：不同级别放入不同的日志文件
	fileOne             // 1：全部放在一个日志文件access中
	fileTwo             // 2：只分access和error两个文件
)

const (
	styleJsonStr     = "json"
	styleJsonMiniStr = "json-mini"
	styleSdxStr      = "sdx"
	styleSdxMiniStr  = "sdx-mini"
)

const (
	styleJson int8 = iota
	styleJsonMini
	styleSdx
	styleSdxMini
)
