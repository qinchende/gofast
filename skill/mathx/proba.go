package mathx

import (
	"math/rand"
	"sync"
	"time"
)

type Maybe struct {
	// rand.New(...) returns a non thread safe object
	r    *rand.Rand
	lock sync.Mutex
}

func NewMaybe() *Maybe {
	return &Maybe{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// 一个随机数小于pb
func (p *Maybe) TrueOnMaybe(pb float64) (truth bool) {
	p.lock.Lock()
	truth = p.r.Float64() < pb
	p.lock.Unlock()
	return
}
