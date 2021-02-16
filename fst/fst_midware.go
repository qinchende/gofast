package fst

import "net/http"

// 添加全局拦截器
func (gft *GoFast) Fit(hds ...IncHandler) *GoFast {
	gft.fitHandlers = append(gft.fitHandlers, hds...)
	ifPanic(len(gft.fitHandlers) >= maxFitLen, "Fit handlers more the 127 error.")
	return gft
}

// 执行下一个拦截器
func (r *Request) NextFit(w http.ResponseWriter) {
	r.fitIdx++
	for r.fitIdx < len(r.gftApp.fitHandlers) {
		r.gftApp.fitHandlers[r.fitIdx](w, r)
		r.fitIdx++
	}
}

func (r *Request) IsAborted() bool {
	return r.fitIdx >= maxFitLen
}

func (r *Request) Abort() {
	r.fitIdx = maxFitLen
}
