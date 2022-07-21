package logx

import (
	"github.com/qinchende/gofast/skill/sysx/host"
	"log"
	"os"
	"path"
	"sync/atomic"
)

var currConfig *LogConfig

func Style() int8 {
	return currConfig.style
}

func MustSetup(c *LogConfig) {
	currConfig = c

	// 必须准备好日志环境，否则启动失败自动退出
	err := setup(currConfig)
	//DefaultWriter = infoLog
	//DefErrorWriter = errorLog
	if err != nil {
		msg := formatWithCaller(err.Error(), 3)
		log.Println(msg)
		output(severeLog, levelFatal, msg, false)
		os.Exit(1)
	}
}

func WithArchiveMillis(millis int) LogOption {
	return func(opts *logOptions) {
		opts.logStackArchiveMills = millis
	}
}

func WithKeepDays(days int) LogOption {
	return func(opts *logOptions) {
		opts.keepDays = days
	}
}

func WithGzip() LogOption {
	return func(opts *logOptions) {
		opts.gzipEnabled = true
	}
}

func SetLevel(level uint32) {
	atomic.StoreUint32(&logLevel, level)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// SetUp sets up the logx. If already set up, just return nil.
// we allow SetUp to be called multiple times, because for example
// we need to allow different service frameworks to initialize logx respectively.
// the same logic for SetUp
func setup(c *LogConfig) error {
	switch c.StyleName {
	case styleSdxStr:
		c.style = StyleSdx
	case styleSdxMiniStr:
		c.style = StyleSdxMini
	case styleJsonMiniStr:
		c.style = StyleJsonMini
	default:
		c.style = StyleJson
	}

	switch c.Mode {
	case consoleMode:
		setupWithConsole(c)
		return nil
	case volumeMode:
		return setupWithVolume(c)
	default:
		return setupWithFiles(c)
	}
}

// 控制台日志模式
func setupWithConsole(c *LogConfig) {
	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)
		writeConsole = true
		setupLogLevel(c)

		// 一般输出
		accessLog = newLogWriter(log.New(os.Stdout, "", flags))
		statLog = accessLog
		warnLog = accessLog
		// 错误输出
		errorLog = newLogWriter(log.New(os.Stderr, "", flags))
		severeLog = newLogWriter(log.New(os.Stderr, "", flags))
		slowLog = newLogWriter(log.New(os.Stderr, "", flags))
		stackLog = newLogWriter(log.New(os.Stderr, "", flags))
		//stackLog = NewLessWriter(errorLog, options.logStackArchiveMills)
	})
}

// 文件日志模式下的初始化工作
func setupWithFiles(c *LogConfig) error {
	if len(c.Path) == 0 {
		return ErrLogPathNotSet
	}

	var opts []LogOption
	var err error

	// 添加 config 参数处理函数
	opts = append(opts, WithArchiveMillis(c.StackArchiveMillis))
	if c.Compress {
		opts = append(opts, WithGzip())
	}
	if c.KeepDays > 0 {
		opts = append(opts, WithKeepDays(c.KeepDays))
	}

	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)
		handleOptions(opts)
		setupLogLevel(c)

		prefix := c.FilePrefix
		if len(c.FilePrefix) > 0 {
			prefix += "."
		}

		accessFilePath := path.Join(c.Path, prefix+accessFilename)
		errorFilePath := path.Join(c.Path, prefix+errorFilename)
		warnFilePath := path.Join(c.Path, prefix+warnFilename)
		severeFilePath := path.Join(c.Path, prefix+severeFilename)
		slowFilePath := path.Join(c.Path, prefix+slowFilename)
		statFilePath := path.Join(c.Path, prefix+statFilename)
		stackFilePath := path.Join(c.Path, prefix+stackFilename)

		// 初始化日志文件, 用 writer-rotate 策略写日志文件
		if accessLog, err = createOutput(accessFilePath); err != nil {
			return
		}
		if c.FileNumber == fileOne {
			errorLog = accessLog
			warnLog = accessLog
			severeLog = accessLog
			slowLog = accessLog
			statLog = accessLog
			stackLog = accessLog
		} else if c.FileNumber == fileTwo {
			if errorLog, err = createOutput(errorFilePath); err != nil {
				return
			}
			warnLog = errorLog
			severeLog = errorLog
			slowLog = errorLog
			statLog = errorLog
			stackLog = errorLog
		} else if c.FileNumber == fileThree {
			if errorLog, err = createOutput(errorFilePath); err != nil {
				return
			}
			if statLog, err = createOutput(statFilePath); err != nil {
				return
			}
			warnLog = errorLog
			severeLog = errorLog
			slowLog = errorLog
			stackLog = errorLog
		} else {
			if warnLog, err = createOutput(warnFilePath); err != nil {
				return
			}
			if errorLog, err = createOutput(errorFilePath); err != nil {
				return
			}
			if severeLog, err = createOutput(severeFilePath); err != nil {
				return
			}
			if slowLog, err = createOutput(slowFilePath); err != nil {
				return
			}
			if statLog, err = createOutput(statFilePath); err != nil {
				return
			}
			if stackLog, err = createOutput(stackFilePath); err != nil {
				return
			}
			//stackLog = NewLessWriter(errorLog, options.logStackArchiveMills)
		}
	})
	return err
}

// 日志存储
func setupWithVolume(c *LogConfig) error {
	if len(c.ServiceName) == 0 {
		return ErrLogServiceNameNotSet
	}
	c.Path = path.Join(c.Path, c.ServiceName, host.Hostname())
	return setupWithFiles(c)
}

func createOutput(path string) (WriterCloser, error) {
	if len(path) == 0 {
		return nil, ErrLogPathNotSet
	}
	rr := DefaultRotateRule(path, backupFileDelimiter, options.keepDays, options.gzipEnabled)
	return NewRotateLogger(path, rr, options.gzipEnabled)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func setupLogLevel(c *LogConfig) {
	switch c.Level {
	case levelInfo:
		SetLevel(InfoLevel)
	case levelError:
		SetLevel(ErrorLevel)
	case levelSevere:
		SetLevel(SevereLevel)
	}
}

func shouldLog(level uint32) bool {
	return atomic.LoadUint32(&logLevel) <= level
}

func handleOptions(opts []LogOption) {
	for _, opt := range opts {
		opt(&options)
	}
}
