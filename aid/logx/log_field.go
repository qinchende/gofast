// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

//// 专为收集请求日志而设计
//type ReqRecord struct {
//	Record
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
//			r.Record.init()
//			return r
//		},
//	}
//	_reqRecordDefValue ReqRecord
//)
//
//func getReqRecordFromPool() *ReqRecord {
//	r := reqRecordPool.Get().(*ReqRecord)
//	r.buf = pool.GetBytes()
//	r.bs = *r.buf
//	return r
//}
//
//func putReqRecordToPool(r *ReqRecord) {
//	*r.buf = r.bs
//	pool.FreeBytes(r.buf)
//	r.buf = nil
//	r.bs = nil
//	reqRecordPool.Put(r)
//}
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func newReqRecord(w io.Writer, label string) *ReqRecord {
//	r := getReqRecordFromPool()
//
//	// Record记录的数据
//	r.Time = timex.NowDur()
//	r.Label = label
//	r.iow = w
//	r.out = r
//
//	return r
//}
//
//func InfoReq() *ReqRecord {
//	if Def.ShowInfo() {
//		return newReqRecord(Def.WReq, LabelReq)
//	}
//	return nil
//}
//
//func (r *ReqRecord) Output(msg string) {
//
//	r.bs = jde.AppendStrField(r.bs, "msg", msg)
//	r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
//	r.bs = append(r.bs, "\n"...)
//
//	if _, err := r.iow.Write(r.bs); err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, "logx: write req-record error: %s\n", err.Error())
//	}
//	putReqRecordToPool(r)
//}
