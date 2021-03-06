package mid

import (
	"context"
	"github.com/qinchende/gofast/fst"
	"net/http"
	"time"
)

const timeoutMsg = "[GoFast]Request Timeout."

func doBadthing(done chan bool) {
	time.Sleep(time.Second)
	done <- true
}

// 设置请求处理的超时时间，单位是毫秒
func ReqTimeout(dur time.Duration) fst.IncHandler {
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

		//done := make(chan bool)
		//go f(done)
		//select {
		//case <-done:
		//	fmt.Println("done")
		//	return nil
		//case <-time.After(dur):
		//	// 超时退出
		//	w.ResW.WriteHeader(http.StatusServiceUnavailable)
		//	fst.RaisePanic(timeoutMsg)
		//	return
		//}

	}
}

// 设置请求处理的超时时间，单位是毫秒
func ReqTimeoutBack(dur time.Duration) fst.IncHandler {
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
			w.ResW.WriteHeader(http.StatusServiceUnavailable)
			fst.RaisePanic(timeoutMsg)
		}
	}
}

//
//func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		if duration > 0 {
//			return http.TimeoutHandler(next, duration, timeoutMsg)
//		} else {
//			return next
//		}
//	}
//}
