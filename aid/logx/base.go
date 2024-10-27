// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

// 日志级别的设定，自动显示对应级别的日志
const (
	LevelStack int8 = -1   // [err, warn, info, debug, stack]
	LevelDebug int8 = iota // [err, warn, info, debug]
	LevelInfo              // [err, warn, info]
	LevelWarn              // [err, warn]
	LevelError             // [err]
)

const (
	callerSkipDepth = 4 // 这里的4最好别动，刚好能打印出错误发生的地方。

	// 日志输出媒介
	toConsole = "console"
	toFile    = "file"
	toVolume  = "volume"

	// 用于区分日志 Label
	labelInfo  = "info"  // 0
	labelDebug = "debug" // 1
	labelWarn  = "warn"  // 2
	labelError = "err"   // 4
	labelStack = "stack" // 8
	labelStat  = "stat"  // 32
	labelSlow  = "slow"  // 64
	labelReq   = "req"   // 64
	labelTimer = "timer" // 128	// 定时器执行的任务日志，一般为定时脚本准备
)

type LogConfig struct {
	AppName   string `v:"def=AppName"`                                     // 应用名称
	HostName  string `v:"def=HostName"`                                    // 运行终端编号
	LogMedium string `v:"def=console,enum=console|file|volume"`            // 记录存储媒介
	LogLevel  string `v:"def=info,enum=stack|debug|info|warn|err|pic"`     // 记录日志的级别
	LogStyle  string `v:"def=sdx,enum=custom|sdx|sdx-json|elk|prometheus"` // 日志样式
	LogStat   bool   `v:"def=true"`                                        // 是否记录统计信息

	FileFolder   string `v:""`                    // 日志文件夹路径
	FilePrefix   string `v:""`                    // 日志文件名统一前缀(默认是AppName)
	FileSplit    uint16 `v:"def=0,range=[0:255]"` // 日志拆分(比如32: info+stat; 64: info+timer; 160: info+stat+timer)
	FileKeepDays int    `v:"def=30"`              // 日志文件保留天数
	FileGzip     bool   `v:"def=false"`           // 是否Gzip压缩日志文件
	// FileStackArchiveMillis int  `v:"def=100"`   // 日志文件堆栈毫秒数

	iLevel int8 // 日志级别
	iStyle int8 // 日志样式类型
}
