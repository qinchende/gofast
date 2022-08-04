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
		if result := recover(); result != nil {
			// TODO: 这里要分两种异常，一种是常规的错误异常，一种是非预测性的系统异常
			if err, ok := result.(cst.GFError); ok {
				c.AbortFaiStr(fmt.Sprint("GFError: ", err))
			} else {
				logx.Stack(c.ReqRaw)
				logx.StackF("%s", debug.Stack())

				c.AbortDirect(http.StatusInternalServerError, fmt.Sprint("panic: ", result))
			}
		}
	}()

	c.Next()
}
