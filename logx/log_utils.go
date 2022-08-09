// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/timex"
	"log"
	"runtime"
	"strconv"
	"strings"
)

var (
	timeFormat     = "2006-01-02 15:04:05"
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
	log.SetPrefix("[GoFast]")    // 前置字符串加上特定标记
	log.SetFlags(log.Lmsgprefix) // 取消前置字符串
	log.SetFlags(log.LstdFlags)  // 设置成日期+时间 格式

	if useStyle == true {
		switch myCnf.logStyle {
		case LogStyleSdx:
			// fmt.Sprint("[", getTimestampMini(), "][", logLevel, "]: ", info)
			sb := strings.Builder{}
			sb.Grow(len(info) + 26)
			sb.WriteByte('[')
			sb.WriteString(getTimestampMini())
			sb.WriteString("][")
			sb.WriteString(logLevel)
			sb.WriteString("]: ")
			sb.WriteString(info)
			outputDirectBuilder(w, &sb)
			return
		case LogStyleSdxMini:
		case LogStyleJsonMini:
		case LogStyleJson:
			logWrap := logEntry{
				Timestamp: getTimestamp(),
				Level:     logLevel,
				Content:   info,
			}
			if content, err := jsonx.Marshal(logWrap); err != nil {
				outputDirectString(w, err.Error())
			} else {
				outputDirectBytes(w, content)
			}
			return
		}
	}
	outputDirectString(w, info)
}

func outputDirectString(w WriterCloser, str string) {
	if w == nil {
		log.Println(str)
	} else {
		_ = w.Writeln(str)
	}
}

// 推荐使用bytes版本
func outputDirectBytes(w WriterCloser, bytes []byte) {
	if w == nil {
		log.Println(bytes)
	} else {
		_ = w.WritelnBytes(bytes)
	}
}

// 推荐使用strings.Builder版本
func outputDirectBuilder(w WriterCloser, sb *strings.Builder) {
	if w == nil {
		log.Println(sb.String())
	} else {
		_ = w.WritelnBuilder(sb)
	}
}
