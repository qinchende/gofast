// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/aid/bag"
	"github.com/qinchende/gofast/core/cst"
	"net/http"
	"time"
)

type Field struct {
	Key string
	Val any
}

//
//type ZapField struct {
//	Key       string
//	Type      uint8
//	Integer   int64
//	String    string
//	Interface interface{}
//}

//// Event represents a log event. It is instanced by one of the level method of
//// Logger and finalized by the Msg or Msgf method.
//type Event struct {
//	buf       []byte
//	w         LevelWriter
//	level     Level
//	done      func(msg string)
//	stack     bool            // enable error stack trace
//	ch        []Hook          // hooks from context
//	skipFrame int             // The number of additional frames to skip when printing the caller.
//	ctx       context.Context // Optional Go context for event
//}

type Event struct {
	bf *[]byte

	Time  time.Duration
	App   string
	Host  string
	Label string
	Msg   string
}

type ReqRecord struct {
	RawReq     *http.Request
	StatusCode int
	Method     string
	RequestURI string
	UserAgent  string
	RemoteAddr string

	Event
	TimeStamp  time.Duration
	Latency    time.Duration
	Pms        cst.SuperKV
	BodySize   int
	ResData    []byte
	CarryItems bag.CarryList
}

//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// Str adds the field key with val as a string to the *Event context.
//func (e *Event) Str(key, val string) *Event {
//	*e.bf = append(*e.bf, key...,  ":"..., val...)
//	return e
//}
//
//// Strs adds the field key with vals as a []string to the *Event context.
//func (e *Event) Strs(key string, vals []string) *Event {
//	e.buf = enc.AppendStrings(enc.AppendKey(e.buf, key), vals)
//	return e
//}
//
//// Stringer adds the field key with val.String() (or null if val is nil)
//// to the *Event context.
//func (e *Event) Stringer(key string, val fmt.Stringer) *Event {
//	e.buf = enc.AppendStringer(enc.AppendKey(e.buf, key), val)
//	return e
//}
