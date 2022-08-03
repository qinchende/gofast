package logx

import (
	"errors"
	"github.com/qinchende/gofast/skill/sysx/host"
	"log"
	"os"
	"path"
	"sync/atomic"
)

var (
	ErrLogNotInitialized    = errors.New("log not initialized")
	ErrLogServiceNameNotSet = errors.New("log service name must be set")
)

// 必须准备好日志环境，否则启动失败自动退出
func MustSetup(cnf *LogConfig) {
	myCnf = cnf
	if len(myCnf.FilePrefix) > 0 {
		myCnf.FilePrefix += "."
	}

	if err := setup(myCnf); err != nil {
		stackInfo := formatWithCaller(err.Error(), 3)
		log.Println(stackInfo)
		output(stackLog, typeStack, stackInfo, false)
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
	}

	switch c.LogStyle {
	case styleSdxStr:
		c.logStyle = StyleSdx
	case styleSdxMiniStr:
		c.logStyle = StyleSdxMini
	case styleJsonMiniStr:
		c.logStyle = StyleJsonMini
	default:
		c.logStyle = StyleJson
	}

	switch c.PrintMedium {
	case printMediumFile:
		return setupWithFiles(c)
	case printMediumVolume:
		return setupWithVolume(c)
	default:
		return setupWithConsole(c)
	}
}

// 控制台日志模式
func setupWithConsole(c *LogConfig) error {
	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)
		writeConsole = true
		//setupLogLevel(c)

		// 一般输出
		infoLog = newLogWriter(log.New(os.Stdout, "", flags))
		debugLog = infoLog
		statLog = infoLog
		slowLog = infoLog
		// 错误输出
		warnLog = newLogWriter(log.New(os.Stderr, "", flags))
		errorLog = warnLog
		stackLog = warnLog
	})
	return nil
}

// 分卷存储文件
func setupWithVolume(c *LogConfig) error {
	if len(c.AppName) == 0 {
		return ErrLogServiceNameNotSet
	}
	c.FilePath = path.Join(c.FilePath, c.AppName, host.Hostname())
	return setupWithFiles(c)
}

// 文件日志模式下的初始化工作
func setupWithFiles(c *LogConfig) error {
	if len(c.FilePath) == 0 {
		return errors.New("log path must be set")
	}
	var err error

	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)

		debugFilePath := logFilePath(typeDebug)
		infoFilePath := logFilePath(typeInfo)
		warnFilePath := logFilePath(typeWarn)
		errorFilePath := logFilePath(typeError)
		stackFilePath := logFilePath(typeStack)
		statFilePath := logFilePath(typeStat)
		slowFilePath := logFilePath(typeSlow)

		// 初始化日志文件, 用 writer-rotate 策略写日志文件
		infoLog = createOutput(infoFilePath)
		if c.FileNumber == fileOne {
			debugLog = infoLog
			warnLog = infoLog
			errorLog = infoLog
			stackLog = infoLog
			statLog = infoLog
			slowLog = infoLog
		} else if c.FileNumber == fileTwo {
			errorLog = createOutput(errorFilePath)
			warnLog = errorLog
			slowLog = errorLog
			statLog = errorLog
			stackLog = errorLog
		} else if c.FileNumber == fileThree {
			errorLog = createOutput(errorFilePath)
			statLog = createOutput(statFilePath)
			warnLog = errorLog
			stackLog = errorLog
			slowLog = statLog
			slowLog = statLog
		} else {
			debugLog = createOutput(debugFilePath)
			warnLog = createOutput(warnFilePath)
			errorLog = createOutput(errorFilePath)
			stackLog = createOutput(stackFilePath)
			statLog = createOutput(statFilePath)
			slowLog = createOutput(slowFilePath)
		}
	})
	return err
}

func logFilePath(logType string) string {
	return path.Join(myCnf.FilePath, myCnf.FilePrefix+logType+".log")
}

func createOutput(path string) WriterCloser {
	rr := DefaultRotateRule(path, backupFileDelimiter, myCnf.FileKeepDays, myCnf.FileGzip)
	wr, err := NewRotateLogger(path, rr, myCnf.FileGzip)
	if err != nil {
		panic(err)
	}
	return wr
}
