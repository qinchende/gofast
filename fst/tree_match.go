// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "net/url"

type matchRoute struct {
	ptrNode *radixMiniNode
	params  *routeParams
	rts     bool // 判断是否可以 RedirectTrailingSlash 找到节点
}

// 在一个函数（作用域）中解决路由匹配的问题，加快匹配速度
func (n *radixMiniNode) matchRoute(fstMem *fstMemSpace, path string, mr *matchRoute, unescape bool) {
	var pLen uint8
	var pNode *radixMiniNode

nextLoop:
	pLen = uint8(len(path))

	// PartA
	// +++++++++++++++++++++++++++++++++++++++++++++
	// 如果当前节点是 模糊匹配 节点，可能是 : 或 *
	if n.nType >= param {
		keyName := fstMem.treeChars[n.matchStart+1 : n.matchStart+uint16(n.matchLen)]

		switch n.nType {
		case catchAll:
			goto mathRestPath
		case param:
			// 找第一个 '/'
			pos := uint8(0)
			hasSlash := false
			for ; pos < pLen; pos++ {
				if path[pos] == '/' {
					hasSlash = true
					break
				}
			}
			// 完全匹配后面的所有字符，这和通配符*逻辑一样了
			if !hasSlash {
				goto mathRestPath
			} else if pos == 0 {
				// 参数匹配：居然一个字符都没有匹配到，直接返回，没找到
				return
			}

			pVal := path[:pos]
			if unescape {
				if v, err := url.QueryUnescape(pVal); err == nil {
					pVal = v
				}
			}
			*mr.params = append(*mr.params, UrlParam{Key: keyName, Value: pVal})
			// 匹配后面的节点，后面肯定只能是一个 '/' 开头的节点
			path = path[pos:]
			goto matchChildNode
		}
	mathRestPath:
		if unescape {
			if v, err := url.QueryUnescape(path); err == nil {
				path = v
			}
		}
		// 说明完全匹配当前url段
		*mr.params = append(*mr.params, UrlParam{Key: keyName, Value: path})
		mr.ptrNode = n
		return
	}

	// PartB
	// 如果当前节点不是 模糊匹配
	// +++++++++++++++++++++++++++++++++++++++++++++
	// 1.1 长度差异，直接不可能
	if pLen < n.matchLen {
		// 有可能 path + '/' 能够匹配到一个路由，去匹配一下重定向的逻辑
		goto checkRTS
	}
	// 1.2 比对每一个字符
	for i := uint16(0); i < uint16(n.matchLen); i++ {
		if path[i] == fstMem.treeChars[n.matchStart+i] {
			continue
		} else {
			// 当前普通节点的字符序列都不能匹配完，匹配直接就失败了
			return
		}
	}
	// 2. 当前节点所有字符都匹配成功，要开始查找下一个可能的节点
	// 2.1 如果完全匹配了，而且当前节点对应一个路由处理函数，已经找到节点
	if pLen == n.matchLen && n.hdsItemIdx != -1 {
		mr.ptrNode = n
		return
	}
	//// 2.2 当前节点没有子节点了，无法匹配，但是有可能重定向
	//if n.childLen <= 0 {
	//	return
	//}
	// 2.3 查找可能的子节点，再次循环，匹配后面的路径
	path = path[n.matchLen:]

matchChildNode:
	pLen = uint8(len(path))
	for id := uint16(0); id < uint16(n.childLen); id++ {
		pNode = &fstMem.allRadixMiniNodes[n.childStart+id]
		// 子节点是 模糊匹配节点（:*） | 首字符匹配的普通节点，就走这个逻辑
		if pNode.nType >= param || (pLen > 0 && fstMem.treeChars[pNode.matchStart] == path[0]) {
			n = pNode
			goto nextLoop
		}
	}

checkRTS:
	// 是否允许重定向，是否有 RedirectTrailingSlash 可以匹配
	if fstMem.myApp.WebConfig.RedirectTrailingSlash {
		if pLen == 0 {
			for id := uint16(0); id < uint16(n.childLen); id++ {
				pNode = &fstMem.allRadixMiniNodes[n.childStart+id]
				mr.rts = pNode.hdsItemIdx != -1 && pNode.matchLen == 1 && fstMem.treeChars[pNode.matchStart] == '/'
				if mr.rts {
					break
				}
			}
		} else if path == "/" {
			mr.rts = n.hdsItemIdx != -1
		} else {
			idx := n.matchStart + uint16(n.matchLen) - 1
			mr.rts = n.hdsItemIdx != -1 && fstMem.treeChars[idx] == '/' && path == fstMem.treeChars[n.matchStart:idx]
		}
	}
}
