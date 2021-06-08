package logx

import (
	"github.com/qinchende/gofast/skill/sysx"
	"io"
	"log"
	"os"
	"path"
	"sync/atomic"
)

// 自定义 logger
type logWriter struct {
	logger *log.Logger
}

func newLogWriter(logger *log.Logger) logWriter {
	return logWriter{
		logger: logger,
	}
}

func (lw logWriter) Close() error {
	return nil
}

func (lw logWriter) Write(data []byte) (int, error) {
	lw.logger.Print(string(data))
	return len(data), nil
}

func (lw logWriter) WriteString(data string) (int, error) {
	lw.logger.Print(data)
	return len(data), nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func MustSetup(c LogConfig) {
	currConfig = &c

	// 必须准备好日志环境，否则启动失败自动退出
	err := setup(currConfig)
	DefaultWriter = infoLog
	DefaultErrorWriter = errorLog
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
	case _styleSdx:
		c.style = styleSdx
	case _styleSdxMini:
		c.style = styleSdxMini
	default:
		c.style = styleJson
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
		stackLog = NewLessWriter(errorLog, options.logStackArchiveMills)
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

		accessFile := path.Join(c.Path, accessFilename)
		errorFile := path.Join(c.Path, errorFilename)
		severeFile := path.Join(c.Path, severeFilename)
		slowFile := path.Join(c.Path, slowFilename)
		statFile := path.Join(c.Path, statFilename)

		if infoLog, err = createOutput(accessFile); err != nil {
			return
		}
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
		stackLog = NewLessWriter(errorLog, options.logStackArchiveMills)
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

func createOutput(path string) (io.WriteCloser, error) {
	if len(path) == 0 {
		return nil, ErrLogPathNotSet
	}

	return NewLogger(path, DefaultRotateRule(path, backupFileDelimiter, options.keepDays,
		options.gzipEnabled), options.gzipEnabled)
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
