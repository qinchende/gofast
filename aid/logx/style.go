// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"errors"
	"io"
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

var (
	Formatter        func(w io.WriteCloser, level string, data any)
	CustomOutputFunc func(logLevel string, data any) string
	RequestsLog      func(p *ReqRecord, flag int8) string
	TimeToStr        func(tm time.Time) string
	WriteRecord      func(r *Record, flag int8) string
	WriteReqRecord   func(r *ReqRecord, flag int8) string
)

// 将名称字符串转换成整数类型，提高判断性能
func (l *Logger) initStyle() error {
	switch l.cnf.LogStyle {

	case styleSdxStr:
		l.iStyle = StyleSdx
		{
			Formatter = outputSdxStyle
			RequestsLog = buildSdxReqLog
		}
	case styleCdoStr:
		l.iStyle = StyleCdo
		//Formatter =
	case styleJsonStr:
		l.iStyle = StyleJson
		Formatter = outputPrometheusStyle
	case styleCustomStr:
		l.iStyle = StyleCustom
		//Formatter = outputCustomStyle
	default:
		//Formatter = outputDirect
		return errors.New("item LogStyle not match")
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Output(w io.WriteCloser, logLevel string, data any) {
	Formatter(w, logLevel, data)
}

//// 打印请求日志，可以指定不同的输出样式
//func RequestsLog(p *ReqLogEntity, flag int8) {
//	switch cnf.iStyle {
//	case StyleSdx:
//		InfoDirect()
//	case StyleCdo:
//		InfoDirect(buildElkReqLog(p, flag))
//	case StyleJson:
//		InfoDirect(buildPrometheusReqLog(p, flag))
//	case StyleCustom:
//		InfoDirect(buildCustomReqLog(p, flag))
//	default:
//	}
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 日志的输出，最后都要到这个方法进行输出
//
//	func output(w WriterCloser, logLevel string, data any, useStyle bool) {
//		// 自定义了 sdx 这种输出样式，否则就是默认的 json 样式
//		//log.SetPrefix("[GoFast]")    // 前置字符串加上特定标记
//		//log.SetFlags(log.Lmsgprefix) // 取消前置字符串
//		//log.SetFlags(log.LstdFlags)  // 设置成日期+时间 格式
//
//		if useStyle == true {
//			switch cnf.iStyle {
//			case StyleCustom:
//				outputCustomStyle(w, logLevel, data)
//			case StyleSdx:
//				outputSdxStyle(w, logLevel, data)
//			case LogStyleSdxJson:
//				outputSdxJsonStyle(w, logLevel, data)
//			case LogStyleELK:
//				outputElkStyle(w, logLevel, data)
//			case StyleJson:
//				outputPrometheusStyle(w, logLevel, data)
//			default:
//				outputDirectString(w, lang.ToString(data))
//			}
//		} else {
//			outputDirectString(w, lang.ToString(data))
//		}
//	}
