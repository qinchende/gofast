// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"net/http"
)

// 添加一组全局拦截器
func (gft *GoFast) Fits(gftFunc fitRegFunc) *GoFast {
	return gftFunc(gft)
}

// 添加单个全局拦截器
func (gft *GoFast) Fit(hds IncHandler) *GoFast {
	if hds != nil {
		gft.fitHandlers = append(gft.fitHandlers, hds)
		ifPanic(len(gft.fitHandlers) >= maxFitLen, "Fit handlers more the 127 error.")
	}
	return gft
}

// 执行下一个拦截器
func (gft *GoFast) NextFit(w http.ResponseWriter, r *http.Request) {
	gft.fitIdx++
	for gft.fitIdx < len(gft.fitHandlers) {
		gft.fitHandlers[gft.fitIdx](w, r)
		gft.fitIdx++
	}
}

func (gft *GoFast) IsAborted() bool {
	return gft.fitIdx >= maxFitLen
}

func (gft *GoFast) AbortFit() {
	gft.fitIdx = maxFitLen
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 添加一组全局的中间件函数
func (gft *GoFast) RegHandlers(gftFunc fitRegFunc) *GoFast {
	return gftFunc(gft)
}
