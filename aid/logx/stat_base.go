// Copyright 2022 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
package logx

type LogStat struct {
	Type   uint8
	Fields []string
}

// 我们把系统收集统计信息日志的一些常量维护在这里
var (
	LogStatSysMonitor  = &LogStat{Type: 1, Fields: []string{"cpu", "mem", "gor", "gc"}}
	LogStatRouteReq    = &LogStat{Type: 2, Fields: []string{"accept", "timeout", "drop", "qps", "ave", "max"}}
	LogStatCpuUsage    = &LogStat{Type: 3, Fields: []string{"cpu", "total", "pass", "drop"}}
	LogStatBreakerOpen = &LogStat{Type: 4, Fields: []string{}}
)
