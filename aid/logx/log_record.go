// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/store/jde"
	"io"
	"os"
	"sync"
	"time"
)

var (
	_recordDef Record
	recordPool = &sync.Pool{New: func() interface{} { return &Record{} }}
)

type (
	Record struct {
		//TDur  time.Duration
		//Label string
		//Message string

		myL *Logger
		iow io.Writer
		out RecordWriter

		//fls []Field // 用来记录key-value
		buf *[]byte // 来自于全局缓存
		bs  []byte  // 用来辅助上面的bf指针，防止24个字节的切片对象堆分配

		isGroup bool // 当前是否处于分组阶段
	}
)

func (r *Record) SetWriter(w io.Writer) *Record {
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
	//r.Label = label
	return r
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func getRecord() *Record {
	r := recordPool.Get().(*Record)
	r.buf = pool.GetBytesMin(512)
	r.bs = *r.buf
	return r
}

func backRecord(r *Record) {
	*r.buf = r.bs
	pool.FreeBytes(r.buf)
	r.buf = nil
	r.bs = nil
	recordPool.Put(r)
}

func (l *Logger) newRecord(w io.Writer, label string) *Record {
	r := getRecord()
	//r.Label = label
	r.myL = l
	r.iow = w
	r.out = r
	r.myL.StyleBegin(r, label)
	return r
}
func (r *Record) reuse() {
	r.bs = r.bs[:0]
}

func (r *Record) write() {
	// 此时才是日志记录时间
	//r.TDur = timex.NowDur()
	// 合成最后的输出内容
	data := r.myL.StyleSummary(r)

	if _, err := r.iow.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "logx: write error: %s\n", err.Error())
	}
}

func (r *Record) endWithMsg(msg string) {
	if r.isGroup {
		r.GEnd()
	}

	r.bs = jde.AppendStrField(r.bs, fMessage, msg)
	r.bs = r.bs[:len(r.bs)-1] // 去掉最后面一个逗号

	r.out.write()
	backRecord(r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 可以先输出一条完整的日志，但是不回收Record，而是继续下一条
func (r *Record) Next() *Record {
	if r != nil {
		if r.isGroup {
			r.GEnd()
		}
		r.out.write()
		r.reuse()
	}
	return r
}

func (r *Record) End() {
	if r != nil {
		if r.isGroup {
			r.GEnd()
		}
		r.out.write()
		backRecord(r)
	}
}

func (r *Record) Msg(msg string) {
	if r != nil {
		r.endWithMsg(msg)
	}
}

// MsgF虽然方便，但不推荐使用
func (r *Record) MsgF(str string, v ...any) {
	if r != nil {
		r.endWithMsg(fmt.Sprintf(str, v...))
	}
}

func (r *Record) MsgFunc(createMsg func() string) {
	if r != nil {
		r.endWithMsg(createMsg())
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) Group(k string) *Record {
	if r != nil {
		if r.isGroup {
			r.GEnd()
		}
		r.isGroup = true
		r.bs = r.myL.StyleGroupBegin(r.bs, k)
	}
	return r
}

func (r *Record) GEnd() *Record {
	if r != nil && r.isGroup {
		r.isGroup = false
		r.bs = r.myL.StyleGroupEnd(r.bs)
	}
	return r
}

// +++++ int
func (r *Record) Int(k string, v int) *Record {
	if r != nil {
		r.bs = jde.AppendIntField(r.bs, k, v)
	}
	return r
}

func (r *Record) Ints(k string, v []int) *Record {
	if r != nil {
		r.bs = jde.AppendIntListField(r.bs, k, v)
	}
	return r
}

func (r *Record) I8s(k string, v []int8) *Record {
	if r != nil {
		r.bs = jde.AppendIntListField(r.bs, k, v)
	}
	return r
}
func (r *Record) I16s(k string, v []int16) *Record {
	if r != nil {
		r.bs = jde.AppendIntListField(r.bs, k, v)
	}
	return r
}
func (r *Record) I32s(k string, v []int32) *Record {
	if r != nil {
		r.bs = jde.AppendIntListField(r.bs, k, v)
	}
	return r
}
func (r *Record) I64s(k string, v []int64) *Record {
	if r != nil {
		r.bs = jde.AppendIntListField(r.bs, k, v)
	}
	return r
}

// +++++ uint
func (r *Record) Uint(k string, v uint64) *Record {
	if r != nil {
		r.bs = jde.AppendUintField(r.bs, k, v)
	}
	return r
}

func (r *Record) Uints(k string, v []uint) *Record {
	if r != nil {
		r.bs = jde.AppendUintListField(r.bs, k, v)
	}
	return r
}

func (r *Record) U8s(k string, v []uint8) *Record {
	if r != nil {
		r.bs = jde.AppendUintListField(r.bs, k, v)
	}
	return r
}

func (r *Record) U16s(k string, v []uint16) *Record {
	if r != nil {
		r.bs = jde.AppendUintListField(r.bs, k, v)
	}
	return r
}

func (r *Record) U32s(k string, v []uint32) *Record {
	if r != nil {
		r.bs = jde.AppendUintListField(r.bs, k, v)
	}
	return r
}

func (r *Record) U64s(k string, v []uint64) *Record {
	if r != nil {
		r.bs = jde.AppendUintListField(r.bs, k, v)
	}
	return r
}

// +++++ bool
func (r *Record) Bool(k string, v bool) *Record {
	if r != nil {
		r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}

// +++++ float
func (r *Record) F32(k string, v float32) *Record {
	if r != nil {
		r.bs = jde.AppendF32Field(r.bs, k, v)
	}
	return r
}

func (r *Record) F64(k string, v float64) *Record {
	if r != nil {
		r.bs = jde.AppendF64Field(r.bs, k, v)
	}
	return r
}

// +++++ string
func (r *Record) Str(k, v string) *Record {
	if r != nil {
		r.bs = jde.AppendStrField(r.bs, k, v)
	}
	return r
}

func (r *Record) Strs(k string, v []string) *Record {
	if r != nil {
		r.bs = jde.AppendStrListField(r.bs, k, v)
	}
	return r
}

// +++++ time.Time
func (r *Record) Time(k string, v time.Time) *Record {
	if r != nil {
		r.bs = jde.AppendTimeField(r.bs, k, v, timeFormat)
	}
	return r
}

//func (r *Record) Times(k string, v []time.Time) *Record {
//	if r != nil {
//		r.bs = jde.AppendStrListField(r.bs, k, v)
//	}
//	return r
//}

// +++++ struct
func (r *Record) Obj(k string, v any) *Record {
	if r != nil {
		//r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}

// +++++ any
func (r *Record) Any(k string, v any) *Record {
	if r != nil {
		//r.bs = jde.AppendBoolField(r.bs, k, v)
	}
	return r
}
