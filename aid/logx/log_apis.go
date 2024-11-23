// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.

// Note: 这里debug,info,stat,slow采用直接调用output的方式，而warn,error,stack采用封装调用的方式是特意设计的。
// 提取封装函数再调用能简化代码，但都采用封装调用的方式，很有可能条件不满足，大量的fmt.Sprint函数做无用功。
package logx

// Default logger
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Trace() *Record {
	return Def.Trace()
}

func Debug() *Record {
	return Def.Debug()
}

func Info() *Record {
	return Def.Info()
}

func InfoReq() *Record {
	return Def.InfoReq()
}

func InfoTimer() *Record {
	return Def.InfoTimer()
}

func InfoStat() *Record {
	return Def.InfoStat()
}

func Warn() *Record {
	return Def.Warn()
}

func WarnSlow() *Record {
	return Def.WarnSlow()
}

func Err() *Record {
	return Def.Err()
}

func ErrPanic() *Record {
	return Def.ErrPanic()
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
	newL.r.bs = make([]byte, 0, cap(l.r.bs))
	copy(newL.r.bs, l.r.bs)
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
		return NewRecord(l, l.WStack, LabelTrace)
	}
	return nil
}

func (l *Logger) Debug() *Record {
	if l.ShowDebug() {
		return NewRecord(l, l.WDebug, LabelDebug)
	}
	return nil
}

func (l *Logger) Info() *Record {
	if l.ShowInfo() {
		return NewRecord(l, l.WInfo, LabelInfo)
	}
	return nil
}

func (l *Logger) InfoReq() *Record {
	if l.ShowInfo() {
		return NewRecord(l, l.WReq, LabelReq)
	}
	return nil
}

func (l *Logger) InfoTimer() *Record {
	if l.ShowInfo() {
		return NewRecord(l, l.WTimer, LabelTimer)
	}
	return nil
}

func (l *Logger) InfoStat() *Record {
	if l.ShowStat() {
		return NewRecord(l, l.WStat, LabelStat)
	}
	return nil
}

func (l *Logger) Warn() *Record {
	if l.ShowWarn() {
		return NewRecord(l, l.WWarn, LabelWarn)
	}
	return nil
}

func (l *Logger) WarnSlow() *Record {
	if l.ShowSlow() {
		return NewRecord(l, l.WSlow, LabelSlow)
	}
	return nil
}

func (l *Logger) Err() *Record {
	if l.ShowErr() {
		return NewRecord(l, l.WErr, LabelErr)
	}
	return nil
}

func (l *Logger) ErrPanic() *Record {
	if l.ShowErr() {
		return NewRecord(l, l.WErr, LabelPanic)
	}
	return nil
}
