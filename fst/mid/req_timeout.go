// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"context"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"net/http"
	"time"
)

// 超时之后的返回内容
var midTimeoutBody = "<html><head><title>Timeout</title></head><body><h1>Timeout</h1></body></html>"

// 方式一：标准库
func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if duration > 0 {
			return http.TimeoutHandler(next, duration, midTimeoutBody)
		} else {
			return next
		}
	}
}

// 方式二：设置请求处理的超时时间，单位是毫秒
// TODO：注意下面的注释
// NOTE：这个只是简易的方案，存在不严谨的情况。比如返回结果render了一部分，结果G被Cancel掉了。上面的标准库处理了这个问题。
// 方式一却自定义了 response write 加入了 输出缓存，返回结果全部好了之后才会一次性 render 给客户端。
func Timeout(useTimeout bool) fst.CtxHandler {
	// Debug模式不设置超时
	if useTimeout == false {
		return nil
	}

	return func(c *fst.Context) {
		rt := RConfigs[c.RouteIdx]
		ctxTimeout, cancelCtx := context.WithTimeout(c.ReqRaw.Context(), time.Duration(rt.Timeout)*time.Millisecond)
		defer cancelCtx()

		panicChan := make(chan any, 1)
		finishChan := make(chan struct{})

		// 启动的协程 没有办法杀死，唯一的办法只能用通道通知他，让协程自己退出
		// 如果协程一直不退出，将会一直占用协程的堆栈内存，并且一直处于GMP的待处理队列，影响整体性能
		// 如果外面超时退出，这里的调用处理还是会继续走下去的。（即有可能调用者看到的错误提示，但是后台却是处理成功的）
		go func() {
			defer func() {
				// 其实这个是执行不到的，因为执行链上的 Recovery 函数把异常吃掉了
				if pic := recover(); pic != nil {
					// NOTE：这里必须使用带缓冲的通道，否则本G可能因为父G的提前退出，而卡死在这里，导致G泄露
					panicChan <- pic
					// TODO：无法确定上面的异常是否传递出去，下面的日志还是需要打印的。
					logx.ErrorF("ReqTimeout-Panic: ", pic)
				}
			}()
			// 不管超不超时，本次请求都会执行完毕，或者等到自己超时退出
			c.Next()
			// 执行完成之后通知 主协程，我做完了，你可以退出了
			close(finishChan)
		}()

		// 任何一个先触发都会执行，并结束当前函数，不会两个以上都触发
		select {
		case pic := <-panicChan:
			// 子G发生异常，抛出传递给上层G
			panic(pic)
			return
		case <-finishChan:
			// 正常退出
			//log.Println("I am back.")
			return
		case <-ctxTimeout.Done():
			c.IsTimeout = true
			c.AbortString(http.StatusGatewayTimeout, midTimeoutBody)
			return
		}
	}
}
