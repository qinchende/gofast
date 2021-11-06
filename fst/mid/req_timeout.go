package mid

import (
	"context"
	"fmt"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"log"
	"net/http"
	"time"
)

var midTimeoutMsg = "Request Timeout. Over %d millisecond."

// ++++++++++++++++++++++ add by cd.net 2021.10.14
// 总说：如果中间件拦截器超时退出，那么fst模块中的 request content 对象 就会被缓冲池回首。
// 此时业务逻辑层代码在执行完IO阻塞调用之后，后面的逻辑大概率会抛出异常，因为只要需要用到上下文对象时就是nil

// 方式一：标准库
func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if duration > 0 {
			return http.TimeoutHandler(next, duration, midTimeoutMsg)
		} else {
			return next
		}
	}
}

//
//// 方式二：设置请求处理的超时时间，单位是毫秒
//// TODO：注意下面的注释
//// NOTE：这个只是简易的方案，存在不严谨的情况。比如返回结果render了一部分，结果G被Cancel掉了。上面的标准库处理了这个问题。
//// 方式一却自定义了 response write 加入了 输出缓存，返回结果全部好了之后才会一次性 render 给客户端。
//func ReqTimeout(dur time.Duration) fst.IncHandler {
//	// Debug模式不设置超时
//	if logx.IsDebugging() {
//		return nil
//	}
//
//	// 默认所有请求超时是 3 秒钟
//	if dur <= 0*time.Second {
//		dur = 3 * time.Second
//	}
//	midTimeoutMsg = fmt.Sprintf(midTimeoutMsg, dur/time.Millisecond)
//
//	return func(w *fst.GFResponse, r *http.Request) {
//		ctx, cancelCtx := context.WithTimeout(r.Context(), dur)
//		defer cancelCtx()
//
//		panicChan := make(chan interface{}, 1)
//		finishChan := make(chan struct{})
//
//		// 启动的协程 没有办法杀死，唯一的办法只能用通道通知他，让协程自己退出
//		// 如果协程一直不退出，将会一直占用协程的堆栈内存，并且一直处于GMP的待处理队列，影响整体性能
//		go func() {
//			defer func() {
//				if p := recover(); p != nil {
//					// NOTE：这里必须使用带缓冲的通道，否则本G可能因为父G的提前退出，而卡死在这里，导致G泄露
//					panicChan <- p
//					// TODO：无法确定上面的异常是否传递出去，下面的日志还是需要打印的。
//					log.Println("ReqTimeout panic: ", p)
//				}
//			}()
//			// 不管超不超时，本次请求都会执行完毕，或者等到自己超时退出
//			w.NextFit(r)
//			// 执行完成之后通知 主协程，我做完了，你可以退出了
//			close(finishChan)
//		}()
//
//		// 任何一个先触发都会执行，并结束当前函数
//		select {
//		case pic := <-panicChan:
//			// 子G发生异常，抛出传递给上层G
//			panic(pic)
//		case <-finishChan:
//			// 正常退出
//			//log.Println("I am back.")
//			return
//		case <-ctx.Done():
//			// 超时退出
//			w.ResWrap.WriteHeader(http.StatusServiceUnavailable)
//			fst.RaisePanic(midTimeoutMsg)
//			return
//		}
//		// 下面这种写法，无法解决批量cancel所有子孙 goroutine 的情况。
//		//case <-time.After(dur):
//		//	// 超时退出
//		//	//w.ResWrap.WriteHeader(http.StatusServiceUnavailable)
//		//	fst.RaisePanic(midTimeoutMsg)
//		//	return
//	}
//}

//// 方式三
//// ++++++++++++++++++ add by chende 2021.10.13
//// NOTE：完善方式二，使其达到方式一的效果，同时满足本自定义框架的特点。
//// 这种方式有个问题，就是用不了标准库中 buffer 的 responseWrite 。这样还不如使用方式二
//func ReqTimeoutSuper(dur time.Duration) fst.IncHandler {
//	// Debug模式不设置超时
//	if logx.IsDebugging() {
//		return nil
//	}
//	// 默认所有请求超时是 3 秒钟
//	if dur <= 0*time.Second {
//		dur = 3 * time.Second
//	}
//	midTimeoutMsg = fmt.Sprintf(midTimeoutMsg, dur/time.Millisecond)
//
//	return func(w *fst.GFResponse, r *http.Request) {
//		twHandler := http.TimeoutHandler(w, dur, midTimeoutMsg)
//		twHandler.ServeHTTP(w.ResWrap, r)
//	}
//}

// 方式二：设置请求处理的超时时间，单位是毫秒
// TODO：注意下面的注释
// NOTE：这个只是简易的方案，存在不严谨的情况。比如返回结果render了一部分，结果G被Cancel掉了。上面的标准库处理了这个问题。
// 方式一却自定义了 response write 加入了 输出缓存，返回结果全部好了之后才会一次性 render 给客户端。
func ReqTimeout(dur time.Duration) fst.CtxHandler {
	// Debug模式不设置超时
	if logx.IsDebugging() {
		return nil
	}

	// 默认所有请求超时是 3 秒钟
	if dur <= 0*time.Second {
		dur = 3 * time.Second
	}
	midTimeoutMsg = fmt.Sprintf(midTimeoutMsg, dur/time.Millisecond)

	return func(ctx *fst.Context) {
		ctx2, cancelCtx := context.WithTimeout(ctx.ReqRaw.Context(), dur)
		defer cancelCtx()

		panicChan := make(chan interface{}, 1)
		finishChan := make(chan struct{})

		// 启动的协程 没有办法杀死，唯一的办法只能用通道通知他，让协程自己退出
		// 如果协程一直不退出，将会一直占用协程的堆栈内存，并且一直处于GMP的待处理队列，影响整体性能
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// NOTE：这里必须使用带缓冲的通道，否则本G可能因为父G的提前退出，而卡死在这里，导致G泄露
					panicChan <- p
					// TODO：无法确定上面的异常是否传递出去，下面的日志还是需要打印的。
					log.Println("ReqTimeout panic: ", p)
				}
			}()
			// 不管超不超时，本次请求都会执行完毕，或者等到自己超时退出
			ctx.Next()
			// 执行完成之后通知 主协程，我做完了，你可以退出了
			close(finishChan)
		}()

		// 任何一个先触发都会执行，并结束当前函数
		select {
		case pic := <-panicChan:
			// 子G发生异常，抛出传递给上层G
			panic(pic)
		case <-finishChan:
			// 正常退出
			//log.Println("I am back.")
			return
		case <-ctx2.Done():
			// 超时退出
			ctx.ResWrap.WriteHeader(http.StatusServiceUnavailable)
			fst.RaisePanic(midTimeoutMsg)
			return
		}
	}
}
