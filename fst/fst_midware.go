// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "net/http"

// 用于封装，框架自定义一组
func (gft *GoFast) Apply(inject injectFunc) *GoFast {
	return inject(gft)
}

// 添加单个全局拦截器
func (gft *GoFast) Fit(hds FitFunc) *GoFast {
	if hds != nil {
		gft.fitHandlers = append(gft.fitHandlers, hds)
		ifPanic(uint8(len(gft.fitHandlers)) >= maxFits, "Fit handlers more the 255 error.")
	}
	return gft
}

// 将下一级 context 的处理函数，加入fitHandlers 执行链的最后面
func (gft *GoFast) bindContextFit(handler http.HandlerFunc) {
	// 倒序加入
	for i := len(gft.fitHandlers) - 1; i >= 0; i-- {
		handler = gft.fitHandlers[i](handler)
	}
	gft.fitEnter = handler
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 添加一组全局的中间件函数
//func (gft *GoFast) RegHandlers(inject injectFunc) *GoFast {
//	return inject(gft)
//}
//
//// 添加一组全局拦截器
//func (gft *GoFast) RegFits(inject injectFunc) *GoFast {
//	return inject(gft)
//}

//
//// 创建日志模板
//func (gft *GoFast) NewMetricsProject() *stat.Metrics {
//	name := gft.Name
//	if len(name) <= 0 {
//		name = gft.Addr
//	}
//	return stat.NewMetrics(name)
//}
