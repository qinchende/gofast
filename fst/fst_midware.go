// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"net/http"
)

// GoFast提供的拦截器全家桶
func (gft *GoFast) RegFits(gftFunc goFastRegFunc) *GoFast {
	return gftFunc(gft)
}

// 添加全局拦截器
func (gft *GoFast) Fit(hds ...IncHandler) *GoFast {
	gft.fitHandlers = append(gft.fitHandlers, hds...)
	ifPanic(len(gft.fitHandlers) >= maxFitLen, "Fit handlers more the 127 error.")
	return gft
}

// 执行下一个拦截器
func (w *GFResponse) NextFit(r *http.Request) {
	w.fitIdx++
	for w.fitIdx < len(w.gftApp.fitHandlers) {
		w.gftApp.fitHandlers[w.fitIdx](w, r)
		w.fitIdx++
	}
}

func (w *GFResponse) IsAborted() bool {
	return w.fitIdx >= maxFitLen
}

func (w *GFResponse) AbortFit() {
	w.fitIdx = maxFitLen
}
