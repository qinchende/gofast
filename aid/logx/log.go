// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/aid/conf"
	"github.com/qinchende/gofast/aid/sysx/host"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

var myLogger *Logger

// 指定LogConfig初始化默认日志记录器
func SetupDefault(cnf *LogConfig) {
	if myLogger != nil {
		myLogger.Warn().Msg("logx: default logger already existed")
		return
	}
	if cnf == nil {
		cnf = &LogConfig{}
		_ = conf.LoadFromJson(cnf, []byte("{}"))
	}
	myLogger = NewMust(cnf)
}

// 可以自定义默认日志记录器
func SetDefault(l *Logger) {
	if l == nil {
		myLogger = l
	}
}

func NewMust(cnf *LogConfig) *Logger {
	l, err := New(cnf)
	if err != nil {
		msg := msgWithStack(err.Error())
		_, _ = fmt.Fprintf(os.Stderr, "logx: NewMust error: %s\n", msg)
		os.Exit(1)
	}
	return l
}

func New(cnf *LogConfig) (l *Logger, err error) {
	l = &Logger{cnf: cnf}
	if err = l.initLogger(); err != nil {
		return nil, err
	}
	return l, nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (l *Logger) initLogger() error {
	switch l.cnf.LogLevel {
	case labelStack:
		l.iLevel = LevelStack
	case labelDebug:
		l.iLevel = LevelDebug
	case labelInfo:
		l.iLevel = LevelInfo
	case labelWarn:
		l.iLevel = LevelWarn
	case labelErr:
		l.iLevel = LevelErr
	case labelDiscard:
		l.iLevel = LevelDiscard
	default:
		return errors.New("Wrong LogLevel of config")
	}

	if err := l.initStyle(); err != nil {
		return err
	}

	switch l.cnf.LogMedium {
	case toConsole:
		return l.setupForConsole()
	case toFile:
		return l.setupForFiles()
	case toVolume:
		return l.setupForVolume()
	case toCustom:
		return nil
	default:
		return errors.New("Wrong LogMedium of config")
	}
}

// 第一种：打印在console
func (l *Logger) setupForConsole() error {
	l.initOnce.Do(func() {
		w1 := os.Stdout
		l.ioStack = w1
		l.ioDebug = w1
		l.ioInfo = w1
		l.ioReq = w1
		l.ioTimer = w1
		l.ioStat = w1

		w2 := os.Stderr
		l.ioWarn = w2
		l.ioSlow = w2
		l.ioErr = w2
		l.ioPanic = w2
	})
	return nil
}

// 第二种：文件日志模式下的初始化工作
func (l *Logger) setupForFiles() error {
	c := l.cnf
	if len(c.FilePath) == 0 {
		return errors.New("log file folder must be set")
	}
	l.initOnce.Do(func() {
		// 初始化日志文件, 用 writer-rotate 策略写日志文件
		l.ioInfo = l.createFile(labelInfo)
		// os.Stderr + os.Stdout + os.Stdin (将标准输出重定向到文件中)
		*os.Stdout = *l.ioInfo.(*RotateWriter).fp
		*os.Stderr = *os.Stdout
		log.SetOutput(l.ioInfo) // 这里不用写了，系统自带的Logger系统默认用的就是 os.stdout 和 os.stderr

		fStep := 0
		fiNames := strings.Split(c.FileSplit, "|")

		if fiNames[fStep] != "debug" {
			l.ioDebug = l.createFile(labelDebug)
		} else {
			l.ioDebug = l.ioInfo
		}
	})

	return nil
}

func (l *Logger) createWriterFile(path string) io.WriteCloser {
	rr := DefDailyRotateRule(path, backupFileDelimiter, l.cnf.FileKeepDays, l.cnf.FileGzip)
	wr, err := NewRotateWriter(path, rr, l.cnf.FileGzip)
	if err != nil {
		panic(err)
	}
	return wr
}

func (l *Logger) createFile(label string) io.WriteCloser {
	if l.cnf.FileName == "" {
		l.cnf.FileName = "[FileName]"
	}
	filePath := path.Join(l.cnf.FilePath, l.cnf.FileName+"."+label+".log")
	return l.createWriterFile(filePath)
}

// 第三种：分卷存储文件（其实也是写文件，但是更严格的分层文件夹。）
func (l *Logger) setupForVolume() error {
	c := l.cnf
	if len(c.AppName) == 0 {
		return errors.New("log config item [AppName] must be set")
	}
	c.FilePath = path.Join(c.FilePath, c.AppName, host.Hostname())
	return l.setupForFiles()
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func CloseFiles() error {
//	if myCnf.LogMedium == toConsole {
//		return nil
//	}
//
//	if ioDebug != nil {
//		if err := ioDebug.Close(); err != nil {
//			return err
//		}
//	}
//	if ioInfo != nil {
//		if err := ioInfo.Close(); err != nil {
//			return err
//		}
//	}
//	if ioWarn != nil {
//		if err := ioWarn.Close(); err != nil {
//			return err
//		}
//	}
//	if ioErr != nil {
//		if err := ioErr.Close(); err != nil {
//			return err
//		}
//	}
//	if ioStack != nil {
//		if err := ioStack.Close(); err != nil {
//			return err
//		}
//	}
//	if ioStat != nil {
//		if err := ioStat.Close(); err != nil {
//			return err
//		}
//	}
//	if ioSlow != nil {
//		if err := ioSlow.Close(); err != nil {
//			return err
//		}
//	}
//	if ioTimer != nil {
//		if err := ioTimer.Close(); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func Disable() {
//	initOnce.Do(func() {
//		//atomic.StoreUint32(&initialized, 1)
//
//		//ioInfo = iox.NopCloser(ioutil.Discard)
//		//ioErr = iox.NopCloser(ioutil.Discard)
//		//ioSlow = iox.NopCloser(ioutil.Discard)
//		//ioStat = iox.NopCloser(ioutil.Discard)
//		//ioStack = ioutil.Discard
//	})
//}
