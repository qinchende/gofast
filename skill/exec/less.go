package exec

import (
	"time"

	"github.com/qinchende/gofast/skill/syncx"
	"github.com/qinchende/gofast/skill/timex"
)

type Less struct {
	threshold time.Duration
	lastTime  *syncx.AtomicDuration
}

func NewLess(threshold time.Duration) *Less {
	return &Less{
		threshold: threshold,
		lastTime:  syncx.NewAtomicDuration(),
	}
}

func (le *Less) DoOrDiscard(execute func()) bool {
	now := timex.Now()
	lastTime := le.lastTime.Load()
	if lastTime == 0 || lastTime+le.threshold < now {
		le.lastTime.Set(now)
		execute()
		return true
	}

	return false
}
