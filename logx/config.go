// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

var currConfig *LogConfig

type LogConfig struct {
	ServiceName        string `cnf:",NA"`
	Mode               string `cnf:",def=console,enum=console|file|volume"`
	Level              string `cnf:",def=info,enum=info|error|severe"` // 记录日志的级别
	Path               string `cnf:",def=logs"`                           // 日志文件路径
	FilePrefix         string `cnf:",NA"`                               // 日志文件名统一前缀
	FileNumber         int8   `cnf:",def=0,range=[0:2]"`                  // 日志文件数量
	Compress           bool   `cnf:",NA"`
	KeepDays           int    `cnf:",NA"`
	StackArchiveMillis int    `cnf:",def=100"`
	NeedCpuMem         bool   `cnf:",def=true"`
	StyleName          string `cnf:",def=sdx,enum=json|json-mini|sdx|sdx-mini"`
	style              int8   `inner:",NA"` // 日志模板样式
}

const (
	fileAll int8 = iota // 默认0：不同级别放入不同的日志文件
	fileOne             // 1：全部放在一个日志文件access中
	fileTwo             // 2：只分access和error两个文件
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
