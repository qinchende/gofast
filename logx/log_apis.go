// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license

// Note: 这里debug,info,stat,slow采用直接调用output的方式，而warn,error,stack采用封装调用的方式是特意设计的。
// 提取封装函数再调用能简化代码，但都采用封装调用的方式，很有可能条件不满足，大量的fmt.Sprint函数做无用功。
package logx

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Debug(v ...any) {
	if myCnf.logLevel <= LogLevelDebug {
		output(debugLog, fmt.Sprint(v...), levelDebug, true)
	}
}

func DebugStr(v string) {
	if myCnf.logLevel <= LogLevelDebug {
		output(debugLog, v, levelDebug, true)
	}
}

func DebugF(format string, v ...any) {
	if myCnf.logLevel <= LogLevelDebug {
		output(debugLog, fmt.Sprintf(format, v...), levelDebug, true)
	}
}

func DebugDirectStr(v string) {
	if myCnf.logLevel <= LogLevelDebug {
		output(debugLog, v, levelDebug, false)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Info(v ...any) {
	if myCnf.logLevel <= LogLevelInfo {
		output(infoLog, fmt.Sprint(v...), levelInfo, true)
	}
}

func InfoStr(v string) {
	if myCnf.logLevel <= LogLevelInfo {
		output(infoLog, v, levelInfo, true)
	}
}

func InfoF(format string, v ...any) {
	if myCnf.logLevel <= LogLevelInfo {
		output(infoLog, fmt.Sprintf(format, v...), levelInfo, true)
	}
}

// 直接打印所给的数据
func InfoDirectStr(v string) {
	if myCnf.logLevel <= LogLevelInfo {
		output(infoLog, v, levelInfo, false)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Warn(v ...any) {
	warnSync(fmt.Sprint(v...), true)
}

func WarnStr(v string) {
	warnSync(v, true)
}

func WarnF(format string, v ...any) {
	warnSync(fmt.Sprintf(format, v...), true)
}

func Error(v ...any) {
	errorSync(fmt.Sprint(v...), callerInnerDepth, true)
}

func ErrorStr(v string) {
	errorSync(v, callerInnerDepth, true)
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

func StackStr(v string) {
	stackSync(v, true)
}

func StackF(format string, v ...any) {
	stackSync(fmt.Sprintf(format, v...), true)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Stat(v ...any) {
	if myCnf.LogStats {
		output(statLog, fmt.Sprint(v...), levelStat, true)
	}
}

func StatStr(v string) {
	if myCnf.LogStats {
		output(statLog, v, levelStat, true)
	}
}

func StatF(format string, v ...any) {
	if myCnf.LogStats {
		output(statLog, fmt.Sprintf(format, v...), levelStat, true)
	}
}

func Slow(v ...any) {
	if myCnf.LogStats {
		output(slowLog, fmt.Sprint(v...), levelSlow, true)
	}
}

func SlowStr(v string) {
	if myCnf.LogStats {
		output(slowLog, v, levelSlow, true)
	}
}

func SlowF(format string, v ...any) {
	if myCnf.LogStats {
		output(slowLog, fmt.Sprintf(format, v...), levelSlow, true)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inner call apis
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
