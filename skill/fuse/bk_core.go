package fuse

import (
	"github.com/qinchende/gofast/skill/lang"
)

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
	//b.throttle = newLoggedThrottle(b.name, newGoogleThrottle())
	b.throttle = newGoogleThrottle(b.name)

	return &b
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (cb *autoBreaker) Name() string {
	return cb.name
}

func (cb *autoBreaker) Allow() (Promise, error) {
	return cb.throttle.allow()
}

func (cb *autoBreaker) Do(req funcReq) error {
	return cb.throttle.doReq(req, nil, defAcceptable)
}

func (cb *autoBreaker) DoWithAcceptable(req funcReq, acceptable Acceptable) error {
	return cb.throttle.doReq(req, nil, acceptable)
}

func (cb *autoBreaker) DoWithFallback(req funcReq, fallback funcFallback) error {
	return cb.throttle.doReq(req, fallback, defAcceptable)
}

func (cb *autoBreaker) DoWithFallbackAcceptable(req funcReq, fallback funcFallback, acceptable Acceptable) error {
	return cb.throttle.doReq(req, fallback, acceptable)
}

func defAcceptable(err error) bool {
	return err == nil
}
