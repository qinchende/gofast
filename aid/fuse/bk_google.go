// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fuse

import (
	"github.com/qinchende/gofast/aid/randx"
)

type autoBreaker struct {
	throttle
	name    string
	errWin  *errorWindow
	markLog bool
}

func NewGBreaker(name string, markLog bool) Breaker {
	b := autoBreaker{
		name:    name,
		errWin:  new(errorWindow),
		markLog: markLog,
	}
	if len(b.name) == 0 {
		b.name = randx.Rand()
	}
	b.throttle = newGoogleThrottle()
	return &b
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ab *autoBreaker) Name() string {
	return ab.name
}

func (ab *autoBreaker) Errors(join string) string {
	return ab.errWin.ErrorsJoin(join)
}

func (ab *autoBreaker) Accept() {
	ab.throttle.markValue(1)
}

func (ab *autoBreaker) Reject(reason string) {
	if ab.markLog {
		ab.errWin.add(reason)
	}
	ab.throttle.markValue(0)
}

func (ab *autoBreaker) Allow() error {
	return ab.throttle.allow()
}

func (ab *autoBreaker) AcceptValue(v float64) {
	ab.throttle.markValue(v)
}

func (ab *autoBreaker) RejectValue(v float64, reason string) {
	if ab.markLog {
		ab.errWin.add(reason)
	}
	ab.throttle.markValue(v)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (ab *autoBreaker) Do(req funcReq) error {
//	return ab.logError(ab.throttle.doReq(req, nil, defAcceptable))
//}
//
//func (ab *autoBreaker) DoWithAcceptable(req funcReq, cpt funcAcceptable) error {
//	return ab.logError(ab.throttle.doReq(req, nil, cpt))
//}
//
//func (ab *autoBreaker) DoWithFallback(req funcReq, fb funcFallback) error {
//	return ab.logError(ab.throttle.doReq(req, fb, defAcceptable))
//}
//
//func (ab *autoBreaker) DoWithFallbackAcceptable(req funcReq, fb funcFallback, cpt funcAcceptable) error {
//	return ab.logError(ab.throttle.doReq(req, fb, cpt))
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// utils
//func defAcceptable(err error) bool {
//	return err == nil
//}
