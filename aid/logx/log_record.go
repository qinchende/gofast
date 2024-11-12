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

//type Field struct {
//	Key string
//	Val any
//}

type Record struct {
	Time  time.Duration `json:"ts"`
	Label string        `json:"lb"`
	//Msg   string

	log *Logger
	iow io.WriteCloser
	out LogBuilder
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
func newRecord(w io.WriteCloser, label string) *Record {
	r := getRecordFromPool()
	r.Label = label
	r.iow = w
	r.out = r
	return r
}

func (r *Record) Output(msg string) {
	r.Time = timex.NowDur() // 此时才是日志记录时间

	r.bs = jde.AppendStrField(r.bs, "msg", msg)
	r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
	r.bs = append(r.bs, "\n"...)

	if _, err := r.iow.Write(r.bs); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logx: write record error: %s\n", err.Error())
	}
	putRecordToPool(r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) SetLabel(label string) *Record {
	if r == nil {
		return nil
	}
	r.Label = label
	return r
}

func (r *Record) Send() {
	if r != nil {
		r.out.Output("")
	}
}

func (r *Record) Msg(msg string) {
	if r != nil {
		r.out.Output(msg)
	}
}

// MsgF虽然方便，但不推荐使用
func (r *Record) MsgF(str string, v ...any) {
	if r != nil {
		r.out.Output(fmt.Sprintf(str, v...))
	}
}

func (r *Record) MsgFunc(createMsg func() string) {
	if r != nil {
		r.out.Output(createMsg())
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) Str(k, v string) *Record {
	if r == nil {
		return nil
	}
	r.bs = jde.AppendStrField(r.bs, k, v)
	return r
}

func (r *Record) Int(k string, v int) *Record {
	if r == nil {
		return nil
	}
	r.bs = jde.AppendIntField(r.bs, k, v)
	return r
}

func (r *Record) Bool(k string, v bool) *Record {
	if r == nil {
		return nil
	}
	r.bs = jde.AppendBoolField(r.bs, k, v)
	return r
}
