// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"time"
)

//type logSdxJsonEntry struct {
//	T string
//	L string
//	D any
//}
type logSdxJsonEntry [3]any

func outputSdxJsonStyle(w WriterCloser, logLevel string, data any) {
	logWrap := logSdxJsonEntry{
		time.Now().Format(timeFormatMini),
		logLevel,
		data,
	}
	if content, err := jsonx.Marshal(logWrap); err != nil {
		outputDirectString(w, err.Error())
	} else {
		outputDirectBytes(w, content)
	}
}

func buildSdxJsonReqLog(p *ReqLogEntity, flag int8) string {
	return ""
}
