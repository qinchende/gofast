// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"errors"
	"github.com/qinchende/gofast/skill/sysx/host"
	"log"
	"os"
	"path"
	"sync"
)

var (
	debugLog WriterCloser
	infoLog  WriterCloser
	warnLog  WriterCloser
	errorLog WriterCloser
	stackLog WriterCloser
	slowLog  WriterCloser
	statLog  WriterCloser

	initOnce sync.Once
	myCnf    *LogConfig
)

// 必须准备好日志环境，否则启动失败自动退出
func MustSetup(cnf *LogConfig) {
	myCnf = cnf
	if len(myCnf.FilePrefix) > 0 {
		myCnf.FilePrefix += "."
	} else if len(myCnf.AppName) > 0 {
		myCnf.FilePrefix = myCnf.AppName + "."
	}

	if err := setup(myCnf); err != nil {
		info := formatWithCaller(err.Error(), callerInnerDepth)
		log.Println(info)
		output(stackLog, info, levelStack, true)
		os.Exit(1)
	}
}

func setup(c *LogConfig) error {
	switch c.LogLevel {
	case "debug":
		c.logLevel = LogLevelDebug
	case "info":
		c.logLevel = LogLevelInfo
	case "warn":
		c.logLevel = LogLevelWarn
	case "error":
		c.logLevel = LogLevelError
	case "stack":
		c.logLevel = LogLevelStack
	default:
		return errors.New("item LogLevel not match")
	}

	switch c.LogStyle {
	case styleSdxStr:
		c.logStyle = LogStyleSdx
	case styleSdxMiniStr:
		c.logStyle = LogStyleSdxMini
	case styleJsonMiniStr:
		c.logStyle = LogStyleJsonMini
	case styleJsonStr:
		c.logStyle = LogStyleJson
	default:
		return errors.New("item LogStyle not match")
	}

	switch c.LogMedium {
	case logMediumConsole:
		return setupWithConsole(c)
	case logMediumFile:
		return setupWithFiles(c)
	case logMediumVolume:
		return setupWithVolume(c)
	default:
		return errors.New("item LogMedium not match")
	}
}

// 第一种：打印在console
func setupWithConsole(c *LogConfig) error {
	initOnce.Do(func() {
		infoLog = newLogWriter(log.New(os.Stdout, "", 0))
		debugLog = infoLog
		statLog = infoLog
		slowLog = infoLog
		warnLog = newLogWriter(log.New(os.Stderr, "", 0))
		errorLog = warnLog
		stackLog = warnLog
	})
	return nil
}

// 第二种：文件日志模式下的初始化工作
func setupWithFiles(c *LogConfig) error {
	if len(c.FileFolder) == 0 {
		return errors.New("log file folder must be set")
	}
	initOnce.Do(func() {
		debugFilePath := logFilePath(levelDebug)
		infoFilePath := logFilePath(levelInfo)
		warnFilePath := logFilePath(levelWarn)
		errorFilePath := logFilePath(levelError)
		stackFilePath := logFilePath(levelStack)
		statFilePath := logFilePath(levelStat)
		slowFilePath := logFilePath(levelSlow)

		// 初始化日志文件, 用 writer-rotate 策略写日志文件
		infoLog = createFileWriter(infoFilePath)
		// os.Stderr + os.Stdout + os.Stdin (将标准输出重定向到文件中)
		*os.Stdout = *infoLog.(*RotateLogger).fp
		*os.Stderr = *os.Stdout
		//log.SetOutput(infoLog) // 这里不用写了，系统自带的Logger系统默认用的就是 os.stdout 和 os.stderr

		if c.FileNumber == fileOne { // all in info file
			debugLog = infoLog
			warnLog = infoLog
			errorLog = infoLog
			stackLog = infoLog
			statLog = infoLog
			slowLog = infoLog
		} else if c.FileNumber == fileTwo { // split info and stat files
			debugLog = infoLog
			warnLog = infoLog
			errorLog = infoLog
			stackLog = infoLog
			statLog = createFileWriter(statFilePath)
			slowLog = statLog
		} else if c.FileNumber == fileThree { // split info error stat files
			debugLog = infoLog
			errorLog = createFileWriter(errorFilePath)
			warnLog = errorLog
			stackLog = errorLog
			statLog = createFileWriter(statFilePath)
			slowLog = statLog
		} else if c.FileNumber == fileAll { // split every files
			debugLog = createFileWriter(debugFilePath)
			warnLog = createFileWriter(warnFilePath)
			errorLog = createFileWriter(errorFilePath)
			stackLog = createFileWriter(stackFilePath)
			statLog = createFileWriter(statFilePath)
			slowLog = createFileWriter(slowFilePath)
		}
	})

	return nil
}

func logFilePath(logType string) string {
	return path.Join(myCnf.FileFolder, myCnf.FilePrefix+logType+".log")
}

func createFileWriter(path string) WriterCloser {
	rr := DefDailyRotateRule(path, backupFileDelimiter, myCnf.FileKeepDays, myCnf.FileGzip)
	wr, err := NewRotateLogger(path, rr, myCnf.FileGzip)
	if err != nil {
		panic(err)
	}
	return wr
}

// 第三种：分卷存储文件（其实也是写文件，但是更严格的分层文件夹。）
func setupWithVolume(c *LogConfig) error {
	if len(c.AppName) == 0 {
		return errors.New("log config item [AppName] must be set")
	}
	c.FileFolder = path.Join(c.FileFolder, c.AppName, host.Hostname())
	return setupWithFiles(c)
}
