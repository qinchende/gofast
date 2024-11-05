// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/core/lang"
	"log"
	"runtime"
	"strconv"
	"strings"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 消息前面加上堆栈信息
func msgWithStack(msg string) string {
	return msgWithCaller(msg, callerSkipDepth)
}

func msgWithCaller(msg string, skip int) string {
	var bf strings.Builder

	_, file, line, ok := runtime.Caller(skip)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		bf.WriteString(short)
		bf.WriteByte(':')
		bf.WriteString(strconv.Itoa(line))
	}

	if bf.Len() > 0 {
		bf.WriteByte(' ')
	}
	bf.WriteString(msg)
	return bf.String()
}

func outputDirect(w WriterCloser, logLevel string, data any) {
	if w == nil {
		log.Println(lang.ToString(data))
	} else {
		_ = w.Writeln(lang.ToString(data))
	}
}

func outputDirectString(w WriterCloser, str string) {
	if w == nil {
		log.Println(str)
	} else {
		_ = w.Writeln(str)
	}
}

// 不推荐使用 bytes 版本
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
