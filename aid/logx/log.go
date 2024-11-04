// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"errors"
	"github.com/qinchende/gofast/aid/sysx/host"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

// 每种分类可以有单独输出到不同的日志文件
var (
	ioStack WriterCloser
	ioDebug WriterCloser
	ioInfo  WriterCloser
	ioReq   WriterCloser
	ioTimer WriterCloser
	ioStat  WriterCloser
	ioWarn  WriterCloser
	ioSlow  WriterCloser
	ioErr   WriterCloser
	ioPanic WriterCloser

	initOnce sync.Once
	myCnf    *LogConfig
)

// 必须准备好日志环境，否则启动失败自动退出
func MustSetup(cnf *LogConfig) {
	if err := Setup(cnf); err != nil {
		data := formatWithCaller(err.Error(), callerSkipDepth)
		if ioErr != nil {
			output(ioErr, labelErr, data)
		} else {
			log.Println(data)
		}
		os.Exit(1)
	}
}

func Setup(cnf *LogConfig) error {
	myCnf = cnf

	if len(myCnf.FileName) > 0 {
		myCnf.FileName += "."
	} else if len(myCnf.AppName) > 0 {
		myCnf.FileName = myCnf.AppName + "."
	}

	return initLogger(myCnf)
}

func initLogger(c *LogConfig) error {
	switch c.LogLevel {
	case labelStack:
		c.iLevel = LevelStack
	case labelDebug:
		c.iLevel = LevelDebug
	case labelInfo:
		c.iLevel = LevelInfo
	case labelWarn:
		c.iLevel = LevelWarn
	case labelErr:
		c.iLevel = LevelErr
	case labelDiscard:
		c.iLevel = LevelDiscard
	default:
		return errors.New("Wrong LogLevel by config")
	}

	if err := initStyle(c); err != nil {
		return err
	}

	switch c.LogMedium {
	case toConsole:
		return setupWithConsole(c)
	case toFile:
		return setupWithFiles(c)
	case toVolume:
		return setupWithVolume(c)
	default:
		return errors.New("Wrong LogMedium by config")
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 第一种：打印在console
func setupWithConsole(c *LogConfig) error {
	initOnce.Do(func() {
		ioInfo = newLogWriter(log.New(os.Stdout, "", 0))
		ioDebug = ioInfo
		ioStat = ioInfo
		ioSlow = ioInfo
		ioTimer = ioInfo
		ioWarn = newLogWriter(log.New(os.Stderr, "", 0))
		ioErr = ioWarn
		ioStack = ioWarn
	})
	return nil
}

// 第二种：文件日志模式下的初始化工作
func setupWithFiles(c *LogConfig) error {
	if len(c.FilePath) == 0 {
		return errors.New("log file folder must be set")
	}
	initOnce.Do(func() {
		// 初始化日志文件, 用 writer-rotate 策略写日志文件
		ioInfo = createFile(labelInfo)
		// os.Stderr + os.Stdout + os.Stdin (将标准输出重定向到文件中)
		*os.Stdout = *ioInfo.(*RotateLogger).fp
		*os.Stderr = *os.Stdout
		//log.SetOutput(ioInfo) // 这里不用写了，系统自带的Logger系统默认用的就是 os.stdout 和 os.stderr

		fStep := 0
		fiNames := strings.Split(c.FileSplit, "|")

		if fiNames[fStep] != "debug" {
			ioDebug = createFile(labelDebug)
		} else {
			ioDebug = ioInfo
		}
		//if c.FileSplit&2 != 0 {
		//	ioWarn = createFile(labelWarn)
		//} else {
		//	ioWarn = ioInfo
		//}
		//if c.FileSplit&4 != 0 {
		//	ioErr = createFile(labelErr)
		//} else {
		//	ioErr = ioWarn
		//}
		//if c.FileSplit&8 != 0 {
		//	ioStack = createFile(labelStack)
		//} else {
		//	ioStack = ioErr
		//}
		//if c.FileSplit&32 != 0 {
		//	ioStat = createFile(labelStat)
		//} else {
		//	ioStat = ioStack
		//}
		//if c.FileSplit&64 != 0 {
		//	ioSlow = createFile(labelSlow)
		//} else {
		//	ioSlow = ioStat
		//}
		//if c.FileSplit&128 != 0 {
		//	ioTimer = createFile(labelTimer)
		//} else {
		//	ioTimer = ioSlow
		//}
	})

	return nil
}

func logFilePath(logType string) string {
	return path.Join(myCnf.FilePath, myCnf.FileName+logType+".log")
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
	c.FilePath = path.Join(c.FilePath, c.AppName, host.Hostname())
	return setupWithFiles(c)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func CloseFiles() error {
	if myCnf.LogMedium == toConsole {
		return nil
	}

	if ioDebug != nil {
		if err := ioDebug.Close(); err != nil {
			return err
		}
	}
	if ioInfo != nil {
		if err := ioInfo.Close(); err != nil {
			return err
		}
	}
	if ioWarn != nil {
		if err := ioWarn.Close(); err != nil {
			return err
		}
	}
	if ioErr != nil {
		if err := ioErr.Close(); err != nil {
			return err
		}
	}
	if ioStack != nil {
		if err := ioStack.Close(); err != nil {
			return err
		}
	}
	if ioStat != nil {
		if err := ioStat.Close(); err != nil {
			return err
		}
	}
	if ioSlow != nil {
		if err := ioSlow.Close(); err != nil {
			return err
		}
	}
	if ioTimer != nil {
		if err := ioTimer.Close(); err != nil {
			return err
		}
	}

	return nil
}

func Disable() {
	initOnce.Do(func() {
		//atomic.StoreUint32(&initialized, 1)

		//ioInfo = iox.NopCloser(ioutil.Discard)
		//ioErr = iox.NopCloser(ioutil.Discard)
		//ioSlow = iox.NopCloser(ioutil.Discard)
		//ioStat = iox.NopCloser(ioutil.Discard)
		//ioStack = ioutil.Discard
	})
}
