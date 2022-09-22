// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// node的类型 0.默认 1.根 2.参数 3.通配
// >=2 则可以认为是匹配了参数或者通配
const (
	normal   uint8 = iota // 0. 默认值，普通字符匹配节点
	root                  // 1. 根节点 root
	param                 // 2. 参数 ：
	catchAll              // 3. 通配 *
)

// method路由树 对应数组（切片）
type methodTrees []*methodTree

// 每种method对应一颗独立的路由树
type methodTree struct {
	method     string
	root       *radixNode
	miniRoot   *radixMiniNode
	nodeCt     uint16
	nodeStrLen uint16
}

func (trees methodTrees) getTreeRoot(method string) *radixNode {
	tree := trees.getTree(method)
	if tree != nil {
		return tree.root
	}
	return nil
}

func (trees methodTrees) getTreeMiniRoot(method string) *radixMiniNode {
	tree := trees.getTree(method)
	if tree != nil {
		return tree.miniRoot
	}
	return nil
}

func (trees methodTrees) getTree(method string) *methodTree {
	for _, tree := range trees {
		if tree.method == method {
			return tree
		}
	}
	return nil
}

// 注册一条完整的路由
func (tr *methodTree) regRoute(path string, ri *RouteItem) {
	if tr.root == nil {
		tr.root = &radixNode{nType: root}
		tr.nodeCt++
		tr.root.bindSegment(tr, path, ri)
		return
	}
	tr.root.regSegment(tr, nil, path, ri)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func commonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

func isBeginWildcard(src string) bool {
	switch src[0] {
	case ':', '*':
		return true
	}
	return false
}
