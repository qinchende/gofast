// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"errors"
	"sync"
	"time"
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

const (
	InfoLevel   = iota // InfoLevel logs everything
	ErrorLevel         // ErrorLevel includes errors, slows, stacks
	SevereLevel        // SevereLevel only log severe messages
)

const (
	fileAll   int8 = iota // 默认0：不同级别放入不同的日志文件
	fileOne               // 1：全部放在一个日志文件access中
	fileTwo               // 2：只分access和error两个文件
	fileThree             // 3：只分access和error和stat三个文件
)

// 日志样式类型
const (
	StyleJson int8 = iota
	StyleJsonMini
	StyleSdx
	StyleSdxMini
)

// 日志样式名称
const (
	styleJsonStr     = "json"
	styleJsonMiniStr = "json-mini"
	styleSdxStr      = "sdx"
	styleSdxMiniStr  = "sdx-mini"
)

const (
	timeFormat     = "2006-01-02T15:04:05.000Z07"
	timeFormatMini = "01-02 15:04:05"

	accessFilename = "access.log"
	errorFilename  = "error.log"
	warnFilename   = "warn.log"
	severeFilename = "severe.log"
	slowFilename   = "slow.log"
	statFilename   = "stat.log"
	stackFilename  = "stack.log"

	consoleMode = "console"
	volumeMode  = "volume"

	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelWarn   = "warn"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"

	callerInnerDepth = 5
	flags            = 0x0
)

var (
	ErrLogPathNotSet        = errors.New("log path must be set")
	ErrLogNotInitialized    = errors.New("log not initialized")
	ErrLogServiceNameNotSet = errors.New("log service name must be set")

	writeConsole bool
	logLevel     uint32
	// 6个不同等级的日志输出
	accessLog WriterCloser
	errorLog  WriterCloser
	warnLog   WriterCloser
	severeLog WriterCloser
	slowLog   WriterCloser
	statLog   WriterCloser
	stackLog  WriterCloser

	once        sync.Once
	initialized uint32
	options     logOptions
)

type (
	LogConfig struct {
		AppName            string `v:""`
		PrintMedium        string `v:"def=console,enum=console|file|volume"`
		LogLevel           string `v:"def=info,enum=info|error|severe"` // 记录日志的级别
		FilePath           string `v:"def=_logs_"`                      // 日志文件路径
		FilePrefix         string `v:""`                                // 日志文件名统一前缀(默认是AppName)
		FileNumber         int8   `v:"def=0,range=[0:3]"`               // 日志文件数量
		Compress           bool   `v:"def=false"`
		KeepDays           int    `v:"def=0"`
		StackArchiveMillis int    `v:"def=100"`
		StyleName          string `v:"def=sdx,enum=json|json-mini|sdx|sdx-mini"` // 日志样式
		styleType          int8   `v:""`                                         // 日志样式类型
	}

	logEntry struct {
		Timestamp string `json:"@timestamp"`
		Level     string `json:"lv"`
		Duration  string `json:"duration,omitempty"`
		Content   string `json:"ct"`
	}

	logOptions struct {
		gzipEnabled          bool
		logStackArchiveMills int
		keepDays             int
	}

	LogOption func(options *logOptions)

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
