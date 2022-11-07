package syncx

import (
	"sync/atomic"
)

type LazyCounter struct {
	Max  int32
	Curr int32
}

func (ct *LazyCounter) TryBorrow() bool {
	if ct.Curr > ct.Max {
		return false
	}
	atomic.AddInt32(&ct.Curr, 1)
	return true
}

func (ct *LazyCounter) Return() error {
	atomic.AddInt32(&ct.Curr, -1)
	if ct.Curr < 0 {
		atomic.AddInt32(&ct.Curr, 1)
		return ErrCounterEmpty
	}
	return nil
}
