package fst

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
	"strings"
)

func debugPrintRoute(gft *GoFast, ri *RouteItem) {
	if !gft.IsDebugging() {
		return
	}

	nuHandlers := len(ri.eHds)
	lastHdsIdx := ri.eHds[nuHandlers-1]
	fun := ri.group.myApp.fstMem.allCtxHandlers[lastHdsIdx]

	logx.DebugF("%-6s %-25s --> %s (%d hds)\n", ri.method, ri.fullPath, lang.NameOfFunc(fun), nuHandlers)
}

func debugPrintRouteTree(gft *GoFast, strTree *strings.Builder) {
	if gft.IsDebugging() {
		logx.DebugDirect(strTree.String())
	}
}
