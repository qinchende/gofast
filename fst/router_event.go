// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// Note：该设计给将来预留了足够的扩展空间
// 请求生命周期，设计了不同点的事件类型，这样可以自由 加入 hook
const (
	EAfterMatch = "onAfterMatch" // 初步匹配路由之后，调用这个做更进一步的自定义Check检查
	EBefore     = "onBefore"
	EAfter      = "onAfter"
	EBeforeSend = "onBeforeSend"
	EAfterSend  = "onAfterSend"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 所有注册的 Context handlers 都要通过此函数来注册
// 包括 RouteGroup 和 RouteItem
func (re *routeEvents) regCtxHandler(fstMem *fstMemSpace, eType string, hds []CtxHandler) (*routeEvents, uint16) {
	if len(hds) == 0 || hds[0] == nil {
		return re, 0
	}

	// 如果 hds 里面的有为 nil 的函数，丢弃掉
	tHds := make([]CtxHandler, 0, len(hds))
	for i := range hds {
		if hds[i] != nil {
			tHds = append(tHds, hds[i])
		}
	}

	switch eType {
	case EAfterMatch:
		re.eAfterMatchHds = append(re.eAfterMatchHds, addCtxHandlers(fstMem, tHds)...)
	case EBefore:
		re.eBeforeHds = append(re.eBeforeHds, addCtxHandlers(fstMem, tHds)...)
	case EAfter:
		re.eAfterHds = append(re.eAfterHds, addCtxHandlers(fstMem, tHds)...)
	case EBeforeSend:
		re.eBeforeSendHds = append(re.eBeforeSendHds, addCtxHandlers(fstMem, tHds)...)
	case EAfterSend:
		re.eAfterSendHds = append(re.eAfterSendHds, addCtxHandlers(fstMem, tHds)...)
	default:
		panic("Event type error, can't find this type.")
	}

	return re, uint16(len(tHds))
}

// SpecialRouteGroup
func (gft *GoFast) SpecialBefore(hds ...CtxHandler) {
	gft.specialGroup.Before(hds...)
}

func (gft *GoFast) SpecialAfter(hds ...CtxHandler) {
	gft.specialGroup.After(hds...)
}

// RouteGroup
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gp *RouteGroup) regGroupCtxHandler(eType string, hds []CtxHandler) *RouteGroup {
	_, ct := gp.regCtxHandler(gp.app.fstMem, eType, hds)
	gp.selfHdsLen += ct // 记录分组中一共加入的 各种处理 函数个数
	return gp
}

func (gp *RouteGroup) Before(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EBefore, hds)
}

func (gp *RouteGroup) B(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EBefore, hds)
}

func (gp *RouteGroup) After(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EAfter, hds)
}

func (gp *RouteGroup) A(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EAfter, hds)
}

func (gp *RouteGroup) BeforeSend(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EBeforeSend, hds)
}

func (gp *RouteGroup) AfterSend(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EAfterSend, hds)
}

func (gp *RouteGroup) AfterMatch(hds ...CtxHandler) *RouteGroup {
	return gp.regGroupCtxHandler(EAfterMatch, hds)
}

// RouteItem
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ri *RouteItem) regItemCtxHandler(eType string, hds []CtxHandler) *RouteItem {
	ri.regCtxHandler(ri.group.app.fstMem, eType, hds)
	return ri
}

// 注册节点的所有事件
func (ri *RouteItem) Before(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EBefore, hds)
}

func (ri *RouteItem) B(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EBefore, hds)
}

func (ri *RouteItem) After(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EAfter, hds)
}

func (ri *RouteItem) A(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EAfter, hds)
}

func (ri *RouteItem) BeforeSend(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EBeforeSend, hds)
}

func (ri *RouteItem) AfterSend(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EAfterSend, hds)
}

// 路由匹配到之后，没等执行中间件就走这个逻辑，可以返回标记，中断后面的中间件
func (ri *RouteItem) AfterMatch(hds ...CtxHandler) *RouteItem {
	return ri.regItemCtxHandler(EAfterMatch, hds)
}

// RouteItemAttrs
func (ri *RouteItem) Bind(ra RouteAttrs) *RouteItem {
	ra.BindRoute(ri)
	return ri
}

// RouteItems
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ris RouteItems) Bind(ra RouteAttrs) RouteItems {
	for i := range ris {
		if i == 0 {
			ra.BindRoute(ris[i])
		} else {
			ra.Clone().BindRoute(ris[i])
		}
	}
	return ris
}

func (ris RouteItems) Before(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].regItemCtxHandler(EBefore, hds)
	}
	return ris
}

func (ris RouteItems) B(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].regItemCtxHandler(EBefore, hds)
	}
	return ris
}

func (ris RouteItems) After(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].regItemCtxHandler(EAfter, hds)
	}
	return ris
}

func (ris RouteItems) A(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].regItemCtxHandler(EAfter, hds)
	}
	return ris
}

func (ris RouteItems) BeforeSend(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].regItemCtxHandler(EBeforeSend, hds)
	}
	return ris
}

func (ris RouteItems) AfterSend(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].regItemCtxHandler(EAfterSend, hds)
	}
	return ris
}

func (ris RouteItems) AfterMatch(hds ...CtxHandler) RouteItems {
	for i := range ris {
		ris[i].regItemCtxHandler(EAfterMatch, hds)
	}
	return ris
}
