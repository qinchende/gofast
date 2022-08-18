// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/tools"
	"github.com/qinchende/gofast/logx"
)

// 添加一项携带消息体的Basket，日志系统会打印出这些信息
func (c *Context) AddMsgBasket(msg string) {
	if len(c.Baskets) > maxCtxBaskets {
		logx.Error("current request context baskets is out of range.")
		return
	}
	b := &tools.Basket{
		Msg:  msg,
		Type: BasketTypeMsg,
		Meta: nil,
	}
	c.Baskets = append(c.Baskets, b)
}

// 取出只作为消息传递的篮子
func (c *Context) MsgBaskets() tools.Baskets {
	return c.Baskets.ByType(BasketTypeMsg)
}
