// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"errors"
)

const (
	timeFormat     = "2006-01-02 15:04:05"
	timeFormatMini = "01-02 15:04:05"
)

// NOTE：系统内置了几个系列的请求日志输出样式，比如custom、sdx、elk、prometheus等，当然你也可以自定义输出样式

// 日志样式类型
const (
	LogStyleCustom int8 = iota
	LogStyleSdx
	LogStyleELK
	LogStyleJson
)

// 日志样式名称
const (
	styleSdxStr    = "sdx"
	styleJsonStr   = "json"
	styleELKStr    = "elk"
	styleCustomStr = "custom" // 自定义
)

var Formatter func(w WriterCloser, level string, data any)

// 将名称字符串转换成整数类型，提高判断性能
func initStyle(c *LogConfig) error {
	switch c.LogStyle {

	case styleSdxStr:
		c.iStyle = LogStyleSdx
		Formatter = outputSdxStyle
	case styleELKStr:
		c.iStyle = LogStyleELK
		Formatter = outputElkStyle
	case styleJsonStr:
		c.iStyle = LogStyleJson
		Formatter = outputPrometheusStyle
	case styleCustomStr:
		c.iStyle = LogStyleCustom
		Formatter = outputCustomStyle
	default:
		Formatter = outputDirect
		return errors.New("item LogStyle not match")
	}
	return nil
}

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
//			switch myCnf.iStyle {
//			case LogStyleCustom:
//				outputCustomStyle(w, logLevel, data)
//			case LogStyleSdx:
//				outputSdxStyle(w, logLevel, data)
//			case LogStyleSdxJson:
//				outputSdxJsonStyle(w, logLevel, data)
//			case LogStyleELK:
//				outputElkStyle(w, logLevel, data)
//			case LogStyleJson:
//				outputPrometheusStyle(w, logLevel, data)
//			default:
//				outputDirectString(w, lang.ToString(data))
//			}
//		} else {
//			outputDirectString(w, lang.ToString(data))
//		}
//	}
func output(w WriterCloser, logLevel string, data any) {
	Formatter(w, logLevel, data)
}

// 打印请求日志，可以指定不同的输出样式
func RequestsLog(p *ReqLogEntity, flag int8) {
	switch myCnf.iStyle {
	case LogStyleSdx:
		InfoDirect(buildSdxReqLog(p, flag))
	case LogStyleELK:
		InfoDirect(buildElkReqLog(p, flag))
	case LogStyleJson:
		InfoDirect(buildPrometheusReqLog(p, flag))
	case LogStyleCustom:
		InfoDirect(buildCustomReqLog(p, flag))
	default:
	}
}
