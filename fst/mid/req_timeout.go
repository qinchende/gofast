package mid

import (
	"context"
	"github.com/qinchende/gofast/fst"
	"net/http"
	"time"
)

const timeoutMsg = "Request Timeout."

// 设置请求处理的超时时间，单位是毫秒
func ReqTimeout(dur time.Duration) fst.IncHandler {
	// 默认所有请求超时是 30 秒钟，再长就直接异常返回
	if dur == 0 {
		dur = 30 * time.Second
	}
	// TODO：Debug模式就不设置超时了
	//if logx.IsDebugging() {
	//	dur = -1
	//}

	return func(w *fst.GFResponse, r *http.Request) {
		if dur <= 0 {
			return
		}

		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		// 启动的协程 没有办法杀死，唯一的办法只能用通道通知他，让协程自己退出
		// 如果协程一直不退出，将会一直占用协程的堆栈内存，并且一直处于GMP的待处理队列，影响整体性能
		go func(finish chan struct{}) {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 不管超不超时，本次请求都会执行完毕，或者等到自己超时退出
			w.NextFit(r)
			// 执行完成之后通知 主协程，我做完了，你可以退出了
			close(finish)
		}(done)

		select {
		case pic := <-panicChan:
			panic(pic)
		case <-done:
			// 正常退出
			//log.Println("I am back.")
		case <-time.After(dur):
			// 超时退出
			//w.ResW.WriteHeader(http.StatusServiceUnavailable)
			fst.RaisePanic(timeoutMsg)
			return
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 设置请求处理的超时时间，单位是毫秒
func ReqTimeoutCtx(dur time.Duration) fst.IncHandler {
	// 默认所有请求超时是 30 秒钟，再长就直接异常返回
	if dur == 0 {
		dur = 30 * time.Second
	}
	// TODO：Debug 模式就不设置超时了
	//if logx.IsDebugging() {
	//	dur = -1
	//}

	return func(w *fst.GFResponse, r *http.Request) {
		if dur <= 0 {
			return
		}

		ctx, cancelCtx := context.WithTimeout(r.Context(), dur)
		defer cancelCtx()

		//r = r.WithContext(ctx)
		done := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			w.NextFit(r)
			close(done)
		}()

		// 主 goroutines 进入下面三种情况的等待，任何一种满足都退出
		select {
		case pic := <-panicChan:
			panic(pic)
		case <-done:
			// 正常结束
		case <-ctx.Done():
			// 超时退出
			//w.ResW.WriteHeader(http.StatusServiceUnavailable)
			fst.RaisePanic(timeoutMsg)
		}
	}
}

func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if duration > 0 {
			return http.TimeoutHandler(next, duration, timeoutMsg)
		} else {
			return next
		}
	}
}
