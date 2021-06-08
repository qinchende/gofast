package fst

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
	"strings"
)

func debugPrintRoute(ri *RouterItem) {
	if logx.IsDebugging() {
		nuHandlers := len(ri.eHds)
		lastHdsIdx := ri.eHds[nuHandlers-1]
		fun := ri.group.gftApp.fstMem.allCtxHandlers[lastHdsIdx]

		handlerName := lang.NameOfFunc(fun)
		if logx.DebugPrintRouteFunc == nil {
			logx.DebugPrint("%-6s %-25s --> %s (%d handlers)\n", ri.method, ri.fullPath, handlerName, nuHandlers)
		} else {
			logx.DebugPrintRouteFunc(ri.method, ri.fullPath, handlerName, nuHandlers)
		}
	}
}

func debugPrintRouteTree(strTree *strings.Builder) {
	if logx.IsDebugging() {
		logx.Info(strTree)
	}
}
