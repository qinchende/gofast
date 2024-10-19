// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/core/cst"
	"net/http"
)

func (gft *GoFast) Apply(apply appSelfFunc) *GoFast {
	return apply(gft)
}

// 添加单个全局拦截器
func (gft *GoFast) UseHttpHandler(hds HttpHandler) *GoFast {
	if hds != nil {
		gft.httpHandlers = append(gft.httpHandlers, hds)
		cst.PanicIf(uint8(len(gft.httpHandlers)) >= maxHttpHandlers, "Http handlers more the 255.")
	}
	return gft
}

// 将下一级 context 的处理函数，加入 httpHandlers 执行链的最后面
func (gft *GoFast) bindContextHandler(handler http.HandlerFunc) {
	// 倒序加入原始http中间件
	for i := len(gft.httpHandlers) - 1; i >= 0; i-- {
		handler = gft.httpHandlers[i](handler)
	}
	gft.httpEnter = handler
}
