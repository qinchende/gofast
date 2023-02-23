// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

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
