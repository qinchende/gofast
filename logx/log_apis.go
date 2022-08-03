package logx

import (
	"fmt"
	"os"
	"runtime/debug"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 严格输出指定字符串，不做任何修饰格式化
//func strictInfoSync(msg string) {
//	outputString(infoLog, msg)
//}

// 直接打印所给的数据
func Print(v ...any) {
	outputString(infoLog, fmt.Sprint(v...))
}

// 直接打印所给的数据
func PrintF(format string, v ...any) {
	outputString(infoLog, fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func Error(v ...any) {
	ErrorCaller(1, v...)
}

func ErrorF(format string, v ...any) {
	ErrorCallerF(1, format, v...)
}

func Fatal(v ...any) {
	ErrorCaller(1, v...)
	os.Exit(1)
}

func FatalF(format string, v ...any) {
	ErrorCallerF(1, format, v...)
	os.Exit(1)
}

func ErrorCaller(callDepth int, v ...any) {
	errorSync(fmt.Sprint(v...), callDepth+callerInnerDepth)
}

func ErrorCallerF(callDepth int, format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callDepth+callerInnerDepth)
}

func ErrorStack(v ...any) {
	stackSync(fmt.Sprint(v...))
}

func ErrorStackF(format string, v ...any) {
	stackSync(fmt.Sprintf(format, v...))
}

func Warn(v ...any) {
	warnSync(fmt.Sprint(v...))
}

func WarnF(format string, v ...any) {
	warnSync(fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	infoSync(fmt.Sprint(v...), false)
}

func InfoSkipLine(v ...any) {
	infoSync(fmt.Sprint(v...), true)
}

func InfoF(format string, v ...any) {
	infoSync(fmt.Sprintf(format, v...), false)
}

//func Severe(v ...any) {
//	severeSync(fmt.Sprint(v...))
//}
//
//func SevereF(format string, v ...any) {
//	severeSync(fmt.Sprintf(format, v...))
//}

func Slow(v ...any) {
	slowSync(fmt.Sprint(v...))
}

func SlowF(format string, v ...any) {
	slowSync(fmt.Sprintf(format, v...))
}

func Stat(v ...any) {
	statSync(fmt.Sprint(v...))
}

func StatF(format string, v ...any) {
	statSync(fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func errorSync(msg string, callDepth int) {
	if myCnf.logLevel <= LogLevelError {
		outputError(errorLog, msg, callDepth)
	}
}

func infoSync(msg string, newLine bool) {
	if myCnf.logLevel <= LogLevelInfo {
		output(infoLog, typeInfo, msg, newLine)
	}
}

func warnSync(msg string) {
	if myCnf.logLevel <= LogLevelInfo {
		output(warnLog, typeWarn, msg, false)
	}
}

//
//func severeSync(msg string) {
//	if myCnf.logLevel <= SevereLevel {
//		output(severeLog, typeSevere, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())), false)
//	}
//}

func slowSync(msg string) {
	if myCnf.logLevel <= LogLevelError {
		output(slowLog, typeSlow, msg, false)
	}
}

func stackSync(msg string) {
	if myCnf.logLevel <= LogLevelError {
		output(stackLog, typeError, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())), false)
	}
}

func statSync(msg string) {
	if myCnf.logLevel <= LogLevelInfo {
		output(statLog, typeStat, msg, false)
	}
}

// 向对应的文件（描述符）写入日志记录
func logBytes(buf []byte) {
	_, _ = fmt.Fprint(infoLog, buf)
}

func logString(text string) {
	infoSync(text, false)
}

func outputError(lwt WriterCloser, msg string, callDepth int) {
	content := formatWithCaller(msg, callDepth)
	output(lwt, typeError, content, false)
}
