// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

import (
	"errors"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/store/jde"
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

const (
	timeFormat     = "2006-01-02 15:04:05"
	timeFormatMini = "01-02 15:04:05"
)

// 将名称字符串转换成整数类型，提高判断性能
func (l *Logger) initStyle() error {
	switch l.cnf.LogStyle {

	case styleSdxStr:
		l.iStyle = StyleSdx
		l.StyleSummary = SdxSummary
		l.StyleGroupBegin = SdxGroupBegin
		l.StyleGroupEnd = SdxGroupEnd
	case styleCdoStr:
		l.iStyle = StyleCdo
		l.StyleSummary = SdxSummary
		l.StyleGroupBegin = SdxGroupBegin
		l.StyleGroupEnd = SdxGroupEnd
	case styleJsonStr:
		l.iStyle = StyleJson
		l.StyleSummary = SdxSummary
		l.StyleGroupBegin = SdxGroupBegin
		l.StyleGroupEnd = SdxGroupEnd
	case styleCustomStr:
		l.iStyle = StyleCustom
	default:
		return errors.New("item LogStyle not match")
	}
	return nil
}

// Sdx-style
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func SdxSummary(r *Record) []byte {
	bs := r.bs

	// 每条日志的第一行样式
	bf := bs[len(bs):]
	bf = timex.ToTime(r.Time).AppendFormat(bf, timeFormat)
	bf = append(bf, " ["...)
	bf = append(bf, r.Label...)
	bf = append(bf, "]: {"...)
	bf = append(bf, r.log.bs...)
	//bf = append(bf, "}\n"...)

	bf = append(bf, bs...)
	bf = append(bf, "}\n"...)
	return bf
}

func SdxGroupBegin(bs []byte, k string) []byte {
	if len(bs) > 0 {
		bs = append(bs, "}\n    "...)
	}
	bs = jde.AppendStr(bs, k)
	bs = append(bs, ": {"...)
	return bs
}

func SdxGroupEnd(bs []byte) []byte {
	bs = append(bs, "}"...)
	return bs
}

// JSON-style
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
