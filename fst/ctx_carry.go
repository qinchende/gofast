// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/tips"
	"github.com/qinchende/gofast/logx"
)

// 将来可以扩展，携带各种数据类型
const (
	carryTypeAny     tips.CarryType = 0
	carryTypePrivate tips.CarryType = 1 << 0
	carryTypePublic  tips.CarryType = 1 << 1
	carryTypeMsg     tips.CarryType = 1 << 2 // 传递消息
)

// 添加一条消息，日志系统会打印出这些传递信息
func (c *Context) CarryAddMsg(msg string) {
	if len(c.CarryItems) > maxCtxCarryLen {
		logx.Error("current request context carry list is out of range.")
		return
	}
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
