// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/aid/bag"
	"github.com/qinchende/gofast/core/cst"
	"net/http"
	"time"
)

type Record struct {
	Time  time.Duration
	App   string
	Host  string
	Label string
	Msg   string
}

type ReqRecord struct {
	Record
	RawReq     *http.Request
	TimeStamp  time.Duration
	Latency    time.Duration
	ClientIP   string
	StatusCode int
	Pms        cst.SuperKV
	BodySize   int
	ResData    []byte
	CarryItems bag.CarryList
}

type Field struct {
	Key string
	Val any
}

type ZapField struct {
	Key       string
	Type      uint8
	Integer   int64
	String    string
	Interface interface{}
}
