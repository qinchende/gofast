// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

// 请求生命周期，设计了不同点的事件类型，这样可以自由 加入 hook
const (
	EValid     = "onValid"
	EBefore    = "onBefore"
	EHandler   = "onHandler"
	EAfter     = "onAfter"
	ESend      = "onSend"
	EAfterSend = "onAfterSend"

	EPreHandler       = "onPreHandler"
	EPreSerialization = "onPreSerialization"
	ERequest          = "onRequest"
	EPreParsing       = "onPreParsing"
	ETimeout          = "onTimeout"
	EError            = "onError"
)

// 分组只能通过这里注册事件处理函数
// 所有注册的 Context handlers 都要通过此函数来注册
func (gp *RouterGroup) eventRegister(eType string, hds CtxHandlers) *RouterGroup {
	ifPanic(len(hds) <= 0, "there must be at least one handler")

	switch eType {
	case EHandler:
		gp.eHds = append(gp.eHds, addCtxHandlers(hds)...)
	case EBefore:
		gp.eBeforeHds = append(gp.eBeforeHds, addCtxHandlers(hds)...)
	case EAfter:
		gp.eAfterHds = append(gp.eAfterHds, addCtxHandlers(hds)...)

	//case ERequest:
	//	gp.eRequestHds = append(gp.eRequestHds, addCtxHandlers(hds)...)
	//case EPreParsing:
	//	gp.ePreHandlerHds = append(gp.ePreHandlerHds, addCtxHandlers(hds)...)
	case EValid:
		gp.eValidHds = append(gp.eValidHds, addCtxHandlers(hds)...)
	//case EPreHandler:
	//	gp.ePreHandlerHds = append(gp.ePreHandlerHds, addCtxHandlers(hds)...)
	//case EPreSerialization:
	//	gp.ePreSerializationHds = append(gp.ePreSerializationHds, addCtxHandlers(hds)...)
	case ESend:
		gp.eSendHds = append(gp.eSendHds, addCtxHandlers(hds)...)
	//case EAfterSend:
	//	gp.eAfterSendHds = append(gp.eAfterSendHds, addCtxHandlers(hds)...)
	case EAfterSend:
		gp.eAfterSendHds = append(gp.eAfterSendHds, addCtxHandlers(hds)...)
	//case ETimeout:
	//	gp.eTimeoutHds = append(gp.eTimeoutHds, addCtxHandlers(hds)...)
	//case EError:
	//	gp.eErrorHds = append(gp.eErrorHds, addCtxHandlers(hds)...)
	default:
		panic("Event type error, can't find this type.")
	}

	// 记录分组中一共加入的 处理 函数个数
	gp.selfHdsLen += uint16(len(hds))
	return gp
}

// 注册节点的所有事件
func (gp *RouterGroup) Before(hds ...CtxHandler) *RouterGroup {
	return gp.eventRegister(EBefore, hds)
}

func (gp *RouterGroup) After(hds ...CtxHandler) *RouterGroup {
	return gp.eventRegister(EAfter, hds)
}

func (gp *RouterGroup) Valid(hds ...CtxHandler) *RouterGroup {
	return gp.eventRegister(EValid, hds)
}

func (gp *RouterGroup) Send(hds ...CtxHandler) *RouterGroup {
	return gp.eventRegister(ESend, hds)
}

func (gp *RouterGroup) AfterSend(hds ...CtxHandler) *RouterGroup {
	return gp.eventRegister(EAfterSend, hds)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 所有注册的 Context handlers 都要通过此函数来注册
func (ri *RouterItem) itemRegister(eType string, hds CtxHandlers) *RouterItem {
	ifPanic(len(hds) <= 0, "there must be at least one handler")

	switch eType {
	case EHandler:
		ri.eHds = append(ri.eHds, addCtxHandlers(hds)...)
	case EBefore:
		ri.eBeforeHds = append(ri.eBeforeHds, addCtxHandlers(hds)...)
	case EAfter:
		ri.eAfterHds = append(ri.eAfterHds, addCtxHandlers(hds)...)

	//case ERequest:
	//	ri.eRequestHds = append(ri.eRequestHds, addCtxHandlers(hds)...)
	//case EPreParsing:
	//	ri.ePreHandlerHds = append(ri.ePreHandlerHds, addCtxHandlers(hds)...)
	case EValid:
		ri.eValidHds = append(ri.eValidHds, addCtxHandlers(hds)...)
	//case EPreHandler:
	//	ri.ePreHandlerHds = append(ri.ePreHandlerHds, addCtxHandlers(hds)...)
	//case EPreSerialization:
	//	ri.ePreSerializationHds = append(ri.ePreSerializationHds, addCtxHandlers(hds)...)
	case ESend:
		ri.eSendHds = append(ri.eSendHds, addCtxHandlers(hds)...)
	//case EAfterSend:
	//	ri.eAfterSendHds = append(ri.eAfterSendHds, addCtxHandlers(hds)...)
	case EAfterSend:
		ri.eAfterSendHds = append(ri.eAfterSendHds, addCtxHandlers(hds)...)
	//case ETimeout:
	//	ri.eTimeoutHds = append(ri.eTimeoutHds, addCtxHandlers(hds)...)
	//case EError:
	//	ri.eErrorHds = append(ri.eErrorHds, addCtxHandlers(hds)...)
	default:
		panic("Event type error, can't find this type.")
	}

	return ri
}

// 注册节点的所有事件
func (ri *RouterItem) Before(hds ...CtxHandler) *RouterItem {
	return ri.itemRegister(EBefore, hds)
}

func (ri *RouterItem) After(hds ...CtxHandler) *RouterItem {
	return ri.itemRegister(EAfter, hds)
}

func (ri *RouterItem) Valid(hds ...CtxHandler) *RouterItem {
	return ri.itemRegister(EValid, hds)
}

func (ri *RouterItem) Send(hds ...CtxHandler) *RouterItem {
	return ri.itemRegister(ESend, hds)
}

func (ri *RouterItem) AfterSend(hds ...CtxHandler) *RouterItem {
	return ri.itemRegister(EAfterSend, hds)
}
