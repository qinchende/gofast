// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/aid/bag"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/fst"
	"io"
	"net/http"
	"sync"
	"time"
)

func LoggerDemo(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	r := logx.Info()

	r.Int("Latency", int(timex.SdxNowDur()-c.EnterTime))
	r.Str("RemoteAddr", c.ClientIP())
	r.Int("BodySize", len(c.Res.WrittenData()))
	r.Str("RemoteAddr", c.ClientIP())

	r.Send()
}

func Logger(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := InfoReq()

	p.RawReq = c.Req.Raw
	p.Pms = c.Pms
	p.RemoteAddr = c.ClientIP()
	p.StatusCode = c.Res.Status()
	p.ResData = c.Res.WrittenData()
	p.BodySize = len(p.ResData)

	// 内部错误信息一般不返回给调用者，但是需要打印日志信息
	p.CarryItems = c.CarryMsgItems()

	// Stop timer
	p.TimeStamp = timex.SdxNowDur()
	p.Latency = p.TimeStamp - c.EnterTime

	p.Send()
}

func LoggerMini(c *fst.Context) {
	// 执行完后面的请求，再打印日志
	c.Next()

	// 请求处理完，并成功返回了，接下来就是打印请求日志
	p := InfoReq()

	p.RawReq = c.Req.Raw
	p.RemoteAddr = c.ClientIP()
	p.StatusCode = c.Res.Status()
	p.ResData = c.Res.WrittenData()
	p.BodySize = len(p.ResData)

	// Stop timer
	p.TimeStamp = timex.SdxNowDur()
	p.Latency = p.TimeStamp - c.EnterTime

	p.Send()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 专为收集请求日志而设计
type ReqRecord struct {
	logx.Record

	RawReq     *http.Request
	StatusCode int
	Method     string
	RequestURI string
	UserAgent  string
	RemoteAddr string
	TimeStamp  time.Duration
	Latency    time.Duration
	Pms        cst.SuperKV
	BodySize   int
	ResData    []byte
	CarryItems bag.CarryList
}

var (
	reqRecordPool = &sync.Pool{
		New: func() interface{} {
			r := &ReqRecord{}
			//r.Record.init()
			return r
		},
	}
	_reqRecordDefValue ReqRecord
)

func getReqRecordFromPool() *ReqRecord {
	r := reqRecordPool.Get().(*ReqRecord)
	//r.bf = pool.GetBytes()
	//r.bs = *r.bf
	return r
}

func putReqRecordToPool(r *ReqRecord) {
	//*r.bf = r.bs
	//pool.FreeBytes(r.bf)
	//r.bf = nil
	//r.bs = nil
	reqRecordPool.Put(r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func newReqRecord(w io.WriteCloser, label string) *ReqRecord {
	r := getReqRecordFromPool()

	// Record记录的数据
	r.Time = timex.NowDur()
	r.Label = label
	//r.iow = w
	//r.out = r

	return r
}

func InfoReq() *ReqRecord {
	if logx.DefLogger.ShowInfo() {
		return newReqRecord(logx.DefLogger.WReq, logx.LabelReq)
	}
	return nil
}

func (r *ReqRecord) output(msg string) {

	//r.bs = jde.AppendStrField(r.bs, "msg", msg)
	//r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
	//r.bs = append(r.bs, "\n"...)
	//
	//if _, err := r.iow.Write(r.bs); err != nil {
	//	_, _ = fmt.Fprintf(os.Stderr, "logx: write req-record error: %s\n", err.Error())
	//}
	putReqRecordToPool(r)
}
