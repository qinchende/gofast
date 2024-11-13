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
)

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

func (r *Record) output(msg string) {
	r.Time = timex.NowDur() // 此时才是日志记录时间

	r.bs = jde.AppendStrField(r.bs, "msg", msg)
	r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号
	r.bs = append(r.bs, "\n"...)

	// 合成最后的输出结果
	r.bs = r.log.StyleFunc(r.log, r.bs)

	if _, err := r.iow.Write(r.bs); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logx: write record error: %s\n", err.Error())
	}
	putRecordToPool(r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) SetWriter(w io.WriteCloser) *Record {
	if r == nil {
		return nil
	}
	r.iow = w
	return r
}

func (r *Record) SetLabel(label string) *Record {
	if r == nil {
		return nil
	}
	r.Label = label
	return r
}

// 可以先输出一条完整的日志，但是不回收Record，而是继续下一条
func (r *Record) Flush() *Record {
	return r
}

func (r *Record) Send() {
	if r != nil {
		r.out.output("")
	}
}

func (r *Record) Msg(msg string) {
	if r != nil {
		r.out.output(msg)
	}
}

// MsgF虽然方便，但不推荐使用
func (r *Record) MsgF(str string, v ...any) {
	if r != nil {
		r.out.output(fmt.Sprintf(str, v...))
	}
}

func (r *Record) MsgFunc(createMsg func() string) {
	if r != nil {
		r.out.output(createMsg())
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) Str(k, v string) *Record {
	if r != nil {
		r.bs = jde.AppendStrField(r.bs, k, v)
	}
	return r
}

func (r *Record) Int(k string, v int) *Record {
	if r != nil {
		r.bs = jde.AppendIntField(r.bs, k, v)
	}
	return r
}

func (r *Record) Bool(k string, v bool) *Record {
	if r != nil {
		r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}

func (r *Record) F32(k string, v float32) *Record {
	if r != nil {
		//r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}

func (r *Record) F64(k string, v float64) *Record {
	if r != nil {
		//r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}

func (r *Record) Obj(k string, v any) *Record {
	if r != nil {
		//r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}

func (r *Record) Any(k string, v any) *Record {
	if r != nil {
		//r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}
