// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

type radixNode struct {
	match      string
	indices    string
	children   []*radixNode
	wildChild  bool
	nType      uint8
	routerItem *RouterItem
}

func (n *radixNode) regSegment(mTree *methodTree, seg string, ri *RouterItem) {
	isWildSeg := isBeginWildcard(seg)
	// 如果当前节点已经是通配符，而加入的新路由不是 *,: 开头，报错
	ifPanic(n.nType >= param && !isWildSeg, spf("已存在: %s，与新加入: %s冲突\n", n.match, seg))

	if isWildSeg {
		// 当前节点是否是通配符节点
		if n.nType >= param {
			n.matchSameWildcards(mTree, seg, ri)
		} else {
			n.addWildChild(mTree, seg, ri)
		}
		return
	}

	// 找到和当前节点共同的前缀
	i := commonPrefix(seg, n.match)

	// TODO: 当一个字符都没有匹配的情况下，如何呢？
	// ...

	// 这个时候需要拆分当前 n
	if i < len(n.match) {
		// 新建splitSub，作为当前n节点的一个子节点
		splitSub := &radixNode{
			match:      n.match[i:],
			indices:    n.indices,
			children:   n.children,
			routerItem: n.routerItem,
			nType:      n.nType,
			wildChild:  n.wildChild,
		}
		mTree.nodeCt++
		n.children = []*radixNode{splitSub}
		n.indices = string([]byte{n.match[i]})
		n.match = seg[:i]
		n.routerItem = nil
		n.wildChild = false
	}

	// 匹配之后剩余的 seg, 要成为 n 的一个子节点
	if i < len(seg) {
		seg = seg[i:]
		// fix by sdx on 2021.01.06
		// n.regSegment(mTree, seg, ri)
		n.addNormalChild(mTree, seg, ri)
		return
	}

	// 到这里 seg 和 n.match 肯定一样
	// 而且当前节点已经被匹配过，就证明重复了，需要报错
	ifPanic(n.routerItem != nil, spf("此路由：%s 已经存在", n.match))
	n.routerItem = ri
}

// n.nType必须是通配型，seg必须以通配符(:,*)打头
// seg: :name/buy/something
// seg: :name
func (n *radixNode) matchSameWildcards(mTree *methodTree, seg string, ri *RouterItem) {
	ifPanic(seg[0] == '*', spf("通配符%s已经存在，不能再注册%s", n.match, seg))

	part := seg
	var hasSlash bool
	for i := 0; i < len(seg); i++ {
		if seg[i] == '/' {
			hasSlash = true
			part = seg[:i]
			seg = seg[i:]
			break
		}
	}
	ifPanic(n.match != part, spf("通配路由 %s 和 %s 参数名称不一致\n", n.match, part))
	//ifPanic(!hasSlash, spf("路径: %s 和已存在的路由匹配冲突\n", seg))
	// 证明刚好完整匹配参数
	if !hasSlash {
		ifPanic(n.routerItem != nil, spf("路由已经存在：%s\n", n.match))
		n.routerItem = ri
		return
	}

	// 下一个节点有没有都是问题。
	n.addNormalChild(mTree, seg, ri)
}

// 只能添加 wildcard path
func (n *radixNode) addWildChild(mTree *methodTree, seg string, ri *RouterItem) {
	ifPanic(n.children != nil || n.routerItem != nil, spf("已存在%s，不可能再加入%s", n.match, seg))

	emptySub := &radixNode{}
	mTree.nodeCt++
	n.wildChild = true
	n.children = append(n.children, emptySub)
	emptySub.bindSegment(mTree, seg, ri)
}

// 只能添加 非 wildcard path
func (n *radixNode) addNormalChild(mTree *methodTree, seg string, ri *RouterItem) {
	ifPanic(isBeginWildcard(seg), spf("当前路由'%s'与已注册的'%s'冲突\n", seg, n.match))

	// 如果子节点是通配符节点，直接循环下一次
	if n.wildChild {
		n.children[0].regSegment(mTree, seg, ri)
		return
	}

	c := seg[0]
	// 查询当前节点的子节点是否有可能已存在的分支
	lenIdx := len(n.indices)
	for i := 0; i < lenIdx; i++ {
		if c == n.indices[i] {
			n.children[i].regSegment(mTree, seg, ri)
			return
		}
	}

	// 如果没有查找到已有子节点，当前剩余路径就是一个新的子节点
	n.indices += string([]byte{c})
	emptySub := &radixNode{}
	mTree.nodeCt++
	n.children = append(n.children, emptySub)
	emptySub.bindSegment(mTree, seg, ri)
}

// 在当前 空node 中绑定当前 path
// 如果这个 path 含有通配符，需要依次匹配所有这些通配符, 可能的情况：
// path: /xxx/xxx
// path: /xxx/:name/age
// path: wx/:name/xxx/*others
// path: /:name/xxx/*others
// path: :name/xxx/*others ->可能吗？不可能是这种！# 如果出现就是出错了（2021.01.06）
func (n *radixNode) bindSegment(mTree *methodTree, path string, ri *RouterItem) {
	var wildCt, slashCt, lastI int
	var wildParts []string

	for i, sLen := 0, len(path); i < sLen; i++ {
		switch path[i] {
		case ':', '*':
			wildCt++
			if i == 0 {
				slashCt++
			} else if path[i-1] == '/' {
				slashCt++
			}
			// 如果是第一次匹配到通配符
			if wildCt == 1 && i != 0 {
				wildParts = append(wildParts, path[:i])
			} else {
				wildParts = splitSlashPath(wildParts, path[lastI:i])
			}
			lastI = i
		}
	}

	if wildCt == 0 {
		wildParts = append(wildParts, path) // 如果没有通配符，就是整个字符串
	} else {
		wildParts = splitSlashPath(wildParts, path[lastI:]) // 最后一个段记录下来
	}
	ifPanic(wildCt != slashCt, "路由配置有误，通配符标识符只能成段出现")

	// 依次串联起所有的节点
	pNode := n
	lenParts := len(wildParts)
	var tSeg string
	for i := 0; i < lenParts; i++ {
		tSeg = wildParts[i]
		ifPanic(tSeg[0] == '*' && i < (lenParts-1), "通配符*只能出现在最后一段")

		if i == 0 {
			if isBeginWildcard(tSeg) {
				ifPanic(len(tSeg) < 2, "通配符必须设置参数名称")
				pNode.nType = param
				if tSeg[0] == '*' {
					pNode.nType = catchAll
				}
			}
			pNode.match = tSeg
			mTree.nodeStrLen += uint16(len(tSeg))
			pNode.children = nil
			pNode.indices = ""
			continue
		}
		newNode := &radixNode{
			match: tSeg,
		}
		mTree.nodeCt++
		mTree.nodeStrLen += uint16(len(tSeg))
		pNode.children = []*radixNode{newNode}
		if isBeginWildcard(tSeg) {
			ifPanic(len(tSeg) < 2, "通配符必须设置参数名称")
			pNode.wildChild = true
			newNode.nType = param
			if tSeg[0] == '*' {
				newNode.nType = catchAll
			}
		} else {
			pNode.indices += string([]byte{tSeg[0]})
		}
		pNode = newNode
	}
	pNode.routerItem = ri
}

// 将下面这一组解析出的通配符继续拆开
// 例如：
// wild: /xxx/
// wild: :name/print/
func splitSlashPath(parts []string, wild string) []string {
	var isMatch bool
	for i, sLen := 1, len(wild); i < sLen; i++ {
		if wild[i] == '/' {
			isMatch = true
			parts = append(parts, wild[:i], wild[i:])
			break
		}
	}
	if isMatch != true {
		parts = append(parts, wild)
	}
	return parts
}
