// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/wit"
)

// 通过类型区分，可以让Context.Values传递各种出现频率并不高的数据，防止Context对象内存占用的扩张
// 将来可以扩展，携带各种数据类型
const (
	carryDef uint = iota
	carryPrivate
	carryPublic
	carryPanicFunc
	carryLogs       // 传递消息
	carryFormCache  // formCache  url.Values  // the parsed form data from POST, PATCH, or PUT body parameters.
	carryQueryCache // queryCache url.Values  // param query result from c.Req.URL.Query()
	
	maxCarryLen int = 8 // 限制 Context.Values 中值的数量
)

// +++++++++++++++++++++++++++++++++++++
// Log ext msg
// +++++++++++++++++++++++++++++++++++++
// 添加一条消息，日志系统会打印出这些传递信息
func (c *Context) LogStr(key string, msg string) {
	c.checkCarrySize()
	c.Values = append(c.Values, wit.KVItemGroup{
		Group:  carryLogs,
		KVItem: wit.Str(key, msg),
	})
}

func (c *Context) LogItem(item wit.KVItem) {
	c.checkCarrySize()
	c.Values = append(c.Values, wit.KVItemGroup{
		Group:  carryLogs,
		KVItem: item,
	})
}

// 取出只作为消息传递的项
func (c *Context) LogItems() wit.KVListGroup {
	return c.Values.ByGroup(carryLogs)
}

// +++++++++++++++++++++++++++++++++++++
// PanicPet func
// +++++++++++++++++++++++++++++++++++++
func (c *Context) SetPanicPet(fn PanicHandler) {
	c.checkCarrySize()
	c.Values = append(c.Values, wit.KVItemGroup{
		KVItem: wit.KVItem{Val: fn},
		Group:  carryPanicFunc,
	})
}

func (c *Context) GetPanicPet() PanicHandler {
	if it := c.Values.FirstOne(carryPanicFunc); it != nil {
		return it.Val.(PanicHandler)
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++
// url FormCache
// +++++++++++++++++++++++++++++++++++++
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

// +++++++++++++++++++++++++++++++++++++
// url QueryCache
// +++++++++++++++++++++++++++++++++++++
func (c *Context) setQueryCache(val cst.WebKV) {
	c.checkCarrySize()
	c.Values = append(c.Values, wit.KVItemGroup{
		KVItem: wit.KVItem{Val: val},
		Group:  carryQueryCache,
	})
}
func (c *Context) getQueryCache() cst.WebKV {
	if it := c.Values.FirstOne(carryQueryCache); it != nil {
		return it.Val.(cst.WebKV)
	}
	return nil
}

// 控制context.CarryList的长度，这个结构要通过sync.Pool复用，内存占用会只增不减
func (c *Context) checkCarrySize() {
	c.PanicIf(len(c.Values) > maxCarryLen, "Carry list is out of range.")
}
