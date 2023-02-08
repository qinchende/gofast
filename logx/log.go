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
	infoLog  WriterCloser
	debugLog WriterCloser
	warnLog  WriterCloser
	errorLog WriterCloser
	stackLog WriterCloser
	slowLog  WriterCloser
	statLog  WriterCloser
	timerLog WriterCloser

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
		output(stackLog, info, typeStack, true)
		os.Exit(1)
	}
}

func setup(c *LogConfig) error {
	switch c.LogLevel {
	case "debug":
		c.logLevelInt8 = LogLevelDebug
	case "info":
		c.logLevelInt8 = LogLevelInfo
	case "warn":
		c.logLevelInt8 = LogLevelWarn
	case "error":
		c.logLevelInt8 = LogLevelError
	case "stack":
		c.logLevelInt8 = LogLevelStack
	default:
		return errors.New("item LogLevel not match")
	}

	if err := initStyle(c); err != nil {
		return err
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
		timerLog = infoLog
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
		// 初始化日志文件, 用 writer-rotate 策略写日志文件
		infoLog = createFile(typeInfo)
		// os.Stderr + os.Stdout + os.Stdin (将标准输出重定向到文件中)
		*os.Stdout = *infoLog.(*RotateLogger).fp
		*os.Stderr = *os.Stdout
		//log.SetOutput(infoLog) // 这里不用写了，系统自带的Logger系统默认用的就是 os.stdout 和 os.stderr

		if c.FileSplit&1 != 0 {
			debugLog = createFile(typeDebug)
		} else {
			debugLog = infoLog
		}
		if c.FileSplit&2 != 0 {
			warnLog = createFile(typeWarn)
		} else {
			warnLog = infoLog
		}
		if c.FileSplit&4 != 0 {
			errorLog = createFile(typeError)
		} else {
			errorLog = warnLog
		}
		if c.FileSplit&8 != 0 {
			stackLog = createFile(typeStack)
		} else {
			stackLog = errorLog
		}
		if c.FileSplit&32 != 0 {
			statLog = createFile(typeStat)
		} else {
			statLog = stackLog
		}
		if c.FileSplit&64 != 0 {
			slowLog = createFile(typeSlow)
		} else {
			slowLog = statLog
		}
		if c.FileSplit&128 != 0 {
			timerLog = createFile(typeTimer)
		} else {
			timerLog = slowLog
		}
	})

	return nil
}

func logFilePath(logType string) string {
	return path.Join(myCnf.FileFolder, myCnf.FilePrefix+logType+".log")
}

func createWriterFile(path string) WriterCloser {
	rr := DefDailyRotateRule(path, backupFileDelimiter, myCnf.FileKeepDays, myCnf.FileGzip)
	wr, err := NewRotateLogger(path, rr, myCnf.FileGzip)
	if err != nil {
		panic(err)
	}
	return wr
}

func createFile(logType string) WriterCloser {
	return createWriterFile(logFilePath(logType))
}

// 第三种：分卷存储文件（其实也是写文件，但是更严格的分层文件夹。）
func setupWithVolume(c *LogConfig) error {
	if len(c.AppName) == 0 {
		return errors.New("log config item [AppName] must be set")
	}
	c.FileFolder = path.Join(c.FileFolder, c.AppName, host.Hostname())
	return setupWithFiles(c)
}
