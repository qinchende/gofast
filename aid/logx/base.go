// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/aid/bag"
	"github.com/qinchende/gofast/core/cst"
	"net/http"
	"time"
)

type LogConfig struct {
	AppName      string `v:"def=AppName"`                                     // 应用名称
	HostName     string `v:"def=HostName"`                                    // 主机名称
	LogMedium    string `v:"def=console,enum=console|file|volume"`            // 记录存储媒介
	LogLevel     string `v:"def=info,enum=discard|debug|info|warn|err|panic"` // 记录日志的级别
	LogStyle     string `v:"def=sdx,enum=sdx|json|cdo|custom"`                // 日志样式
	LogStat      bool   `v:"def=true"`                                        // 是否记录统计信息
	FilePath     string `v:"def=_logs_"`                                      // 文件路径
	FileName     string `v:"def=[AppName]"`                                   // 文件名称(默认是AppName)
	FileSplit    string `v:"def=info,range=[0:255]"`                          // 日志拆分(比如: stack; info|warn; info|timer)
	FileKeepDays int    `v:"def=30,range=[0:3650]"`                           // 日志文件保留天数
	FileGzip     bool   `v:"def=false"`                                       // 是否Gzip压缩日志文件
	// FileStackArchiveMillis int  `v:"def=100"`   // 日志文件堆栈毫秒数

	iLevel int8 // 日志级别
	iStyle int8 // 日志样式类型
}

// 日志级别的设定，自动输出对应级别的日志
const (
	LevelDiscard int8 = -8
	LevelStack   int8 = -1
	LevelDebug   int8 = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic

	// 用于区分日志分类的 Label
	labelStack = "stack" //
	labelDebug = "debug" //
	labelInfo  = "info"  //
	labelReq   = "req"   // 请求日志
	labelTimer = "timer" // 定时器执行的任务日志，一般为定时脚本准备
	labelStat  = "stat"  //
	labelWarn  = "warn"  //
	labelSlow  = "slow"  //
	labelErr   = "err"   //
	labelPanic = "panic" //
)

const (
	callerSkipDepth = 4 // 这里的4最好别动，刚好能打印出错误发生的地方。

	// 日志输出媒介
	toConsole = "console"
	toFile    = "file"
	toVolume  = "volume"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type Field struct {
	Key string
	Val any
}

//type Entry struct {
//	Level      Level
//	Time       time.Time
//	LoggerName string
//	Message    string
//	Caller     EntryCaller
//	Stack      string
//}

// 日志参数实体
type ReqLogEntity struct {
	RawReq     *http.Request
	TimeStamp  time.Duration
	Latency    time.Duration
	ClientIP   string
	StatusCode int
	Pms        cst.SuperKV
	BodySize   int
	ResData    []byte
	CarryItems bag.CarryList
}
