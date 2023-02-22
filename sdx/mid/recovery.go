// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
	"net/http"
	"runtime/debug"
)

// 截获异常，防止程序崩溃。
func Recovery(c *fst.Context) {
	defer func() {
		if pic := recover(); pic != nil {
			c.ExecAfterPanicHandlers() // 异常Render之前需要特殊处理的逻辑

			// 可能需要重定向异常结果的返回
			if c.PanicPet != nil {
				switch ret := c.PanicPet.(type) {
				case fst.Ret: // 返回指定错误提示
					c.CarryAddMsg(lang.ToString(pic))
					c.AbortFai(ret.Code, ret.Msg, ret.Data)
					return
				case fst.PanicFunc, *fst.PanicFunc: // 执行自定义异常函数，比如变量初始化等
					c.PanicPet.Callback()
				}
			}

			// 异常分类: 1.模拟返回错误信息 2.模拟返回错误编码 3.主动的error异常 4.非预测性的系统异常
			switch info := pic.(type) {
			case cst.GFFaiString:
				c.AbortFai(0, string(info), nil)
			case cst.GFFaiInt:
				c.AbortFai(int(info), "", nil)
			case cst.GFError:
				c.AbortFai(0, info.Error(), nil)
			default:
				// TODO-important: 非预期的异常，将会作为熔断的判断依据（业务逻辑不要随意使用系统panic，请用框架panic）
				logx.Stacks(c.ReqRaw)
				logx.StackF("%s", debug.Stack())
				c.AbortDirect(http.StatusInternalServerError, info)
			}
		}
	}()

	c.Next()
}
