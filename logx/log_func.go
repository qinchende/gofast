package logx

import (
	"encoding/json"
	"fmt"
	"github.com/qinchende/gofast/skill/timex"
	"log"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"
)

func Alert(v string) {
	output(errorLog, levelAlert, v)
}

func Close() error {
	if writeConsole {
		return nil
	}

	if atomic.LoadUint32(&initialized) == 0 {
		return ErrLogNotInitialized
	}

	atomic.StoreUint32(&initialized, 0)

	if accessLog != nil {
		if err := accessLog.Close(); err != nil {
			return err
		}
	}

	if warnLog != nil {
		if err := warnLog.Close(); err != nil {
			return err
		}
	}

	if errorLog != nil {
		if err := errorLog.Close(); err != nil {
			return err
		}
	}

	if severeLog != nil {
		if err := severeLog.Close(); err != nil {
			return err
		}
	}

	if slowLog != nil {
		if err := slowLog.Close(); err != nil {
			return err
		}
	}

	if statLog != nil {
		if err := statLog.Close(); err != nil {
			return err
		}
	}

	return nil
}

func Disable() {
	once.Do(func() {
		atomic.StoreUint32(&initialized, 1)

		//infoLog = iox.NopCloser(ioutil.Discard)
		//errorLog = iox.NopCloser(ioutil.Discard)
		//severeLog = iox.NopCloser(ioutil.Discard)
		//slowLog = iox.NopCloser(ioutil.Discard)
		//statLog = iox.NopCloser(ioutil.Discard)
		//stackLog = ioutil.Discard
	})
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 严格输出指定字符串，不做任何修饰格式化
//func strictInfoSync(msg string) {
//	outputString(infoLog, msg)
//}

// 直接打印所给的数据
func Print(v ...any) {
	outputString(accessLog, fmt.Sprint(v...))
}

// 直接打印所给的数据
func Printf(format string, v ...any) {
	outputString(accessLog, fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func Error(v ...any) {
	ErrorCaller(1, v...)
}

func Errorf(format string, v ...any) {
	ErrorCallerf(1, format, v...)
}

func ErrorCaller(callDepth int, v ...any) {
	errorSync(fmt.Sprint(v...), callDepth+callerInnerDepth)
}

func ErrorCallerf(callDepth int, format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callDepth+callerInnerDepth)
}

func ErrorStack(v ...any) {
	// there is newline in stack string
	stackSync(fmt.Sprint(v...))
}

func ErrorStackf(format string, v ...any) {
	// there is newline in stack string
	stackSync(fmt.Sprintf(format, v...))
}

func Warn(v ...any) {
	warnSync(fmt.Sprint(v...))
}

func Warnf(format string, v ...any) {
	warnSync(fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	infoSync(fmt.Sprint(v...))
}

func Infof(format string, v ...any) {
	infoSync(fmt.Sprintf(format, v...))
}

func Severe(v ...any) {
	severeSync(fmt.Sprint(v...))
}

func Severef(format string, v ...any) {
	severeSync(fmt.Sprintf(format, v...))
}

func Slow(v ...any) {
	slowSync(fmt.Sprint(v...))
}

func Slowf(format string, v ...any) {
	slowSync(fmt.Sprintf(format, v...))
}

func Stat(v ...any) {
	statSync(fmt.Sprint(v...))
}

func Statf(format string, v ...any) {
	statSync(fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func errorSync(msg string, callDepth int) {
	if shouldLog(ErrorLevel) {
		outputError(errorLog, msg, callDepth)
	}
}

func formatWithCaller(msg string, callDepth int) string {
	var buf strings.Builder

	caller := getCaller(callDepth)
	if len(caller) > 0 {
		buf.WriteString(caller)
		buf.WriteByte(' ')
	}

	buf.WriteString(msg)

	return buf.String()
}

func getCaller(callDepth int) string {
	var buf strings.Builder

	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		buf.WriteString(short)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
	}

	return buf.String()
}

func getTimestamp() string {
	return timex.Time().Format(timeFormat)
}

func getTimestampMini() string {
	return timex.Time().Format(timeFormatMini)
}

func infoSync(msg string) {
	if shouldLog(InfoLevel) {
		output(accessLog, levelInfo, msg)
	}
}

func warnSync(msg string) {
	if shouldLog(InfoLevel) {
		output(warnLog, levelWarn, msg)
	}
}

func severeSync(msg string) {
	if shouldLog(SevereLevel) {
		output(severeLog, levelSevere, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func slowSync(msg string) {
	if shouldLog(ErrorLevel) {
		output(slowLog, levelSlow, msg)
	}
}

func stackSync(msg string) {
	if shouldLog(ErrorLevel) {
		//output(stackLog, levelError, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
		output(errorLog, levelError, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func statSync(msg string) {
	if shouldLog(InfoLevel) {
		output(statLog, levelStat, msg)
	}
}

// 向对应的文件（描述符）写入日志记录
//func logBytes(buf []byte) {
//	_, _ = fmt.Fprint(DefaultWriter, buf)
//}
//
//func logString(text string) {
//	infoSync(text)
//}

func outputError(lwt WriterCloser, msg string, callDepth int) {
	content := formatWithCaller(msg, callDepth)
	output(lwt, levelError, content)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 日志的输出，最后都要落脚到这个方法
func output(lwt WriterCloser, level, msg string) {
	// 自定义了 sdx 这种输出样式，否则就是默认的 json 样式
	//log.SetPrefix("[GoFast]")    // 前置字符串加上特定标记
	//log.SetFlags(log.Lmsgprefix) // 取消前置字符串
	//log.SetFlags(log.LstdFlags) // 设置成日期+时间 格式

	// TODO: 打印日志，套用不同的日志模板
	switch currConfig.style {
	case StyleSdx:
		outputString(lwt, fmt.Sprint("[", getTimestampMini(), "][", level, "]: ", msg))

	case StyleSdxMini:
		outputString(lwt, msg)

	case StyleJsonMini:
		outputJson(lwt, msg)

	default:
		info := logEntry{
			Timestamp: getTimestamp(),
			Level:     level,
			Content:   msg,
		}
		outputJson(lwt, info)
	}
}

func outputJson(lwt WriterCloser, info any) {
	if content, err := json.Marshal(info); err != nil {
		log.Println(err.Error())
	} else if atomic.LoadUint32(&initialized) == 0 || lwt == nil {
		log.Println(string(content))
	} else {
		_, _ = lwt.Write(append(content, '\n'))
	}
}

func outputString(lwt WriterCloser, info string) {
	if atomic.LoadUint32(&initialized) == 0 || lwt == nil {
		log.Println(info)
	} else {
		_ = lwt.Writeln(info)
	}
}
