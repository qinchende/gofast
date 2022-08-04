package logx

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Debug(v ...any) {
	debugSync(fmt.Sprint(v...))
}

func DebugDirect(v ...any) {
	debugSyncDirect(fmt.Sprint(v...))
}

func DebugF(format string, v ...any) {
	debugSync(fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 直接打印所给的数据
func InfoDirect(v ...any) {
	infoSyncDirect(fmt.Sprint(v...))
}

func InfoDirectF(format string, v ...any) {
	infoSyncDirect(fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	infoSync(fmt.Sprint(v...))
}

func InfoF(format string, v ...any) {
	infoSync(fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Warn(v ...any) {
	warnSync(fmt.Sprint(v...))
}

func WarnF(format string, v ...any) {
	warnSync(fmt.Sprintf(format, v...))
}

func Error(v ...any) {
	errorSync(fmt.Sprint(v...), callerInnerDepth)
}

func ErrorF(format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callerInnerDepth)
}

func ErrorFatal(v ...any) {
	errorSync(fmt.Sprint(v...), callerInnerDepth)
	os.Exit(1)
}

func ErrorFatalF(format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callerInnerDepth)
	os.Exit(1)
}

func Stack(v ...any) {
	stackSync(fmt.Sprint(v...))
}

func StackF(format string, v ...any) {
	stackSync(fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Stat(v ...any) {
	statSync(fmt.Sprint(v...))
}

func StatF(format string, v ...any) {
	statSync(fmt.Sprintf(format, v...))
}

func Slow(v ...any) {
	slowSync(fmt.Sprint(v...))
}

func SlowF(format string, v ...any) {
	slowSync(fmt.Sprintf(format, v...))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inner call apis
func debugSync(msg string) {
	if myCnf.logLevel <= LogLevelDebug {
		output(debugLog, typeDebug, msg)
	}
}

func debugSyncDirect(msg string) {
	if myCnf.logLevel <= LogLevelDebug {
		outputDirect(debugLog, msg)
	}
}

func infoSync(msg string) {
	if myCnf.logLevel <= LogLevelInfo {
		output(infoLog, typeInfo, msg)
	}
}

func infoSyncDirect(msg string) {
	if myCnf.logLevel <= LogLevelInfo {
		outputDirect(infoLog, msg)
	}
}

func warnSync(msg string) {
	if myCnf.logLevel <= LogLevelWarn {
		output(warnLog, typeWarn, msg)
	}
}

func errorSync(msg string, callDepth int) {
	if myCnf.logLevel <= LogLevelError {
		content := formatWithCaller(msg, callDepth)
		output(errorLog, typeError, content)
	}
}

func stackSync(msg string) {
	if myCnf.logLevel <= LogLevelStack {
		output(stackLog, typeStack, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())))
	}
}

func statSync(msg string) {
	if myCnf.LogStats {
		output(statLog, typeStat, msg)
	}
}

func slowSync(msg string) {
	if myCnf.LogStats {
		output(slowLog, typeSlow, msg)
	}
}
