// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/aid/timex"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/store/jde"
	"io"
	"os"
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

	w   io.WriteCloser
	out RecordOutput
	bf  *[]byte
	bs  []byte // 用来辅助上面的bf指针，防止24个字节的切片对象堆分配
}

var (
	recordPool = &sync.Pool{
		New: func() interface{} {
			r := &Record{}
			r.init()
			return r
		},
	}
	_recordDefValue Record
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
func NewRecord(w io.WriteCloser, label string) *Record {
	r := getRecordFromPool()
	r.Time = timex.NowDur()
	r.Label = label
	r.w = w
	r.out = r
	return r
}

func (r *Record) output(msg string) {
	r.bs = jde.AppendStrField(r.bs, "msg", msg)
	r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
	r.bs = append(r.bs, "\n"...)

	if _, err := r.w.Write(r.bs); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logx: write record error: %s\n", err.Error())
	}
	putRecordToPool(r)
}
