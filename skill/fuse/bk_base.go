package fuse

import "errors"

var ErrServiceUnavailable = errors.New("auto-breaker is open")

type (
	// Acceptable is the func to check if the error can be accepted.
	funcAcceptable func(err error) bool
	funcFallback   func(err error) error
	funcReq        func() error

	// A Breaker represents a circuit breaker.
	Breaker interface {
		// Name returns the name of the Breaker.
		Name() string
		// Allow checks if the request is allowed.
		// If allowed, a promise will be returned, the caller needs to call promise.Accept()
		// on success, or call promise.Reject() on failure.
		// If not allow, ErrServiceUnavailable will be returned.
		Allow() error

		Accept()              // allow successful.
		Reject(reason string) // allow failed.

		Do(req funcReq) error
		DoWithAcceptable(req funcReq, cpt funcAcceptable) error
		DoWithFallback(req funcReq, fb funcFallback) error
		DoWithFallbackAcceptable(req funcReq, fb funcFallback, cpt funcAcceptable) error
	}

	throttle interface {
		allow() error
		doReq(req funcReq, fb funcFallback, cpt funcAcceptable) error
		markSuc()
		markFai()
	}
)
