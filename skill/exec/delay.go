package exec

import (
	"github.com/qinchende/gofast/skill/gmp"
	"sync"
	"time"
)

type Delay struct {
	fn        func()
	delay     time.Duration
	triggered bool
	lock      sync.Mutex
}

func NewDelay(fn func(), delay time.Duration) *Delay {
	return &Delay{
		fn:    fn,
		delay: delay,
	}
}

func (de *Delay) Trigger() {
	de.lock.Lock()
	defer de.lock.Unlock()

	if de.triggered {
		return
	}

	de.triggered = true
	gmp.GoSafe(func() {
		timer := time.NewTimer(de.delay)
		defer timer.Stop()
		<-timer.C

		// set triggered to false before calling fn to ensure no triggers are missed.
		de.lock.Lock()
		de.triggered = false
		de.lock.Unlock()
		de.fn()
	})
}
