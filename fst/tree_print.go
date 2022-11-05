// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"fmt"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
	"strings"
)

// 打印 路由树
func (gft *GoFast) printRouteTrees() {
	strTree := new(strings.Builder)
	strTree.WriteString("\n+++++++++++++++The route tree:\n")

	for _, tree := range gft.routerTrees {
		printTree(tree, strTree)
	}
	strTree.WriteString("++++++++++++++++++++++++++++++\n")
	// 打印到控制台
	debugPrintRouteTree(gft, strTree)
}

func printTree(tree *methodTree, strTree *strings.Builder) {
	if tree.root == nil {
		return
	}
	strTree.WriteString("(")
	strTree.WriteString(tree.method)
	strTree.WriteString(")\n")
	tree.root.prettyPrint(strTree, "", true)
}

func (n *radixNode) prettyPrint(str *strings.Builder, prefix string, isTail bool) {
	str.WriteString(prefix)

	nextPrefix := prefix
	if isTail {
		str.WriteString("└── ")
		nextPrefix += "    "
	} else {
		str.WriteString("├── ")
		nextPrefix += "│   "
	}

	// 要显示的节点内容
	str.WriteString(n.match)
	curLen := len([]rune(prefix)) + len([]rune(n.match))
	// 缩进最大是160字符占位符
	retract := 60 - curLen
	for i := 1; i <= 5; i++ {
		if retract >= 0 {
			break
		}
		retract += 20
	}
	if retract < 0 {
		retract = 0
	}
	strFmt := "%-" + fmt.Sprintf("%ds", retract)
	str.WriteString(fmt.Sprintf(strFmt, ""))
	// [优先级，动态匹配参数数量，handler数量，所有子节点首字符]
	// [4-0-0-im]
	//genPrintNode(str, []string{fmt.Sprint(n.priority), fmt.Sprint(n.maxParams), fmt.Sprint(len(n.handlers)), n.indices})
	//genPrintNode(str, []string{fmt.Sprint(len(n.hdsItem)), n.indices})
	genPrintNode(str, []string{fmt.Sprint(n.leafItem != nil), n.indices})
	//genPrintNode(str, []string{n.indices})

	chLen := len(n.children)
	for i := 0; i < chLen-1; i++ {
		n.children[i].prettyPrint(str, nextPrefix, false)
	}
	if chLen > 0 {
		n.children[chLen-1].prettyPrint(str, nextPrefix, true)
	}
}

// 打印节点对应的相关值
func genPrintNode(str *strings.Builder, arr []string) {
	str.WriteString(" [")
	if arr != nil && len(arr) > 0 {
		if arr[0] == "true" {
			str.WriteString("1")
		} else {
			str.WriteString("0")
		}
	}
	for i := 1; i < len(arr); i++ {
		if arr[i] != "" {
			str.WriteString("-")
			str.WriteString(arr[i])
		}
	}
	str.WriteString("]\n")
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// debug print
func debugPrintRoute(gft *GoFast, ri *RouteItem) {
	if !gft.IsDebugging() {
		return
	}

	nuHandlers := len(ri.eHds)
	lastHdsIdx := ri.eHds[nuHandlers-1]
	fun := ri.group.myApp.fstMem.allCtxHandlers[lastHdsIdx]

	logx.DebugF("%-6s %-25s %s (%d hds)", ri.method, ri.fullPath, lang.NameOfFunc(fun), nuHandlers)
}

func debugPrintRouteTree(gft *GoFast, strTree *strings.Builder) {
	if gft.IsDebugging() {
		logx.DebugDirect(strTree.String())
	}
}
