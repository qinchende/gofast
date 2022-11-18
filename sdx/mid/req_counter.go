// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/sdx/gate"
	"net/http"
)

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
