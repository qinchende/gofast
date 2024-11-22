// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"fmt"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/fst"
	"time"
)

func LoggerDemo(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	r := logx.Info()
	r.Int("Latency", int64(timex.SdxNowDur()-c.EnterTime))
	r.Str("RemoteAddr", c.ClientIP())
	r.Int("BodySize", int64(len(c.Res.WrittenData())))
	r.Str("RemoteAddr", c.ClientIP())
	r.Send()
}

//[GET] /admin/sdx (127.0.0.1/08-23 15:41:07) [200/63/0]
//B: {}
//P: {"nowTime":"2024-08-22T15:41:07+08:00"}
//R: {"status":"fai","code":0,"msg":"","data":{"data":"handle sdx"}}

func Logger(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	r := logx.InfoReq().
		Str("Method", c.Req.Raw.Method).
		Str("Url", c.Req.Raw.URL.RawPath).
		Str("IP", c.ClientIP()).
		Str("Mark", fmt.Sprintf("%d/%d/%d", c.Res.Status(), len(c.Res.WrittenData()), (timex.SdxNowDur()-c.EnterTime)/time.Millisecond))

	r.Group("B").GEnd()
	r.Group("P").GEnd()
	r.Group("R").GEnd()

	r.Send()
}

func LoggerMini(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	r := logx.InfoReq().
		Str("Method", c.Req.Raw.Method).
		Str("Url", c.Req.Raw.URL.RawPath).
		Str("IP", c.ClientIP()).
		Str("Mark", fmt.Sprintf("%d/%d/%d", c.Res.Status(), len(c.Res.WrittenData()), (timex.SdxNowDur()-c.EnterTime)/time.Millisecond))

	r.Group("B").GEnd()
	r.Group("P").GEnd()
	r.Group("R").GEnd()

	r.Send()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 专为收集请求日志而设计
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

//var (
//	reqRecordPool = &sync.Pool{
//		New: func() interface{} {
//			r := &ReqRecord{}
//			//r.Record.init()
//			return r
//		},
//	}
//	_reqRecordDefValue ReqRecord
//)
//
//func getReqRecordFromPool() *ReqRecord {
//	r := reqRecordPool.Get().(*ReqRecord)
//	//r.bf = pool.GetBytes()
//	//r.bs = *r.bf
//	return r
//}
//
//func putReqRecordToPool(r *ReqRecord) {
//	//*r.bf = r.bs
//	//pool.FreeBytes(r.bf)
//	//r.bf = nil
//	//r.bs = nil
//	reqRecordPool.Put(r)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func newReqRecord(w io.WriteCloser, label string) *ReqRecord {
//	r := getReqRecordFromPool()
//
//	// Record记录的数据
//	r.Time = timex.NowDur()
//	r.Label = label
//	//r.iow = w
//	//r.out = r
//
//	return r
//}
//
//func InfoReq() *ReqRecord {
//	if logx.Def.ShowInfo() {
//		return newReqRecord(logx.Def.WReq, logx.LabelReq)
//	}
//	return nil
//}
//
//func (r *ReqRecord) output(msg string) {
//
//	//r.bs = jde.AppendStrField(r.bs, "msg", msg)
//	//r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
//	//r.bs = append(r.bs, "\n"...)
//	//
//	//if _, err := r.iow.Write(r.bs); err != nil {
//	//	_, _ = fmt.Fprintf(os.Stderr, "logx: write req-record error: %s\n", err.Error())
//	//}
//	putReqRecordToPool(r)
//}
