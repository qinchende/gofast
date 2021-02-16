// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 请求生命周期，设计了不同点的事件类型，这样可以自由 加入 hook
const (
	EPreValid  = "onPreValid"
	EBefore    = "onBefore"
	EHandler   = "onHandler"
	EAfter     = "onAfter"
	EPreSend   = "onPreSend"
	EAfterSend = "onAfterSend"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 所有注册的 Context handlers 都要通过此函数来注册
// 包括 RouterGroup 和 RouterItem
func (re *routeEvents) regCtxHandler(fstMem *fstMemSpace, eType string, hds CtxHandlers) *routeEvents {
	ifPanic(len(hds) <= 0, "there must be at least one handler")

	switch eType {
	case EPreValid:
		re.ePreValidHds = append(re.ePreValidHds, addCtxHandlers(fstMem, hds)...)
	case EBefore:
		re.eBeforeHds = append(re.eBeforeHds, addCtxHandlers(fstMem, hds)...)
	//case EHandler:
	//	re.eHds = append(re.eHds, addCtxHandlers(hds)...)
	case EAfter:
		re.eAfterHds = append(re.eAfterHds, addCtxHandlers(fstMem, hds)...)
	case EPreSend:
		re.ePreSendHds = append(re.ePreSendHds, addCtxHandlers(fstMem, hds)...)
	case EAfterSend:
		re.eAfterSendHds = append(re.eAfterSendHds, addCtxHandlers(fstMem, hds)...)

	default:
		panic("Event type error, can't find this type.")
	}

	return re
}

// RouterGroup
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gp *RouterGroup) regGroupCtxHandler(eType string, hds CtxHandlers) *RouterGroup {
	gp.regCtxHandler(gp.gftApp.fstMem, eType, hds)
	// 记录分组中一共加入的 处理 函数个数
	gp.selfHdsLen += uint16(len(hds))
	return gp
}

//// TODO: 需要注册拦截器处理链
//func (gp *RouterGroup) Use(f func()) *RouterGroup {
//	//return gp.regGroupCtxHandler(EBefore, before)
//	return nil
//}

func (gp *RouterGroup) Before(hds ...CtxHandler) *RouterGroup {
	return gp.regGroupCtxHandler(EBefore, hds)
}

func (gp *RouterGroup) After(hds ...CtxHandler) *RouterGroup {
	return gp.regGroupCtxHandler(EAfter, hds)
}

func (gp *RouterGroup) PreValid(hds ...CtxHandler) *RouterGroup {
	return gp.regGroupCtxHandler(EPreValid, hds)
}

func (gp *RouterGroup) PreSend(hds ...CtxHandler) *RouterGroup {
	return gp.regGroupCtxHandler(EPreSend, hds)
}

func (gp *RouterGroup) AfterSend(hds ...CtxHandler) *RouterGroup {
	return gp.regGroupCtxHandler(EAfterSend, hds)
}

// RouterItem
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ri *RouterItem) regItemCtxHandler(eType string, hds CtxHandlers) *RouterItem {
	ri.regCtxHandler(ri.parent.gftApp.fstMem, eType, hds)
	return ri
}

// 注册节点的所有事件
func (ri *RouterItem) Before(hds ...CtxHandler) *RouterItem {
	return ri.regItemCtxHandler(EBefore, hds)
}

func (ri *RouterItem) After(hds ...CtxHandler) *RouterItem {
	return ri.regItemCtxHandler(EAfter, hds)
}

func (ri *RouterItem) PreValid(hds ...CtxHandler) *RouterItem {
	return ri.regItemCtxHandler(EPreValid, hds)
}

func (ri *RouterItem) PreSend(hds ...CtxHandler) *RouterItem {
	return ri.regItemCtxHandler(EPreSend, hds)
}

func (ri *RouterItem) AfterSend(hds ...CtxHandler) *RouterItem {
	return ri.regItemCtxHandler(EAfterSend, hds)
}
