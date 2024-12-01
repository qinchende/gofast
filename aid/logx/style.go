// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

import (
	"errors"
	"github.com/qinchende/gofast/store/jde"
	"time"
)

// NOTE：系统内置了几个系列的请求日志输出样式，当然你也可以自定义输出样式
// 当前内置三种：sdx-style、json-style、cdo-style
const (
	StyleSdx int8 = iota
	StyleJson
	StyleCdo
	StyleCustom
)

// 日志样式名称
const (
	styleSdxStr    = "sdx"
	styleJsonStr   = "json"
	styleCdoStr    = "cdo"
	styleCustomStr = "custom"
)

// 将名称字符串转换成整数类型，提高判断性能s
func (l *Logger) initStyle() error {
	switch l.cnf.LogStyle {
	case styleSdxStr:
		l.iStyle = StyleSdx
		l.LogBegin = SdxBegin
		l.LogEnd = SdxEnd
		l.GroupBegin = SdxGroupBegin
		l.GroupEnd = SdxGroupEnd
	case styleJsonStr:
		l.iStyle = StyleJson
		l.LogBegin = JsonBegin
		l.LogEnd = JsonEnd
		l.GroupBegin = JsonGroupBegin
		l.GroupEnd = JsonGroupEnd
	case styleCdoStr: // TODO: need to impl
		l.iStyle = StyleCdo
		l.LogBegin = JsonBegin
		l.LogEnd = JsonEnd
		l.GroupBegin = JsonGroupBegin
		l.GroupEnd = JsonGroupEnd
	case styleCustomStr:
		l.iStyle = StyleCustom
		l.LogBegin = CustomBegin
		l.LogEnd = CustomEnd
		l.GroupBegin = CustomGroupBegin
		l.GroupEnd = CustomGroupEnd
	default:
		return errors.New("item LogStyle not match")
	}
	return nil
}

// Sdx-style
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func SdxBegin(bs []byte, label string) []byte {
	//bs = append(time.Now().AppendFormat(append(bs, '\n'), timeFormat), " ["...)
	bs = append(time.Now().AppendFormat(bs, timeFormat), " ["...)
	return append(append(bs, label...), "]: "...)
}

func SdxEnd(bs []byte) []byte {
	if bs[len(bs)-1] == ',' {
		bs = bs[:len(bs)-1]
	}
	return append(bs, "\n"...)
}

func SdxGroupBegin(bs []byte, k string) []byte {
	return append(jde.AppendStrNoQuotes(append(bs, "\n  "...), k), ": {"...)
}

func SdxGroupEnd(bs []byte) []byte {
	if len(bs) > 0 && bs[len(bs)-1] == ',' {
		bs = bs[:len(bs)-1]
	}
	return append(bs, "},"...)
}

// JSON-style
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func JsonBegin(bs []byte, label string) []byte {
	bs = jde.AppendTimeField(append(bs, '{'), fTimeStamp, time.Now(), timeFormat)
	return jde.AppendStrField(bs, fLabel, label)
}

func JsonEnd(bs []byte) []byte {
	if bs[len(bs)-1] == ',' {
		bs = bs[:len(bs)-1]
	}
	return append(bs, "}\n"...)
}

func JsonGroupBegin(bs []byte, k string) []byte {
	return append(jde.AppendKey(bs, k), '{')
}

func JsonGroupEnd(bs []byte) []byte {
	if len(bs) > 0 && bs[len(bs)-1] == ',' {
		bs = bs[:len(bs)-1]
	}
	return append(bs, "},"...)
}

// Custom-style
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func CustomBegin(bs []byte, label string) []byte {
	return bs
}

func CustomEnd(bs []byte) []byte {
	return bs
}

func CustomGroupBegin(bs []byte, k string) []byte {
	return bs
}

func CustomGroupEnd(bs []byte) []byte {
	return bs
}
