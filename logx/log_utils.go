package logx

import (
	"fmt"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/stringx"
	"github.com/qinchende/gofast/skill/timex"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	timeFormat     = "2006-01-02T15:04:05.000Z07"
	timeFormatMini = "01-02 15:04:05"
)

func getTimestamp() string {
	return timex.Time().Format(timeFormat)
}

func getTimestampMini() string {
	return timex.Time().Format(timeFormatMini)
}

// 消息前面加上堆栈信息
func formatWithCaller(msg string, callDepth int) string {
	callerBuf := getCaller(callDepth)
	if callerBuf.Len() > 0 {
		callerBuf.WriteByte(' ')
	}
	callerBuf.WriteString(msg)
	return callerBuf.String()
}

func getCaller(callDepth int) *strings.Builder {
	var buf strings.Builder

	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		buf.WriteString(short)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
	}
	return &buf
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 日志的输出，最后都要到这个方法进行输出
func output(w WriterCloser, level, str string, newLine bool) {
	// 自定义了 sdx 这种输出样式，否则就是默认的 json 样式
	//log.SetPrefix("[GoFast]")    // 前置字符串加上特定标记
	//log.SetFlags(log.Lmsgprefix) // 取消前置字符串
	//log.SetFlags(log.LstdFlags) // 设置成日期+时间 格式

	// TODO: 打印日志，套用不同的日志模板
	switch myCnf.logStyle {
	case StyleSdx:
		str = fmt.Sprint("[", getTimestampMini(), "][", level, "]: ", str)
	case StyleSdxMini:
	case StyleJsonMini:
	case StyleJson:
		info := logEntry{
			Timestamp: getTimestamp(),
			Level:     level,
			Content:   str,
		}
		if content, err := jsonx.Marshal(info); err != nil {
			str = err.Error()
		} else {
			str = stringx.BytesToString(content)
		}
	}

	if newLine {
		str = "\n" + str
	}
	outputString(w, str)
}

func outputString(w WriterCloser, str string) {
	if atomic.LoadUint32(&initialized) == 0 || w == nil {
		log.Println(str)
	} else {
		_ = w.Writeln(str)
	}
}
