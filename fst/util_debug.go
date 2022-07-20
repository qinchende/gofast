package fst

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
	"strings"
)

func debugPrintRoute(ri *RouteItem) {
	if !logx.IsDebugging() {
		return
	}

	nuHandlers := len(ri.eHds)
	lastHdsIdx := ri.eHds[nuHandlers-1]
	fun := ri.group.myApp.fstMem.allCtxHandlers[lastHdsIdx]

	logx.DebugPrint("%-6s %-25s --> %s (%d hds)\n", ri.method, ri.fullPath, lang.NameOfFunc(fun), nuHandlers)
}

func debugPrintRouteTree(strTree *strings.Builder) {
	if logx.IsDebugging() {
		logx.Print(strTree)
	}
}
