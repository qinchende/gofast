// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
)

// 异常处理逻辑的接口定义
type panicHandler interface {
	Callback()
}

type PanicFunc struct {
	Func func()
}

func (pw PanicFunc) Callback() { pw.Func() }

//func NewPanicPet(fn func()) *PanicFunc {
//	return &PanicFunc{Func: fn}
//}

// 取出只作为消息传递的项
func (c *Context) PanicCatch(ret any) {
	if pic := recover(); pic != nil {
		c.CarryAddMsg(lang.ToString(pic))
		switch ret.(type) {
		case string:
			c.FaiMsg(ret.(string))
		case *cst.Ret:
			c.FaiRet(ret.(*cst.Ret))
		default:
			c.FaiMsg(lang.ToString(ret))
		}
	}
}

// 如果有错误就引发异常，同时返回指定的结果
func (c *Context) PanicIfErr(err error, ret any) {
	if err == nil {
		return
	}
	if ret == nil {
		cst.Panic(err)
		return
	}
	c.CarryAddMsg(err.Error())
	cst.Panic(ret)
}

func (c *Context) PanicIf(ifTrue bool, ret any) {
	cst.PanicIf(ifTrue, ret)
}

func (c *Context) PanicString(ret string) {
	cst.PanicString(ret)
}

func (c *Context) Panic(ret any) {
	cst.Panic(ret)
}
