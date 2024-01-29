// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "github.com/qinchende/gofast/cst"

// 实现Session存储时，处理失败就抛出异常
type SessionKeeper interface {
	// Session值相关方法
	Get(string) (string, bool) // 获取某个key的值
	GetValues() cst.WebKV      // 获取当前session中的所有键值对
	Set(string, string)        // 设置session的一组kv值
	SetValues(cst.WebKV)       // cst.WebKV类型的session数据
	Del(string)                // 删除某个key
	SetUid(string)             // 设置当前用户的唯一标识ID
	GetUid() string            // 获取用户唯一标识ID

	// Session控制方法
	Token() string    // SessionID
	TokenIsNew() bool // SessionID is new?
	ExpireS(uint32)   // 设置过期时间秒
	Save()            // 保存session数据
	Destroy()         // 销毁当前session数据
	Recreate()        // 重新创建session需要的基础数据

	//Saved() bool // 是否已保存
}
