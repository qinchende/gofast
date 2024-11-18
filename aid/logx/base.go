// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

import (
	"io"
)

type LogConfig struct {
	AppName   string `v:"def=App"`                                     // 应用名称
	HostName  string `v:"def=Host"`                                    // 运行主机名称
	LogLevel  string `v:"def=INF,enum=TRC|DBG|INF|WRN|ERR|disable"`    // 记录日志的级别
	LogStyle  string `v:"def=json,enum=sdx|json|cdo|custom"`           // 日志样式
	LogMedium string `v:"def=console,enum=console|file|volume|custom"` // 记录存储媒介

	// 当LogMedium为 file 或 volume 有效
	FilePath     string `v:"def=_logs_"`            // 文件路径
	FileName     string `v:"def=[AppName]"`         // 文件名称(默认是AppName)
	FileSplit    string `v:"def=no,range=[0:255]"`  // 日志拆分(比如: no; stack|info|req|warn|err)
	FileKeepDays int    `v:"def=30,range=[0:3650]"` // 日志文件保留天数
	FileGzip     bool   `v:"def=false"`             // 是否Gzip压缩日志文件
	// FileStackArchiveMillis int  `v:"def=100"`   // 日志文件堆栈毫秒数

	// 控制1
	DiscardIO    bool `v:"def=false"` // 禁止IO输出
	DisableTimer bool `v:"def=false"` // 是否记录定时器日志
	DisableReq   bool `v:"def=false"` // 是否记录请求日志
	DisableStat  bool `v:"def=false"` // 是否记录统计数据
	DisableSlow  bool `v:"def=false"` // 是否记录慢日志
	// 控制2
	DisableMark  bool   `v:"def=false"` // 是否打印应用标记
	CdoGroupSize uint16 `v:"def=100"`   // Cdo编码时分页大小
}

const (
	// 日志级别的设定，主要是用来控制日志输出的多少
	// Note: 日志级别不需要太多，如果你觉得自己需要，多半都是打印日志的逻辑出问题了
	// 默认下面几个日志级别，足够了。可以自定义扩展，但我不推荐
	LevelTrace   int8 = -8  // 1
	LevelDebug   int8 = -4  // 2
	LevelInfo    int8 = 0   // 3
	LevelWarn    int8 = 4   // 4
	LevelErr     int8 = 8   // 5
	LevelDisable int8 = 127 // 6 丢弃日志

	// 日志分类Label。这和日志级别是不同的概念，Label可以有很多，但是日志级别不要太多
	LabelTrace   = "TRC"     // 1
	LabelDebug   = "DBG"     // 2
	LabelInfo    = "INF"     // 3
	LabelReq     = "REQ"     // 3 请求日志
	LabelTimer   = "TIM"     // 3 定时器执行的任务日志，一般为定时脚本准备
	LabelStat    = "STA"     // 3 运行状态日志
	LabelWarn    = "WRN"     // 4
	LabelSlow    = "SLO"     // 4 慢日志
	LabelErr     = "ERR"     // 5
	LabelPanic   = "PIC"     // 5
	LabelDisable = "disable" // 6
)

//var (
//	labels = [11]string{"trace", "debug", "info", "req", "timer", "stat", "warn", "slow", "err", "panic", "discard"}
//)
//iTrace   int8 = 0
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

const (
	callerSkipDepth = 4 // 这里的4最好别动，刚好能打印出错误发生的地方。

	// 日志输出媒介
	toConsole = "console"
	toFile    = "file"
	toVolume  = "volume"
	toCustom  = "custom"

	fMessage   = "msg"
	fTimeStamp = "ts"
	fLabel     = "label"
	fApp       = "app"
	fHost      = "host"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type (
	// 对象自定义输出方法，实现此接口用来自定义处理敏感信息
	Printer interface {
		Print([]byte) []byte
	}

	RecordWriter interface {
		write()
	}

	LogRecord struct {
		r Record
	}

	Logger struct {
		// *LogConfig
		// 自己也需要集成一个记录器
		// App  string `json:"app"`
		// Host string `json:"host"`

		// 每种分类可以单独输出到不同的介质
		WStack io.Writer
		WDebug io.Writer
		WInfo  io.Writer
		WReq   io.Writer
		WTimer io.Writer
		WStat  io.Writer
		WWarn  io.Writer
		WSlow  io.Writer
		WErr   io.Writer
		WPanic io.Writer
		//WDiscard io.Writer

		// 指定下面的方法即可自定义输出日志样式
		StyleSummary    func(r *Record) []byte
		StyleGroupBegin func([]byte, string) []byte
		StyleGroupEnd   func([]byte) []byte

		LogRecord

		// initOnce sync.Once
		cnf *LogConfig

		iLevel int8 // 日志级别
		iStyle int8 // 日志样式类型
	}
)
