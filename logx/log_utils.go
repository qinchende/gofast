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
)

var (
	timeFormat     = "2006-01-02T15:04:05.000Z07"
	timeFormatMini = "01-02 15:04:05"
)

type logEntry struct {
	Timestamp string `json:"@timestamp"`
	Level     string `json:"lv"`
	Duration  string `json:"duration,omitempty"`
	Content   string `json:"ct"`
}

func LogStyleType() int8 {
	return myCnf.logStyle
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
func output(w WriterCloser, info string, logLevel string, useStyle bool) {
	// 自定义了 sdx 这种输出样式，否则就是默认的 json 样式
	//log.SetPrefix("[GoFast]")    // 前置字符串加上特定标记
	//log.SetFlags(log.Lmsgprefix) // 取消前置字符串
	//log.SetFlags(log.LstdFlags) // 设置成日期+时间 格式

	if useStyle == true {
		switch myCnf.logStyle {
		case LogStyleSdx:
			info = fmt.Sprint("[", getTimestampMini(), "][", logLevel, "]: ", info)
		case LogStyleSdxMini:
		case LogStyleJsonMini:
		case LogStyleJson:
			logWrap := logEntry{
				Timestamp: getTimestamp(),
				Level:     logLevel,
				Content:   info,
			}
			if content, err := jsonx.Marshal(logWrap); err != nil {
				info = err.Error()
			} else {
				info = stringx.BytesToString(content)
			}
		}
	}
	outputDirect(w, info)
}

func outputDirect(w WriterCloser, str string) {
	if w == nil {
		log.Println(str)
	} else {
		_ = w.Writeln(str)
	}
}
