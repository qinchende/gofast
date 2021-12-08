// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/logx"
	"time"
)

// 固定分钟作为统计周期
var LogInterval = time.Minute

type (
	FuncGetPath func(id uint16) string // 获取当前请求对应的路径
	// 每个请求需要消耗 2个字长 16字节的空间
	ReqItem struct {
		LossTime time.Duration // 单次请求耗时
		RouteIdx uint16        // 当前请求对应路由树节点的index，这用来单独统计不同route
		Drop     bool          // 是否是一个被丢弃的请求（熔断或者资源超限拒绝处理）
	}

	reqItems struct {
		items    []ReqItem
		duration time.Duration
		drops    int
	}

	// 存放所有请求的处理时间，作为统计的容器
	reqContainer struct {
		getPath  FuncGetPath
		name     string
		pid      int
		duration time.Duration // 本容器中所有请求的总耗时
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
			rc.duration += item.LossTime
		}
	}
	return true
}

// 返回当前容器中的所有数据，同时重置容器
func (rc *reqContainer) RemoveAll() interface{} {
	ret := reqItems{
		items:    rc.items,
		duration: rc.duration,
		drops:    rc.drops,
	}
	rc.items = nil
	rc.duration = 0
	rc.drops = 0

	return ret
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
		report.Path = rc.getPath(ret.items[0].RouteIdx)
	}

	log(report)
}

func log(report *PrintInfo) {
	// writeReport(report)
	logx.Statf("(%s) | %s - qps: %.1f/s", report.Name, report.Path, report.ReqsPerSecond)
}

//func (c *metricsContainer) Execute(v interface{}) {
//	pair := v.(tasksDurationPair)
//	tasks := pair.tasks
//	duration := pair.duration
//	drops := pair.drops
//	size := len(tasks)
//	report := &StatReport{
//		Name:          c.name,
//		Timestamp:     time.Now().Unix(),
//		Pid:           c.pid,
//		ReqsPerSecond: float32(size) / float32(logInterval/time.Second),
//		Drops:         drops,
//	}
//
//	if size > 0 {
//		report.Average = float32(duration/time.Millisecond) / float32(size)
//
//		fiftyPercent := size >> 1
//		if fiftyPercent > 0 {
//			top50pTasks := topK(tasks, fiftyPercent)
//			medianTask := top50pTasks[0]
//			report.Median = float32(medianTask.Duration) / float32(time.Millisecond)
//			tenPercent := fiftyPercent / 5
//			if tenPercent > 0 {
//				top10pTasks := topK(tasks, tenPercent)
//				task90th := top10pTasks[0]
//				report.Top90th = float32(task90th.Duration) / float32(time.Millisecond)
//				onePercent := tenPercent / 10
//				if onePercent > 0 {
//					top1pTasks := topK(top10pTasks, onePercent)
//					task99th := top1pTasks[0]
//					report.Top99th = float32(task99th.Duration) / float32(time.Millisecond)
//					pointOnePercent := onePercent / 10
//					if pointOnePercent > 0 {
//						topPointOneTasks := topK(top1pTasks, pointOnePercent)
//						task99Point9th := topPointOneTasks[0]
//						report.Top99p9th = float32(task99Point9th.Duration) / float32(time.Millisecond)
//					} else {
//						report.Top99p9th = getTopDuration(top1pTasks)
//					}
//				} else {
//					mostDuration := getTopDuration(top10pTasks)
//					report.Top99th = mostDuration
//					report.Top99p9th = mostDuration
//				}
//			} else {
//				mostDuration := getTopDuration(tasks)
//				report.Top90th = mostDuration
//				report.Top99th = mostDuration
//				report.Top99p9th = mostDuration
//			}
//		} else {
//			mostDuration := getTopDuration(tasks)
//			report.Median = mostDuration
//			report.Top90th = mostDuration
//			report.Top99th = mostDuration
//			report.Top99p9th = mostDuration
//		}
//	}
//
//	log(report)
//}

//func log(report *StatReport) {
//	writeReport(report)
//	if logEnabled.True() {
//		logx.Statf("(%s) - qps: %.1f/s, drops: %d, avg time: %.1fms, med: %.1fms, "+
//			"90th: %.1fms, 99th: %.1fms, 99.9th: %.1fms",
//			report.Name, report.ReqsPerSecond, report.Drops, report.Average, report.Median,
//			report.Top90th, report.Top99th, report.Top99p9th)
//	}
//}
