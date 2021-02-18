package fst

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
)

func DebugPrintRoute(httpMethod, absolutePath string, hds CtxHandlers) {
	if logx.IsDebugging() {
		nuHandlers := len(hds)
		handlerName := lang.NameOfFunc(hds.Last())
		if logx.DebugPrintRouteFunc == nil {
			logx.DebugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
		} else {
			logx.DebugPrintRouteFunc(httpMethod, absolutePath, handlerName, nuHandlers)
		}
	}
}
