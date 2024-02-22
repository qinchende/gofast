// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/store/jde"
	"time"
)

// 待实现
type logPrometheusEntry struct {
	Timestamp string `pms:"@timestamp"`
	Level     string `pms:"lv"`
	Duration  string `pms:"duration"`
	Content   any    `pms:"ct"`
}

func outputPrometheusStyle(w WriterCloser, logLevel string, data any) {
	logWrap := logPrometheusEntry{
		Timestamp: time.Now().Format(timeFormat),
		Level:     logLevel,
		Content:   data,
	}
	if content, err := jde.EncodeToBytes(logWrap); err != nil {
		outputDirectString(w, err.Error())
	} else {
		outputDirectBytes(w, content)
	}
}

func buildPrometheusReqLog(p *ReqLogEntity, flag int8) string {
	return ""
}
