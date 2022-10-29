package fuse

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/proc"
)

const (
	numHistoryReasons = 5
	timeFormat        = "15:04:05"
)

var ErrServiceUnavailable = errors.New("auto-breaker is open")

type autoBreaker struct {
	name string
	throttle
}

// NewBreaker returns a Breaker object. opts can be used to customize the Breaker.
func NewBreaker(name string) Breaker {
	b := autoBreaker{name: name}
	if len(b.name) == 0 {
		b.name = lang.Rand()
	}
	b.throttle = newLoggedThrottle(b.name, newGoogleBreaker())

	return &b
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (cb *autoBreaker) Allow() (Promise, error) {
	return cb.throttle.allow()
}

func (cb *autoBreaker) Do(req func() error) error {
	return cb.throttle.doReq(req, nil, defaultAcceptable)
}

func (cb *autoBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return cb.throttle.doReq(req, nil, acceptable)
}

func (cb *autoBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return cb.throttle.doReq(req, fallback, defaultAcceptable)
}

func (cb *autoBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error,
	acceptable Acceptable) error {
	return cb.throttle.doReq(req, fallback, acceptable)
}

func (cb *autoBreaker) Name() string {
	return cb.name
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func defaultAcceptable(err error) bool {
	return err == nil
}

type loggedThrottle struct {
	name string
	bkThrottle
	errWin *errorWindow
}

func newLoggedThrottle(name string, t bkThrottle) loggedThrottle {
	return loggedThrottle{
		name:       name,
		bkThrottle: t,
		errWin:     new(errorWindow),
	}
}

func (lt loggedThrottle) allow() (Promise, error) {
	promise, err := lt.bkThrottle.allow()
	return promiseWithReason{
		promise: promise,
		errWin:  lt.errWin,
	}, lt.logError(err)
}

func (lt loggedThrottle) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	return lt.logError(lt.bkThrottle.doReq(req, fallback, func(err error) bool {
		accept := acceptable(err)
		if !accept {
			lt.errWin.add(err.Error())
		}
		return accept
	}))
}

func (lt loggedThrottle) logError(err error) error {
	if err == ErrServiceUnavailable {
		// if circuit open, not possible to have empty error window
		Report(fmt.Sprintf(
			"proc(%s/%d), callee: %s, breaker is open and requests dropped\nlast errors:\n%s",
			proc.ProcessName(), proc.Pid(), lt.name, lt.errWin))
	}

	return err
}
