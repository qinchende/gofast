// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"io"
	"time"
)

type LogConfig struct {
	AppName   string `v:"def=App"`                                         // 应用名称
	HostName  string `v:"def=Host"`                                        // 运行主机名称
	LogLevel  string `v:"def=info,enum=trace|debug|info|warn|err|discard"` // 记录日志的级别
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
	LevelDiscard int8 = 127 // 6 禁用日志

	// 用于区分日志分类的 Label。这和日志级别是不同的概念
	LabelTrace   = "trace"   // 1
	LabelDebug   = "debug"   // 2
	LabelInfo    = "info"    // 3
	LabelReq     = "req"     // 3 请求日志
	LabelTimer   = "timer"   // 3 定时器执行的任务日志，一般为定时脚本准备
	LabelStat    = "stat"    // 3 运行状态日志
	LabelWarn    = "warn"    // 4
	LabelSlow    = "slow"    // 4 慢日志
	LabelErr     = "err"     // 5
	LabelPanic   = "panic"   // 5
	LabelDiscard = "discard" // 6
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
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type (
	LogBuilder interface {
		output(msg string)
	}

	Logger struct {
		//*LogConfig

		// 自己也需要集成一个记录器
		Record
		//App  string `json:"app"`
		//Host string `json:"host"`

		// 每种分类可以单独输出到不同的介质
		WStack io.WriteCloser
		WDebug io.WriteCloser
		WInfo  io.WriteCloser
		WReq   io.WriteCloser
		WTimer io.WriteCloser
		WStat  io.WriteCloser
		WWarn  io.WriteCloser
		WSlow  io.WriteCloser
		WErr   io.WriteCloser
		WPanic io.WriteCloser
		//WDiscard io.WriteCloser

		StyleFunc func(*Logger, []byte) []byte

		//initOnce sync.Once
		cnf *LogConfig

		iLevel int8 // 日志级别
		iStyle int8 // 日志样式类型
	}

	Field struct {
		Key string
		Val any
	}

	Record struct {
		Time  time.Duration `json:"ts"`
		Label string        `json:"lb"`
		//Msg   string

		log *Logger
		iow io.WriteCloser
		out LogBuilder
		bf  *[]byte
		bs  []byte // 用来辅助上面的bf指针，防止24个字节的切片对象堆分配
		//fls []Field // 用来记录key-value
	}
)
