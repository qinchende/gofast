// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/aid/bag"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/pool"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// 专为收集请求日志而设计
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

var (
	reqRecordPool = &sync.Pool{
		New: func() interface{} {
			r := &ReqRecord{}
			//r.Record.init()
			return r
		},
	}
)

func getReqRecordFromPool() *ReqRecord {
	r := reqRecordPool.Get().(*ReqRecord)
	r.pBuf = pool.GetBytes()
	r.bs = *r.pBuf
	return r
}

func putReqRecordToPool(r *ReqRecord) {
	*r.pBuf = r.bs
	pool.FreeBytes(r.pBuf)
	r.pBuf = nil
	r.bs = nil
	reqRecordPool.Put(r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func newReqRecord(l *Logger, w io.Writer, label string) *ReqRecord {
	r := getReqRecordFromPool()

	r.myL = l
	r.iow = w
	r.out = r
	r.bs = append(l.FnLogBegin(r.bs, label), l.r.bs...)

	return r
}

func InfoReqX() *ReqRecord {
	if Def.ShowInfo() {
		return newReqRecord(Def, Def.WReq, LabelReq)
	}
	return nil
}

func (r *ReqRecord) write() {
	data := r.myL.FnLogEnd(r.bs)
	if _, err := r.iow.Write(data); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logx: write req-record error: %s\n", err.Error())
	}
	putReqRecordToPool(r)
}
