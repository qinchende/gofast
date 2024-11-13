// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license

// Note: 这里debug,info,stat,slow采用直接调用output的方式，而warn,error,stack采用封装调用的方式是特意设计的。
// 提取封装函数再调用能简化代码，但都采用封装调用的方式，很有可能条件不满足，大量的fmt.Sprint函数做无用功。
package logx

// Default logger
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Stack() *Record {
	return DefLogger.Trace()
}

func Debug() *Record {
	return DefLogger.Debug()
}

func Info() *Record {
	return DefLogger.Info()
}

func InfoTimer() *Record {
	return DefLogger.InfoTimer()
}

func InfoStat() *Record {
	return DefLogger.InfoStat()
}

func Warn() *Record {
	return DefLogger.Warn()
}

func WarnSlow() *Record {
	return DefLogger.WarnSlow()
}

func Err() *Record {
	return DefLogger.Err()
}

func ErrPanic() *Record {
	return DefLogger.ErrPanic()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Logger instance methods
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// @@++@@
func (l *Logger) Clone() *Logger {
	if l == nil {
		return nil
	}
	newL := *l
	newL.bs = make([]byte, 0, cap(l.bs))
	copy(newL.bs, l.bs)
	return &newL
}

// @@++@@
func (l *Logger) ShowStack() bool {
	return l.iLevel <= LevelTrace
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
func (l *Logger) Trace() *Record {
	if l.ShowStack() {
		return newRecord(l.WStack, LabelTrace)
	}
	return nil
}

func (l *Logger) Debug() *Record {
	if l.ShowDebug() {
		return newRecord(l.WDebug, LabelDebug)
	}
	return nil
}

func (l *Logger) Info() *Record {
	if l.ShowInfo() {
		return newRecord(l.WInfo, LabelInfo)
	}
	return nil
}

func (l *Logger) InfoTimer() *Record {
	if l.ShowInfo() {
		return newRecord(l.WTimer, LabelTimer)
	}
	return nil
}

func (l *Logger) InfoStat() *Record {
	if l.ShowStat() {
		return newRecord(l.WStat, LabelStat)
	}
	return nil
}

func (l *Logger) Warn() *Record {
	if l.ShowWarn() {
		return newRecord(l.WWarn, LabelWarn)
	}
	return nil
}

func (l *Logger) WarnSlow() *Record {
	if l.ShowSlow() {
		return newRecord(l.WSlow, LabelSlow)
	}
	return nil
}

func (l *Logger) Err() *Record {
	if l.ShowErr() {
		return newRecord(l.WErr, LabelErr)
	}
	return nil
}

func (l *Logger) ErrPanic() *Record {
	if l.ShowErr() {
		return newRecord(l.WErr, LabelPanic)
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func Trace(v string) {
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
//		output(WDebug, LabelDebug, v)
//	}
//}
//
//func Debugs(v ...any) {
//	if cnf.iLevel <= LevelDebug {
//		output(WDebug, LabelDebug, fmt.Sprint(v...))
//	}
//}
//
//func DebugF(format string, v ...any) {
//	if cnf.iLevel <= LevelDebug {
//		output(WDebug, LabelDebug, fmt.Sprintf(format, v...))
//	}
//}
//
//func DebugDirect(v string) {
//	if cnf.iLevel <= LevelDebug {
//		output(WDebug, LabelDebug, v)
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//	func Info(v string) {
//		if cnf.iLevel <= LevelInfo {
//			output(WInfo, LabelInfo, v)
//		}
//	}
//
//	func InfoKV(v cst.KV) {
//		if cnf.iLevel <= LevelInfo {
//			output(WInfo, LabelInfo, v)
//		}
//	}
//
//	func Infos(v ...any) {
//		if cnf.iLevel <= LevelInfo {
//			output(WInfo, LabelInfo, fmt.Sprint(v...))
//		}
//	}
//
//	func InfoF(format string, v ...any) {
//		if cnf.iLevel <= LevelInfo {
//			output(WInfo, LabelInfo, fmt.Sprintf(format, v...))
//		}
//	}
//
// // 直接打印所给的数据
//
//	func InfoDirect(v string) {
//		if cnf.iLevel <= LevelInfo {
//			output(WInfo, LabelInfo, v)
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
//		output(WStat, LabelStat, v)
//	}
//}
//
//func StatKV(data cst.KV) {
//	if !cnf.DisableStat {
//		output(WStat, LabelStat, data)
//	}
//}

//
//func Stats(v ...any) {
//	if cnf.EnableStat {
//		output(WStat, LabelStat, fmt.Sprint(v...))
//	}
//}
//
//func StatF(format string, v ...any) {
//	if cnf.EnableStat {
//		output(WStat, LabelStat, fmt.Sprintf(format, v...))
//	}
//}
//
//// +++
//func Slow(v string) {
//	if cnf.EnableStat {
//		output(WSlow, LabelSlow, v)
//	}
//}
//
//func Slows(v ...any) {
//	if cnf.EnableStat {
//		output(WSlow, LabelSlow, fmt.Sprint(v...))
//	}
//}
//
//func SlowF(format string, v ...any) {
//	if cnf.EnableStat {
//		output(WSlow, LabelSlow, fmt.Sprintf(format, v...))
//	}
//}
//
//// +++
//func Timer(v string) {
//	if cnf.EnableStat {
//		output(WTimer, LabelTimer, v)
//	}
//}
//
//func TimerKV(data cst.KV) {
//	if cnf.EnableStat {
//		output(WTimer, LabelTimer, data)
//	}
//}
//
//func Timers(v ...any) {
//	if cnf.EnableStat {
//		output(WTimer, LabelTimer, fmt.Sprint(v...))
//	}
//}
//
//func TimerF(format string, v ...any) {
//	if cnf.EnableStat {
//		output(WTimer, LabelTimer, fmt.Sprintf(format, v...))
//	}
//}
//
//func TimerError(v string) {
//	if cnf.EnableStat {
//		output(WErr, LabelErr, msgWithStack(v))
//	}
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// inner call apis
//func warnSync(msg string) {
//	if ShowWarn() {
//		output(WWarn, LabelWarn, msg)
//	}
//}
//
//func errorSync(msg string, skip int) {
//	if ShowErr() {
//		output(WErr, LabelErr, msgWithCaller(msg, skip))
//	}
//}
//
//func stackSync(msg string) {
//	if ShowStack() {
//		output(WStack, LabelTrace, fmt.Sprintf("MSG: %s Trace: %s", msg, debug.Trace()))
//	}
//}
