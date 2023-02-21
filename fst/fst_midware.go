// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "net/http"

// 用于封装，框架自定义一组
func (gft *GoFast) UseGlobal(inject injectFunc) *GoFast {
	return inject(gft)
}

// 添加单个全局拦截器
func (gft *GoFast) UseHttpHandler(hds HttpHandler) *GoFast {
	if hds != nil {
		gft.httpHandlers = append(gft.httpHandlers, hds)
		GFPanicIf(uint8(len(gft.httpHandlers)) >= maxHttpHandlers, "Http handlers more the 255.")
	}
	return gft
}

// 将下一级 context 的处理函数，加入fitHandlers 执行链的最后面
func (gft *GoFast) bindContextHandler(handler http.HandlerFunc) {
	// 倒序加入
	for i := len(gft.httpHandlers) - 1; i >= 0; i-- {
		handler = gft.httpHandlers[i](handler)
	}
	gft.httpEnter = handler
}
