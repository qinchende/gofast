// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

// 如果采用Custom模式，必须外部指定下面这两个方法
var CustomOutputFunc func(info, logLevel string) string
var CustomReqLogFunc func(p *ReqLogEntity, flag int8) string

func outputCustomStyle(w WriterCloser, info, logLevel string) {
	outputDirectString(w, CustomOutputFunc(info, logLevel))
}

func buildCustomReqLog(p *ReqLogEntity, flag int8) string {
	return CustomReqLogFunc(p, flag)
}
