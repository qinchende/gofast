// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
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

var DefLogger *Logger

// 指定LogConfig初始化默认日志记录器
func SetupDefault(cnf *LogConfig) {
	if DefLogger != nil {
		Warn().Msg("logx: default logger already existed")
		return
	}
	if cnf == nil {
		cnf = &LogConfig{}
		_ = conf.LoadFromJson(cnf, []byte("{}"))
	}
	DefLogger = NewMust(cnf)
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
	case LabelTrace:
		l.iLevel = LevelTrace
	case LabelDebug:
		l.iLevel = LevelDebug
	case LabelInfo:
		l.iLevel = LevelInfo
	case LabelWarn:
		l.iLevel = LevelWarn
	case LabelErr:
		l.iLevel = LevelErr
	case LabelDiscard:
		l.iLevel = LevelDiscard
	default:
		return errors.New("Wrong LogLevel of config")
	}

	if err := l.initStyle(); err != nil {
		return err
	}

	// 全局内容
	l.Str(fApp, l.cnf.AppName)
	l.Str(fHost, l.cnf.HostName)

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
	w1 := os.Stdout
	l.WStack = w1
	l.WDebug = w1
	l.WInfo = w1
	l.WReq = w1
	l.WTimer = w1
	l.WStat = w1

	w2 := os.Stdout
	l.WWarn = w2
	l.WSlow = w2
	l.WErr = w2
	l.WPanic = w2
	return nil
}

// 第二种：文件日志模式下的初始化工作
func (l *Logger) setupForFiles() error {
	c := l.cnf
	if len(c.FilePath) == 0 {
		return errors.New("log file folder must be set")
	}
	// 初始化日志文件, 用 writer-rotate 策略写日志文件
	l.WInfo = l.createFile(LabelInfo)
	// os.Stderr + os.Stdout + os.Stdin (将标准输出重定向到文件中)
	*os.Stdout = *l.WInfo.(*RotateWriter).fp
	*os.Stderr = *os.Stdout
	log.SetOutput(l.WInfo) // 这里不用写了，系统自带的Logger系统默认用的就是 os.stdout 和 os.stderr

	fStep := 0
	fiNames := strings.Split(c.FileSplit, "|")

	if fiNames[fStep] != "debug" {
		l.WDebug = l.createFile(LabelDebug)
	} else {
		l.WDebug = l.WInfo
	}

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
//	if WDebug != nil {
//		if err := WDebug.Close(); err != nil {
//			return err
//		}
//	}
//	if WInfo != nil {
//		if err := WInfo.Close(); err != nil {
//			return err
//		}
//	}
//	if WWarn != nil {
//		if err := WWarn.Close(); err != nil {
//			return err
//		}
//	}
//	if WErr != nil {
//		if err := WErr.Close(); err != nil {
//			return err
//		}
//	}
//	if WStack != nil {
//		if err := WStack.Close(); err != nil {
//			return err
//		}
//	}
//	if WStat != nil {
//		if err := WStat.Close(); err != nil {
//			return err
//		}
//	}
//	if WSlow != nil {
//		if err := WSlow.Close(); err != nil {
//			return err
//		}
//	}
//	if WTimer != nil {
//		if err := WTimer.Close(); err != nil {
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
//		//WInfo = iox.NopCloser(ioutil.Discard)
//		//WErr = iox.NopCloser(ioutil.Discard)
//		//WSlow = iox.NopCloser(ioutil.Discard)
//		//WStat = iox.NopCloser(ioutil.Discard)
//		//WStack = ioutil.Discard
//	})
//}
