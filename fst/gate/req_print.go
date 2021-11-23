// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/logx"
	"time"
)

// 统计的间隔时间
var LogInterval = time.Minute

type (
	FuncGetPath func(id uint16) string // 获取当前请求对应的路径

	ReqItem struct {
		RouterIdx uint16        // 当前请求对应路由树节点的index
		Duration  time.Duration // 请求耗时
		Drop      bool          // 是否是一个被丢弃的请求（熔断或者资源超限拒绝处理）
	}

	reqItems struct {
		items    []ReqItem
		duration time.Duration
		drops    int
	}

	reqContainer struct {
		getPath  FuncGetPath
		name     string
		pid      int
		duration time.Duration
		items    []ReqItem
		drops    int
	}

	PrintInfo struct {
		Name          string  `json:"name"`
		Path          string  `json:"path"`
		Timestamp     int64   `json:"tm"`
		Pid           int     `json:"pid"`
		PerDur        int     `json:"ptime"`
		ReqsPerSecond float32 `json:"qps"`
		Drops         int     `json:"drops"`
		Average       float32 `json:"avg"`
		Median        float32 `json:"med"`
		Top90th       float32 `json:"t90"`
		Top99th       float32 `json:"t99"`
		Top99p9th     float32 `json:"t99p9"`
	}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 添加统计项目
// 如果这里返回true，意味着要立刻刷新当前所有统计数据，这个开关用户自定义输出日志
func (rc *reqContainer) AddItem(v interface{}) bool {
	if item, ok := v.(ReqItem); ok {
		if item.Drop {
			rc.drops++
		} else {
			rc.items = append(rc.items, item)
			rc.duration += item.Duration
		}
	}
	return true
}

// 返回当前容器中的所有数据，同时清空容器
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

// 执行，
func (rc *reqContainer) Execute(items interface{}) {
	ret := items.(reqItems)
	//items := ret.items
	duration := ret.duration
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
		report.Average = float32(duration/time.Millisecond) / (float32(size))
		report.Path = rc.getPath(ret.items[0].RouterIdx)
	}

	log(report)
}

func log(report *PrintInfo) {
	// writeReport(report)
	logx.Statf("(%s) | %s - qps: %.1f/s", report.Name, report.Path, report.ReqsPerSecond)
}
