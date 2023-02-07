// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license

// Note: 这里debug,info,stat,slow采用直接调用output的方式，而warn,error,stack采用封装调用的方式是特意设计的。
// 提取封装函数再调用能简化代码，但都采用封装调用的方式，很有可能条件不满足，大量的fmt.Sprint函数做无用功。
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"os"
	"runtime/debug"
)

func ShowDebug() bool {
	return myCnf.logLevelInt8 <= LogLevelDebug
}

func ShowInfo() bool {
	return myCnf.logLevelInt8 <= LogLevelInfo
}

func ShowWarn() bool {
	return myCnf.logLevelInt8 <= LogLevelWarn
}

func ShowError() bool {
	return myCnf.logLevelInt8 <= LogLevelError
}

func ShowStack() bool {
	return myCnf.logLevelInt8 <= LogLevelStack
}

func ShowStat() bool {
	return myCnf.LogStats
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Debug(v string) {
	if myCnf.logLevelInt8 <= LogLevelDebug {
		output(debugLog, levelDebug, v, true)
	}
}

func Debugs(v ...any) {
	if myCnf.logLevelInt8 <= LogLevelDebug {
		output(debugLog, levelDebug, fmt.Sprint(v...), true)
	}
}

func DebugF(format string, v ...any) {
	if myCnf.logLevelInt8 <= LogLevelDebug {
		output(debugLog, levelDebug, fmt.Sprintf(format, v...), true)
	}
}

func DebugDirect(v string) {
	if myCnf.logLevelInt8 <= LogLevelDebug {
		output(debugLog, levelDebug, v, false)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Info(v string) {
	if myCnf.logLevelInt8 <= LogLevelInfo {
		output(infoLog, levelInfo, v, true)
	}
}

func InfoKV(v cst.KV) {
	if myCnf.logLevelInt8 <= LogLevelInfo {
		output(infoLog, levelInfo, v, true)
	}
}

func Infos(v ...any) {
	if myCnf.logLevelInt8 <= LogLevelInfo {
		output(infoLog, levelInfo, fmt.Sprint(v...), true)
	}
}

func InfoF(format string, v ...any) {
	if myCnf.logLevelInt8 <= LogLevelInfo {
		output(infoLog, levelInfo, fmt.Sprintf(format, v...), true)
	}
}

// 直接打印所给的数据
func InfoDirect(v string) {
	if myCnf.logLevelInt8 <= LogLevelInfo {
		output(infoLog, levelInfo, v, false)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Warn(v string) {
	warnSync(v, true)
}

func Warns(v ...any) {
	warnSync(fmt.Sprint(v...), true)
}

func WarnF(format string, v ...any) {
	warnSync(fmt.Sprintf(format, v...), true)
}

func Error(v string) {
	errorSync(v, callerInnerDepth, true)
}

func Errors(v ...any) {
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

func Stack(v string) {
	stackSync(v, true)
}

func Stacks(v ...any) {
	stackSync(fmt.Sprint(v...), true)
}

func StackF(format string, v ...any) {
	stackSync(fmt.Sprintf(format, v...), true)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Stat(v string) {
	if myCnf.LogStats {
		output(statLog, levelStat, v, true)
	}
}

func StatKV(data cst.KV) {
	if myCnf.LogStats {
		output(statLog, levelStat, data, true)
	}
}

func Stats(v ...any) {
	if myCnf.LogStats {
		output(statLog, levelStat, fmt.Sprint(v...), true)
	}
}

func StatF(format string, v ...any) {
	if myCnf.LogStats {
		output(statLog, levelStat, fmt.Sprintf(format, v...), true)
	}
}

func Slow(v string) {
	if myCnf.LogStats {
		output(slowLog, levelSlow, v, true)
	}
}

func Slows(v ...any) {
	if myCnf.LogStats {
		output(slowLog, levelSlow, fmt.Sprint(v...), true)
	}
}

func SlowF(format string, v ...any) {
	if myCnf.LogStats {
		output(slowLog, levelSlow, fmt.Sprintf(format, v...), true)
	}
}

func Script(v string) {
	if myCnf.LogStats {
		output(timerLog, levelTimer, v, true)
	}
}

func ScriptKV(data cst.KV) {
	if myCnf.LogStats {
		output(timerLog, levelTimer, data, true)
	}
}

func Scripts(v ...any) {
	if myCnf.LogStats {
		output(timerLog, levelTimer, fmt.Sprint(v...), true)
	}
}

func ScriptF(format string, v ...any) {
	if myCnf.LogStats {
		output(timerLog, levelTimer, fmt.Sprintf(format, v...), true)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inner call apis
func warnSync(msg string, useStyle bool) {
	if myCnf.logLevelInt8 <= LogLevelWarn {
		output(warnLog, levelWarn, msg, useStyle)
	}
}

func errorSync(msg string, callDepth int, useStyle bool) {
	if myCnf.logLevelInt8 <= LogLevelError {
		output(errorLog, levelError, formatWithCaller(msg, callDepth), useStyle)
	}
}

func stackSync(msg string, useStyle bool) {
	if myCnf.logLevelInt8 <= LogLevelStack {
		output(stackLog, levelStack, fmt.Sprintf("%s\n%s", msg, string(debug.Stack())), useStyle)
	}
}
