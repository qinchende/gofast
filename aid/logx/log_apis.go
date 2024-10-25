// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license

// Note: 这里debug,info,stat,slow采用直接调用output的方式，而warn,error,stack采用封装调用的方式是特意设计的。
// 提取封装函数再调用能简化代码，但都采用封装调用的方式，很有可能条件不满足，大量的fmt.Sprint函数做无用功。
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/core/cst"
	"os"
	"runtime/debug"
)

func ShowDebug() bool {
	return myCnf.iLevel <= LevelDebug
}

func ShowInfo() bool {
	return myCnf.iLevel <= LevelInfo
}

func ShowWarn() bool {
	return myCnf.iLevel <= LevelWarn
}

func ShowError() bool {
	return myCnf.iLevel <= LevelError
}

func ShowStack() bool {
	return myCnf.iLevel <= LevelStack
}

func ShowStat() bool {
	return myCnf.LogStat
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Debug(v string) {
	if myCnf.iLevel <= LevelDebug {
		output(debugLog, labelDebug, v, true)
	}
}

func Debugs(v ...any) {
	if myCnf.iLevel <= LevelDebug {
		output(debugLog, labelDebug, fmt.Sprint(v...), true)
	}
}

func DebugF(format string, v ...any) {
	if myCnf.iLevel <= LevelDebug {
		output(debugLog, labelDebug, fmt.Sprintf(format, v...), true)
	}
}

func DebugDirect(v string) {
	if myCnf.iLevel <= LevelDebug {
		output(debugLog, labelDebug, v, false)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Info(v string) {
	if myCnf.iLevel <= LevelInfo {
		output(infoLog, labelInfo, v, true)
	}
}

func InfoKV(v cst.KV) {
	if myCnf.iLevel <= LevelInfo {
		output(infoLog, labelInfo, v, true)
	}
}

func Infos(v ...any) {
	if myCnf.iLevel <= LevelInfo {
		output(infoLog, labelInfo, fmt.Sprint(v...), true)
	}
}

func InfoF(format string, v ...any) {
	if myCnf.iLevel <= LevelInfo {
		output(infoLog, labelInfo, fmt.Sprintf(format, v...), true)
	}
}

// 直接打印所给的数据
func InfoDirect(v string) {
	if myCnf.iLevel <= LevelInfo {
		output(infoLog, labelInfo, v, false)
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

// +++
func Error(v string) {
	errorSync(v, callerSkipDepth, true)
}

func Errors(v ...any) {
	errorSync(fmt.Sprint(v...), callerSkipDepth, true)
}

func ErrorF(format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callerSkipDepth, true)
}

func ErrorFatal(v ...any) {
	errorSync(fmt.Sprint(v...), callerSkipDepth, true)
	os.Exit(1)
}

func ErrorFatalF(format string, v ...any) {
	errorSync(fmt.Sprintf(format, v...), callerSkipDepth, true)
	os.Exit(1)
}

// +++
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
	if myCnf.LogStat {
		output(statLog, labelStat, v, true)
	}
}

func StatKV(data cst.KV) {
	if myCnf.LogStat {
		output(statLog, labelStat, data, true)
	}
}

func Stats(v ...any) {
	if myCnf.LogStat {
		output(statLog, labelStat, fmt.Sprint(v...), true)
	}
}

func StatF(format string, v ...any) {
	if myCnf.LogStat {
		output(statLog, labelStat, fmt.Sprintf(format, v...), true)
	}
}

// +++
func Slow(v string) {
	if myCnf.LogStat {
		output(slowLog, labelSlow, v, true)
	}
}

func Slows(v ...any) {
	if myCnf.LogStat {
		output(slowLog, labelSlow, fmt.Sprint(v...), true)
	}
}

func SlowF(format string, v ...any) {
	if myCnf.LogStat {
		output(slowLog, labelSlow, fmt.Sprintf(format, v...), true)
	}
}

// +++
func Timer(v string) {
	if myCnf.LogStat {
		output(timerLog, labelTimer, v, true)
	}
}

func TimerKV(data cst.KV) {
	if myCnf.LogStat {
		output(timerLog, labelTimer, data, true)
	}
}

func Timers(v ...any) {
	if myCnf.LogStat {
		output(timerLog, labelTimer, fmt.Sprint(v...), true)
	}
}

func TimerF(format string, v ...any) {
	if myCnf.LogStat {
		output(timerLog, labelTimer, fmt.Sprintf(format, v...), true)
	}
}

func TimerError(v string) {
	if myCnf.LogStat {
		output(errorLog, labelError, formatWithCaller(v, callerSkipDepth), true)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inner call apis
func warnSync(msg string, useStyle bool) {
	if myCnf.iLevel <= LevelWarn {
		output(warnLog, labelWarn, msg, useStyle)
	}
}

func errorSync(msg string, callDepth int, useStyle bool) {
	if myCnf.iLevel <= LevelError {
		output(errorLog, labelError, formatWithCaller(msg, callDepth), useStyle)
	}
}

func stackSync(msg string, useStyle bool) {
	if myCnf.iLevel <= LevelStack {
		output(stackLog, labelStack, fmt.Sprintf("MSG: %s Stack: %s", msg, debug.Stack()), useStyle)
	}
}
