// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/fst"
	"strconv"
	"time"
)

// [GET] /admin/sdx (127.0.0.1/08-23 15:41:07) [200/63/0]
// P: {"nowTime":"2024-08-22T15:41:07+08:00"}
// R: {"status":"fai","code":0,"msg":"","data":{"data":"handle sdx"}}
// E: {}
func Logger(c *fst.Context) {
	c.Next() // 执行完后面的请求，再打印日志

	r := logx.InfoReq()
	appendReqBaseInfo(r, c)

	r.Group("P").Any("pms", c.Pms).GroupEnd()
	r.Group("R").Json("res", c.Res.WrittenData()).GroupEnd()
	r.Group("E").GroupEnd() // carry items

	r.Send()
}

// [GET] /admin/sdx (127.0.0.1/08-23 15:41:07) [200/63/0]
func LoggerMini(c *fst.Context) {
	c.Next()

	r := logx.InfoReq()
	appendReqBaseInfo(r, c)
	r.Send()
}

func appendReqBaseInfo(r *logx.Record, c *fst.Context) {
	bf := r.GetBuf()

	//// 请求处理完，并成功返回了，接下来就是打印请求日志
	//r := logx.InfoReq().
	//	Str("Method", c.Req.Raw.Method).
	//	Str("Url", c.Req.Raw.URL.RawPath).
	//	Str("IP", c.ClientIP()).
	//	Str("Mark", fmt.Sprintf("%d/%d/%d", c.Res.Status(), len(c.Res.WrittenData()), timeMillSeconds)
	// pbs := pool.GetBytesMin(256)
	// bf := *pbs
	bf = append(bf, "\"A\":\"["...)
	bf = append(bf, c.Req.Raw.Method...)
	bf = append(bf, "] "...)
	bf = append(bf, c.Req.Raw.URL.Path...) // 或者 jde.AppendStrNoQuotes()
	bf = append(bf, " ("...)
	bf = append(bf, c.ClientIP()...)
	bf = append(bf, ") ["...)
	bf = append(bf, strconv.Itoa(c.Res.Status())...)
	bf = append(bf, '/')
	bf = append(bf, strconv.Itoa(len(c.Res.WrittenData()))...)
	bf = append(bf, '/')
	bf = append(bf, strconv.Itoa(int((timex.SdxNowDur()-c.EnterTime)/time.Millisecond))...)
	bf = append(bf, "]\","...)
	// r.Str("A", lang.B2S(bf))
	// pool.FreeBytes(pbs)

	r.SetBuf(bf)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 通过模板构造字符串可能性能更好。
//func buildSdxReqLog(p *ReqRecord, flag int8) string {
//	// 需要用Mini版本
//	if flag > 0 {
//		return buildSdxReqLogMini(p)
//	}
//
//	formatStr := `
//[%s] %s (%s/%s) [%d/%d/%d]
//  B: %s
//  P: %s
//  R: %s%s
//`
//	// 最长打印出 1024个字节的结果
//	tLen := p.BodySize
//	if tLen > 1024 {
//		tLen = 1024
//	}
//
//	// 这个时候可以随意改变 p.Pms ，这是请求最后一个执行的地方了
//	reqParams := []byte("{}")
//	reqBaseParams := []byte("{}")
//
//	// 当熔断降载的时候，还没有进入c.Pms的处理逻辑，c.Pms为nil
//	if p.Pms != nil {
//		// 1. 请求核心参数
//		var basePms = make(cst.KV)
//		if tok, ok := p.Pms.Get("tok"); ok {
//			basePms["tok"] = tok
//			p.Pms.Del("tok")
//		}
//		reqBaseParams, _ = jsonx.Marshal(basePms)
//
//		// 2. 请求的其它参数
//		reqParams, _ = jsonx.Marshal(p.Pms)
//	} else if p.RawReq.Form != nil {
//		reqParams, _ = jsonx.Marshal(p.RawReq.Form)
//	}
//
//	return fmt.Sprintf(formatStr,
//		p.RawReq.Method,
//		p.RawReq.URL.Path,
//		p.RemoteAddr,
//		timex.ToTime(p.TimeStamp).Format(timeFormatMini),
//		p.StatusCode,
//		p.BodySize,
//		p.Latency/time.Millisecond,
//		reqBaseParams,
//		reqParams,
//		(p.ResData)[:tLen],
//		buildCarryInfos(p.CarryItems),
//	)
//}
//
//func buildSdxReqLogMini(p *ReqRecord) string {
//	return ""
//	//	formatStr := `
//	//[%s] %s (%s/%s) [%d/%d/%d] %s
//	//`
//	//	// 最长打印出 1024个字节的结果
//	//	tLen := p.BodySize
//	//	if tLen > 1024 {
//	//		tLen = 1024
//	//	}
//	//
//	//	return fmt.Sprintf(formatStr,
//	//		p.RawReq.Method,
//	//		p.RawReq.URL.Path,
//	//		p.ClientIP,
//	//		timex.ToTime(p.TimeStamp).Format(timeFormatMini),
//	//		p.StatusCode,
//	//		p.BodySize,
//	//		p.Latency/time.Millisecond,
//	//		(p.ResData)[:tLen],
//	//	)
//}
//
//// 所有错误合并成字符串
//func buildCarryInfos(bs bag.CarryList) string {
//	if len(bs) == 0 {
//		return ""
//	}
//
//	var buf strings.Builder
//	buf.Grow(len(bs[0].Msg) + 10)
//
//	buf.WriteString("\n  E: ")
//	infos := bs.CollectMessages()
//	for i, str := range infos {
//		if i != 0 {
//			buf.WriteString("\n     ")
//		}
//		buf.WriteString(strconv.Itoa(i))
//		buf.WriteString(". ")
//		buf.WriteString(str)
//	}
//	return buf.String()
//}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//// 专为收集请求日志而设计
//type ReqRecord struct {
//	logx.Record
//
//	RawReq     *http.Request
//	StatusCode int
//	Method     string
//	RequestURI string
//	UserAgent  string
//	RemoteAddr string
//	TimeStamp  time.Duration
//	Latency    time.Duration
//	Pms        cst.SuperKV
//	BodySize   int
//	ResData    []byte
//	CarryItems bag.CarryList
//}
//
//var (
//	reqRecordPool = &sync.Pool{
//		New: func() interface{} {
//			r := &ReqRecord{}
//			//r.Record.init()
//			return r
//		},
//	}
//)
//
//func getReqRecordFromPool() *ReqRecord {
//	r := reqRecordPool.Get().(*ReqRecord)
//
//	return r
//}
//
//func putReqRecordToPool(r *ReqRecord) {
//
//	reqRecordPool.Put(r)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func newReqRecord(l *logx.Logger, w io.Writer, label string) *ReqRecord {
//	r := getReqRecordFromPool()
//	return r
//}
//
//func (r *ReqRecord) write() {
//	putReqRecordToPool(r)
//}
