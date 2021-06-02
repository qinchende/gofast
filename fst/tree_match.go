// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

type matchResult struct {
	ptrNode *radixMiniNode
	params  Params
	needRTS bool // 是否需要做 RedirectTrailingSlash 的检测
	rts     bool // 是否可以通过重定向，URL最后加或减一个 ‘/’ 访问到有处理函数的节点
}

// 在一个函数（作用域）中解决路由匹配的问题，加快匹配速度
func (n *radixMiniNode) matchRoute(fstMem *fstMemSpace, path string, mr *matchResult) {
nextLoop:
	var pLen = uint8(len(path))
	var pNode *radixMiniNode

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

			mr.params = append(mr.params, Param{Key: keyName, Value: path[:pos]})
			// 匹配后面的节点，后面肯定只能是一个 '/' 开头的节点
			path = path[pos:]
			// 看看子节点有没有一个能匹配上'/'字符，有就进入下一轮循环
			for id := uint16(0); id < uint16(n.childLen); id++ {
				n = &fstMem.allRadixMiniNodes[n.childStart+id]
				if fstMem.treeChars[n.matchStart] == path[0] {
					goto nextLoop
				}
			}
			// 没有匹配就要返回没匹配到路由了，不过这里可以看看是否能重定向
			goto checkRTS
		}

	mathRestPath:
		// 说明完全匹配当前url段
		mr.params = append(mr.params, Param{Key: keyName, Value: path})
		mr.ptrNode = n
		return
	}

	// PartB
	// 如果当前节点不是 模糊匹配
	// +++++++++++++++++++++++++++++++++++++++++++++
	// 1.1 长度差异，直接不可能
	if pLen < n.matchLen {
		return
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
	for id := uint16(0); id < uint16(n.childLen); id++ {
		pNode = &fstMem.allRadixMiniNodes[n.childStart+id]
		// 子节点是 模糊匹配节点（:*） | 首字符匹配的普通节点，就走这个逻辑
		if pNode.nType >= param || fstMem.treeChars[pNode.matchStart] == path[0] {
			n = pNode
			goto nextLoop
		}
	}
	// 2.4 上面没有能匹配子串的节点，说明请求url，没有找到能匹配的路由
	// 但是可以看是否有 RedirectTrailingSlash 可以匹配
checkRTS:
	if mr.needRTS {
		mr.rts = path == "/" && n.hdsItemIdx != -1
	}
}
