// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
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
		Name() string              // Name returns the name of the Breaker.
		Errors(join string) string // reason

		Allow() error         //
		Accept()              // allow successful.
		Reject(reason string) // allow failed.

		AcceptValue(v float64)                // allow successful.
		RejectValue(v float64, reason string) // allow failed.

		//Do(req funcReq) error
		//DoWithAcceptable(req funcReq, cpt funcAcceptable) error
		//DoWithFallback(req funcReq, fb funcFallback) error
		//DoWithFallbackAcceptable(req funcReq, fb funcFallback, cpt funcAcceptable) error
	}

	throttle interface {
		allow() error
		doReq(req funcReq, fb funcFallback, cpt funcAcceptable) error
		markValue(val float64)
	}
)
