// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"net/http"
	"runtime/debug"
)

// 截获异常，防止程序崩溃。
func Recovery(c *fst.Context) {
	defer func() {
		if pic := recover(); pic != nil {
			// TODO: 这里要分三种异常，1.模拟返回错误 2.常规的错误异常 3.非预测性的系统异常
			switch info := pic.(type) {
			case cst.GFFaiString:
				c.AbortFaiStr(string(info))
			case cst.GFError:
				c.AbortFaiStr(fmt.Sprint("GFError: ", info))
			default:
				logx.Stacks(c.ReqRaw)
				logx.StackF("%s", debug.Stack())
				c.AbortDirect(http.StatusInternalServerError, fmt.Sprint("panic: ", info))
			}
		}
	}()

	c.Next()
}
