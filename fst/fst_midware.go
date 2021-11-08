// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 添加一组全局拦截器
func (gft *GoFast) RegFits(gftFunc fitRegFunc) *GoFast {
	return gftFunc(gft)
}

// 添加单个全局拦截器
func (gft *GoFast) Fit(hds IncMiddlewareFunc) *GoFast {
	if hds != nil {
		gft.fitHandlers = append(gft.fitHandlers, hds)
		ifPanic(uint8(len(gft.fitHandlers)) >= maxFits, "Fit handlers more the 255 error.")
	}
	return gft
}

//
//// 执行下一个拦截器
//func (gft *GoFast) NextFit(w http.ResponseWriter, r *http.Request) {
//	for gft.fitIdx < uint8(len(gft.fitHandlers)) {
//		gft.fitHandlers[gft.fitIdx](w, r)
//		gft.fitIdx++
//	}
//}

// 这是构造链式中间件的关键函数
func applyIncMiddleware(h IncHandler, middleware ...IncMiddlewareFunc) IncHandler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func (gft *GoFast) IsAborted() bool {
	return gft.fitIdx >= maxFits
}

func (gft *GoFast) AbortFit() {
	gft.fitIdx = maxFits
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 添加一组全局的中间件函数
func (gft *GoFast) RegHandlers(gftFunc fitRegFunc) *GoFast {
	return gftFunc(gft)
}
