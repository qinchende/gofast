// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"github.com/qinchende/gofast/skill/sysx"
	"net/http"
)

func HttpHighMemProtect(kp *gate.RequestKeeper, pos uint16) fst.HttpHandler {
	// 前提是必须启动系统资源自动监控
	if kp == nil || sysx.MonitorStarted == false {
		return nil
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// 执行后面的处理函数
			next(w, r)
		}
	}
}
