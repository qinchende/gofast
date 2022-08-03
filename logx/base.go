// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"sync"
	"time"
)

const (
	LogLevelDebug int8 = iota // LogLevelDebug logs [everything]
	LogLevelInfo              // LogLevelInfo logs [info, warn, error, stack]
	LogLevelWarn              // LogLevelError includes [warn, error, stack]
	LogLevelError             // LogLevelError includes [error, stack]
	LogLevelStack             // LogLevelError includes [stack]
)

const (
	fileAll   int8 = iota // 默认0：不同级别放入不同的日志文件
	fileOne               // 1：全部放在一个日志文件access中
	fileTwo               // 2：只分access和error两个文件
	fileThree             // 3：只分access和error和stat三个文件
)

const (
	// 日志样式类型
	StyleJson int8 = iota
	StyleJsonMini
	StyleSdx
	StyleSdxMini
)

const (
	// 日志样式名称
	styleJsonStr     = "json"
	styleJsonMiniStr = "json-mini"
	styleSdxStr      = "sdx"
	styleSdxMiniStr  = "sdx-mini"

	printMediumConsole = "console"
	printMediumFile    = "file"
	printMediumVolume  = "volume"

	// 多钟不同的log分类
	typeDebug = "debug"
	typeInfo  = "info"
	typeWarn  = "warn"
	typeError = "error"
	typeStack = "stack"
	// 几种统计日志
	typeStat = "stat"
	typeSlow = "slow"

	callerInnerDepth = 5
	flags            = 0x0
)

var (
	debugLog WriterCloser
	infoLog  WriterCloser
	warnLog  WriterCloser
	errorLog WriterCloser
	stackLog WriterCloser
	slowLog  WriterCloser
	statLog  WriterCloser

	writeConsole bool
	once         sync.Once
	initialized  uint32
	//options     logOptions

	myCnf *LogConfig
)

type (
	LogConfig struct {
		AppName     string `v:""`
		PrintMedium string `v:"def=console,enum=console|file|volume"`
		LogLevel    string `v:"def=info,enum=debug|info|warn|error|stack"` // 记录日志的级别
		LogStyle    string `v:"def=sdx,enum=json|json-mini|sdx|sdx-mini"`  // 日志样式
		LogStats    bool   `v:"def=true"`                                  // 是否打印统计信息

		FilePath   string `v:"def=_logs_"`        // 日志文件路径
		FilePrefix string `v:""`                  // 日志文件名统一前缀(默认是AppName)
		FileNumber int8   `v:"def=0,range=[0:3]"` // 日志文件数量

		FileKeepDays           int  `v:"def=0"`     // 日志文件保留天数
		FileStackArchiveMillis int  `v:"def=100"`   // 日志文件堆栈毫秒数
		FileGzip               bool `v:"def=false"` // 是否Gzip压缩日志文件

		logLevel int8 // 日志级别
		logStyle int8 // 日志样式类型
	}

	logEntry struct {
		Timestamp string `json:"@timestamp"`
		Level     string `json:"lv"`
		Duration  string `json:"duration,omitempty"`
		Content   string `json:"ct"`
	}

	//logOptions struct {
	//	gzipEnabled          bool
	//	logStackArchiveMills int
	//	keepDays             int
	//}
	//
	//LogOption func(options *logOptions)

	Logger interface {
		Error(...any)
		ErrorF(string, ...any)
		Info(...any)
		InfoF(string, ...any)
		Slow(...any)
		Slowf(string, ...any)
		WithDuration(time.Duration) Logger
	}
)

//const (
//	green   = "\033[97;42m"
//	white   = "\033[90;47m"
//	yellow  = "\033[90;43m"
//	red     = "\033[97;41m"
//	blue    = "\033[97;44m"
//	magenta = "\033[97;45m"
//	cyan    = "\033[97;46m"
//	Reset   = "\033[0m"
//)
//
//var (
////DefWriter      io.Writer = os.Stdout
////DefErrorWriter io.Writer = os.Stderr
//)

func Style() int8 {
	return myCnf.logStyle
}
