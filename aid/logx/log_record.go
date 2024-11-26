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

type Record struct {
	myL *Logger      // 日志记录器
	iow io.Writer    // 最终日志写入的介质
	out RecordWriter // Build日志数据

	pBuf    *[]byte // 来自于全局缓存
	bs      []byte  // 用来辅助上面的bf指针，防止24个字节的切片对象堆分配
	isGroup bool    // 当前是否处于分组阶段
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var (
	//_recordDefVal Record
	recordPool = &sync.Pool{New: func() interface{} { return &Record{} }}
)

func getRecord() *Record {
	r := recordPool.Get().(*Record)
	// 最少2048字节，这里可以改进算法，动态计算最优大小，或者按Label分类给大小
	r.pBuf = pool.GetBytesMin(2048)
	r.bs = *r.pBuf
	return r
}

func backRecord(r *Record) {
	*r.pBuf = r.bs
	pool.FreeBytes(r.pBuf)
	r.pBuf = nil
	r.bs = nil
	recordPool.Put(r)
}

func newRecord(l *Logger, w io.Writer, label string) *Record {
	r := getRecord()
	r.myL = l
	r.iow = w
	r.out = r
	r.bs = append(l.LogBegin(r.bs, label), l.r.bs...)
	return r
}

func (r *Record) reuse(w io.Writer, label string) {
	r.iow = w
	r.bs = append(r.myL.LogBegin(r.bs[0:0], label), r.myL.r.bs...)
}

func (r *Record) write() {
	data := r.myL.LogEnd(r.bs)
	if _, err := r.iow.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "logx: write error: %s\n", err.Error())
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) SetWriter(w io.Writer) *Record {
	if r == nil {
		return nil
	}
	r.iow = w
	return r
}

func (r *Record) GetBuf() []byte {
	if r != nil {
		return r.bs
	}
	return nil
}

func (r *Record) SetBuf(bs []byte) *Record {
	if r != nil {
		r.bs = bs
	}
	return r
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) Reuse(w io.Writer, level int8, label string) *Record {
	if r != nil && r.myL.iLevel <= level {
		r.reuse(w, label)
	}
	return r
}

func (r *Record) Trace() *Record {
	if r != nil && r.myL.ShowTrace() {
		r.reuse(r.myL.WTrace, LabelTrace)
	}
	return r
}

func (r *Record) Debug() *Record {
	if r != nil && r.myL.ShowDebug() {
		r.reuse(r.myL.WDebug, LabelDebug)
	}
	return r
}

func (r *Record) Info() *Record {
	if r != nil && r.myL.ShowInfo() {
		r.reuse(r.myL.WInfo, LabelInfo)
	}
	return r
}

func (r *Record) InfoTimer() *Record {
	if r != nil && r.myL.ShowInfo() {
		r.reuse(r.myL.WTimer, LabelTimer)
	}
	return r
}

func (r *Record) InfoReq() *Record {
	if r != nil && r.myL.ShowInfo() {
		r.reuse(r.myL.WReq, LabelReq)
	}
	return r
}

func (r *Record) InfoStat() *Record {
	if r != nil && r.myL.ShowInfo() {
		r.reuse(r.myL.WStat, LabelStat)
	}
	return r
}

func (r *Record) Warn() *Record {
	if r != nil && r.myL.ShowWarn() {
		r.reuse(r.myL.WWarn, LabelWarn)
	}
	return r
}

func (r *Record) WarnSlow() *Record {
	if r != nil && r.myL.ShowWarn() {
		r.reuse(r.myL.WSlow, LabelSlow)
	}
	return r
}

func (r *Record) Error() *Record {
	if r != nil && r.myL.ShowErr() {
		r.reuse(r.myL.WErr, LabelErr)
	}
	return r
}

func (r *Record) ErrPanic() *Record {
	if r != nil && r.myL.ShowErr() {
		r.reuse(r.myL.WPanic, LabelPanic)
	}
	return r
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (r *Record) Send() {
	r.sendWithMessage("")
}

func (r *Record) SendMsg(msg string) {
	r.sendWithMessage(msg)
}

// MsgF虽然方便，但不推荐使用
func (r *Record) SendMsgF(str string, v ...any) {
	r.sendWithMessage(fmt.Sprintf(str, v...))
}

func (r *Record) SendMsgFunc(createMsg func() string) {
	r.sendWithMessage(createMsg())
}

func (r *Record) sendWithMessage(msg string) {
	if r != nil {
		if r.isGroup {
			r.GEnd()
		}
		if len(msg) > 0 {
			bf := jde.AppendStrField(r.bs, fMessage, msg)
			r.bs = bf[:len(bf)-1] // 去掉最后面一个逗号
		}
		r.out.write()
		backRecord(r)
	}
}

// 可以先输出一条完整的日志，但是不回收Record，而是继续打印下一条
func (r *Record) SendBack() *Record {
	if r != nil {
		if r.isGroup {
			r.GEnd()
		}
		r.out.write()
	}
	return r
}

func (r *Record) Group(k string) *Record {
	if r != nil {
		if r.isGroup {
			r.GEnd()
		}
		r.isGroup = true
		r.bs = r.myL.GroupBegin(r.bs, k)
	}
	return r
}

func (r *Record) GEnd() *Record {
	if r != nil && r.isGroup {
		r.isGroup = false
		r.bs = r.myL.GroupEnd(r.bs)
	}
	return r
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++++ int
func (r *Record) Int(k string, v int) *Record {
	if r != nil {
		r.bs = jde.AppendIntField(r.bs, k, v)
	}
	return r
}

func (r *Record) I64(k string, v int64) *Record {
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
func (r *Record) Uint(k string, v uint) *Record {
	if r != nil {
		r.bs = jde.AppendUintField(r.bs, k, v)
	}
	return r
}

func (r *Record) U64(k string, v uint64) *Record {
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

func (r *Record) Bools(k string, v []bool) *Record {
	if r != nil {
		r.bs = jde.AppendBoolListField(r.bs, k, v)
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

func (r *Record) F32s(k string, v []float32) *Record {
	if r != nil {
		r.bs = jde.AppendF32sField(r.bs, k, v)
	}
	return r
}

func (r *Record) F64(k string, v float64) *Record {
	if r != nil {
		r.bs = jde.AppendF64Field(r.bs, k, v)
	}
	return r
}

func (r *Record) F64s(k string, v []float64) *Record {
	if r != nil {
		r.bs = jde.AppendF64sField(r.bs, k, v)
	}
	return r
}

// +++++ string
func (r *Record) Append(v []byte) *Record {
	if r != nil {
		r.bs = append(r.bs, v...)
	}
	return r
}

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

func (r *Record) Times(k string, v []time.Time) *Record {
	if r != nil {
		r.bs = jde.AppendTimeListField(r.bs, k, v, timeFormat)
	}
	return r
}

// +++++ error
func (r *Record) Err(v error) *Record {
	if r != nil && v != nil {
		r.bs = jde.AppendStrField(r.bs, fError, v.Error())
	}
	return r
}

// +++++ struct
func (r *Record) Obj(k string, v ObjEncoder) *Record {
	if r != nil {
		bf := jde.AppendKey(r.bs, k)
		if v == nil {
			r.bs = append(bf, "null"...)
			return r
		}

		r.bs = append(bf, '{')
		v.EncodeLogx(r)

		bf = r.bs
		if bf[len(bf)-1] == ',' {
			bf = bf[:len(bf)-1]
		}
		r.bs = append(bf, "},"...)
	}
	return r
}

// 任意struct切片类型转换成LogX输出的切片类型
func ToObjs[T ObjEncoder](list []T) []ObjEncoder {
	arr := make([]ObjEncoder, len(list))
	for idx := range list {
		arr[idx] = list[idx]
	}
	return arr
}

func (r *Record) Objs(k string, v []ObjEncoder) *Record {
	if r != nil && v != nil {
		bf := jde.AppendKey(r.bs, k)
		if len(v) == 0 {
			r.bs = append(bf, "[],"...)
		} else {
			bf = append(bf, '[')
			for idx := range v {
				r.bs = append(bf, '{')
				v[idx].EncodeLogx(r)
				bf = r.bs
				if bf[len(bf)-1] == ',' {
					bf = bf[:len(bf)-1]
				}
				bf = append(bf, "},"...)
			}
			r.bs = append(bf[:len(bf)-1], "],"...)
		}
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
