// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"fmt"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/core/lang"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/fst"
	"strconv"
	"time"
)

// [GET] /admin/sdx (127.0.0.1/08-23 15:41:07) [200/63/0]
// B: {}
// P: {"nowTime":"2024-08-22T15:41:07+08:00"}
// R: {"status":"fai","code":0,"msg":"","data":{"data":"handle sdx"}}
func Logger(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	r := logx.InfoReq()
	//bs := r.GetBuf()

	pbs := pool.GetBytesMin(256)
	bf := *pbs
	bf = append(bf, "["...)
	bf = append(bf, c.Req.Raw.Method...)
	bf = append(bf, "] "...)
	bf = append(bf, c.Req.Raw.URL.Path...)
	bf = append(bf, " ("...)
	bf = append(bf, c.ClientIP()...)
	bf = append(bf, ") ["...)
	bf = append(bf, strconv.Itoa(c.Res.Status())...)
	bf = append(bf, '/')
	bf = append(bf, strconv.Itoa(len(c.Res.WrittenData()))...)
	bf = append(bf, '/')
	bf = append(bf, strconv.Itoa(int((timex.SdxNowDur()-c.EnterTime)/time.Millisecond))...)
	bf = append(bf, ']')
	r.Str("A", lang.B2S(bf))
	pool.FreeBytes(pbs)

	//r.SetBuf(bf)

	//// 请求处理完，并成功返回了，接下来就是打印请求日志
	//r := logx.InfoReq().
	//	Str("Method", c.Req.Raw.Method).
	//	Str("Url", c.Req.Raw.URL.RawPath).
	//	Str("IP", c.ClientIP()).
	//	Str("Mark", fmt.Sprintf("%d/%d/%d", c.Res.Status(), len(c.Res.WrittenData()), (timex.SdxNowDur()-c.EnterTime)/time.Millisecond))

	r.Group("B").GroupEnd()
	r.Group("P").GroupEnd()
	r.Group("R").Append(c.Res.WrittenData()).GroupEnd()

	r.Send()
}

func LoggerSdxMini() fst.CtxHandler {
	return func(c *fst.Context) {
	}
}

func LoggerSdxMini2(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	r := logx.InfoReq().
		Str("Method", c.Req.Raw.Method).
		Str("Url", c.Req.Raw.URL.RawPath).
		Str("IP", c.ClientIP()).
		Str("Mark", fmt.Sprintf("%d/%d/%d", c.Res.Status(), len(c.Res.WrittenData()), (timex.SdxNowDur()-c.EnterTime)/time.Millisecond))

	r.Group("B").GroupEnd()
	r.Group("P").GroupEnd()
	r.Group("R").GroupEnd()

	r.Send()
}

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
