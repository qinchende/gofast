// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/syncx"
	"net/http"
)

// 限制最大并发连接数，相当于做一个请求资源数量连接池
func MaxConnections(limit int32) fst.FitFunc {
	// 并发数不做限制
	if limit <= 0 {
		return nil
	}

	latch := syncx.Counter{Max: limit}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if latch.TryBorrow() {
				defer func() {
					if err := latch.Return(); err != nil {
						logx.Errorf("Error: MaxConnections return err, info -> %s", err)
					}
				}()
				next(w, r)
			} else {
				logx.Errorf("curr request %d over %d, rejected with code %d", latch.Curr, limit, http.StatusServiceUnavailable)
				// 返回客户端服务器错误
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		}
	}
}
