// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "github.com/qinchende/gofast/cst"

type SessionKeeper interface {
	GetValues() cst.KV
	Get(string) any
	Set(string, any)
	SetKV(cst.KV)
	Del(string)
	Save() error
	Saved() bool
	Expire(int32)
	SidIsNew() bool      // SessionID is new?
	Sid() string         // SessionID
	Destroy()            // 销毁当前session数据
	Recreate(c *Context) // 重新创建session信息
}
