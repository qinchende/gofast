// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

// 日志级别的设定，自动显示对应级别的日志
const (
	LogLevelDebug int8 = iota // LogLevelDebug logs [everything]
	LogLevelInfo              // LogLevelInfo logs [info, warn, error, stack]
	LogLevelWarn              // LogLevelError includes [warn, error, stack]
	LogLevelError             // LogLevelError includes [error, stack]
	LogLevelStack             // LogLevelError includes [stack]
)

// 日志文件拆分成几个分别保存内容
const (
	fileAll   int8 = iota // 默认0：不同级别放入不同的日志文件
	fileOne               // 1：全部放在一个日志文件info中
	fileTwo               // 2：只分info和stat两个文件
	fileThree             // 3：只分info和error和stat三个文件
)

const (
	logMediumConsole = "console"
	logMediumFile    = "file"
	logMediumVolume  = "volume"

	// 多钟不同的log分类
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
	levelStack = "stack"
	// 几种统计日志
	levelStat = "stat"
	levelSlow = "slow"

	callerInnerDepth = 3
)

type LogConfig struct {
	AppName   string `v:""`
	LogMedium string `v:"def=console,enum=console|file|volume"`
	LogLevel  string `v:"def=info,enum=debug|info|warn|error|stack"` // 记录日志的级别
	LogStyle  string `v:"def=sdx,enum=json|json-mini|sdx|sdx-mini"`  // 日志样式
	LogStats  bool   `v:"def=true"`                                  // 是否打印统计信息

	FileFolder string `v:""`                  // 日志文件夹路径
	FilePrefix string `v:""`                  // 日志文件名统一前缀(默认是AppName)
	FileNumber int8   `v:"def=2,range=[0:3]"` // 日志文件数量，默认只拆分成 info and stat (日志 + 统计)

	FileKeepDays int  `v:"def=7"`     // 日志文件保留天数
	FileGzip     bool `v:"def=false"` // 是否Gzip压缩日志文件
	// FileStackArchiveMillis int  `v:"def=100"`   // 日志文件堆栈毫秒数

	logLevel int8 // 日志级别
	logStyle int8 // 日志样式类型
}

//// 日志文件的目标系统
//const (
//	LogTypeConsole    = "console"
//	LogTypeELK        = "elk"
//	LogTypePrometheus = "prometheus"
//)

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

//logOptions struct {
//	gzipEnabled          bool
//	logStackArchiveMills int
//	keepDays             int
//}
//
//LogOption func(options *logOptions)

//Logger interface {
//	Error(...any)
//	ErrorF(string, ...any)
//	Info(...any)
//	InfoF(string, ...any)
//	Slow(...any)
//	Slowf(string, ...any)
//	WithDuration(time.Duration) Logger
//}

//writeConsole bool
//initialized  uint32 // 是否已经完成初始化
