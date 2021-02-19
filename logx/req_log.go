// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"fmt"
	"net/http"
	"time"
)

type LogReqParams struct {
	Request      *http.Request
	TimeStamp    time.Time
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	Method       string
	Path         string
	ErrorMessage string
	isTerm       bool
	BodySize     int
	Keys         map[string]interface{}
}

var GenReqLogString = func(p *LogReqParams) string {
	formatStr := `
[%s] %s (%s/%s) [%d]
  B: %s C: %s
  P: %s
  R: %s
  E: %s
`
	return fmt.Sprintf(formatStr,
		p.Method,
		p.Path,
		p.ClientIP,
		p.TimeStamp.Format("01-02 15:04:05"),
		p.Latency/time.Millisecond,
		"",
		"",
		"",
		"",
		p.ErrorMessage,
	)
}

func WriteReqLog(p *LogReqParams) {
	writeNow(GenReqLogString(p))
}
