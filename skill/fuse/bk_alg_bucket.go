// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fuse

import "github.com/qinchende/gofast/skill/collect"

// 滑动窗口的数据结构，统计请求信息
type BreakBucket struct {
	total   int64   // 总访问次数
	accepts float64 // 处理成功的次数
}

func (bk *BreakBucket) Reset() {
	bk.total = 0
	bk.accepts = 0
}

func (bk *BreakBucket) Add(v float64) {
	bk.total++
	bk.accepts += v
}

func (bk *BreakBucket) AddByFlag(v float64, flag int8) {
}

func (bk *BreakBucket) Discount(past collect.SlideWinBucket) {
	pbk := past.(*BreakBucket)
	bk.total -= pbk.total
	bk.accepts -= pbk.accepts
}
