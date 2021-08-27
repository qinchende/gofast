package logx

import (
	"github.com/qinchende/gofast/skill/sysx"
	"log"
	"os"
	"path"
	"sync/atomic"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func MustSetup(c LogConfig) {
	currConfig = &c

	// 必须准备好日志环境，否则启动失败自动退出
	err := setup(currConfig)
	//DefaultWriter = infoLog
	DefErrorWriter = errorLog
	if err != nil {
		msg := formatWithCaller(err.Error(), 3)
		log.Println(msg)
		output(severeLog, levelFatal, msg)
		os.Exit(1)
	}
}

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
		infoLog = newLogWriter(log.New(os.Stdout, "", flags))
		statLog = infoLog
		// 错误输出
		errorLog = newLogWriter(log.New(os.Stderr, "", flags))
		severeLog = newLogWriter(log.New(os.Stderr, "", flags))
		slowLog = newLogWriter(log.New(os.Stderr, "", flags))
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

		accessFile := path.Join(c.Path, prefix+accessFilename)
		errorFile := path.Join(c.Path, prefix+errorFilename)
		severeFile := path.Join(c.Path, prefix+severeFilename)
		slowFile := path.Join(c.Path, prefix+slowFilename)
		statFile := path.Join(c.Path, prefix+statFilename)

		// 初始化日志文件, 用 writer-rotate 策略写日志文件
		if infoLog, err = createOutput(accessFile); err != nil {
			return
		}
		if c.FileNumber == fileOne {
			severeLog = infoLog
			slowLog = infoLog
			statLog = infoLog
			//stackLog = infoLog
		} else if c.FileNumber == fileTwo {
			if errorLog, err = createOutput(errorFile); err != nil {
				return
			}
			severeLog = errorLog
			slowLog = errorLog
			statLog = errorLog
			//stackLog = errorLog
		} else {
			if errorLog, err = createOutput(errorFile); err != nil {
				return
			}
			if severeLog, err = createOutput(severeFile); err != nil {
				return
			}
			if slowLog, err = createOutput(slowFile); err != nil {
				return
			}
			if statLog, err = createOutput(statFile); err != nil {
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
	c.Path = path.Join(c.Path, c.ServiceName, sysx.Hostname())
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

func SetLevel(level uint32) {
	atomic.StoreUint32(&logLevel, level)
}

func shouldLog(level uint32) bool {
	return atomic.LoadUint32(&logLevel) <= level
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

func handleOptions(opts []LogOption) {
	for _, opt := range opts {
		opt(&options)
	}
}
