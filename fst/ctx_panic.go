// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "github.com/qinchende/gofast/skill/lang"

// 异常处理逻辑的接口定义
type PanicHandler interface {
	Callback()
}

type PanicFunc struct {
	Func func()
}

func (pw PanicFunc) Callback() { pw.Func() }

func NewPanicPet(fn func()) *PanicFunc {
	return &PanicFunc{Func: fn}
}

// 取出只作为消息传递的项
func (c *Context) PanicCatch(ret any) {
	if pic := recover(); pic != nil {
		c.CarryAddMsg(lang.ToString(pic))
		switch ret.(type) {
		case string:
			c.FaiMsg(ret.(string))
		case *Ret:
			c.FaiRet(ret.(*Ret))
		default:
			c.FaiMsg(lang.ToString(ret))
		}
	}
}
