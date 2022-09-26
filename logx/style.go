package logx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/tools"
	"net/http"
	"time"
)

const (
	timeFormat     = "2006-01-02 15:04:05"
	timeFormatMini = "01-02 15:04:05"
)

// 日志样式类型
const (
	LogStyleJson int8 = iota
	LogStyleJsonMini
	LogStyleSdx
	LogStyleSdxMini
)

// 日志样式名称
const (
	styleJsonStr     = "json"
	styleJsonMiniStr = "json-mini"
	styleSdxStr      = "sdx"
	styleSdxMiniStr  = "sdx-mini"
)

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
	MsgBaskets tools.Baskets
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 打印请求日志，可以指定不同的输出样式
func PrintReqLog(p *ReqLogEntity) {
	switch myCnf.logStyle {
	case LogStyleSdx:
		InfoDirect(buildSdxReqLog(p))
	case LogStyleSdxMini:
		InfoDirect(buildSdxReqLogMini(p))
	case LogStyleJson:
		InfoDirect(buildSdxReqLog(p))
	case LogStyleJsonMini:
		InfoDirect(buildSdxReqLog(p))
	default:
	}
}

func PrintReqLogMini(p *ReqLogEntity) {
	InfoDirect(buildSdxReqLogMini(p))
}

// 日志的输出，最后都要到这个方法进行输出
func output(w WriterCloser, info string, logLevel string, useStyle bool) {
	// 自定义了 sdx 这种输出样式，否则就是默认的 json 样式
	//log.SetPrefix("[GoFast]")    // 前置字符串加上特定标记
	//log.SetFlags(log.Lmsgprefix) // 取消前置字符串
	//log.SetFlags(log.LstdFlags)  // 设置成日期+时间 格式

	if useStyle == true {
		switch myCnf.logStyle {
		case LogStyleSdx:
			outputSdxStyle(w, info, logLevel)
		case LogStyleSdxMini:
			outputSdxStyle(w, info, logLevel)
		case LogStyleJson:
			outputJsonStyle(w, info, logLevel)
		case LogStyleJsonMini:
			outputJsonStyle(w, info, logLevel)
		default:
			outputDirectString(w, info)
		}
	} else {
		outputDirectString(w, info)
	}
}

func initStyle(c *LogConfig) error {
	switch c.LogStyle {
	case styleSdxStr:
		c.logStyle = LogStyleSdx
	case styleSdxMiniStr:
		c.logStyle = LogStyleSdxMini
	case styleJsonMiniStr:
		c.logStyle = LogStyleJsonMini
	case styleJsonStr:
		c.logStyle = LogStyleJson
	default:
		return errors.New("item LogStyle not match")
	}
	return nil
}

//func StyleType() int8 {
//	return myCnf.logStyle
//}
