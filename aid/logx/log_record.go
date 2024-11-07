// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/aid/bag"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/store/jde"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Field struct {
	Key string
	Val any
}

type Record struct {
	Time  time.Duration
	App   string
	Host  string
	Label string
	//Msg   string

	w  io.Writer
	bf *[]byte
	bs []byte // 用来辅助上面的bf指针，防止24个字节的切片对象堆分配
}

type ReqRecord struct {
	RawReq     *http.Request
	StatusCode int
	Method     string
	RequestURI string
	UserAgent  string
	RemoteAddr string

	Record
	TimeStamp  time.Duration
	Latency    time.Duration
	Pms        cst.SuperKV
	BodySize   int
	ResData    []byte
	CarryItems bag.CarryList
}

var (
	recordPool = &sync.Pool{
		New: func() interface{} {
			r := &Record{}
			r.init()
			return r
		},
	}
	reqRecordPool = &sync.Pool{
		New: func() interface{} {
			r := &ReqRecord{}
			r.Record.init()
			return r
		},
	}

	_recordDefValue    Record
	_reqRecordDefValue ReqRecord
)

func (r *Record) init() {
	r.bf = pool.GetBytes()
	r.bs = *r.bf
	r.App = myCnf.AppName
	r.Host = myCnf.HostName
}

func getRecordFromPool() *Record {
	r := recordPool.Get().(*Record)
	r.bf = pool.GetBytes()
	r.bs = *r.bf
	//r.bs = r.bs[0:0]
	return r
}

func putRecordToPool(r *Record) {
	*r.bf = r.bs
	pool.FreeBytes(r.bf)
	r.bf = nil
	r.bs = nil
	recordPool.Put(r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func NewRecord(w io.Writer, label string) *Record {
	r := getRecordFromPool()
	r.Time = timex.NowDur()
	r.Label = label
	r.w = w
	return r
}

func (r *Record) output(msg string) {
	r.bs = jde.AppendStrField(r.bs, "msg", msg)
	r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
	r.bs = append(r.bs, "\n"...)

	if _, err := r.w.Write(r.bs); err != nil {
		log.Println("Panic to write log, error: " + err.Error())
		panic(err)
	}
	putRecordToPool(r)
}
