// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license

// Note: 这里debug,info,stat,slow采用直接调用output的方式，而warn,error,stack采用封装调用的方式是特意设计的。
// 提取封装函数再调用能简化代码，但都采用封装调用的方式，很有可能条件不满足，大量的fmt.Sprint函数做无用功。
package logx

// Default logger
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Stack() *Record {
	return myLogger.Stack()
}

func Debug() *Record {
	return myLogger.Debug()
}

func Info() *Record {
	return myLogger.Info()
}

func InfoTimer() *Record {
	return myLogger.InfoTimer()
}

func InfoStat() *Record {
	return myLogger.InfoStat()
}

func Warn() *Record {
	return myLogger.Warn()
}

func WarnSlow() *Record {
	return myLogger.WarnSlow()
}

func Err() *Record {
	return myLogger.Err()
}

func ErrPanic() *Record {
	return myLogger.ErrPanic()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Logger Methods
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (l *Logger) ShowStack() bool {
	return l.iLevel <= LevelStack
}

func (l *Logger) ShowDebug() bool {
	return l.iLevel <= LevelDebug
}

func (l *Logger) ShowInfo() bool {
	return l.iLevel <= LevelInfo
}

func (l *Logger) ShowWarn() bool {
	return l.iLevel <= LevelWarn
}

func (l *Logger) ShowErr() bool {
	return l.iLevel <= LevelErr
}

func (l *Logger) ShowStat() bool {
	return l.ShowInfo() && !l.cnf.DisableStat
}

func (l *Logger) ShowSlow() bool {
	return l.ShowWarn() && !l.cnf.DisableSlow
}

// @@++@@
func (l *Logger) Stack() *Record {
	if l.ShowStack() {
		return newRecord(l.ioStack, labelStack)
	}
	return nil
}

func (l *Logger) Debug() *Record {
	if l.ShowDebug() {
		return newRecord(l.ioDebug, labelDebug)
	}
	return nil
}

func (l *Logger) Info() *Record {
	if l.ShowInfo() {
		return newRecord(l.ioInfo, labelInfo)
	}
	return nil
}

func (l *Logger) InfoTimer() *Record {
	if l.ShowInfo() {
		return newRecord(l.ioTimer, labelTimer)
	}
	return nil
}

func (l *Logger) InfoStat() *Record {
	if l.ShowStat() {
		return newRecord(l.ioStat, labelStat)
	}
	return nil
}

func (l *Logger) Warn() *Record {
	if l.ShowWarn() {
		return newRecord(l.ioWarn, labelWarn)
	}
	return nil
}

func (l *Logger) WarnSlow() *Record {
	if l.ShowSlow() {
		return newRecord(l.ioSlow, labelSlow)
	}
	return nil
}

func (l *Logger) Err() *Record {
	if l.ShowErr() {
		return newRecord(l.ioErr, labelErr)
	}
	return nil
}

func (l *Logger) ErrPanic() *Record {
	if l.ShowErr() {
		return newRecord(l.ioErr, labelPanic)
	}
	return nil
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
//	if cnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, v)
//	}
//}
//
//func Debugs(v ...any) {
//	if cnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, fmt.Sprint(v...))
//	}
//}
//
//func DebugF(format string, v ...any) {
//	if cnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, fmt.Sprintf(format, v...))
//	}
//}
//
//func DebugDirect(v string) {
//	if cnf.iLevel <= LevelDebug {
//		output(ioDebug, labelDebug, v)
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//	func Info(v string) {
//		if cnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, v)
//		}
//	}
//
//	func InfoKV(v cst.KV) {
//		if cnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, v)
//		}
//	}
//
//	func Infos(v ...any) {
//		if cnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, fmt.Sprint(v...))
//		}
//	}
//
//	func InfoF(format string, v ...any) {
//		if cnf.iLevel <= LevelInfo {
//			output(ioInfo, labelInfo, fmt.Sprintf(format, v...))
//		}
//	}
//
// // 直接打印所给的数据
//
//	func InfoDirect(v string) {
//		if cnf.iLevel <= LevelInfo {
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
//func Stat(v string) {
//	if !cnf.DisableStat {
//		output(ioStat, labelStat, v)
//	}
//}
//
//func StatKV(data cst.KV) {
//	if !cnf.DisableStat {
//		output(ioStat, labelStat, data)
//	}
//}

//
//func Stats(v ...any) {
//	if cnf.EnableStat {
//		output(ioStat, labelStat, fmt.Sprint(v...))
//	}
//}
//
//func StatF(format string, v ...any) {
//	if cnf.EnableStat {
//		output(ioStat, labelStat, fmt.Sprintf(format, v...))
//	}
//}
//
//// +++
//func Slow(v string) {
//	if cnf.EnableStat {
//		output(ioSlow, labelSlow, v)
//	}
//}
//
//func Slows(v ...any) {
//	if cnf.EnableStat {
//		output(ioSlow, labelSlow, fmt.Sprint(v...))
//	}
//}
//
//func SlowF(format string, v ...any) {
//	if cnf.EnableStat {
//		output(ioSlow, labelSlow, fmt.Sprintf(format, v...))
//	}
//}
//
//// +++
//func Timer(v string) {
//	if cnf.EnableStat {
//		output(ioTimer, labelTimer, v)
//	}
//}
//
//func TimerKV(data cst.KV) {
//	if cnf.EnableStat {
//		output(ioTimer, labelTimer, data)
//	}
//}
//
//func Timers(v ...any) {
//	if cnf.EnableStat {
//		output(ioTimer, labelTimer, fmt.Sprint(v...))
//	}
//}
//
//func TimerF(format string, v ...any) {
//	if cnf.EnableStat {
//		output(ioTimer, labelTimer, fmt.Sprintf(format, v...))
//	}
//}
//
//func TimerError(v string) {
//	if cnf.EnableStat {
//		output(ioErr, labelErr, msgWithStack(v))
//	}
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// inner call apis
//func warnSync(msg string) {
//	if ShowWarn() {
//		output(ioWarn, labelWarn, msg)
//	}
//}
//
//func errorSync(msg string, skip int) {
//	if ShowErr() {
//		output(ioErr, labelErr, msgWithCaller(msg, skip))
//	}
//}
//
//func stackSync(msg string) {
//	if ShowStack() {
//		output(ioStack, labelStack, fmt.Sprintf("MSG: %s Stack: %s", msg, debug.Stack()))
//	}
//}
