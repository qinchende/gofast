package fuse

import (
	"github.com/qinchende/gofast/skill/collect"
	"github.com/qinchende/gofast/skill/syncx"
	"time"
)

type (
	// A Promise interface is returned by Shedder.Allow to let callers tell
	// whether the processing request is successful or not.
	Promise interface {
		// Pass lets the caller tell that the call is successful.
		Pass()
		// Fail lets the caller tell that the call is failed.
		Fail()
	}

	// Shedder is the interface that wraps the Allow method.
	Shedder interface {
		// Allow returns the Promise if allowed, otherwise ErrServiceOverloaded.
		Allow() (Promise, error)
	}

	// ShedderOption lets caller customize the Shedder.
	ShedderOption func(opts *shedderOptions)

	shedderOptions struct {
		window       time.Duration
		buckets      int
		cpuThreshold float64
	}

	adaptiveShedder struct {
		cpuThreshold    float64
		windows         int64
		flying          int64
		avgFlying       float64
		avgFlyingLock   syncx.SpinLock
		dropTime        *syncx.AtomicDuration
		droppedRecently *syncx.AtomicBool
		passCounter     *collect.RollingWindow
		rtCounter       *collect.RollingWindow
	}
)
