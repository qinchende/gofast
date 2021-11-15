// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package check

import "time"

var LogInterval = time.Minute

type (
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
