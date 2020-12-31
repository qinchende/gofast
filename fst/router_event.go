package fst

//
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

// 所有注册的 Context handlers 都要通过此函数来注册
func (gp *RouterGroup) eventRegister(eType string, hds CHandlers) *RouterGroup {
	ifPanic(len(hds) <= 0, "there must be at least one handler")

	switch eType {
	case EHandler:
		gp.eHds = append(gp.eHds, addCHandlers(hds)...)
	case EBefore:
		gp.eBeforeHds = append(gp.eBeforeHds, addCHandlers(hds)...)
	case EAfter:
		gp.eAfterHds = append(gp.eAfterHds, addCHandlers(hds)...)

	//case ERequest:
	//	gp.eRequestHds = append(gp.eRequestHds, addCHandlers(hds)...)
	//case EPreParsing:
	//	gp.ePreHandlerHds = append(gp.ePreHandlerHds, addCHandlers(hds)...)
	case EValid:
		gp.eValidHds = append(gp.eValidHds, addCHandlers(hds)...)
	//case EPreHandler:
	//	gp.ePreHandlerHds = append(gp.ePreHandlerHds, addCHandlers(hds)...)
	//case EPreSerialization:
	//	gp.ePreSerializationHds = append(gp.ePreSerializationHds, addCHandlers(hds)...)
	case ESend:
		gp.eSendHds = append(gp.eSendHds, addCHandlers(hds)...)
	//case EAfterSend:
	//	gp.eAfterSendHds = append(gp.eAfterSendHds, addCHandlers(hds)...)
	case EAfterSend:
		gp.eAfterSendHds = append(gp.eAfterSendHds, addCHandlers(hds)...)
	//case ETimeout:
	//	gp.eTimeoutHds = append(gp.eTimeoutHds, addCHandlers(hds)...)
	//case EError:
	//	gp.eErrorHds = append(gp.eErrorHds, addCHandlers(hds)...)
	default:
		panic("Event type error, can't find this type.")
	}

	return gp
}

// 注册节点的所有事件
func (gp *RouterGroup) Before(hds ...CHandler) *RouterGroup {
	return gp.eventRegister(EBefore, hds)
}

func (gp *RouterGroup) After(hds ...CHandler) *RouterGroup {
	return gp.eventRegister(EAfter, hds)
}

func (gp *RouterGroup) Valid(hds ...CHandler) *RouterGroup {
	return gp.eventRegister(EValid, hds)
}

func (gp *RouterGroup) Send(hds ...CHandler) *RouterGroup {
	return gp.eventRegister(ESend, hds)
}

func (gp *RouterGroup) AfterSend(hds ...CHandler) *RouterGroup {
	return gp.eventRegister(EAfterSend, hds)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 所有注册的 Context handlers 都要通过此函数来注册
func (ri *RouterItem) itemRegister(eType string, hds CHandlers) *RouterItem {
	ifPanic(len(hds) <= 0, "there must be at least one handler")

	switch eType {
	case EHandler:
		ri.eHds = append(ri.eHds, addCHandlers(hds)...)
	case EBefore:
		ri.eBeforeHds = append(ri.eBeforeHds, addCHandlers(hds)...)
	case EAfter:
		ri.eAfterHds = append(ri.eAfterHds, addCHandlers(hds)...)

	//case ERequest:
	//	ri.eRequestHds = append(ri.eRequestHds, addCHandlers(hds)...)
	//case EPreParsing:
	//	ri.ePreHandlerHds = append(ri.ePreHandlerHds, addCHandlers(hds)...)
	case EValid:
		ri.eValidHds = append(ri.eValidHds, addCHandlers(hds)...)
	//case EPreHandler:
	//	ri.ePreHandlerHds = append(ri.ePreHandlerHds, addCHandlers(hds)...)
	//case EPreSerialization:
	//	ri.ePreSerializationHds = append(ri.ePreSerializationHds, addCHandlers(hds)...)
	case ESend:
		ri.eSendHds = append(ri.eSendHds, addCHandlers(hds)...)
	//case EAfterSend:
	//	ri.eAfterSendHds = append(ri.eAfterSendHds, addCHandlers(hds)...)
	case EAfterSend:
		ri.eAfterSendHds = append(ri.eAfterSendHds, addCHandlers(hds)...)
	//case ETimeout:
	//	ri.eTimeoutHds = append(ri.eTimeoutHds, addCHandlers(hds)...)
	//case EError:
	//	ri.eErrorHds = append(ri.eErrorHds, addCHandlers(hds)...)
	default:
		panic("Event type error, can't find this type.")
	}

	return ri
}

// 注册节点的所有事件
func (ri *RouterItem) Before(hds ...CHandler) *RouterItem {
	return ri.itemRegister(EBefore, hds)
}

func (ri *RouterItem) After(hds ...CHandler) *RouterItem {
	return ri.itemRegister(EAfter, hds)
}

func (ri *RouterItem) Valid(hds ...CHandler) *RouterItem {
	return ri.itemRegister(EValid, hds)
}

func (ri *RouterItem) Send(hds ...CHandler) *RouterItem {
	return ri.itemRegister(ESend, hds)
}

func (ri *RouterItem) AfterSend(hds ...CHandler) *RouterItem {
	return ri.itemRegister(EAfterSend, hds)
}
