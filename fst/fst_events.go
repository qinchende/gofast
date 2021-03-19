// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

const (
	//ERoute = "onRoute"
	EReady = "onReady" // server 接收正式请求之前调用
	EClose = "onClose" // server 关闭退出之前调用
)

type appEvents struct {
	eReadyHds AppHandlers
	//eRouteHds AppHandlers
	eCloseHds AppHandlers
}

func (gft *GoFast) execAppHandlers(hds AppHandlers) {
	for i, hLen := 0, len(hds); i < hLen; i++ {
		hds[i](gft)
	}
}

func (gft *GoFast) On(eType string, handles ...AppHandler) {
	switch eType {
	case EReady:
		gft.eReadyHds = append(gft.eReadyHds, handles...)
	//case ERoute:
	//	gft.eRouteHds = append(gft.eRouteHds, handles...)
	case EClose:
		gft.eCloseHds = append(gft.eCloseHds, handles...)
	default:
		panic("Server event type error, can't find this type.")
	}
}

func (gft *GoFast) OnReady(hds ...AppHandler) {
	gft.On(EReady, hds...)
}

func (gft *GoFast) OnClose(hds ...AppHandler) {
	gft.On(EClose, hds...)
}
