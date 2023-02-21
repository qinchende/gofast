// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "github.com/qinchende/gofast/cst"

// 实现Session存储时，处理失败就抛出异常
type SessionKeeper interface {
	GetValues() cst.KV   // 获取当前session中的所有键值对
	Get(string) any      // 获取某个key的值
	Set(string, any)     // 设置session的一组kv值
	SetKV(cst.KV)        // cst.KV 类型的session数据
	Del(string)          // 删除某个session中的key
	Save()               // 保存session数据
	Saved() bool         // 是否已保存
	ExpireS(int32)       // 设置过期时间秒
	SidIsNew() bool      // SessionID is new?
	Sid() string         // SessionID
	Destroy()            // 销毁当前session数据
	Recreate(c *Context) // 重新创建session信息
}
