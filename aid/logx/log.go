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
	"time"
)

var Def *Logger

func NewDefaultConfig() *LogConfig {
	cnf := &LogConfig{}
	_ = conf.LoadFromJson(cnf, []byte("{}"))
	return cnf
}

// 指定LogConfig初始化默认日志记录器
func SetupDefault(cnf *LogConfig) {
	if Def != nil {
		Warn().Msg("logx: default logger already existed")
		return
	}
	Def = NewMust(cnf)
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
	if cnf == nil {
		cnf = NewDefaultConfig()
	}
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
	case LabelDisable:
		l.iLevel = LevelDisable
	default:
		return errors.New("Wrong LogLevel of config")
	}

	if err := l.initStyle(); err != nil {
		return err
	}

	// 全局内容
	if len(l.cnf.AppName) > 0 {
		l.Str(fApp, l.cnf.AppName)
	}
	if len(l.cnf.HostName) > 0 {
		l.Str(fHost, l.cnf.HostName)
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
	w1 := io.Writer(os.Stdout)
	w2 := io.Writer(os.Stderr)
	if l.cnf.DiscardIO {
		w1 = io.Discard
		w2 = io.Discard
	}

	l.WStack = w1
	l.WDebug = w1
	l.WInfo = w1
	l.WReq = w1
	l.WTimer = w1
	l.WStat = w1

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

func (l *Logger) createWriterFile(path string) io.Writer {
	rr := DefDailyRotateRule(path, backupFileDelimiter, l.cnf.FileKeepDays, l.cnf.FileGzip)
	wr, err := NewRotateWriter(path, rr, l.cnf.FileGzip)
	if err != nil {
		panic(err)
	}
	return wr
}

func (l *Logger) createFile(label string) io.Writer {
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++ int
func (lr *TopRecord) Int(k string, v int) *TopRecord {
	lr.r.Int(k, v)
	return lr
}

func (lr *TopRecord) I64(k string, v int64) *TopRecord {
	lr.r.I64(k, v)
	return lr
}

func (lr *TopRecord) Ints(k string, v []int) *TopRecord {
	lr.r.Ints(k, v)
	return lr
}

func (lr *TopRecord) I8s(k string, v []int8) *TopRecord {
	lr.r.I8s(k, v)
	return lr
}

func (lr *TopRecord) I16s(k string, v []int16) *TopRecord {
	lr.r.I16s(k, v)
	return lr
}

func (lr *TopRecord) I32s(k string, v []int32) *TopRecord {
	lr.r.I32s(k, v)
	return lr
}

func (lr *TopRecord) I64s(k string, v []int64) *TopRecord {
	lr.r.I64s(k, v)
	return lr
}

// +++ uint
func (lr *TopRecord) Uint(k string, v uint) *TopRecord {
	lr.r.Uint(k, v)
	return lr
}

func (lr *TopRecord) U64(k string, v uint64) *TopRecord {
	lr.r.U64(k, v)
	return lr
}

func (lr *TopRecord) Uints(k string, v []uint) *TopRecord {
	lr.r.Uints(k, v)
	return lr
}

func (lr *TopRecord) U8s(k string, v []uint8) *TopRecord {
	lr.r.U8s(k, v)
	return lr
}

func (lr *TopRecord) U16s(k string, v []uint16) *TopRecord {
	lr.r.U16s(k, v)
	return lr
}

func (lr *TopRecord) U32s(k string, v []uint32) *TopRecord {
	lr.r.U32s(k, v)
	return lr
}

func (lr *TopRecord) U64s(k string, v []uint64) *TopRecord {
	lr.r.U64s(k, v)
	return lr
}

// +++ bool
func (lr *TopRecord) Bool(k string, v bool) *TopRecord {
	lr.r.Bool(k, v)
	return lr
}

func (lr *TopRecord) Bools(k string, v []bool) *TopRecord {
	lr.r.Bools(k, v)
	return lr
}

// +++ float
func (lr *TopRecord) F32(k string, v float32) *TopRecord {
	lr.r.F32(k, v)
	return lr
}

func (lr *TopRecord) F32s(k string, v []float32) *TopRecord {
	lr.r.F32s(k, v)
	return lr
}

func (lr *TopRecord) F64(k string, v float64) *TopRecord {
	lr.r.F64(k, v)
	return lr
}

func (lr *TopRecord) F64s(k string, v []float64) *TopRecord {
	lr.r.F64s(k, v)
	return lr
}

// +++ string
func (lr *TopRecord) Str(k, v string) *TopRecord {
	lr.r.Str(k, v)
	return lr
}

func (lr *TopRecord) Strs(k string, v []string) *TopRecord {
	lr.r.Strs(k, v)
	return lr
}

// +++ time.Times
func (lr *TopRecord) Time(k string, v time.Time) *TopRecord {
	lr.r.Time(k, v)
	return lr
}

func (lr *TopRecord) Times(k string, v []time.Time) *TopRecord {
	lr.r.Times(k, v)
	return lr
}

// +++ struct
func (lr *TopRecord) Obj(k string, v ObjEncoder) *TopRecord {
	lr.r.Obj(k, v)
	return lr
}

func (lr *TopRecord) Objs(k string, v []ObjEncoder) *TopRecord {
	lr.r.Objs(k, v)
	return lr
}

// +++ any
func (lr *TopRecord) Any(k string, v any) *TopRecord {
	lr.r.Any(k, v)
	return lr
}
