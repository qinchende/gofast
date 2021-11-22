// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

const (
	EBeforeBuildRoutes = "onBeforeBuildRoutes" // 开始重构路由树前
	EAfterBuildRoutes  = "onAfterBuildRoutes"  // 重构路由树后
	EReady             = "onReady"             // server 接收正式请求之前调用
	EClose             = "onClose"             // server 关闭退出之前调用
)

type appEvents struct {
	eBeforeBuildRoutesHds []AppHandler
	eAfterBuildRoutesHds  []AppHandler
	eReadyHds             []AppHandler
	eCloseHds             []AppHandler
}

func (gft *GoFast) execAppHandlers(hds []AppHandler) {
	for i, hLen := 0, len(hds); i < hLen; i++ {
		hds[i](gft)
	}
}

func (gft *GoFast) On(eType string, handles ...AppHandler) {
	switch eType {
	case EBeforeBuildRoutes:
		gft.eBeforeBuildRoutesHds = append(gft.eBeforeBuildRoutesHds, handles...)
	case EAfterBuildRoutes:
		gft.eAfterBuildRoutesHds = append(gft.eAfterBuildRoutesHds, handles...)
	case EReady:
		gft.eReadyHds = append(gft.eReadyHds, handles...)
	case EClose:
		gft.eCloseHds = append(gft.eCloseHds, handles...)
	default:
		panic("Server event type error, can't find this type.")
	}
}

func (gft *GoFast) OnBeforeBuildRoutes(hds ...AppHandler) {
	gft.On(EBeforeBuildRoutes, hds...)
}

func (gft *GoFast) OnAfterBuildRoutes(hds ...AppHandler) {
	gft.On(EAfterBuildRoutes, hds...)
}

func (gft *GoFast) OnReady(hds ...AppHandler) {
	gft.On(EReady, hds...)
}

func (gft *GoFast) OnClose(hds ...AppHandler) {
	gft.On(EClose, hds...)
}
