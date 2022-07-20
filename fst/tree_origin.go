// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 目前一共 73 字节
type radixNode struct {
	match     string       // 16字节
	indices   string       // 16字节
	children  []*radixNode // 24字节
	wildChild bool         // 1字节
	nType     uint8        // 8字节
	leafItem  *RouteItem   // 8字节
}

// n.nType必须是通配型，seg必须以通配符(:,*)打头
// seg: :name/buy/something
// seg: :name
func (n *radixNode) matchSameWildcards(mTree *methodTree, seg string, ri *RouteItem) {
	if seg[0] == '*' {
		panic(spf("通配符%s已经存在，不能再注册%s", n.match, seg))
	}

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
	if n.match != part {
		panic(spf("通配路由 %s 和 %s 参数名称不一致\n", n.match, part))
	}
	// 证明刚好完整匹配参数
	if !hasSlash {
		if n.leafItem != nil {
			panic(spf("路由已经存在：%s\n", n.match))
		}
		n.leafItem = ri
		return
	}

	// 下一个节点有没有都是问题。
	n.addNormalChild(mTree, seg, ri)
}

// 只能添加 wildcard path
func (n *radixNode) addWildChild(mTree *methodTree, seg string, ri *RouteItem) {
	// 如果子节点是通配符节点，直接循环下一次
	if n.wildChild {
		n.children[0].regSegment(mTree, n, seg, ri)
		return
	}

	// TODO: 有可能需要加入已有的同样是模糊匹配到的节点，这个时候应该匹配
	if n.children != nil {
		panic(spf("已存在%s，不可能再加入%s", n.match, seg))
	}

	emptySub := &radixNode{}
	mTree.nodeCt++
	n.wildChild = true
	n.children = append(n.children, emptySub)
	emptySub.bindSegment(mTree, seg, ri)
}

// 只能添加 非 wildcard path
func (n *radixNode) addNormalChild(mTree *methodTree, seg string, ri *RouteItem) {
	// 如果子节点是通配符节点，直接循环下一次
	if n.wildChild {
		n.children[0].regSegment(mTree, n, seg, ri)
		return
	}

	c := seg[0]
	// 查询当前节点的子节点是否有可能已存在的分支
	lenIdx := len(n.indices)
	for i := 0; i < lenIdx; i++ {
		if c == n.indices[i] {
			n.children[i].regSegment(mTree, n, seg, ri)
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

// 在当前节点下，解析指定的路由段
// regSegment 和 bindSegment 的区别：是否有现存的n节点
func (n *radixNode) regSegment(mTree *methodTree, nParent *radixNode, seg string, ri *RouteItem) {
	isWildSeg := isBeginWildcard(seg)
	// 如果当前节点已经是通配符，而加入的新路由不是 *,: 开头，报错
	if n.nType >= param && !isWildSeg {
		panic(spf("已存在: %s，与新加入: %s冲突\n", n.match, seg))
	}

	// 1. 新来的是一个通配段
	if isWildSeg {
		// 当前节点是否是通配符节点
		if n.nType >= param {
			n.matchSameWildcards(mTree, seg, ri)
		} else {
			n.addWildChild(mTree, seg, ri)
		}
		return
	}

	// 2. 新来的不是一个通配段，需要找共同的前缀
	i := commonPrefix(seg, n.match)

	// 2.1 没有共同的前缀，他们只可能是兄弟节点。
	// 新来的只可能绑定在 n 节点的父节点上
	if i == 0 {
		IfPanic(nParent == nil, "解析路由错误，找不到父节点")
		nParent.regSegment(mTree, nil, seg, ri)
		return
	}

	// 2.2 有共同的前缀
	// 2.2.1 如果n的当前匹配大于共同前缀。提取出共同的前缀，剩余部分当做是新的n，继续走后面流程
	// 2.2.2 如果n被全部匹配，啥也不做。新来的只能是和n一样，或者是子节点。
	if i < len(n.match) {
		// 新建splitSub，作为当前n节点的一个子节点
		splitSub := &radixNode{
			match:     n.match[i:],
			indices:   n.indices,
			children:  n.children,
			leafItem:  n.leafItem,
			nType:     n.nType,
			wildChild: n.wildChild,
		}
		mTree.nodeCt++
		n.children = []*radixNode{splitSub}
		n.indices = string([]byte{n.match[i]})
		n.match = seg[:i]
		n.leafItem = nil
		n.wildChild = false
	}

	// 2.3.1 匹配之后剩余的 seg, 要成为 n 的一个子节点
	if i < len(seg) {
		seg = seg[i:]
		// fix by sdx on 2021.01.06
		// 子节点有两种可能，即是否带通配符
		isWildSeg = isBeginWildcard(seg)
		if isWildSeg {
			n.addWildChild(mTree, seg, ri)
		} else {
			n.addNormalChild(mTree, seg, ri)
		}
		return
	}

	// 2.3.2 到这里 seg 和 n.match 肯定一样
	// 如果当前节点已经被匹配过，就证明重复了，需要报错
	if n.leafItem != nil {
		panic(spf("此路由：%s 已经存在", n.match))
	}
	n.leafItem = ri
}

// 在当前 空node 中绑定当前 path | 强调当前是空Node
// 如果这个 path 含有通配符，需要依次匹配所有这些通配符, 可能的情况：
// path: /xxx/xxx
// path: /xxx/:name/age
// path: wx/:name/xxx/*others
// path: /:name/xxx/*others
// path: :name/xxx/*others ->可能吗？可能的
func (n *radixNode) bindSegment(mTree *methodTree, path string, ri *RouteItem) {
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
			if wildCt == 1 {
				if i != 0 {
					wildParts = append(wildParts, path[:i])
				}
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
	IfPanic(wildCt != slashCt, "路由配置有误，通配符标识符只能成段出现")

	// 依次串联起所有的节点
	pNode := n
	lenParts := len(wildParts)
	var tSeg string
	for i := 0; i < lenParts; i++ {
		tSeg = wildParts[i]
		IfPanic(tSeg[0] == '*' && i < (lenParts-1), "通配符*只能出现在最后一段")

		if i == 0 {
			if isBeginWildcard(tSeg) {
				IfPanic(len(tSeg) < 2, "通配符必须设置参数名称")
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
			IfPanic(len(tSeg) < 2, "通配符必须设置参数名称")
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
	pNode.leafItem = ri
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
