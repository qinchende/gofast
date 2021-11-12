// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"bytes"
	"fmt"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/render"
	"github.com/qinchende/gofast/logx"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// 默认的异常处理函数，到这里了，证明err肯定发生了。
func recoveryHandler(ctx *fst.Context, err interface{}) {
	ctx.AbortWithStatus(http.StatusInternalServerError)
	ctx.ErrorN(err)

	// 默认返回 JSON 格式的结果
	jsonData := fst.NewRenderKV("fai", fmt.Sprintf("%v", err), 0)
	// todo: 这里可能也会抛出异常
	_ = render.WriteJSON(ctx.ResWrap, jsonData)
}

func Recovery() fst.CtxHandler {
	var logger *log.Logger
	if logx.DefErrorWriter != nil {
		logger = log.New(logx.DefErrorWriter, "\n\n\x1b[31m", log.LstdFlags)
	}

	return func(ctx *fst.Context) {
		defer func() {
			// 没有捕获错误，啥也不用做
			err := recover()
			if err == nil {
				return
			}

			// 如果是框架主动panic，只打印简单错误日志
			// 说明是程序自己触发的异常，带有一定的预见性(类似退出当前协程等)。
			if _, ok := err.(fst.GFPanic); ok {
				recoveryHandler(ctx, err.(error))
				ctx.AbortBehind()
				return
			}

			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
						strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			if logger != nil {
				stackInfo := stackCollect(3)
				httpRequest, _ := httputil.DumpRequest(ctx.ReqRaw, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				headersToStr := strings.Join(headers, "\r\n")
				if brokenPipe {
					logger.Printf("%s\n%s%s", err, headersToStr, logx.Reset)
				} else if logx.IsDebugging() {
					logx.Infof("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
						time.Now().Format("2006/01/02 - 15:04:05"), headersToStr, err, stackInfo, logx.Reset)
				} else {
					logger.Printf("[Recovery] %s panic recovered:\n%s\n%s%s",
						time.Now().Format("2006/01/02 - 15:04:05"), err, stackInfo, logx.Reset)
				}
			}
			if brokenPipe {
				// If the connection is dead, we can't write a status to it.
				ctx.ErrorN(err) // nolint: errcheck
			} else {
				recoveryHandler(ctx, err)
			}
			ctx.AbortBehind()
		}()

		// 有 defer 函数，所以这里的 ctx.Next() 有意义。
		ctx.Next()
	}
}

//
//// 默认的异常处理函数，到这里了，证明err肯定发生了。
//func recoveryHandler(w *fst.GFResponse, err interface{}) {
//	w.AbortWithStatus(http.StatusInternalServerError)
//	w.ErrorN(err)
//
//	// 默认返回 JSON 格式的结果
//	jsonData := fst.NewRenderKV("fai", fmt.Sprintf("%v", err), 0)
//	// todo: 这里可能也会抛出异常
//	_ = render.WriteJSON(w.ResWrap, jsonData)
//}

//func Recovery() http.HandlerFunc {
//	var logger *log.Logger
//	if logx.DefErrorWriter != nil {
//		logger = log.New(logx.DefErrorWriter, "\n\n\x1b[31m", log.LstdFlags)
//	}
//
//	return func(w *fst.GFResponse, r *http.Request) {
//		defer func() {
//			// 没有捕获错误，啥也不用做
//			err := recover()
//			if err == nil {
//				return
//			}
//
//			// 如果是框架主动panic，只打印简单错误日志
//			// 说明是程序自己触发的异常，带有一定的预见性(类似退出当前协程等)。
//			if _, ok := err.(fst.GFPanic); ok {
//				recoveryHandler(w, err.(error))
//				w.AbortFit()
//				return
//			}
//
//			// Check for a broken connection, as it is not really a
//			// condition that warrants a panic stack trace.
//			var brokenPipe bool
//			if ne, ok := err.(*net.OpError); ok {
//				if se, ok := ne.Err.(*os.SyscallError); ok {
//					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
//						strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
//						brokenPipe = true
//					}
//				}
//			}
//			if logger != nil {
//				stackInfo := stackCollect(3)
//				httpRequest, _ := httputil.DumpRequest(r, false)
//				headers := strings.Split(string(httpRequest), "\r\n")
//				for idx, header := range headers {
//					current := strings.Split(header, ":")
//					if current[0] == "Authorization" {
//						headers[idx] = current[0] + ": *"
//					}
//				}
//				headersToStr := strings.Join(headers, "\r\n")
//				if brokenPipe {
//					logger.Printf("%s\n%s%s", err, headersToStr, logx.Reset)
//				} else if logx.IsDebugging() {
//					logx.Infof("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
//						time.Now().Format("2006/01/02 - 15:04:05"), headersToStr, err, stackInfo, logx.Reset)
//				} else {
//					logger.Printf("[Recovery] %s panic recovered:\n%s\n%s%s",
//						time.Now().Format("2006/01/02 - 15:04:05"), err, stackInfo, logx.Reset)
//				}
//			}
//			if brokenPipe {
//				// If the connection is dead, we can't write a status to it.
//				w.ErrorN(err) // nolint: errcheck
//			} else {
//				recoveryHandler(w, err)
//			}
//			w.AbortFit()
//		}()
//
//		w.NextFit(r)
//	}
//}

// 下面这些方法用于收集异常堆栈信息
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// stack returns a nicely formatted stack frame, skipping skip frames.
func stackCollect(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", stackFuncName(pc), stackSource(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func stackSource(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func stackFuncName(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
