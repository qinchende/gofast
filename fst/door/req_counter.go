// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package door

import (
	"github.com/qinchende/gofast/logx"
	"time"
)

type (
	ReqItem struct {
		Duration  time.Duration // 任务耗时
		RouterIdx int16         // 路由节点的index
		Drop      bool          // 是否是一个丢弃的任务
	}

	reqItems struct {
		items    []ReqItem
		duration time.Duration
		drops    int
	}

	reqContainer struct {
		//name     string
		pid      int
		items    []ReqItem
		duration time.Duration
		drops    int
	}
)

// 添加统计项目
func (rc *reqContainer) AddItem(v interface{}) bool {
	if item, ok := v.(ReqItem); ok {
		if item.Drop {
			rc.drops++
		} else {
			rc.items = append(rc.items, item)
			rc.duration += item.Duration
		}
	}
	return false
}

// 执行
func (rc *reqContainer) Execute(items interface{}) {
	ret := items.(reqItems)
	//items := pair.items
	//duration := pair.duration
	drops := ret.drops
	size := len(ret.items)
	report := &PrintInfo{
		Name:          "Door.PrintInfo",
		Timestamp:     time.Now().Unix(),
		Pid:           rc.pid,
		ReqsPerSecond: float32(size) / float32(LogInterval/time.Second),
		Drops:         drops,
	}
	if size > 0 {
		//report.PerDur = (ret.duration / time.Millisecond) / size
		report.Path = Keeper.getPath(ret.items[0].RouterIdx)
	}

	log(report)
}

// 返回当前容器中的所有数据
func (rc *reqContainer) RemoveAll() interface{} {
	items := rc.items
	duration := rc.duration
	drops := rc.drops
	rc.items = nil
	rc.duration = 0
	rc.drops = 0

	return reqItems{
		items:    items,
		duration: duration,
		drops:    drops,
	}
}

func log(report *PrintInfo) {
	// writeReport(report)
	logx.Statf("(%s) | %s - qps: %.1f/s", report.Name, report.Path, report.ReqsPerSecond)
}
