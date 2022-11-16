// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"net/http"
)

// 请求分析，针对不同路由分别执行熔断策略
func Breaker(kp *gate.RequestKeeper) fst.CtxHandler {
	if kp == nil {
		return nil
	}

	return func(c *fst.Context) {
		// 检查是否允许本次访问通过，主要是滑动窗口判断是否达到熔断条件
		brk := kp.Breakers[c.RouteIdx]
		err := brk.Allow()
		// 有错误信息返回，证明本次请求被熔断，接下来：
		// 1. 本次记入丢弃请求统计  2. 打印错误信息  3. 返回服务器出错
		if err != nil {
			kp.CountRouteDrop(c.RouteIdx)
			// 有可能会连续疯狂的熔断，确认要打印所有信息吗？这里先不打印
			//r := c.ReqRaw
			//logx.ErrorF("[http] break, %s - %s - %s", r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent())
			c.AbortDirect(http.StatusServiceUnavailable, midFusingBody)
			// 返回之后，后面的 defer 和 c.Next() 都不会执行。
			return
		}

		defer func() {
			status := c.ResWrap.Status()
			// 5xx 以下的错误被认为是正常返回。否认就是服务器错误，被认定是拒绝服务
			if status < http.StatusInternalServerError {
				brk.Accept() // 一次正常请求
			} else {
				brk.Reject(http.StatusText(status)) // 一次异常返回
			}
		}()

		// 执行后面的处理函数
		c.Next()
	}
}
