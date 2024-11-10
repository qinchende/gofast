// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/aid/bag"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/store/jde"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type ReqRecord struct {
	Record

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

func (r *ReqRecord) Close() error {
	return nil
}

var (
	reqRecordPool = &sync.Pool{
		New: func() interface{} {
			r := &ReqRecord{}
			r.Record.init()
			return r
		},
	}
	_reqRecordDefValue ReqRecord
)

func getReqRecordFromPool() *ReqRecord {
	r := reqRecordPool.Get().(*ReqRecord)
	r.bf = pool.GetBytes()
	r.bs = *r.bf
	return r
}

func putReqRecordToPool(r *ReqRecord) {
	*r.bf = r.bs
	pool.FreeBytes(r.bf)
	r.bf = nil
	r.bs = nil
	reqRecordPool.Put(r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func newReqRecord(w io.WriteCloser, label string) *ReqRecord {
	r := getReqRecordFromPool()

	// Record记录的数据
	r.Time = timex.NowDur()
	r.Label = label
	r.w = w
	r.out = r

	return r
}

func InfoReq() *ReqRecord {
	if ShowInfo() {
		return newReqRecord(ioReq, labelReq)
	}
	return nil
}

func (r *ReqRecord) Output(msg string) {

	r.bs = jde.AppendStrField(r.bs, "msg", msg)
	r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
	r.bs = append(r.bs, "\n"...)

	if _, err := r.w.Write(r.bs); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logx: write req-record error: %s\n", err.Error())
	}
	putReqRecordToPool(r)
}
