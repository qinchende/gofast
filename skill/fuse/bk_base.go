package fuse

import "errors"

var ErrServiceUnavailable = errors.New("auto-breaker is open")

type (
	// Acceptable is the func to check if the error can be accepted.
	Acceptable   func(err error) bool
	funcFallback func(err error) error
	funcReq      func() error

	// A Breaker represents a circuit breaker.
	Breaker interface {
		// Name returns the name of the Breaker.
		Name() string

		// Allow checks if the request is allowed.
		// If allowed, a promise will be returned, the caller needs to call promise.Accept()
		// on success, or call promise.Reject() on failure.
		// If not allow, ErrServiceUnavailable will be returned.
		Allow() (Promise, error)

		Do(req funcReq) error
		DoWithAcceptable(req funcReq, acceptable Acceptable) error
		DoWithFallback(req funcReq, fallback funcFallback) error
		DoWithFallbackAcceptable(req funcReq, fallback funcFallback, acceptable Acceptable) error
	}

	Promise interface {
		Accept()              // allow successful.
		Reject(reason string) // allow failed.
	}

	throttle interface {
		allow() (Promise, error)
		doReq(req funcReq, fallback funcFallback, acceptable Acceptable) error
	}

	//// 来一个简版的接口 ++++++++++++++++++++++++++++++++
	//promiseNoReason interface {
	//	Accept()
	//	Reject()
	//}
	//
	//throttleNoReason interface {
	//	allow() (promiseNoReason, error)
	//	doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error
	//}
)
