package fst

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
)

func DebugPrintRoute(ri *RouterItem, hds CtxHandlers) {
	if logx.IsDebugging() {
		nuHandlers := len(ri.eHds)
		handlerName := lang.NameOfFunc(hds.Last())
		if logx.DebugPrintRouteFunc == nil {
			logx.DebugPrint("%-6s %-25s --> %s (%d handlers)\n", ri.method, ri.fullPath, handlerName, nuHandlers)
		} else {
			logx.DebugPrintRouteFunc(ri.method, ri.fullPath, handlerName, nuHandlers)
		}
	}
}
