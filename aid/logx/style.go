// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

import (
	"errors"
	"github.com/qinchende/gofast/core/cst"
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

const (
	//timeFormatMini = cst.TimeFmtMdHms
	timeFormat = cst.TimeFmtYmdHms
)

// 将名称字符串转换成整数类型，提高判断性能s
func (l *Logger) initStyle() error {
	switch l.cnf.LogStyle {

	case styleSdxStr:
		l.iStyle = StyleSdx
		l.FnLogBegin = SdxBegin
		l.FnLogEnd = SdxEnd
		l.FnGroupBegin = SdxGroupBegin
		l.FnGroupEnd = SdxGroupEnd
	case styleCdoStr:
		l.iStyle = StyleCdo
		l.FnLogBegin = JsonBegin
		l.FnLogEnd = JsonEnd
		l.FnGroupBegin = JsonGroupBegin
		l.FnGroupEnd = JsonGroupEnd
	case styleJsonStr:
		l.iStyle = StyleJson
		l.FnLogBegin = JsonBegin
		l.FnLogEnd = JsonEnd
		l.FnGroupBegin = JsonGroupBegin
		l.FnGroupEnd = JsonGroupEnd
	case styleCustomStr:
		l.iStyle = StyleCustom
	default:
		return errors.New("item LogStyle not match")
	}
	return nil
}

// Sdx-style
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func SdxBegin(r *Record, label string) {
	bf := append(time.Now().AppendFormat(r.bs, timeFormat), " ["...)
	r.bs = append(append(bf, label...), "]: {"...)
}

func SdxEnd(r *Record) []byte {
	bf := r.bs
	if bf[len(bf)-1] == ',' {
		bf = bf[:len(bf)-1]
	}
	return append(bf, "}\n"...)
}

func SdxGroupBegin(bs []byte, k string) []byte {
	bs = append(bs, "\n    \""...)
	bs = jde.AppendStrNoQuotes(bs, k)
	bs = append(bs, "\": {"...)
	return bs
}

func SdxGroupEnd(bs []byte) []byte {
	if len(bs) > 0 && bs[len(bs)-1] == ',' {
		bs = bs[:len(bs)-1]
	}
	bs = append(bs, "},"...)
	return bs
}

// JSON-style
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func JsonBegin(r *Record, label string) {
	r.bs = jde.AppendStrField(jde.AppendTimeField(append(r.bs, '{'), fTimeStamp, time.Now(), timeFormat), fLabel, label)
}

func JsonEnd(r *Record) []byte {
	bf := r.bs
	if bf[len(bf)-1] == ',' {
		bf = bf[:len(bf)-1]
	}
	return append(bf, "}\n"...)
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
