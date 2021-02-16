// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

const (
	EReady = "onReady"
	//ERoute = "onRoute"
	EClose = "onClose"
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
