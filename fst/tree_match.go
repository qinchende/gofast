// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

type matchVal struct {
	ptrNode *radixMiniNode
	params  Params
	tsr     bool // 是否可以通过重定向，URL最后加入一个 ‘/’ 访问到有处理函数的节点
}

// 在一个函数（作用域）中解决路由匹配的问题，避免函数调用的开销
func (n *radixMiniNode) matchRoute(path string, p Params) (ret matchVal) {
	ret.params = p
nextLoop:
	var pLen = uint8(len(path))
	// 当前是通配符节点
	if n.nType >= param {
		// 找第一个 '/'
		pos := uint8(0)
		hasSlash := false
		for ; pos < pLen; pos++ {
			if path[pos] == '/' {
				hasSlash = true
				break
			}
		}
		if hasSlash && pos == 0 {
			return
		}

		keyName := fstMem.treeChars[n.matchStart+1 : n.matchStart+uint16(n.matchLen)]
		if hasSlash {
			switch n.nType {
			case param:
				ret.params = append(ret.params, Param{Key: keyName, Value: path[:pos]})

				// 匹配后面的节点，后面肯定只能是一个 '/' 开头的节点
				path = path[pos:]
				for id := uint8(0); id < n.childLen; id++ {
					n = &fstMem.allRadixMiniNodes[n.childStart+uint16(id)]
					if fstMem.treeChars[n.matchStart] == path[0] {
						goto nextLoop
					}
				}
				return
			case catchAll:
				return
			default:
				return
			}
		} else {
			// 说明完全匹配当前 通配字段
			ret.params = append(ret.params, Param{Key: keyName, Value: path})
			ret.ptrNode = n
			return
		}
	}

	// 如果当前节点不是通配符
	// 1.1 长度差异，直接不可能
	if pLen < n.matchLen {
		return
	}
	// 1.2 比对每一个字符
	for i := uint8(0); i < n.matchLen; i++ {
		if path[i] == fstMem.treeChars[n.matchStart+uint16(i)] {
			continue
		} else {
			// 匹配直接就失败了
			return
		}
	}
	// 2. 当前节点所有字符都匹配成功，要开始查找下一个可能的节点
	// 2.1 如果完全匹配了，已经找到节点
	if n.matchLen == pLen {
		ret.ptrNode = n
		return
	}
	if n.childLen <= 0 {
		return
	}

	path = path[n.matchLen:]
	var pNode *radixMiniNode
	for id := uint8(0); id < n.childLen; id++ {
		pNode = &fstMem.allRadixMiniNodes[n.childStart+uint16(id)]
		if pNode.nType >= param || fstMem.treeChars[pNode.matchStart] == path[0] {
			n = pNode
			goto nextLoop
		}
	}

	// 所有的情况都没有匹配，直接返回
	return
}
