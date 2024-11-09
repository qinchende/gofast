// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

type LogConfig struct {
	AppName   string `v:"def=AppName"`                                     // 应用名称
	HostName  string `v:"def=HostName"`                                    // 主机名称
	LogLevel  string `v:"def=info,enum=stack|debug|info|warn|err|discard"` // 记录日志的级别
	LogStyle  string `v:"def=sdx,enum=sdx|json|cdo|custom"`                // 日志样式
	LogMedium string `v:"def=console,enum=console|file|volume|custom"`     // 记录存储媒介

	// 当LogMedium为 file 或 volume 有效
	FilePath     string `v:"def=_logs_"`            // 文件路径
	FileName     string `v:"def=[AppName]"`         // 文件名称(默认是AppName)
	FileSplit    string `v:"def=no,range=[0:255]"`  // 日志拆分(比如: no; stack|info|req|warn|err)
	FileKeepDays int    `v:"def=30,range=[0:3650]"` // 日志文件保留天数
	FileGzip     bool   `v:"def=false"`             // 是否Gzip压缩日志文件
	// FileStackArchiveMillis int  `v:"def=100"`   // 日志文件堆栈毫秒数

	// 控制1
	DisableTimer bool `v:"def=false"` // 是否记录定时器日志
	DisableReq   bool `v:"def=false"` // 是否记录请求日志
	DisableStat  bool `v:"def=false"` // 是否记录统计数据
	DisableSlow  bool `v:"def=false"` // 是否记录慢日志
	// 控制2
	EnableMini   bool   `v:"def=false"`
	DisableMark  bool   `v:"def=false"` // 是否打印应用标记
	CdoGroupSize uint16 `v:"def=100"`   // Cdo编码时分页大小

	iLevel int8 // 日志级别
	iStyle int8 // 日志样式类型
}

const (
	// 日志级别的设定，自动输出对应级别的日志。主要是用来控制日志输出的多少
	// Note: 日志级别不需要太多，如果你觉得自己需要，多半都是打印日志的逻辑出问题了
	// 默认下面几个日志级别，足够了。当然你是可以扩展的
	LevelStack   int8 = -8  // 1
	LevelDebug   int8 = -4  // 2
	LevelInfo    int8 = 0   // 3
	LevelWarn    int8 = 4   // 4
	LevelErr     int8 = 8   // 5
	LevelDiscard int8 = 127 // 6 禁用日志

	// 用于区分日志分类的 Label。这和日志级别是不同的概念
	labelStack   = "stack"   // 1
	labelDebug   = "debug"   // 2
	labelInfo    = "info"    // 3
	labelTimer   = "timer"   // 3 定时器执行的任务日志，一般为定时脚本准备
	labelReq     = "req"     // 3 请求日志
	labelStat    = "stat"    // 3 运行状态日志
	labelWarn    = "warn"    // 4
	labelSlow    = "slow"    // 4 慢日志
	labelErr     = "err"     // 5
	labelPanic   = "panic"   // 5
	labelDiscard = "discard" // 6

	//iStack   int8 = 0
	//iDebug   int8 = 1
	//iInfo    int8 = 2
	//iReq     int8 = 3
	//iTimer   int8 = 4
	//iStat    int8 = 5
	//iWarn    int8 = 6
	//iSlow    int8 = 7
	//iErr     int8 = 8
	//iPanic   int8 = 9
	//iDiscard int8 = 10
)

var (
// labels = [11]string{"stack", "debug", "info", "req", "timer", "stat", "warn", "slow", "err", "panic", "discard"}
)

const (
	callerSkipDepth = 4 // 这里的4最好别动，刚好能打印出错误发生的地方。

	// 日志输出媒介
	toConsole = "console"
	toFile    = "file"
	toVolume  = "volume"
	toCustom  = "custom"
)
