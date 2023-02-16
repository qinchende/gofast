// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"time"
)

// 待实现
type logElkEntry struct {
	Timestamp string `json:"@timestamp"`
	Level     string `json:"lv"`
	Duration  string `json:"duration,omitempty"`
	Content   any    `json:"ct"`
}

func outputElkStyle(w WriterCloser, logLevel string, data any) {
	logWrap := logElkEntry{
		Timestamp: time.Now().Format(timeFormat),
		Level:     logLevel,
		Content:   data,
	}
	if content, err := jsonx.Marshal(logWrap); err != nil {
		outputDirectString(w, err.Error())
	} else {
		outputDirectBytes(w, content)
	}
}

func buildElkReqLog(p *ReqLogEntity, flag int8) string {
	return ""
}
