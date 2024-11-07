// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license

// Note: 这里debug,info,stat,slow采用直接调用output的方式，而warn,error,stack采用封装调用的方式是特意设计的。
// 提取封装函数再调用能简化代码，但都采用封装调用的方式，很有可能条件不满足，大量的fmt.Sprint函数做无用功。
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/core/cst"
	"runtime/debug"
)

func Stack() *Record {
	if ShowStack() {
		return NewRecord(ioStack, labelStack)
	}
	return nil
}

func Debug() *Record {
	if ShowDebug() {
		return NewRecord(ioDebug, labelDebug)
	}
	return nil
}

func Info() *Record {
	if ShowInfo() {
		return NewRecord(ioInfo, labelInfo)
	}
	return nil
}

func InfoReq() *Record {
	if ShowInfo() {
		return NewRecord(ioInfo, labelReq)
	}
	return nil
}

func InfoTimers() *Record {
	if ShowInfo() {
		return NewRecord(ioInfo, labelTimer)
	}
	return nil
}

func InfoStat() *Record {
	if ShowInfo() {
		return NewRecord(ioInfo, labelStat)
	}
	return nil
}

func Warn() *Record {
	if ShowWarn() {
		return NewRecord(ioWarn, labelWarn)
	}
	return nil
}

func WarnSlow() *Record {
	if ShowWarn() {
		return NewRecord(ioWarn, labelSlow)
	}
	return nil
}

func Err() *Record {
	if ShowErr() {
		return NewRecord(ioErr, labelErr)
	}
	return nil
}

func ErrPanic() *Record {
	if ShowErr() {
		return NewRecord(ioErr, labelPanic)
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func ShowStack() bool {
	return myCnf.iLevel <= LevelStack
}

func ShowDebug() bool {
	return myCnf.iLevel <= LevelDebug
}

func ShowInfo() bool {
	return myCnf.iLevel <= LevelInfo
}

func ShowWarn() bool {
	return myCnf.iLevel <= LevelWarn
}

func ShowErr() bool {
	return myCnf.iLevel <= LevelErr
}

func ShowStat() bool {
	return myCnf.EnableStat && ShowInfo()
}

func ShowSlow() bool {
	return myCnf.EnableSlow && ShowWarn()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func Stack(v string) {
//	stackSync(v)
//}
//
//func Stacks(v ...any) {
//	stackSync(fmt.Sprint(v...))
//}
//
//func StackF(format string, v ...any) {
//	stackSync(fmt.Sprintf(format, v...))
//}
//
//// +++
//func Debug(v string) {
//	if myCnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, v)
//	}
//}
//
//func Debugs(v ...any) {
//	if myCnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, fmt.Sprint(v...))
//	}
//}
//
//func DebugF(format string, v ...any) {
//	if myCnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, fmt.Sprintf(format, v...))
//	}
//}
//
//func DebugDirect(v string) {
//	if myCnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, v)
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//	func Info(v string) {
//		if myCnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, v)
//		}
//	}
//
//	func InfoKV(v cst.KV) {
//		if myCnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, v)
//		}
//	}
//
//	func Infos(v ...any) {
//		if myCnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, fmt.Sprint(v...))
//		}
//	}
//
//	func InfoF(format string, v ...any) {
//		if myCnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, fmt.Sprintf(format, v...))
//		}
//	}
//
// // 直接打印所给的数据
//
//	func InfoDirect(v string) {
//		if myCnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, v)
//		}
//	}
//
// // +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//	func Warn(v string) {
//		warnSync(v)
//	}
//
//	func Warns(v ...any) {
//		warnSync(fmt.Sprint(v...))
//	}
//
//	func WarnF(format string, v ...any) {
//		warnSync(fmt.Sprintf(format, v...))
//	}
//
// // +++
//
//	func Error(v string) {
//		errorSync(v, callerSkipDepth)
//	}
//
//	func Errors(v ...any) {
//		errorSync(fmt.Sprint(v...), callerSkipDepth)
//	}
//
//	func ErrorF(format string, v ...any) {
//		errorSync(fmt.Sprintf(format, v...), callerSkipDepth)
//	}
//
//	func ErrorFatal(v ...any) {
//		errorSync(fmt.Sprint(v...), callerSkipDepth)
//		os.Exit(1)
//	}
//
//	func ErrorFatalF(format string, v ...any) {
//		errorSync(fmt.Sprintf(format, v...), callerSkipDepth)
//		os.Exit(1)
//	}
//
// // +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Stat(v string) {
	if myCnf.EnableStat {
		output(ioStat, labelStat, v)
	}
}

func StatKV(data cst.KV) {
	if myCnf.EnableStat {
		output(ioStat, labelStat, data)
	}
}

//
//func Stats(v ...any) {
//	if myCnf.EnableStat {
//		output(ioStat, labelStat, fmt.Sprint(v...))
//	}
//}
//
//func StatF(format string, v ...any) {
//	if myCnf.EnableStat {
//		output(ioStat, labelStat, fmt.Sprintf(format, v...))
//	}
//}
//
//// +++
//func Slow(v string) {
//	if myCnf.EnableStat {
//		output(ioSlow, labelSlow, v)
//	}
//}
//
//func Slows(v ...any) {
//	if myCnf.EnableStat {
//		output(ioSlow, labelSlow, fmt.Sprint(v...))
//	}
//}
//
//func SlowF(format string, v ...any) {
//	if myCnf.EnableStat {
//		output(ioSlow, labelSlow, fmt.Sprintf(format, v...))
//	}
//}
//
//// +++
//func Timer(v string) {
//	if myCnf.EnableStat {
//		output(ioTimer, labelTimer, v)
//	}
//}
//
//func TimerKV(data cst.KV) {
//	if myCnf.EnableStat {
//		output(ioTimer, labelTimer, data)
//	}
//}
//
//func Timers(v ...any) {
//	if myCnf.EnableStat {
//		output(ioTimer, labelTimer, fmt.Sprint(v...))
//	}
//}
//
//func TimerF(format string, v ...any) {
//	if myCnf.EnableStat {
//		output(ioTimer, labelTimer, fmt.Sprintf(format, v...))
//	}
//}
//
//func TimerError(v string) {
//	if myCnf.EnableStat {
//		output(ioErr, labelErr, msgWithStack(v))
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inner call apis
func warnSync(msg string) {
	if ShowWarn() {
		output(ioWarn, labelWarn, msg)
	}
}

func errorSync(msg string, skip int) {
	if ShowErr() {
		output(ioErr, labelErr, msgWithCaller(msg, skip))
	}
}

func stackSync(msg string) {
	if ShowStack() {
		output(ioStack, labelStack, fmt.Sprintf("MSG: %s Stack: %s", msg, debug.Stack()))
	}
}
