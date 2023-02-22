// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/tips"
	"github.com/qinchende/gofast/skill/lang"
	"net/http"
	"time"
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
	LogStyleSdxJson
	LogStyleELK
	LogStylePrometheus
)

// 日志样式名称
const (
	styleCustomStr     = "custom" // 自定义
	styleSdxStr        = "sdx"
	styleSdxJson       = "sdx-json"
	styleELKStr        = "elk"
	stylePrometheusStr = "prometheus"
)

// 将名称字符串转换成整数类型，提高判断性能
func initStyle(c *LogConfig) error {
	switch c.LogStyle {
	case styleCustomStr:
		c.logStyleInt8 = LogStyleCustom
	case styleSdxStr:
		c.logStyleInt8 = LogStyleSdx
	case styleSdxJson:
		c.logStyleInt8 = LogStyleSdxJson
	case styleELKStr:
		c.logStyleInt8 = LogStyleELK
	case stylePrometheusStr:
		c.logStyleInt8 = LogStylePrometheus
	default:
		return errors.New("item LogStyle not match")
	}
	return nil
}

// 日志参数实体
type ReqLogEntity struct {
	RawReq     *http.Request
	TimeStamp  time.Duration
	Latency    time.Duration
	ClientIP   string
	StatusCode int
	Pms        cst.KV
	BodySize   int
	ResData    []byte
	CarryItems tips.CarryList
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 日志的输出，最后都要到这个方法进行输出
func output(w WriterCloser, logLevel string, data any, useStyle bool) {
	// 自定义了 sdx 这种输出样式，否则就是默认的 json 样式
	//log.SetPrefix("[GoFast]")    // 前置字符串加上特定标记
	//log.SetFlags(log.Lmsgprefix) // 取消前置字符串
	//log.SetFlags(log.LstdFlags)  // 设置成日期+时间 格式

	if useStyle == true {
		switch myCnf.logStyleInt8 {
		case LogStyleCustom:
			outputCustomStyle(w, logLevel, data)
		case LogStyleSdx:
			outputSdxStyle(w, logLevel, data)
		case LogStyleSdxJson:
			outputSdxJsonStyle(w, logLevel, data)
		case LogStyleELK:
			outputElkStyle(w, logLevel, data)
		case LogStylePrometheus:
			outputPrometheusStyle(w, logLevel, data)
		default:
			outputDirectString(w, lang.ToString(data))
		}
	} else {
		outputDirectString(w, lang.ToString(data))
	}
}

// 打印请求日志，可以指定不同的输出样式
func RequestsLog(p *ReqLogEntity, flag int8) {
	switch myCnf.logStyleInt8 {
	case LogStyleCustom:
		InfoDirect(buildCustomReqLog(p, flag))
	case LogStyleSdx:
		InfoDirect(buildSdxReqLog(p, flag))
	case LogStyleSdxJson:
		InfoDirect(buildSdxReqLog(p, flag))
	case LogStyleELK:
		InfoDirect(buildElkReqLog(p, flag))
	case LogStylePrometheus:
		InfoDirect(buildPrometheusReqLog(p, flag))
	default:
	}
}
