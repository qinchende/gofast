// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/logx"
	"github.com/qinchende/gofast/fst"
	"net/http"
)

// 截获异常，防止程序崩溃。
func Recovery(c *fst.Context) {
	defer func() {
		if pic := recover(); pic != nil {
			// 可能需要重定向异常结果的返回
			if c.PanicPet != nil {
				switch ret := c.PanicPet.(type) {
				case *cst.Ret:
					c.CarryMsg(lang.ToString(pic))
					c.AbortRet(ret)
					return
				case fst.PanicFunc, *fst.PanicFunc: // 执行自定义异常函数，比如变量初始化等
					c.PanicPet.Callback()
				}
			}

			// 异常分类: 1.模拟返回错误信息 2.模拟返回错误编码 3.主动的error异常 4.非预测性的系统异常
			switch info := pic.(type) {
			case cst.TypeString:
				c.AbortFai(0, string(info), nil)
			case cst.TypeError:
				c.AbortFai(0, info.Error(), nil)
			case cst.TypeInt:
				c.AbortFai(int(info), "", nil)
			case *cst.Ret:
				c.AbortRet(info)
			case cst.Ret:
				c.AbortRet(&info)
			default:
				// TODO-important: 非预期的异常，比如系统异常
				// 将会作为熔断的判断依据（业务逻辑不要随意使用系统panic，请用框架GFPanic）
				logx.Stacks(c.Req.Raw.RequestURI)
				c.AbortDirect(http.StatusInternalServerError, info)
			}
		}
	}()

	c.Next()
}
