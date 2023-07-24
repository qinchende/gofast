// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/tips"
)

// 将来可以扩展，携带各种数据类型
const (
	carryTypeAny        tips.CarryType = 0
	carryTypePrivate    tips.CarryType = 1 << 0
	carryTypePublic     tips.CarryType = 1 << 1
	carryTypeMsg        tips.CarryType = 1 << 2 // 传递消息
	carryTypeFormCache  tips.CarryType = 1 << 3 // formCache  url.Values  // the parsed form data from POST, PATCH, or PUT body parameters.
	carryTypeQueryCache tips.CarryType = 1 << 4 // queryCache url.Values  // param query result from c.Req.URL.Query()
)

// 添加一条消息，日志系统会打印出这些传递信息
func (c *Context) CarryMsg(msg string) {
	c.checkCarrySize()
	msgItem := &tips.CarryItem{
		Type: carryTypeMsg,
		Msg:  msg,
		Meta: nil,
	}
	c.CarryItems = append(c.CarryItems, msgItem)
}

// 取出只作为消息传递的项
func (c *Context) CarryMsgItems() tips.CarryList {
	return c.CarryItems.ByType(carryTypeMsg)
}

//
//func (c *Context) setFormCache(val url.Values) {
//	c.checkCarrySize()
//	formItem := &tips.CarryItem{
//		Type: carryTypeFormCache,
//		Msg:  "form params",
//		Meta: val,
//	}
//	c.CarryItems = append(c.CarryItems, formItem)
//}
//
//func (c *Context) formCache() url.Values {
//	it := c.CarryItems.FirstOne(carryTypeFormCache)
//	if it == nil {
//		return nil
//	}
//	return it.Meta.(url.Values)
//}

func (c *Context) setQueryCache(val cst.WebKV) {
	c.checkCarrySize()
	queryItem := &tips.CarryItem{
		Type: carryTypeQueryCache,
		Msg:  "form params",
		Meta: val,
	}
	c.CarryItems = append(c.CarryItems, queryItem)
}
func (c *Context) queryCache() cst.WebKV {
	it := c.CarryItems.FirstOne(carryTypeQueryCache)
	if it == nil {
		return nil
	}
	return it.Meta.(cst.WebKV)
}

// 控制context.CarryList的长度，这个结构要通过sync.Pool复用，内存占用会只增不减
func (c *Context) checkCarrySize() {
	c.PanicIf(len(c.CarryItems) > maxCtxCarryLen, "Carry list is out of range.")
}
