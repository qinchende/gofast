// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/timex"
)

// 待实现
type logPrometheusEntry struct {
	Timestamp string `json:"@timestamp"`
	Level     string `json:"lv"`
	Duration  string `json:"duration,omitempty"`
	Content   string `json:"ct"`
}

func outputPrometheusStyle(w WriterCloser, info, logLevel string) {
	logWrap := logPrometheusEntry{
		Timestamp: timex.Time().Format(timeFormat),
		Level:     logLevel,
		Content:   info,
	}
	if content, err := jsonx.Marshal(logWrap); err != nil {
		outputDirectString(w, err.Error())
	} else {
		outputDirectBytes(w, content)
	}
}

func buildPrometheusReqLog(p *ReqLogEntity, flag int8) string {
	return ""
}
