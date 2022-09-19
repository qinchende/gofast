// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// Note：该设计给将来预留了足够的扩展空间
// 请求生命周期，设计了不同点的事件类型，这样可以自由 加入 hook
const (
	EPreBind   = "onPreBind" // 这个事件暂时不用，没有发现有大的必要
	EBefore    = "onBefore"
	EAfter     = "onAfter"
	EPreSend   = "onPreSend"
	EAfterSend = "onAfterSend"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 所有注册的 Context handlers 都要通过此函数来注册
// 包括 RouteGroup 和 RouteItem
func (re *routeEvents) regCtxHandler(fstMem *fstMemSpace, eType string, hds []CtxHandler) (*routeEvents, uint16) {
	if len(hds) == 0 || hds[0] == nil {
		return re, 0
	}

	// 如果 hds 里面的有为 nil 的函数，丢弃掉
	tHds := make([]CtxHandler, len(hds))
	for _, h := range hds {
		if h != nil {
			tHds = append(tHds, h)
		}
	}

	switch eType {
	case EPreBind:
		re.ePreValidHds = append(re.ePreValidHds, addCtxHandlers(fstMem, hds)...)
	case EBefore:
		re.eBeforeHds = append(re.eBeforeHds, addCtxHandlers(fstMem, hds)...)
	case EAfter:
		re.eAfterHds = append(re.eAfterHds, addCtxHandlers(fstMem, hds)...)
	case EPreSend:
		re.ePreSendHds = append(re.ePreSendHds, addCtxHandlers(fstMem, hds)...)
	case EAfterSend:
		re.eAfterSendHds = append(re.eAfterSendHds, addCtxHandlers(fstMem, hds)...)

	default:
		panic("Event type error, can't find this type.")
	}

	return re, uint16(len(tHds))
}

// RouteGroup
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gp *RouteGroup) regGroupCtxHandler(eType string, hds []CtxHandler) *RouteGroup {
	_, ct := gp.regCtxHandler(gp.myApp.fstMem, eType, hds)
	// 记录分组中一共加入的 处理 函数个数
	gp.selfHdsLen += ct
	return gp
}

func (gp *RouteGroup) Before(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EBefore, hds)
}

func (gp *RouteGroup) After(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EAfter, hds)
}

func (gp *RouteGroup) PreBind(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EPreBind, hds)
}

func (gp *RouteGroup) PreSend(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EPreSend, hds)
}

func (gp *RouteGroup) AfterSend(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EAfterSend, hds)
}

// RouteItem
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ri *RouteItem) regItemCtxHandler(eType string, hds []CtxHandler) *RouteItem {
	ri.regCtxHandler(ri.group.myApp.fstMem, eType, hds)
	return ri
}

// 注册节点的所有事件
func (ri *RouteItem) Before(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EBefore, hds)
}

func (ri *RouteItem) After(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EAfter, hds)
}

func (ri *RouteItem) PreValid(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EPreBind, hds)
}

func (ri *RouteItem) PreSend(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EPreSend, hds)
}

func (ri *RouteItem) AfterSend(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EAfterSend, hds)
}

// RouterItemConfig
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ri *RouteItem) Config(rc RouteConfig) *RouteItem {
	rc.AddToList(ri.routeIdx)
	return ri
}
