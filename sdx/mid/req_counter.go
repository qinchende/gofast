// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/skill/syncx"
	"net/http"
)

// 限制最大并发连接数，相当于做一个请求资源数量连接池
func HttpMaxConnections(limit int32) fst.HttpHandler {
	// 并发数不做限制
	if limit <= 0 {
		return nil
	}

	latch := syncx.LazyCounter{Max: limit}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if latch.TryBorrow() {
				defer func() {
					if err := latch.Return(); err != nil {
						logx.ErrorF("Error: MaxConnections return err, info -> %s", err)
					}
				}()
				next(w, r)
			} else {
				logx.ErrorF("req %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
				w.WriteHeader(http.StatusServiceUnavailable) // 返回客户端服务器错误
			}
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 访问计数
func HttpReqCountPos(kp *gate.RequestKeeper, pos uint16) fst.HttpHandler {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			kp.CountExtras(pos)
			next(w, r)
		}
	}
}

func ReqCountPos(kp *gate.RequestKeeper, idx uint16) fst.CtxHandler {
	return func(c *fst.Context) {
		kp.CountExtras(idx)
		c.Next()
	}
}

func ReqCount(kp *gate.RequestKeeper) fst.CtxHandler {
	return func(c *fst.Context) {
		kp.CountRoutePass(c.RouteIdx)
		c.Next()
	}
}
