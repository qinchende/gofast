package logx

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Debug(v ...any) {
	debugSync(fmt.Sprint(v...), true)
}

func DebugF(format string, v ...any) {
	debugSync(fmt.Sprintf(format, v...), true)
}

func DebugDirect(v ...any) {
	debugSync(fmt.Sprint(v...), false)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Info(v ...any) {
	infoSync(fmt.Sprint(v...), true)
}

func InfoF(format string, v ...any) {
	infoSync(fmt.Sprintf(format, v...), true)
}

// 直接打印所给的数据
func InfoDirect(v ...any) {
	infoSync(fmt.Sprint(v...), false)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Warn(v ...any) {
	warnSync(fmt.Sprint(v...), true)
}

func WarnF(format string, v ...any) {
	warnSync(fmt.Sprintf(format, v...), true)
}

func Error(v ...any) {
	errorSync(fmt.Sprint(v...), callerInnerDepth, true)
}

func ErrorF(format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callerInnerDepth, true)
}

func ErrorFatal(v ...any) {
	errorSync(fmt.Sprint(v...), callerInnerDepth, true)
	os.Exit(1)
}

func ErrorFatalF(format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callerInnerDepth, true)
	os.Exit(1)
}

func Stack(v ...any) {
	stackSync(fmt.Sprint(v...), true)
}

func StackF(format string, v ...any) {
	stackSync(fmt.Sprintf(format, v...), true)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Stat(v ...any) {
	statSync(fmt.Sprint(v...), true)
}

func StatF(format string, v ...any) {
	statSync(fmt.Sprintf(format, v...), true)
}

func Slow(v ...any) {
	slowSync(fmt.Sprint(v...), true)
}

func SlowF(format string, v ...any) {
	slowSync(fmt.Sprintf(format, v...), true)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inner call apis
func debugSync(msg string, useStyle bool) {
	if myCnf.logLevel <= LogLevelDebug {
		output(debugLog, msg, levelDebug, useStyle)
	}
}

func infoSync(msg string, useStyle bool) {
	if myCnf.logLevel <= LogLevelInfo {
		output(infoLog, msg, levelInfo, useStyle)
	}
}

func warnSync(msg string, useStyle bool) {
	if myCnf.logLevel <= LogLevelWarn {
		output(warnLog, msg, levelWarn, useStyle)
	}
}

func errorSync(msg string, callDepth int, useStyle bool) {
	if myCnf.logLevel <= LogLevelError {
		output(errorLog, formatWithCaller(msg, callDepth), levelError, useStyle)
	}
}

func stackSync(msg string, useStyle bool) {
	if myCnf.logLevel <= LogLevelStack {
		output(stackLog, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())), levelStack, useStyle)
	}
}

func statSync(msg string, useStyle bool) {
	if myCnf.LogStats {
		output(statLog, msg, levelStat, useStyle)
	}
}

func slowSync(msg string, useStyle bool) {
	if myCnf.LogStats {
		output(slowLog, msg, levelSlow, useStyle)
	}
}
