package logx

import (
	"fmt"
	"github.com/qinchende/gofast/store/jde"
)

func (r *Record) Send() {
	if r != nil {
		r.output("")
	}
}

func (r *Record) Msg(msg string) {
	if r != nil {
		r.output(msg)
	}
}

// MsgF虽然方便，但不推荐使用
func (r *Record) MsgF(str string, v ...any) {
	if r != nil {
		r.output(fmt.Sprintf(str, v...))
	}
}

func (r *Record) MsgFunc(createMsg func() string) {
	if r != nil {
		r.output(createMsg())
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
